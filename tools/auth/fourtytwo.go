package auth

import (
	"context"
	"encoding/json"
	"strconv"

	"golang.org/x/oauth2"
)

var _ Provider = (*Fourtytwo)(nil)

// NameFourtytwo is the unique name of the Fourtytwo provider.
const NameFourtytwo string = "fourtytwo"

// Fourtytwo allows authentication via Fourtytwo OAuth2 provider.
type Fourtytwo struct {
	*baseProvider
}

// NewFourtytwoProvider creates new Fourtytwo provider instance with some defaults.
func NewFourtytwoProvider() *Fourtytwo {
	return &Fourtytwo{&baseProvider{
		ctx: context.Background(),
		scopes: []string{
			"public", // minimal requirement to return the id,
		},
		authUrl:    "https://api.intra.42.fr/oauth/authorize",
		tokenUrl:   "https://api.intra.42.fr/oauth/token",
		userApiUrl: "https://api.intra.42.fr/v2/me",
	}}
}

// FetchAuthUser returns an AuthUser instance based the provider's user api.
//
// API reference: https://api.intra.42.fr/apidoc/2.0.html
func (p *Fourtytwo) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserData(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Id       int    `json:"id"`
		Name     string `json:"displayname"`
		Username string `json:"login"`
		Picture  struct {
			AvartarUrl string `json:"link"`
		} `json:"image"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"active?"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           strconv.Itoa(extracted.Id),
		Name:         extracted.Name,
		Username:     extracted.Username,
		AvatarUrl:    extracted.Picture.AvartarUrl,
		Email:        extracted.Email,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	return user, nil
}
