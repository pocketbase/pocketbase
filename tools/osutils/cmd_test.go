package osutils_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/tools/osutils"
)

func TestYesNoPrompt(t *testing.T) {
	scenarios := []struct {
		stdin    string
		fallback bool
		expected bool
	}{
		{"", false, false},
		{"", true, true},

		// yes
		{"y", false, true},
		{"Y", false, true},
		{"Yes", false, true},
		{"yes", false, true},

		// no
		{"n", true, false},
		{"N", true, false},
		{"No", true, false},
		{"no", true, false},

		// invalid -> no/yes
		{"invalid|no", true, false},
		{"invalid|yes", false, true},
	}

	for _, s := range scenarios {
		t.Run(fmt.Sprintf("%s_%v", s.stdin, s.fallback), func(t *testing.T) {
			stdinread, stdinwrite, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}

			parts := strings.Split(s.stdin, "|")
			for _, p := range parts {
				if _, err := stdinwrite.WriteString(p + "\n"); err != nil {
					t.Fatalf("Failed to write test stdin part %q: %v", p, err)
				}
			}

			if err = stdinwrite.Close(); err != nil {
				t.Fatal(err)
			}

			defer func(oldStdin *os.File) { os.Stdin = oldStdin }(os.Stdin)
			os.Stdin = stdinread

			result := osutils.YesNoPrompt("test", s.fallback)

			if result != s.expected {
				t.Fatalf("Expected %v, got %v", s.expected, result)
			}
		})
	}
}
