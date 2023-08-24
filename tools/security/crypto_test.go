package security_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/tools/security"
)

func TestS256Challenge(t *testing.T) {
	scenarios := []struct {
		code     string
		expected string
	}{
		{"", "47DEQpj8HBSa-_TImW-5JCeuQeRkm5NMpJWZG3hSuFU"},
		{"123", "pmWkWSBCL51Bfkhn79xPuKBKHz__H6B-mY6G9_eieuM"},
	}

	for _, s := range scenarios {
		t.Run(s.code, func(t *testing.T) {
			result := security.S256Challenge(s.code)

			if result != s.expected {
				t.Fatalf("Expected %q, got %q", s.expected, result)
			}
		})
	}
}

func TestMD5(t *testing.T) {
	scenarios := []struct {
		code     string
		expected string
	}{
		{"", "d41d8cd98f00b204e9800998ecf8427e"},
		{"123", "202cb962ac59075b964b07152d234b70"},
	}

	for _, s := range scenarios {
		t.Run(s.code, func(t *testing.T) {
			result := security.MD5(s.code)

			if result != s.expected {
				t.Fatalf("Expected %v, got %v", s.expected, result)
			}
		})
	}
}

func TestSHA256(t *testing.T) {
	scenarios := []struct {
		code     string
		expected string
	}{
		{"", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{"123", "a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3"},
	}

	for _, s := range scenarios {
		t.Run(s.code, func(t *testing.T) {
			result := security.SHA256(s.code)

			if result != s.expected {
				t.Fatalf("Expected %v, got %v", s.expected, result)
			}
		})
	}
}

func TestSHA512(t *testing.T) {
	scenarios := []struct {
		code     string
		expected string
	}{
		{"", "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e"},
		{"123", "3c9909afec25354d551dae21590bb26e38d53f2173b8d3dc3eee4c047e7ab1c1eb8b85103e3be7ba613b31bb5c9c36214dc9f14a42fd7a2fdb84856bca5c44c2"},
	}

	for _, s := range scenarios {
		t.Run(s.code, func(t *testing.T) {
			result := security.SHA512(s.code)

			if result != s.expected {
				t.Fatalf("Expected %v, got %v", s.expected, result)
			}
		})
	}
}
