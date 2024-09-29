package security_test

import (
	"fmt"
	"testing"

	"github.com/pocketbase/pocketbase/tools/security"
)

func TestEncrypt(t *testing.T) {
	scenarios := []struct {
		data        string
		key         string
		expectError bool
	}{
		{"", "", true},
		{"123", "test", true}, // key must be valid 32 char aes string
		{"123", "abcdabcdabcdabcdabcdabcdabcdabcd", false},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%s", i, s.data), func(t *testing.T) {
			result, err := security.Encrypt([]byte(s.data), s.key)

			hasErr := err != nil

			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr %v, got %v (%v)", s.expectError, hasErr, err)
			}

			if hasErr {
				if result != "" {
					t.Fatalf("Expected empty Encrypt result on error, got %q", result)
				}

				return
			}

			// try to decrypt
			decrypted, err := security.Decrypt(result, s.key)
			if err != nil || string(decrypted) != s.data {
				t.Fatalf("Expected decrypted value to match with the data input, got %q (%v)", decrypted, err)
			}
		})
	}
}

func TestDecrypt(t *testing.T) {
	scenarios := []struct {
		cipher       string
		key          string
		expectError  bool
		expectedData string
	}{
		{"", "", true, ""},
		{"123", "test", true, ""}, // key must be valid 32 char aes string
		{"8kcEqilvvYKYcfnSr0aSC54gmnQCsB02SaB8ATlnA==", "abcdabcdabcdabcdabcdabcdabcdabcd", true, ""}, // illegal base64 encoded cipherText
		{"8kcEqilvv+YKYcfnSr0aSC54gmnQCsB02SaB8ATlnA==", "abcdabcdabcdabcdabcdabcdabcdabcd", false, "123"},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%s", i, s.key), func(t *testing.T) {
			result, err := security.Decrypt(s.cipher, s.key)

			hasErr := err != nil

			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr %v, got %v (%v)", s.expectError, hasErr, err)
			}

			if hasErr {
				return
			}

			if str := string(result); str != s.expectedData {
				t.Fatalf("Expected %q, got %q", s.expectedData, str)
			}
		})
	}
}
