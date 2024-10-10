package validators_test

import (
	"fmt"
	"testing"

	"github.com/pocketbase/pocketbase/core/validators"
)

func TestIsRegex(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		val         string
		expectError bool
	}{
		{"", false},
		{`abc`, false},
		{`\w+`, false},
		{`\w*((abc+`, true},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#v", i, s.val), func(t *testing.T) {
			err := validators.IsRegex(s.val)

			hasErr := err != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr to be %v, got %v (%v)", s.expectError, hasErr, err)
			}
		})
	}
}
