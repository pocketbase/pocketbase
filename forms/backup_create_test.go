package forms_test

import (
	"strings"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/tests"
)

func TestBackupCreateValidateAndSubmit(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		name           string
		backupName     string
		expectedErrors []string
	}{
		{
			"invalid length",
			strings.Repeat("a", 97) + ".zip",
			[]string{"name"},
		},
		{
			"valid length + invalid format",
			strings.Repeat("a", 96),
			[]string{"name"},
		},
		{
			"valid length + valid format",
			strings.Repeat("a", 96) + ".zip",
			[]string{},
		},
		{
			"auto generated name",
			"",
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

			form := forms.NewBackupCreate(app)
			form.Name = s.backupName

			result := form.Submit()

			// parse errors
			errs, ok := result.(validation.Errors)
			if !ok && result != nil {
				t.Fatalf("Failed to parse errors %v", result)
				return
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

			// retrieve all created backup files
			files, err := fsys.List("")
			if err != nil {
				t.Fatal("Failed to retrieve backup files")
				return
			}

			if result != nil {
				if total := len(files); total != 0 {
					t.Fatalf("Didn't expected backup files, found %d", total)
				}
				return
			}

			if total := len(files); total != 1 {
				t.Fatalf("Expected 1 backup file, got %d", total)
				return
			}

			if s.backupName == "" {
				prefix := "pb_backup_"
				if !strings.HasPrefix(files[0].Key, prefix) {
					t.Fatalf("Expected the backup file, to have prefix %q: %q", prefix, files[0].Key)
				}
			} else if s.backupName != files[0].Key {
				t.Fatalf("Expected backup file %q, got %q", s.backupName, files[0].Key)
			}
		})
	}
}
