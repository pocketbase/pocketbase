package security

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
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

// MD5 creates md5 hash from the provided plain text.
func MD5(text string) string {
	h := md5.New()
	h.Write([]byte(text))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// SHA256 creates sha256 hash as defined in FIPS 180-4 from the provided text.
func SHA256(text string) string {
	h := sha256.New()
	h.Write([]byte(text))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// SHA512 creates sha512 hash as defined in FIPS 180-4 from the provided text.
func SHA512(text string) string {
	h := sha512.New()
	h.Write([]byte(text))
	return fmt.Sprintf("%x", h.Sum(nil))
}
