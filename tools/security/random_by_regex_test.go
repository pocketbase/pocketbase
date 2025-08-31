package security_test

import (
	"fmt"
	"regexp"
	"regexp/syntax"
	"slices"
	"testing"

	"github.com/pocketbase/pocketbase/tools/security"
)

func TestRandomStringByRegex(t *testing.T) {
	generated := []string{}

	scenarios := []struct {
		pattern     string
		flags       []syntax.Flags
		expectError bool
	}{
		{``, nil, true},
		{`test`, nil, false},
		{`\d+`, []syntax.Flags{syntax.POSIX}, true},
		{`\d+`, nil, false},
		{`\d*`, nil, false},
		{`\d{1,20}`, nil, false},
		{`\d{5}`, nil, false},
		{`\d{0,}-abc`, nil, false},
		{`[a-zA-Z_]*`, nil, false},
		{`[^a-zA-Z]{5,30}`, nil, false},
		{`\w+_abc`, nil, false},
		{`[2-9]{10}-\w+`, nil, false},
		{`(a|b|c)`, nil, false},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%q", i, s.pattern), func(t *testing.T) {
			str, err := security.RandomStringByRegex(s.pattern, s.flags...)

			hasErr := err != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr %v, got %v (%v)", s.expectError, hasErr, err)
			}

			if hasErr {
				return
			}

			r, err := regexp.Compile(s.pattern)
			if err != nil {
				t.Fatal(err)
			}

			if !r.Match([]byte(str)) {
				t.Fatalf("Expected %q to match pattern %v", str, s.pattern)
			}

			if slices.Contains(generated, str) {
				t.Fatalf("The generated string %q already exists in\n%v", str, generated)
			}

			generated = append(generated, str)
		})
	}
}
