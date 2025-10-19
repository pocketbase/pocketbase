// Package jwk implements some common utilities for interacting with JWKs
// (mostly used with OIDC providers).
package jwk

import (
	"context"
	"crypto/ed25519"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type JWK struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	// RS256 (RSA)
	E string `json:"e"`
	N string `json:"n"`
	// Ed25519 (OKP)
	Crv string `json:"crv"`
	X   string `json:"x"`
}

// PublicKey reconstructs and returns the public key from the current JWK.
func (key *JWK) PublicKey() (any, error) {
	switch key.Kty {
	case "RSA":
		// RFC 7518
		// https://datatracker.ietf.org/doc/html/rfc7518#section-6.3
		// https://datatracker.ietf.org/doc/html/rfc7517#appendix-A.1
		exponent, err := base64.RawURLEncoding.DecodeString(strings.TrimRight(key.E, "="))
		if err != nil {
			return nil, err
		}

		modulus, err := base64.RawURLEncoding.DecodeString(strings.TrimRight(key.N, "="))
		if err != nil {
			return nil, err
		}

		return &rsa.PublicKey{
			E: int(big.NewInt(0).SetBytes(exponent).Uint64()),
			N: big.NewInt(0).SetBytes(modulus),
		}, nil
	case "OKP":
		// RFC 8037
		// https://datatracker.ietf.org/doc/html/rfc8037#section-2
		// https://datatracker.ietf.org/doc/html/rfc8037#appendix-A
		if key.Crv != "Ed25519" {
			return nil, fmt.Errorf("unsupported OKP curve (must be Ed25519): %q", key.Crv)
		}

		x, err := base64.RawURLEncoding.DecodeString(strings.TrimRight(key.X, "="))
		if err != nil {
			return nil, err
		}

		if l := len(x); l != ed25519.PublicKeySize {
			return nil, fmt.Errorf("invalid Ed25519 key length: %d", l)
		}

		return ed25519.PublicKey(x), nil
	default:
		return nil, fmt.Errorf("unsupported kty (must be RSA or OKP): %q", key.Kty)
	}
}

// Fetch retrieves the JSON Web Key Set located at jwksURL and returns
// the first key that matches the specified kid.
func Fetch(ctx context.Context, jwksURL string, kid string) (*JWK, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", jwksURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// http.Client.Get doesn't treat non 2xx responses as error
	if res.StatusCode >= 400 {
		return nil, fmt.Errorf(
			"failed to fetch JSON Web Key Set from %s (%d):\n%s",
			jwksURL,
			res.StatusCode,
			string(rawBody),
		)
	}

	jwks := struct {
		Keys []*JWK
	}{}

	err = json.Unmarshal(rawBody, &jwks)
	if err != nil {
		return nil, err
	}

	for _, key := range jwks.Keys {
		if key.Kid == kid {
			return key, nil
		}
	}

	return nil, fmt.Errorf("JWK with kid %q was not found", kid)
}

// ValidateTokenSignature validates the signature of a token with the
// public key retrievied from a remote JWKS.
func ValidateTokenSignature(ctx context.Context, token string, jwksURL string) error {
	// extract the kid token header
	// ---
	t, _, err := jwt.NewParser().ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		return err
	}

	kid, _ := t.Header["kid"].(string)
	if kid == "" {
		return errors.New("missing kid header value")
	}

	// fetch the public key set
	// ---
	key, err := Fetch(ctx, jwksURL, kid)
	if err != nil {
		return err
	}

	// verify the signature
	// ---
	parser := jwt.NewParser(jwt.WithValidMethods([]string{key.Alg}))

	parsedToken, err := parser.Parse(token, func(t *jwt.Token) (any, error) {
		return key.PublicKey()
	})
	if err != nil {
		return err
	}

	if !parsedToken.Valid {
		return errors.New("the parsed token is invalid")
	}

	return nil
}
