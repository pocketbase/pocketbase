package auth

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/golang-jwt/jwt/v5"
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

	// extract the token header params and claims
	// ---
	claims := jwt.MapClaims{}
	t, _, err := jwt.NewParser().ParseUnverified(idToken, claims)
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
	kid, _ := t.Header["kid"].(string)
	err = validateIdTokenSignature(p.ctx, idToken, p.jwksURL, kid)
	if err != nil {
		return nil, err
	}

	return claims, nil
}
