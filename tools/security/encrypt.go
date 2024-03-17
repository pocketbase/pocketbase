package security

import (
	"crypto/aes"
	"crypto/cipher"
	crand "crypto/rand"
	"encoding/base64"
	"io"
)

// Encrypt encrypts "data" with the specified "key" (must be valid 32 char AES key).
//
// This method uses AES-256-GCM block cypher mode.
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

// Decrypt decrypts encrypted text with key (must be valid 32 chars AES key).
//
// This method uses AES-256-GCM block cypher mode.
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
