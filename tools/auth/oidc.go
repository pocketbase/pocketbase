package auth

import (
	"context"
	"encoding/json"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

func init() {
	Providers[NameOIDC] = wrapFactory(NewOIDCProvider)
	Providers[NameOIDC+"2"] = wrapFactory(NewOIDCProvider)
	Providers[NameOIDC+"3"] = wrapFactory(NewOIDCProvider)
}

var _ Provider = (*OIDC)(nil)

// NameOIDC is the unique name of the OpenID Connect (OIDC) provider.
const NameOIDC string = "oidc"

// OIDC allows authentication via OpenID Connect (OIDC) OAuth2 provider.
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
//
// @todo consider adding support for reading the user data from the id_token.
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
		EmailVerified bool   `json:"email_verified"`
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

	if extracted.EmailVerified {
		user.Email = extracted.Email
	}

	return user, nil
}
