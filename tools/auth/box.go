package auth

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

func init() {
	Providers[NameBox] = wrapFactory(NewBoxProvider)
}

var _ Provider = (*Box)(nil)

// NameBox is the unique name of the Box provider.
const NameBox = "box"

// Box is an auth provider for Box.
type Box struct {
	BaseProvider
}

// NewBoxProvider creates a new Box provider instance with some defaults.
func NewBoxProvider() *Box {
	return &Box{BaseProvider{
		ctx:         context.Background(),
		displayName: "Box",
		pkce:        true,
		scopes:      []string{"root_readonly"},
		authURL:     "https://account.box.com/api/oauth2/authorize",
		tokenURL:    "https://api.box.com/oauth2/token",
		userInfoURL: "https://api.box.com/2.0/users/me",
	}}
}

// FetchAuthUser returns an AuthUser instance based on the Box's user API.
//
// API reference: https://developer.box.com/reference/get-users-me/
func (p *Box) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserInfo(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Id        string `json:"id"`
		Name      string `json:"name"`
		Login     string `json:"login"`
		AvatarURL string `json:"avatar_url"`
		Status    string `json:"status"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	if extracted.Status != "active" {
		return nil, fmt.Errorf("Box user account is not active (status: %q)", extracted.Status)
	}

	user := &AuthUser{
		Id:           extracted.Id,
		Name:         extracted.Name,
		AvatarURL:    extracted.AvatarURL,
		Email:        extracted.Login, // Box requires verified email for OAuth authorization
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	return user, nil
}
