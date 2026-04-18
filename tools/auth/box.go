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
		order:       20,
		logo:        `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 40 21.6"><path fill="#0061d5" d="M39.7 19.2q.7 1.2-.2 2.1-1.2.7-2.2-.2l-3.5-4.5-3.4 4.4c-.5.7-1.5.7-2.2.2q-1-.9-.3-2.1l4-5.2-4-5.2c-.5-.7-.3-1.7.3-2.2s1.7-.3 2.2.3l3.4 4.5L37.3 7q.9-1 2.2-.3 1 1 .2 2.2L35.8 14zm-18.2-.6c-2.6 0-4.7-2-4.7-4.6s2.1-4.6 4.7-4.6 4.7 2.1 4.7 4.6a4.7 4.7 0 0 1-4.7 4.6m-13.8 0c-2.6 0-4.7-2-4.7-4.6s2.1-4.6 4.7-4.6 4.7 2.1 4.7 4.6c0 2.6-2.1 4.6-4.7 4.6M21.5 6.4a8 8 0 0 0-6.8 4 8 8 0 0 0-6.9-4q-2.7 0-4.7 1.5V1.5Q3 .2 1.6 0 .1.1 0 1.5v12.6a7.7 7.7 0 0 0 7.7 7.5c3 0 5.6-1.7 6.9-4.1a8 8 0 0 0 6.8 4.1c4.3 0 7.8-3.4 7.8-7.7a7.5 7.5 0 0 0-7.7-7.5"/></svg>`,
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
