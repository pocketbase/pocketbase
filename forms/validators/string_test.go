package validators_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/forms/validators"
)

func TestCompare(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		valA        string
		valB        string
		expectError bool
	}{
		{"", "", false},
		{"", "456", true},
		{"123", "", true},
		{"123", "456", true},
		{"123", "123", false},
	}

	for i, s := range scenarios {
		err := validators.Compare(s.valA)(s.valB)

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("(%d) Expected hasErr to be %v, got %v (%v)", i, s.expectError, hasErr, err)
		}
	}
}
