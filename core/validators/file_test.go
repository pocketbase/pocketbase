package validators_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/core/validators"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

func TestUploadedFileSize(t *testing.T) {
	t.Parallel()

	file, err := filesystem.NewFileFromBytes([]byte("test"), "test.txt")
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		maxBytes    int64
		file        *filesystem.File
		expectError bool
	}{
		{0, nil, false},
		{4, nil, false},
		{3, file, true}, // all test files have "test" as content
		{4, file, false},
		{5, file, false},
	}

	for _, s := range scenarios {
		t.Run(fmt.Sprintf("%d", s.maxBytes), func(t *testing.T) {
			err := validators.UploadedFileSize(s.maxBytes)(s.file)

			hasErr := err != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr to be %v, got %v (%v)", s.expectError, hasErr, err)
			}
		})
	}
}

func TestUploadedFileMimeType(t *testing.T) {
	t.Parallel()

	file, err := filesystem.NewFileFromBytes([]byte("test"), "test.png") // the extension shouldn't matter
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		types       []string
		file        *filesystem.File
		expectError bool
	}{
		{nil, nil, false},
		{[]string{"image/jpeg"}, nil, false},
		{[]string{}, file, true},
		{[]string{"image/jpeg"}, file, true},
		// test files are detected as "text/plain; charset=utf-8" content type
		{[]string{"image/jpeg", "text/plain; charset=utf-8"}, file, false},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%s", i, strings.Join(s.types, ";")), func(t *testing.T) {
			err := validators.UploadedFileMimeType(s.types)(s.file)

			hasErr := err != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr to be %v, got %v (%v)", s.expectError, hasErr, err)
			}
		})
	}
}
