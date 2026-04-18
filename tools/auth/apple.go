package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pocketbase/pocketbase/tools/auth/internal/jwk"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/spf13/cast"
	"golang.org/x/oauth2"
)

func init() {
	Providers[NameApple] = wrapFactory(NewAppleProvider)
}

var _ Provider = (*Apple)(nil)

// NameApple is the unique name of the Apple provider.
const NameApple string = "apple"

// Apple allows authentication via Apple OAuth2.
//
// OIDC differences: https://bitbucket.org/openid/connect/src/master/How-Sign-in-with-Apple-differs-from-OpenID-Connect.md.
type Apple struct {
	BaseProvider

	jwksURL string
}

// NewAppleProvider creates a new Apple provider instance with some defaults.
func NewAppleProvider() *Apple {
	return &Apple{
		BaseProvider: BaseProvider{
			ctx:         context.Background(),
			order:       1,
			logo:        `<svg xmlns="http://www.w3.org/2000/svg" width="256" height="315" preserveAspectRatio="xMidYMid"><path d="M213.8 167c.4 47.6 41.7 63.4 42.2 63.6-.3 1.2-6.6 22.6-21.8 44.8-13 19.1-26.7 38.2-48 38.6-21.1.4-28-12.5-52-12.5s-31.6 12.1-51.5 12.9c-20.7.8-36.4-20.7-49.6-39.8-27-39-47.7-110.3-20-158.4a77 77 0 0 1 65.1-39.4c20.3-.4 39.5 13.6 51.9 13.6s35.7-16.9 60.2-14.4c10.2.4 39 4.2 57.5 31.2-1.5 1-34.4 20-34 59.8M174.2 50.2A69 69 0 0 0 190.6 0c-15.8.6-35 10.5-46.3 23.8-10.2 11.8-19.1 30.6-16.7 48.7 17.6 1.3 35.7-9 46.6-22.3"/></svg>`,
			displayName: "Apple",
			pkce:        true,
			scopes:      []string{"name", "email"},
			authURL:     "https://appleid.apple.com/auth/authorize",
			tokenURL:    "https://appleid.apple.com/auth/token",
		},
		jwksURL: "https://appleid.apple.com/auth/keys",
	}
}

// FetchAuthUser returns an AuthUser instance based on the provided token.
//
// API reference: https://developer.apple.com/documentation/signinwithapple/authenticating-users-with-sign-in-with-apple#Retrieve-the-users-information-from-Apple-ID-servers.
func (p *Apple) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserInfo(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		EmailVerified any    `json:"email_verified"` // could be string or bool
		Email         string `json:"email"`
		Id            string `json:"sub"`

		// not returned at the time of writing and it is usually
		// manually populated in apis.recordAuthWithOAuth2
		Name string `json:"name"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           extracted.Id,
		Name:         extracted.Name,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	if cast.ToBool(extracted.EmailVerified) {
		user.Email = extracted.Email
	}

	return user, nil
}

// FetchRawUserInfo implements Provider.FetchRawUserInfo interface.
//
// Note that Apple doesn't have a UserInfo endpoint and claims about
// the users are included in the id_token (without the name - see #7090).
func (p *Apple) FetchRawUserInfo(token *oauth2.Token) ([]byte, error) {
	idToken, _ := token.Extra("id_token").(string)

	claims, err := p.parseAndVerifyIdToken(idToken)
	if err != nil {
		return nil, err
	}

	return json.Marshal(claims)
}

func (p *Apple) parseAndVerifyIdToken(idToken string) (jwt.MapClaims, error) {
	if idToken == "" {
		return nil, errors.New("empty id_token")
	}

	// extract the token claims
	// ---
	claims := jwt.MapClaims{}
	_, _, err := jwt.NewParser().ParseUnverified(idToken, claims)
	if err != nil {
		return nil, err
	}

	// validate common claims per https://developer.apple.com/documentation/sign_in_with_apple/sign_in_with_apple_rest_api/verifying_a_user#3383769
	// ---
	jwtValidator := jwt.NewValidator(
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithLeeway(idTokenLeeway),
		jwt.WithIssuer("https://appleid.apple.com"),
		jwt.WithAudience(p.clientId),
	)
	err = jwtValidator.Validate(claims)
	if err != nil {
		return nil, err
	}

	// validate id_token signature
	//
	// note: this step could be technically considered optional because we trust
	// the token which is a result of direct TLS communication with the provider
	// (see also https://openid.net/specs/openid-connect-core-1_0.html#IDTokenValidation)
	// ---
	err = jwk.ValidateTokenSignature(p.ctx, idToken, p.jwksURL)
	if err != nil {
		return nil, fmt.Errorf("id_token validation failed: %w", err)
	}

	return claims, nil
}
