package auth

import (
	"encoding/json"

	"golang.org/x/oauth2"
)

var _ Provider = (*Authentik)(nil)

// NameAuthentik is the unique name of the Authentik provider.
const NameAuthentik string = "authentik"

// Authentik allows authentication via Authentik OAuth2.
type Authentik struct {
	*baseProvider
}

// NewAuthentikProvider creates new Authentik provider instance with some defaults.
func NewAuthentikProvider() *Authentik {
	return &Authentik{&baseProvider{
		scopes: []string{
			"openid", // minimal requirement to return the id
			"email",
			"profile",
		},
	}}
}

// FetchAuthUser returns an AuthUser instance based the Authentik's user api.
//
// API reference: https://goauthentik.io/docs/providers/oauth2/
func (p *Authentik) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
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
		AvatarUrl:    extracted.Picture,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	if extracted.EmailVerified {
		user.Email = extracted.Email
	}

	return user, nil
}
