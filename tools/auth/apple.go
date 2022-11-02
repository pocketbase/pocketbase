package auth

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pocketbase/pocketbase/tools/security"
	"golang.org/x/oauth2"
)

var _ Provider = (*Apple)(nil)

// NameApple is the unique name of the Apple provider.
const NameApple string = "apple"

// Apple allows authentication via AzureADEndpoint OAuth2.
type Apple struct {
	*baseProvider
}

// NewAppleProvider creates new Apple AD provider instance with some defaults.
func NewAppleProvider() *Apple {
	return &Apple{&baseProvider{
		scopes:   []string{},
		authUrl:  "https://appleid.apple.com/auth/authorize",
		tokenUrl: "https://appleid.apple.com/auth/token",
	}}
}

func (p *Apple) BuildAuthUrl(state string, opts ...oauth2.AuthCodeOption) string {
	// opts = append(opts, oauth2.SetAuthURLParam("response_mode", "form_post"))
	return p.oauth2Config().AuthCodeURL(state, opts...)
}

// FetchAuthUser returns an AuthUser instance based on the Apple's user api.
func (p *Apple) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("missing id_token")
	}

	claims, err := security.ParseUnverifiedJWT(idToken)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%#v\n", claims)

	user := &AuthUser{
		Email: claims["email"].(string),
		Id:    claims["sub"].(string),
	}

	return user, nil
}

// GenerateClientSecret returns a client secret to be used with apple ouath2 provider
func GenerateAppleClientSecret(signingKey, teamID, clientID, keyID string) (string, error) {
	block, _ := pem.Decode([]byte(signingKey))
	if block == nil {
		return "", errors.New("empty block after decoding")
	}

	privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	jwt.MarshalSingleStringAsArray = false

	now := time.Now()
	claims := &jwt.RegisteredClaims{
		Issuer:    teamID,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour*24*180 - time.Second)), // 180 days
		Audience:  []string{"https://appleid.apple.com"},
		Subject:   clientID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["alg"] = "ES256"
	token.Header["kid"] = keyID

	return token.SignedString(privKey)
}
