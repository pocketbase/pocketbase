package auth

import (
	"context"
	"encoding/json"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

var _ Provider = (*Patreon)(nil)

// NamePatreon is the unique name of the Patreon provider.
const NamePatreon string = "patreon"

// Patreon allows authentication via Patreon OAuth2.
type Patreon struct {
	*baseProvider
}

// NewPatreonProvider creates new Patreon provider instance with some defaults.
func NewPatreonProvider() *Patreon {
	return &Patreon{&baseProvider{
		ctx:         context.Background(),
		displayName: "Patreon",
		pkce:        true,
		scopes:      []string{"identity", "identity[email]"},
		authUrl:     "https://www.patreon.com/oauth2/authorize",
		tokenUrl:    "https://www.patreon.com/api/oauth2/token",
		userApiUrl:  "https://www.patreon.com/api/oauth2/v2/identity?fields%5Buser%5D=full_name,email,vanity,image_url,is_email_verified",
	}}
}

// FetchAuthUser returns an AuthUser instance based on the Patreons's identity api.
//
// API reference:
// https://docs.patreon.com/#get-api-oauth2-v2-identity
// https://docs.patreon.com/#user-v2
func (p *Patreon) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserData(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Data struct {
			Id         string `json:"id"`
			Attributes struct {
				Email           string `json:"email"`
				Name            string `json:"full_name"`
				Username        string `json:"vanity"`
				AvatarUrl       string `json:"image_url"`
				IsEmailVerified bool   `json:"is_email_verified"`
			} `json:"attributes"`
		} `json:"data"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           extracted.Data.Id,
		Username:     extracted.Data.Attributes.Username,
		Name:         extracted.Data.Attributes.Name,
		AvatarUrl:    extracted.Data.Attributes.AvatarUrl,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	if extracted.Data.Attributes.IsEmailVerified {
		user.Email = extracted.Data.Attributes.Email
	}

	return user, nil
}
