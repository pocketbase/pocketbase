package s3

import (
	"net/url"
	"testing"
)

func TestEscapePath(t *testing.T) {
	t.Parallel()

	escaped := escapePath("/ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_.~ !@#$%^&*()+={}[]?><\\|,`'\"/@sub1/@sub2/a/b/c/1/2/3")

	expected := "/ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_.~%20%21%40%23%24%25%5E%26%2A%28%29%2B%3D%7B%7D%5B%5D%3F%3E%3C%5C%7C%2C%60%27%22/%40sub1/%40sub2/a/b/c/1/2/3"

	if escaped != expected {
		t.Fatalf("Expected\n%s\ngot\n%s", expected, escaped)
	}
}

func TestEscapeQuery(t *testing.T) {
	t.Parallel()

	escaped := escapeQuery(url.Values{
		"abc": []string{"123"},
		"/ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_.~ !@#$%^&*()+={}[]?><\\|,`'\"": []string{
			"/ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_.~ !@#$%^&*()+={}[]?><\\|,`'\"",
		},
	})

	expected := "%2FABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_.~%20%21%40%23%24%25%5E%26%2A%28%29%2B%3D%7B%7D%5B%5D%3F%3E%3C%5C%7C%2C%60%27%22=%2FABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_.~%20%21%40%23%24%25%5E%26%2A%28%29%2B%3D%7B%7D%5B%5D%3F%3E%3C%5C%7C%2C%60%27%22&abc=123"

	if escaped != expected {
		t.Fatalf("Expected\n%s\ngot\n%s", expected, escaped)
	}
}
