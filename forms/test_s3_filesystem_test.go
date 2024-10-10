package forms_test

import (
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/tests"
)

func TestS3FilesystemValidate(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		name           string
		filesystem     string
		expectedErrors []string
	}{
		{
			"empty filesystem",
			"",
			[]string{"filesystem"},
		},
		{
			"invalid filesystem",
			"something",
			[]string{"filesystem"},
		},
		{
			"backups filesystem",
			"backups",
			[]string{},
		},
		{
			"storage filesystem",
			"storage",
			[]string{},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			app, _ := tests.NewTestApp()
			defer app.Cleanup()

			form := forms.NewTestS3Filesystem(app)
			form.Filesystem = s.filesystem

			result := form.Validate()

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
		})
	}
}

func TestS3FilesystemSubmitFailure(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// check if validate was called
	{
		form := forms.NewTestS3Filesystem(app)
		form.Filesystem = ""

		result := form.Submit()

		if result == nil {
			t.Fatal("Expected error, got nil")
		}

		if _, ok := result.(validation.Errors); !ok {
			t.Fatalf("Expected validation.Error, got %v", result)
		}
	}

	// check with valid storage and disabled s3
	{
		form := forms.NewTestS3Filesystem(app)
		form.Filesystem = "storage"

		result := form.Submit()

		if result == nil {
			t.Fatal("Expected error, got nil")
		}

		if _, ok := result.(validation.Error); ok {
			t.Fatalf("Didn't expect validation.Error, got %v", result)
		}
	}
}
