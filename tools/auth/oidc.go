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
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/spf13/cast"
	"golang.org/x/oauth2"
)

// idTokenLeeway is the optional leeway for the id_token timestamp fields validation.
//
// It can be changed externally using the PB_ID_TOKEN_LEEWAY env variable
// (the value must be in seconds, e.g. "PB_ID_TOKEN_LEEWAY=60" for 1 minute).
var idTokenLeeway time.Duration = 5 * time.Minute

func init() {
	Providers[NameOIDC] = wrapFactory(NewOIDCProvider)
	Providers[NameOIDC+"2"] = wrapFactory(NewOIDCProvider)
	Providers[NameOIDC+"3"] = wrapFactory(NewOIDCProvider)

	if leewayStr := os.Getenv("PB_ID_TOKEN_LEEWAY"); leewayStr != "" {
		leeway, err := strconv.Atoi(leewayStr)
		if err == nil {
			idTokenLeeway = time.Duration(leeway) * time.Second
		}
	}
}

var _ Provider = (*OIDC)(nil)

// NameOIDC is the unique name of the OpenID Connect (OIDC) provider.
const NameOIDC string = "oidc"

// OIDC allows authentication via OpenID Connect (OIDC) OAuth2 provider.
//
// If specified the user data is fetched from the userInfoURL.
// Otherwise - from the id_token payload.
//
// The provider support the following Extra config options:
//   - "jwksURL" - url to the keys to validate the id_token signature (optional and used only when reading the user data from the id_token)
//   - "issuers" - list of valid issuers for the iss id_token claim (optioanl and used only when reading the user data from the id_token)
type OIDC struct {
	BaseProvider
}

// NewOIDCProvider creates new OpenID Connect (OIDC) provider instance with some defaults.
func NewOIDCProvider() *OIDC {
	return &OIDC{BaseProvider{
		ctx:         context.Background(),
		displayName: "OIDC",
		pkce:        true,
		scopes: []string{
			"openid", // minimal requirement to return the id
			"email",
			"profile",
		},
	}}
}

// FetchAuthUser returns an AuthUser instance based the provider's user api.
//
// API reference: https://openid.net/specs/openid-connect-core-1_0.html#StandardClaims
func (p *OIDC) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserInfo(token)
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
		Username      string `json:"preferred_username"`
		Picture       string `json:"picture"`
		Email         string `json:"email"`
		EmailVerified any    `json:"email_verified"` // see #6657
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           extracted.Id,
		Name:         extracted.Name,
		Username:     extracted.Username,
		AvatarURL:    extracted.Picture,
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

// FetchRawUserInfo implements Provider.FetchRawUserInfo interface method.
//
// It either fetch the data from p.userInfoURL, or if not set - returns the id_token claims.
func (p *OIDC) FetchRawUserInfo(token *oauth2.Token) ([]byte, error) {
	if p.userInfoURL != "" {
		return p.BaseProvider.FetchRawUserInfo(token)
	}

	claims, err := p.parseIdToken(token)
	if err != nil {
		return nil, err
	}

	return json.Marshal(claims)
}

func (p *OIDC) parseIdToken(token *oauth2.Token) (jwt.MapClaims, error) {
	idToken := token.Extra("id_token").(string)
	if idToken == "" {
		return nil, errors.New("empty id_token")
	}

	claims := jwt.MapClaims{}
	t, _, err := jwt.NewParser().ParseUnverified(idToken, claims)
	if err != nil {
		return nil, err
	}

	// validate common claims
	jwtValidator := jwt.NewValidator(
		jwt.WithIssuedAt(),
		jwt.WithLeeway(idTokenLeeway),
		jwt.WithAudience(p.clientId),
	)
	err = jwtValidator.Validate(claims)
	if err != nil {
		return nil, err
	}

	// validate iss (if "issuers" extra config is set)
	issuers := cast.ToStringSlice(p.Extra()["issuers"])
	if len(issuers) > 0 {
		var isIssValid bool
		claimIssuer, _ := claims.GetIssuer()

		for _, issuer := range issuers {
			if security.Equal(claimIssuer, issuer) {
				isIssValid = true
				break
			}
		}

		if !isIssValid {
			return nil, fmt.Errorf("iss must be one of %v, got %#v", issuers, claims["iss"])
		}
	}

	// validate signature (if "jwksURL" extra config is set)
	//
	// note: this step could be technically considered optional because we trust
	// the token which is a result of direct TLS communication with the provider
	// (see also https://openid.net/specs/openid-connect-core-1_0.html#IDTokenValidation)
	jwksURL := cast.ToString(p.Extra()["jwksURL"])
	if jwksURL != "" {
		kid, _ := t.Header["kid"].(string)
		err = validateIdTokenSignature(p.ctx, idToken, jwksURL, kid)
		if err != nil {
			return nil, err
		}
	}

	return claims, nil
}

func validateIdTokenSignature(ctx context.Context, idToken string, jwksURL string, kid string) error {
	// fetch the public key set
	// ---
	if kid == "" {
		return errors.New("missing kid header value")
	}

	key, err := fetchJWK(ctx, jwksURL, kid)
	if err != nil {
		return err
	}

	// decode the key params per RFC 7518 (https://tools.ietf.org/html/rfc7518#section-6.3)
	// and construct a valid publicKey from them
	// ---
	exponent, err := base64.RawURLEncoding.DecodeString(strings.TrimRight(key.E, "="))
	if err != nil {
		return err
	}

	modulus, err := base64.RawURLEncoding.DecodeString(strings.TrimRight(key.N, "="))
	if err != nil {
		return err
	}

	publicKey := &rsa.PublicKey{
		// https://tools.ietf.org/html/rfc7517#appendix-A.1
		E: int(big.NewInt(0).SetBytes(exponent).Uint64()),
		N: big.NewInt(0).SetBytes(modulus),
	}

	// verify the signiture
	// ---
	parser := jwt.NewParser(jwt.WithValidMethods([]string{key.Alg}))

	parsedToken, err := parser.Parse(idToken, func(t *jwt.Token) (any, error) {
		return publicKey, nil
	})
	if err != nil {
		return err
	}

	if !parsedToken.Valid {
		return errors.New("the parsed id_token is invalid")
	}

	return nil
}

type jwk struct {
	Kty string
	Kid string
	Use string
	Alg string
	N   string
	E   string
}

func fetchJWK(ctx context.Context, jwksURL string, kid string) (*jwk, error) {
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
