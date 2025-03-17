package mailer

import (
	"fmt"
	"io"
	"net/mail"
	"strings"
	"testing"
)

func TestAddressesToStrings(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		withName  bool
		addresses []mail.Address
		expected  []string
	}{
		{
			true,
			[]mail.Address{{Name: "John Doe", Address: "test1@example.com"}, {Name: "Jane Doe", Address: "test2@example.com"}},
			[]string{`"John Doe" <test1@example.com>`, `"Jane Doe" <test2@example.com>`},
		},
		{
			true,
			[]mail.Address{{Name: "John Doe", Address: "test1@example.com"}, {Address: "test2@example.com"}},
			[]string{`"John Doe" <test1@example.com>`, `test2@example.com`},
		},
		{
			false,
			[]mail.Address{{Name: "John Doe", Address: "test1@example.com"}, {Name: "Jane Doe", Address: "test2@example.com"}},
			[]string{`test1@example.com`, `test2@example.com`},
		},
	}

	for _, s := range scenarios {
		t.Run(fmt.Sprintf("%v_%v", s.withName, s.addresses), func(t *testing.T) {
			result := addressesToStrings(s.addresses, s.withName)

			if len(s.expected) != len(result) {
				t.Fatalf("Expected\n%v\ngot\n%v", s.expected, result)
			}

			for k, v := range s.expected {
				if v != result[k] {
					t.Fatalf("Expected %d address %q, got %q", k, v, result[k])
				}
			}
		})
	}
}

func TestDetectReaderMimeType(t *testing.T) {
	t.Parallel()

	str := "#!/bin/node\n" + strings.Repeat("a", 10000) // ensure that it is large enough to remain after the signature sniffing

	r, mime, err := detectReaderMimeType(strings.NewReader(str))
	if err != nil {
		t.Fatal(err)
	}

	expectedMime := "text/javascript"
	if mime != expectedMime {
		t.Fatalf("Expected mime %q, got %q", expectedMime, mime)
	}

	raw, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	rawStr := string(raw)

	if rawStr != str {
		t.Fatalf("Expected content\n%s\ngot\n%s", str, rawStr)
	}
}
