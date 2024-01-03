package forms_test

import (
	"archive/zip"
	"bytes"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

func TestBackupUploadValidateAndSubmit(t *testing.T) {
	t.Parallel()

	var zb bytes.Buffer
	zw := zip.NewWriter(&zb)
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}

	f0, _ := filesystem.NewFileFromBytes([]byte("test"), "existing")
	f1, _ := filesystem.NewFileFromBytes([]byte("456"), "nozip")
	f2, _ := filesystem.NewFileFromBytes(zb.Bytes(), "existing")
	f3, _ := filesystem.NewFileFromBytes(zb.Bytes(), "zip")

	scenarios := []struct {
		name           string
		file           *filesystem.File
		expectedErrors []string
	}{
		{
			"missing file",
			nil,
			[]string{"file"},
		},
		{
			"non-zip file",
			f1,
			[]string{"file"},
		},
		{
			"zip file with non-unique name",
			f2,
			[]string{"file"},
		},
		{
			"zip file with unique name",
			f3,
			[]string{},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			app, _ := tests.NewTestApp()
			defer app.Cleanup()

			fsys, err := app.NewBackupsFilesystem()
			if err != nil {
				t.Fatal(err)
			}
			defer fsys.Close()
			// create a dummy backup file to simulate existing backups
			if err := fsys.UploadFile(f0, f0.OriginalName); err != nil {
				t.Fatal(err)
			}

			form := forms.NewBackupUpload(app)
			form.File = s.file

			result := form.Submit()

			// parse errors
			errs, ok := result.(validation.Errors)
			if !ok && result != nil {
				t.Fatalf("Failed to parse errors %v", result)
			}

			// check errors
			if len(errs) > len(s.expectedErrors) {
				t.Fatalf("Expected error keys %v, got %v", s.expectedErrors, errs)
			}
			for _, k := range s.expectedErrors {
				if _, ok := errs[k]; !ok {
					t.Fatalf("Missing expected error key %q in %v", k, errs)
				}
			}

			expectedFiles := []*filesystem.File{f0}
			if result == nil {
				expectedFiles = append(expectedFiles, s.file)
			}

			// retrieve all uploaded backup files
			files, err := fsys.List("")
			if err != nil {
				t.Fatal("Failed to retrieve backup files")
			}

			if len(files) != len(expectedFiles) {
				t.Fatalf("Expected %d files, got %d", len(expectedFiles), len(files))
			}

			for _, ef := range expectedFiles {
				exists := false
				for _, f := range files {
					if f.Key == ef.OriginalName {
						exists = true
						break
					}
				}
				if !exists {
					t.Fatalf("Missing expected backup file %v", ef.OriginalName)
				}
			}
		})
	}
}
