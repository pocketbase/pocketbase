package auth

import (
	"context"
	"encoding/json"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/instagram"
)

var _ Provider = (*Instagram)(nil)

// NameInstagram is the unique name of the Instagram provider.
const NameInstagram string = "instagram"

// Instagram allows authentication via Instagram OAuth2.
type Instagram struct {
	*baseProvider
}

// NewInstagramProvider creates new Instagram provider instance with some defaults.
func NewInstagramProvider() *Instagram {
	return &Instagram{&baseProvider{
		ctx:         context.Background(),
		displayName: "Instagram",
		pkce:        true,
		scopes:      []string{"user_profile"},
		authUrl:     instagram.Endpoint.AuthURL,
		tokenUrl:    instagram.Endpoint.TokenURL,
		userApiUrl:  "https://graph.instagram.com/me?fields=id,username,account_type",
	}}
}

// FetchAuthUser returns an AuthUser instance based on the Instagram's user api.
//
// API reference: https://developers.facebook.com/docs/instagram-basic-display-api/reference/user#fields
func (p *Instagram) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserData(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Id       string `json:"id"`
		Username string `json:"username"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           extracted.Id,
		Username:     extracted.Username,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	return user, nil
}
