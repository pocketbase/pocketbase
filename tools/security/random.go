package security

import (
	cryptoRand "crypto/rand"
	"math/big"
	mathRand "math/rand"
	"time"
)

const defaultRandomAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func init() {
	mathRand.Seed(time.Now().UnixNano())
}

// RandomString generates a cryptographically random string with the specified length.
//
// The generated string matches [A-Za-z0-9]+ and it's transparent to URL-encoding.
func RandomString(length int) string {
	return RandomStringWithAlphabet(length, defaultRandomAlphabet)
}

// RandomStringWithAlphabet generates a cryptographically random string
// with the specified length and characters set.
//
// It panics if for some reason rand.Int returns a non-nil error.
func RandomStringWithAlphabet(length int, alphabet string) string {
	b := make([]byte, length)
	mx := big.NewInt(int64(len(alphabet)))

	for i := range b {
		n, err := cryptoRand.Int(cryptoRand.Reader, mx)
		if err != nil {
			panic(err)
		}
		b[i] = alphabet[n.Int64()]
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
	mx := len(alphabet)

	for i := range b {
		b[i] = alphabet[mathRand.Intn(mx)]
	}

	return string(b)
}
