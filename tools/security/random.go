package security

import (
	"crypto/rand"
)

// RandomString generates a random string with the specified length.
//
// The generated string is cryptographically random and matches
// [A-Za-z0-9]+ (aka. it's transparent to URL-encoding).
func RandomString(length int) string {
	const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	return RandomStringWithAlphabet(length, alphabet)
}

// RandomStringWithAlphabet generates a cryptographically random string
// with the specified length and characters set.
func RandomStringWithAlphabet(length int, alphabet string) string {
	bytes := make([]byte, length)

	rand.Read(bytes)

	for i, b := range bytes {
		bytes[i] = alphabet[b%byte(len(alphabet))]
	}

	return string(bytes)
}
