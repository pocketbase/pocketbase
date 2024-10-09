package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

var _ Provider = (*Zoho)(nil)

// NameZoho is the unique name of the Zoho provider.
const NameZoho string = "zoho"

// Zoho allows authentication via Zoho OAuth2.
type Zoho struct {
	*baseProvider
}

// NewZohoProvider creates new Zoho provider instance with some defaults.
// Zoho claims to support OIDC, but their userinfo endpoint (userApiUrl) https://accounts.zoho.com/oauth/v2/userinfo is not working
func NewZohoProvider() *Zoho {
	return &Zoho{&baseProvider{
		ctx:         context.Background(),
		displayName: "Zoho",
		pkce:        false,
		scopes:      []string{"profile", "email", "openid"},
		authUrl:     "https://accounts.zoho.com/oauth/v2/auth",
		tokenUrl:    "https://accounts.zoho.com/oauth/v2/token",
	}}
}

// FetchAuthUser returns an AuthUser instance based on JWT id_token.
//
// API reference:
// https://www.zoho.com/accounts/protocol/oauth/sign-in-using-zoho.html
func (p *Zoho) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	// Zoho's OIDC implementation is not a standard one, because
	// it uses idToken instead of userApiUrl for user information.

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok || rawIDToken == "" {
		return nil, errors.New("missing id_token")
	}

	claims, err := p.parseIDToken(rawIDToken)
	if err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           claims["sub"].(string),
		Username:     claims["email"].(string), // Zoho doesn't provide username field, using email instead.
		Name:         claims["name"].(string),
		Email:        claims["email"].(string),
		AvatarUrl:    claims["picture"].(string),
		RawUser:      claims,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
	if claims["email_verified"] == true {
		user.Email = claims["email"].(string)
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	return user, nil
}

// -------------------------------------------------------------------

func (p *Zoho) parseIDToken(rawIDToken string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}

	_, _, err := new(jwt.Parser).ParseUnverified(rawIDToken, &claims)
	if err != nil {
		return nil, fmt.Errorf("failed to parse id_token: %w", err)
	}

	// verify issuer
	if !claims.VerifyIssuer("accounts.zoho.com", true) {
		return nil, errors.New("invalid id_token issuer")
	}

	// verify audience
	if !claims.VerifyAudience(p.clientId, true) {
		return nil, errors.New("invalid id_token audience")
	}

	return claims, nil
}
