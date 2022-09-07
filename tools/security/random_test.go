package security_test

import (
	"regexp"
	"testing"

	"github.com/pocketbase/pocketbase/tools/security"
)

func TestRandomString(t *testing.T) {
	generated := []string{}
	reg := regexp.MustCompile(`[a-zA-Z0-9]+`)
	length := 10

	for i := 0; i < 100; i++ {
		result := security.RandomString(length)

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

func TestRandomStringWithAlphabet(t *testing.T) {
	scenarios := []struct {
		alphabet      string
		expectPattern string
	}{
		{"0123456789_", `[0-9_]+`},
		{"abcd", `[abcd]+`},
		{"!@#$%^&*()", `[\!\@\#\$\%\^\&\*\(\)]+`},
	}

	for i, s := range scenarios {
		generated := make([]string, 100)
		length := 10

		for j := 0; j < 100; j++ {
			result := security.RandomStringWithAlphabet(length, s.alphabet)

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
