package tokenizer

import (
	"io"
	"strings"
	"testing"
)

func TestFactories(t *testing.T) {
	expectedContent := "test"

	scenarios := []struct {
		name string
		tk   *Tokenizer
	}{
		{
			"New()",
			New(strings.NewReader(expectedContent)),
		},
		{
			"NewFromString()",
			NewFromString(expectedContent),
		},
		{
			"NewFromBytes()",
			NewFromBytes([]byte(expectedContent)),
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			content, _ := s.tk.r.ReadString(0)

			if content != expectedContent {
				t.Fatalf("Expected reader with content %q, got %q", expectedContent, content)
			}

			if s.tk.keepSeparator != false {
				t.Fatal("Expected keepSeparator false, got true")
			}

			if s.tk.ignoreParenthesis != false {
				t.Fatal("Expected ignoreParenthesis false, got true")
			}

			if len(s.tk.separators) != len(DefaultSeparators) {
				t.Fatalf("Expected \n%v, \ngot \n%v", DefaultSeparators, s.tk.separators)
			}

			for _, r := range s.tk.separators {
				exists := false
				for _, def := range s.tk.separators {
					if r == def {
						exists = true
						break
					}
				}
				if !exists {
					t.Fatalf("Unexpected sepator %s", string(r))
				}
			}
		})
	}
}

func TestScan(t *testing.T) {
	tk := NewFromString("abc, 123.456, (abc)")

	expectedTokens := []string{"abc", "123.456", "(abc)"}

	for _, token := range expectedTokens {
		result, err := tk.Scan()
		if err != nil {
			t.Fatalf("Expected token %q, got error %v", token, err)
		}

		if result != token {
			t.Fatalf("Expected token %q, got error %v", token, result)
		}
	}

	// scan the last character
	token, err := tk.Scan()
	if err != io.EOF {
		t.Fatalf("Expected EOF error, got %v", err)
	}
	if token != "" || err != io.EOF {
		t.Fatalf("Expected empty token, got %q", token)
	}
}

func TestScanAll(t *testing.T) {
	scenarios := []struct {
		name              string
		content           string
		separators        []rune
		keepSeparator     bool
		keepEmptyTokens   bool
		ignoreParenthesis bool
		expectError       bool
		expectTokens      []string
	}{
		{
			name:              "empty string",
			content:           "",
			separators:        DefaultSeparators,
			keepSeparator:     false,
			keepEmptyTokens:   false,
			ignoreParenthesis: false,
			expectError:       false,
			expectTokens:      nil,
		},
		{
			name:              "unbalanced parenthesis",
			content:           `(a,b() c`,
			separators:        DefaultSeparators,
			keepSeparator:     false,
			keepEmptyTokens:   false,
			ignoreParenthesis: false,
			expectError:       true,
			expectTokens:      []string{},
		},
		{
			name:              "unmatching quotes",
			content:           `'asd"`,
			separators:        DefaultSeparators,
			keepSeparator:     false,
			keepEmptyTokens:   false,
			ignoreParenthesis: false,
			expectError:       true,
			expectTokens:      []string{},
		},
		{
			name:              "no separators",
			content:           `a, b, c, d, e 123, "abc"`,
			separators:        nil,
			keepSeparator:     false,
			keepEmptyTokens:   false,
			ignoreParenthesis: false,
			expectError:       false,
			expectTokens:      []string{`a, b, c, d, e 123, "abc"`},
		},
		{
			name: "default separators",
			content: `a, b , c  , d e  , "a,b,  c  " , ,, ,	  (123, 456)
			`,
			separators:        DefaultSeparators,
			keepSeparator:     false,
			keepEmptyTokens:   false,
			ignoreParenthesis: false,
			expectError:       false,
			expectTokens: []string{
				"a",
				"b",
				"c",
				"d e",
				`"a,b,  c  "`,
				`(123, 456)`,
			},
		},
		{
			name:              "keep separators",
			content:           `a, b, c, d  e, "a,b,  c  ",	(123, 456)`,
			separators:        []rune{',', ' '}, // the space should be removed from the cutset
			keepSeparator:     true,
			keepEmptyTokens:   true,
			ignoreParenthesis: false,
			expectError:       false,
			expectTokens: []string{
				"a,",
				" ",
				"b,",
				" ",
				"c,",
				" ",
				"d ",
				" ",
				"e,",
				" ",
				`"a,b,  c  ",`,
				`(123, 456)`,
			},
		},
		{
			name:              "custom separators",
			content:           `a | b c  d &(e + f) &  "g & h" & & &`,
			separators:        []rune{'|', '&'},
			keepSeparator:     false,
			keepEmptyTokens:   false,
			ignoreParenthesis: false,
			expectError:       false,
			expectTokens: []string{
				"a",
				"b c  d",
				"(e + f)",
				`"g & h"`,
			},
		},
		{
			name:              "ignoring parenthesis",
			content:           `a, b, (c,d)`,
			separators:        DefaultSeparators,
			keepSeparator:     false,
			keepEmptyTokens:   false,
			ignoreParenthesis: true,
			expectError:       false,
			expectTokens: []string{
				"a",
				"b",
				"(c",
				"d)",
			},
		},
		{
			name:              "keep empty tokens",
			content:           `a, b, (c, d), ,, , e, , f`,
			separators:        DefaultSeparators,
			keepSeparator:     false,
			keepEmptyTokens:   true,
			ignoreParenthesis: false,
			expectError:       false,
			expectTokens: []string{
				"a",
				"b",
				"(c, d)",
				"",
				"",
				"",
				"e",
				"",
				"f",
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			tk := NewFromString(s.content)

			tk.Separators(s.separators...)
			tk.KeepSeparator(s.keepSeparator)
			tk.KeepEmptyTokens(s.keepEmptyTokens)
			tk.IgnoreParenthesis(s.ignoreParenthesis)

			tokens, err := tk.ScanAll()

			hasErr := err != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr %v, got %v (%v)", s.expectError, hasErr, err)
			}

			if len(tokens) != len(s.expectTokens) {
				t.Fatalf("Expected \n%v (%d), \ngot \n%v (%d)", s.expectTokens, len(s.expectTokens), tokens, len(tokens))
			}

			for _, tok := range tokens {
				exists := false
				for _, def := range s.expectTokens {
					if tok == def {
						exists = true
						break
					}
				}
				if !exists {
					t.Fatalf("Unexpected token %q", tok)
				}
			}
		})
	}
}

func TestTrimCutset(t *testing.T) {
	scenarios := []struct {
		name           string
		separators     []rune
		expectedCutset string
	}{
		{
			"default factory separators",
			nil,
			"\t\n\v\f\r \u0085\u00a0",
		},
		{
			"custom separators",
			[]rune{'\t', ' ', '\r', ','},
			"\n\v\f\u0085\u00a0",
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			tk := NewFromString("")

			if len(s.separators) > 0 {
				tk.Separators(s.separators...)
			}

			if tk.trimCutset != s.expectedCutset {
				t.Fatalf("Expected cutset %q, got %q", s.expectedCutset, tk.trimCutset)
			}
		})
	}
}
