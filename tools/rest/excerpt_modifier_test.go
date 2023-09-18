package rest

import (
	"fmt"
	"testing"

	"github.com/spf13/cast"
)

func TestNewExcerptModifier(t *testing.T) {
	scenarios := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			"no arguments",
			nil,
			true,
		},
		{
			"too many arguments",
			[]string{"12", "false", "something"},
			true,
		},
		{
			"non-numeric max argument",
			[]string{"something"}, // should fallback to 0 which is not allowed
			true,
		},
		{
			"numeric max argument",
			[]string{"12"},
			false,
		},
		{
			"non-bool withEllipsis argument",
			[]string{"12", "something"}, // should fallback to false which is allowed
			false,
		},
		{
			"truthy withEllipsis argument",
			[]string{"12", "t"},
			false,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			m, err := newExcerptModifier(s.args...)

			hasErr := err != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr %v, got %v (%v)", s.expectError, hasErr, err)
			}

			if hasErr {
				if m != nil {
					t.Fatalf("Expected nil modifier, got %v", m)
				}

				return
			}

			var argMax int
			if len(s.args) > 0 {
				argMax = cast.ToInt(s.args[0])
			}

			var argWithEllipsis bool
			if len(s.args) > 1 {
				argWithEllipsis = cast.ToBool(s.args[1])
			}

			if m.max != argMax {
				t.Fatalf("Expected max %d, got %d", argMax, m.max)
			}

			if m.withEllipsis != argWithEllipsis {
				t.Fatalf("Expected withEllipsis %v, got %v", argWithEllipsis, m.withEllipsis)
			}
		})
	}
}

func TestExcerptModifierModify(t *testing.T) {
	// plain text value: "Hello t est12 3 word"
	html := ` <script>var a = 123;</script>   <p>Hello</p><div id="test_id">t   est<b>12
	3</b></div> <h1>word  </h1> `

	plainText := "Hello t est12 3 word"

	scenarios := []struct {
		name     string
		args     []string
		value    string
		expected string
	}{
		// without ellipsis
		{
			"only max < len(plainText)",
			[]string{"2"},
			html,
			plainText[:2],
		},
		{
			"only max = len(plainText)",
			[]string{fmt.Sprint(len(plainText))},
			html,
			plainText,
		},
		{
			"only max > len(plainText)",
			[]string{fmt.Sprint(len(plainText) + 5)},
			html,
			plainText,
		},

		// with ellipsis
		{
			"with ellipsis and max < len(plainText)",
			[]string{"2", "t"},
			html,
			plainText[:2] + "...",
		},
		{
			"with ellipsis and max = len(plainText)",
			[]string{fmt.Sprint(len(plainText)), "t"},
			html,
			plainText,
		},
		{
			"with ellipsis and max > len(plainText)",
			[]string{fmt.Sprint(len(plainText) + 5), "t"},
			html,
			plainText,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			m, err := newExcerptModifier(s.args...)
			if err != nil {
				t.Fatal(err)
			}

			raw, err := m.Modify(s.value)
			if err != nil {
				t.Fatal(err)
			}

			if v := cast.ToString(raw); v != s.expected {
				t.Fatalf("Expected %q, got %q", s.expected, v)
			}
		})
	}
}
