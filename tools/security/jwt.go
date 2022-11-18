package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// ParseUnverifiedJWT parses JWT token and returns its claims
// but DOES NOT verify the signature.
//
// It verifies only the exp, iat and nbf claims.
func ParseUnverifiedJWT(token string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}

	parser := &jwt.Parser{}
	_, _, err := parser.ParseUnverified(token, claims)

	if err == nil {
		err = claims.Valid()
	}

	return claims, err
}

// ParseJWT verifies and parses JWT token and returns its claims.
func ParseJWT(token string, verificationKey string) (jwt.MapClaims, error) {
	parser := jwt.NewParser(jwt.WithValidMethods([]string{"HS256"}))

	parsedToken, err := parser.Parse(token, func(t *jwt.Token) (any, error) {
		return []byte(verificationKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return claims, nil
	}

	return nil, errors.New("Unable to parse token.")
}

// NewToken generates and returns new HS256 signed JWT token.
func NewToken(payload jwt.MapClaims, signingKey string, secondsDuration int64) (string, error) {
	seconds := time.Duration(secondsDuration) * time.Second

	claims := jwt.MapClaims{
		"exp": time.Now().Add(seconds).Unix(),
	}

	for k, v := range payload {
		claims[k] = v
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(signingKey))
}
