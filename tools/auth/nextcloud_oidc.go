package auth

import (
	"encoding/json"

	"golang.org/x/oauth2"
)

var _ Provider = (*Nextcloud)(nil)

// NameNextcloud is the unique name of the Nextcloud OIDC provider.
const NameNextcloud string = "nextcloud"

// Nextcloud allows authentication via OpenID connect.
type Nextcloud struct {
	*baseProvider
}

// NewNextcloudProvider creates new NextcloudOIDC provider instance with some defaults.
func NewNextcloudProvider() *Nextcloud {
	return &Nextcloud{&baseProvider{
		scopes: []string{
			"profile",
			"email",
		},
	}}
}

// FetchAuthUser returns an AuthUser instance based the Nextcloud OIDC app's user api.
//
// API reference: https://github.com/H2CK/oidc/wiki/User-Documentation
func (p *Nextcloud) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserData(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Sub   string
		Name  string
		Email string
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           extracted.Sub,
		Username: 		extracted.Sub,
		Name:         extracted.Name,
		Email:        extracted.Email,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	return user, nil
}
