package forms_test

import (
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/tests"
)

func TestS3FilesystemValidate(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

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
		form := forms.NewTestS3Filesystem(app)
		form.Filesystem = s.filesystem

		result := form.Validate()

		// parse errors
		errs, ok := result.(validation.Errors)
		if !ok && result != nil {
			t.Errorf("[%s] Failed to parse errors %v", s.name, result)
			continue
		}

		// check errors
		if len(errs) > len(s.expectedErrors) {
			t.Errorf("[%s] Expected error keys %v, got %v", s.name, s.expectedErrors, errs)
			continue
		}
		for _, k := range s.expectedErrors {
			if _, ok := errs[k]; !ok {
				t.Errorf("[%s] Missing expected error key %q in %v", s.name, k, errs)
			}
		}
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
