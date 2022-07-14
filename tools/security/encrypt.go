package security

import (
	"crypto/aes"
	"crypto/cipher"
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"strings"
)

// S256Challenge creates base64 encoded sha256 challenge string derived from code.
// The padding of the result base64 string is stripped per [RFC 7636].
//
// [RFC 7636]: https://datatracker.ietf.org/doc/html/rfc7636#section-4.2
func S256Challenge(code string) string {
	h := sha256.New()
	h.Write([]byte(code))
	return strings.TrimRight(base64.URLEncoding.EncodeToString(h.Sum(nil)), "=")
}

// Encrypt encrypts data with key (must be valid 32 char aes key).
func Encrypt(data []byte, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())

	// populates the nonce with a cryptographically secure random sequence
	if _, err := io.ReadFull(crand.Reader, nonce); err != nil {
		return "", err
	}

	cipherByte := gcm.Seal(nonce, nonce, data, nil)

	result := base64.StdEncoding.EncodeToString(cipherByte)

	return result, nil
}

// Decrypt decrypts encrypted text with key (must be valid 32 chars aes key).
func Decrypt(cipherText string, key string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()

	cipherByte, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}

	nonce, cipherByteClean := cipherByte[:nonceSize], cipherByte[nonceSize:]
	return gcm.Open(nil, nonce, cipherByteClean, nil)
}
