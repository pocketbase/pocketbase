package dbutils_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/tools/dbutils"
)

func TestAliasOrIdentifier(t *testing.T) {
	scenarios := []struct {
		value    string
		expected string
	}{
		{"", ""},
		{"abc", "abc"},
		{"abc  ", "abc  "}, // return unmodified
		{"abc.def", "abc.def"},
		{"abc.123 def", "def"},
		{"abc.123 as def.456", "def.456"},
		{"(abc) def", "def"},
		{"(abc) as def", "def"},
		{"abc   def", "def"},
		{"abc as   def", "def"},
		// technically invalid identifier but consistent with the dbx regex matching
		{"a b c d", "d"},
		{"a b c as d", "d"},
	}

	for _, s := range scenarios {
		t.Run(s.value, func(t *testing.T) {
			result := dbutils.AliasOrIdentifier(s.value)

			if result != s.expected {
				t.Fatalf("Expected\n%v\ngot\n%v", s.expected, result)
			}
		})
	}
}
