package validators_test

import (
	"fmt"
	"testing"

	"github.com/pocketbase/pocketbase/core/validators"
)

func Equal(t *testing.T) {
	t.Parallel()

	strA := "abc"
	strB := "abc"
	strC := "123"
	var strNilPtr *string
	var strNilPtr2 *string

	scenarios := []struct {
		valA        any
		valB        any
		expectError bool
	}{
		{nil, nil, false},
		{"", "", false},
		{"", "456", true},
		{"123", "", true},
		{"123", "456", true},
		{"123", "123", false},
		{true, false, true},
		{false, true, true},
		{false, false, false},
		{true, true, false},
		{0, 0, false},
		{0, 1, true},
		{1, 2, true},
		{1, 1, false},
		{&strA, &strA, false},
		{&strA, &strB, false},
		{&strA, &strC, true},
		{"abc", &strA, false},
		{&strA, "abc", false},
		{"abc", &strC, true},
		{"test", 123, true},
		{nil, 123, true},
		{nil, strA, true},
		{nil, &strA, true},
		{nil, strNilPtr, false},
		{strNilPtr, strNilPtr2, false},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%v_%v", i, s.valA, s.valB), func(t *testing.T) {
			err := validators.Equal(s.valA)(s.valB)

			hasErr := err != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr to be %v, got %v (%v)", s.expectError, hasErr, err)
			}
		})
	}
}
