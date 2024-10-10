package tests

import (
	"errors"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// TestValidationErrors checks whether the provided rawErrors are
// instance of [validation.Errors] and contains the expectedErrors keys.
func TestValidationErrors(t *testing.T, rawErrors error, expectedErrors []string) {
	var errs validation.Errors

	if rawErrors != nil && !errors.As(rawErrors, &errs) {
		t.Fatalf("Failed to parse errors, expected to find validation.Errors, got %T\n%v", rawErrors, rawErrors)
	}

	if len(errs) != len(expectedErrors) {
		keys := make([]string, 0, len(errs))
		for k := range errs {
			keys = append(keys, k)
		}
		t.Fatalf("Expected error keys \n%v\ngot\n%v\n%v", expectedErrors, keys, errs)
	}

	for _, k := range expectedErrors {
		if _, ok := errs[k]; !ok {
			t.Fatalf("Missing expected error key %q in %v", k, errs)
		}
	}
}
