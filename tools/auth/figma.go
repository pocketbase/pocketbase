package auth

import (
	"context"
	"encoding/json"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

func init() {
	Providers[NameFigma] = wrapFactory(NewFigmaProvider)
}

var _ Provider = (*Figma)(nil)

// NameFigma is the unique name of the Figma provider.
const NameFigma string = "figma"

// Figma allows authentication via Figma OAuth2.
type Figma struct {
	BaseProvider
}

// NewFigmaProvider creates a new Figma provider instance with some defaults.
func NewFigmaProvider() *Figma {
	// https://www.figma.com/developers/api#oauth2
	// https://developers.figma.com/docs/rest-api/authentication/
	return &Figma{BaseProvider{
		ctx:         context.Background(),
		displayName: "Figma",
		pkce:        true,
		scopes:      []string{"current_user:read"},
		authURL:     "https://www.figma.com/oauth",
		tokenURL:    "https://api.figma.com/v1/oauth/token",
		userInfoURL: "https://api.figma.com/v1/me",
	}}
}

// FetchAuthUser returns an AuthUser instance based on Figma's user api.
//
// API reference: https://developers.figma.com/docs/rest-api/users-endpoints/
func (p *Figma) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserInfo(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Id       string `json:"id"`
		Handle   string `json:"handle"`
		Email    string `json:"email"`
		ImageURL string `json:"img_url"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           extracted.Id,
		Name:         extracted.Handle,
		Email:        extracted.Email,
		AvatarURL:    extracted.ImageURL,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	return user, nil
}
