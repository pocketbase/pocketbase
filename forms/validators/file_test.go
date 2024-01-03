package validators_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pocketbase/pocketbase/forms/validators"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/pocketbase/pocketbase/tools/rest"
)

func TestUploadedFileSize(t *testing.T) {
	t.Parallel()

	data, mp, err := tests.MockMultipartData(nil, "test")
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/", data)
	req.Header.Add("Content-Type", mp.FormDataContentType())

	files, err := rest.FindUploadedFiles(req, "test")
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 1 {
		t.Fatalf("Expected one test file, got %d", len(files))
	}

	scenarios := []struct {
		maxBytes    int
		file        *filesystem.File
		expectError bool
	}{
		{0, nil, false},
		{4, nil, false},
		{3, files[0], true}, // all test files have "test" as content
		{4, files[0], false},
		{5, files[0], false},
	}

	for i, s := range scenarios {
		err := validators.UploadedFileSize(s.maxBytes)(s.file)

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("(%d) Expected hasErr to be %v, got %v (%v)", i, s.expectError, hasErr, err)
		}
	}
}

func TestUploadedFileMimeType(t *testing.T) {
	t.Parallel()

	data, mp, err := tests.MockMultipartData(nil, "test")
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/", data)
	req.Header.Add("Content-Type", mp.FormDataContentType())

	files, err := rest.FindUploadedFiles(req, "test")
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 1 {
		t.Fatalf("Expected one test file, got %d", len(files))
	}

	scenarios := []struct {
		types       []string
		file        *filesystem.File
		expectError bool
	}{
		{nil, nil, false},
		{[]string{"image/jpeg"}, nil, false},
		{[]string{}, files[0], true},
		{[]string{"image/jpeg"}, files[0], true},
		// test files are detected as "text/plain; charset=utf-8" content type
		{[]string{"image/jpeg", "text/plain; charset=utf-8"}, files[0], false},
	}

	for i, s := range scenarios {
		err := validators.UploadedFileMimeType(s.types)(s.file)

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("(%d) Expected hasErr to be %v, got %v (%v)", i, s.expectError, hasErr, err)
		}
	}
}
