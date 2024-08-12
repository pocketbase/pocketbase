package auth

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/spf13/cast"
	"golang.org/x/oauth2"
)

var _ Provider = (*Apple)(nil)

// NameApple is the unique name of the Apple provider.
const NameApple string = "apple"

// Apple allows authentication via Apple OAuth2.
//
// [OIDC differences]: https://bitbucket.org/openid/connect/src/master/How-Sign-in-with-Apple-differs-from-OpenID-Connect.md
type Apple struct {
	*baseProvider

	jwksUrl string
}

// NewAppleProvider creates a new Apple provider instance with some defaults.
func NewAppleProvider() *Apple {
	return &Apple{
		baseProvider: &baseProvider{
			ctx:         context.Background(),
			displayName: "Apple",
			pkce:        true,
			scopes:      []string{"name", "email"},
			authUrl:     "https://appleid.apple.com/auth/authorize",
			tokenUrl:    "https://appleid.apple.com/auth/token",
		},
		jwksUrl: "https://appleid.apple.com/auth/keys",
	}
}

// FetchAuthUser returns an AuthUser instance based on the provided token.
//
// API reference: https://developer.apple.com/documentation/sign_in_with_apple/tokenresponse.
func (p *Apple) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserData(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Id            string `json:"sub"`
		Name          string `json:"name"`
		Email         string `json:"email"`
		EmailVerified any    `json:"email_verified"` // could be string or bool
		User          struct {
			Name struct {
				FirstName string `json:"firstName"`
				LastName  string `json:"lastName"`
			} `json:"name"`
		} `json:"user"`
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

	if user.Name == "" {
		user.Name = strings.TrimSpace(extracted.User.Name.FirstName + " " + extracted.User.Name.LastName)
	}

	return user, nil
}

// FetchRawUserData implements Provider.FetchRawUserData interface.
//
// Apple doesn't have a UserInfo endpoint and claims about users
// are instead included in the "id_token" (https://openid.net/specs/openid-connect-core-1_0.html#id_tokenExample)
func (p *Apple) FetchRawUserData(token *oauth2.Token) ([]byte, error) {
	idToken, _ := token.Extra("id_token").(string)

	claims, err := p.parseAndVerifyIdToken(idToken)
	if err != nil {
		return nil, err
	}

	// Apple only returns the user object the first time the user authorizes the app
	// https://developer.apple.com/documentation/sign_in_with_apple/sign_in_with_apple_js/configuring_your_webpage_for_sign_in_with_apple#3331292
	rawUser, _ := token.Extra("user").(string)
	if rawUser != "" {
		user := map[string]any{}
		err = json.Unmarshal([]byte(rawUser), &user)
		if err != nil {
			return nil, err
		}
		claims["user"] = user
	}

	return json.Marshal(claims)
}

// -------------------------------------------------------------------

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
	if !claims.VerifyIssuer("https://appleid.apple.com", true) {
		return nil, errors.New("iss must be https://appleid.apple.com")
	}

	if !claims.VerifyAudience(p.clientId, true) {
		return nil, errors.New("aud must be the developer's client_id")
	}

	// fetch the public key set
	// ---
	kid, _ := t.Header["kid"].(string)
	if kid == "" {
		return nil, errors.New("missing kid header value")
	}

	key, err := p.fetchJWK(kid)
	if err != nil {
		return nil, err
	}

	// decode the key params per RFC 7518 (https://tools.ietf.org/html/rfc7518#section-6.3)
	// and construct a valid publicKey from them
	// ---
	exponent, err := base64.RawURLEncoding.DecodeString(strings.TrimRight(key.E, "="))
	if err != nil {
		return nil, err
	}

	modulus, err := base64.RawURLEncoding.DecodeString(strings.TrimRight(key.N, "="))
	if err != nil {
		return nil, err
	}

	publicKey := &rsa.PublicKey{
		// https://tools.ietf.org/html/rfc7517#appendix-A.1
		E: int(big.NewInt(0).SetBytes(exponent).Uint64()),
		N: big.NewInt(0).SetBytes(modulus),
	}

	// verify the id_token
	// ---
	parser := jwt.NewParser(jwt.WithValidMethods([]string{key.Alg}))

	parsedToken, err := parser.Parse(idToken, func(t *jwt.Token) (any, error) {
		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return claims, nil
	}

	return nil, errors.New("the parsed id_token is invalid")
}

type jwk struct {
	Kty string
	Kid string
	Use string
	Alg string
	N   string
	E   string
}

func (p *Apple) fetchJWK(kid string) (*jwk, error) {
	req, err := http.NewRequestWithContext(p.ctx, "GET", p.jwksUrl, nil)
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
			"failed to verify the provided id_token (%d):\n%s",
			res.StatusCode,
			string(rawBody),
		)
	}

	jwks := struct {
		Keys []*jwk
	}{}
	if err := json.Unmarshal(rawBody, &jwks); err != nil {
		return nil, err
	}

	for _, key := range jwks.Keys {
		if key.Kid == kid {
			return key, nil
		}
	}

	return nil, fmt.Errorf("jwk with kid %q was not found", kid)
}
