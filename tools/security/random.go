package security

import (
	cryptoRand "crypto/rand"
	mathRand "math/rand/v2"
)

const defaultRandomAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// RandomString generates a cryptographically random string with the specified length.
//
// The generated string matches [A-Za-z0-9]+ and it's transparent to URL-encoding.
func RandomString(length int) string {
	return RandomStringWithAlphabet(length, defaultRandomAlphabet)
}

// RandomStringWithAlphabet generates a cryptographically random string
// with the specified length and characters set.
//
// It panics if for some reason rand.Read returns a non-nil error.
func RandomStringWithAlphabet(length int, alphabet string) string {
	b := make([]byte, length)
	alphaLen := byte(len(alphabet))

	// Compute the largest multiple of alphaLen that fits in a byte
	// to ensure uniform distribution via rejection sampling.
	maxValid := 256 - (256 % int(alphaLen))

	// Read random bytes in bulk instead of one crypto/rand.Int per character.
	// Use rejection sampling to maintain uniform distribution.
	randomBytes := make([]byte, length)
	if _, err := cryptoRand.Read(randomBytes); err != nil {
		panic(err)
	}

	for i := 0; i < length; i++ {
		// Rejection sampling: if the random byte falls in the biased range,
		// request a new one until it doesn't.
		for int(randomBytes[i]) >= maxValid {
			if _, err := cryptoRand.Read(randomBytes[i : i+1]); err != nil {
				panic(err)
			}
		}
		b[i] = alphabet[randomBytes[i]%alphaLen]
	}

	return string(b)
}

// PseudorandomString generates a pseudorandom string with the specified length.
//
// The generated string matches [A-Za-z0-9]+ and it's transparent to URL-encoding.
//
// For a cryptographically random string (but a little bit slower) use RandomString instead.
func PseudorandomString(length int) string {
	return PseudorandomStringWithAlphabet(length, defaultRandomAlphabet)
}

// PseudorandomStringWithAlphabet generates a pseudorandom string
// with the specified length and characters set.
//
// For a cryptographically random (but a little bit slower) use RandomStringWithAlphabet instead.
func PseudorandomStringWithAlphabet(length int, alphabet string) string {
	b := make([]byte, length)
	max := len(alphabet)

	for i := range b {
		b[i] = alphabet[mathRand.IntN(max)]
	}

	return string(b)
}
