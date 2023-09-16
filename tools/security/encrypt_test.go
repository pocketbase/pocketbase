package security_test

import (
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

	for i, scenario := range scenarios {
		result, err := security.Encrypt([]byte(scenario.data), scenario.key)

		if scenario.expectError && err == nil {
			t.Errorf("(%d) Expected error got nil", i)
		}
		if !scenario.expectError && err != nil {
			t.Errorf("(%d) Expected nil got error %v", i, err)
		}

		if scenario.expectError && result != "" {
			t.Errorf("(%d) Expected empty string, got %q", i, result)
		}
		if !scenario.expectError && result == "" {
			t.Errorf("(%d) Expected non empty encrypted result string", i)
		}

		// try to decrypt
		if result != "" {
			decrypted, _ := security.Decrypt(result, scenario.key)
			if string(decrypted) != scenario.data {
				t.Errorf("(%d) Expected decrypted value to match with the data input, got %q", i, decrypted)
			}
		}
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

	for i, scenario := range scenarios {
		result, err := security.Decrypt(scenario.cipher, scenario.key)

		if scenario.expectError && err == nil {
			t.Errorf("(%d) Expected error got nil", i)
		}
		if !scenario.expectError && err != nil {
			t.Errorf("(%d) Expected nil got error %v", i, err)
		}

		resultStr := string(result)
		if resultStr != scenario.expectedData {
			t.Errorf("(%d) Expected %q, got %q", i, scenario.expectedData, resultStr)
		}
	}
}
