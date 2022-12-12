package security_test

import (
	"regexp"
	"testing"

	"github.com/pocketbase/pocketbase/tools/security"
)

func TestRandomString(t *testing.T) {
	testRandomString(t, security.RandomString)
}

func TestRandomStringWithAlphabet(t *testing.T) {
	testRandomStringWithAlphabet(t, security.RandomStringWithAlphabet)
}

func TestPseudorandomString(t *testing.T) {
	testRandomString(t, security.PseudorandomString)
}

func TestPseudorandomStringWithAlphabet(t *testing.T) {
	testRandomStringWithAlphabet(t, security.PseudorandomStringWithAlphabet)
}

// -------------------------------------------------------------------

func testRandomStringWithAlphabet(t *testing.T, randomFunc func(n int, alphabet string) string) {
	scenarios := []struct {
		alphabet      string
		expectPattern string
	}{
		{"0123456789_", `[0-9_]+`},
		{"abcdef123", `[abcdef123]+`},
		{"!@#$%^&*()", `[\!\@\#\$\%\^\&\*\(\)]+`},
	}

	for i, s := range scenarios {
		generated := make([]string, 0, 1000)
		length := 10

		for j := 0; j < 1000; j++ {
			result := randomFunc(length, s.alphabet)

			if len(result) != length {
				t.Fatalf("(%d:%d) Expected the length of the string to be %d, got %d", i, j, length, len(result))
			}

			reg := regexp.MustCompile(s.expectPattern)
			if match := reg.MatchString(result); !match {
				t.Fatalf("(%d:%d) The generated string should have only %s characters, got %q", i, j, s.expectPattern, result)
			}

			for _, str := range generated {
				if str == result {
					t.Fatalf("(%d:%d) Repeating random string - found %q in %q", i, j, result, generated)
				}
			}

			generated = append(generated, result)
		}
	}
}

func testRandomString(t *testing.T, randomFunc func(n int) string) {
	generated := make([]string, 0, 1000)
	reg := regexp.MustCompile(`[a-zA-Z0-9]+`)
	length := 10

	for i := 0; i < 1000; i++ {
		result := randomFunc(length)

		if len(result) != length {
			t.Fatalf("(%d) Expected the length of the string to be %d, got %d", i, length, len(result))
		}

		if match := reg.MatchString(result); !match {
			t.Fatalf("(%d) The generated string should have only [a-zA-Z0-9]+ characters, got %q", i, result)
		}

		for _, str := range generated {
			if str == result {
				t.Fatalf("(%d) Repeating random string - found %q in \n%v", i, result, generated)
			}
		}

		generated = append(generated, result)
	}
}
