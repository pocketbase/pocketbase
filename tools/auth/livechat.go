package auth

import (
	"context"
	"encoding/json"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

var _ Provider = (*Livechat)(nil)

// NameLivechat is the unique name of the Livechat provider.
const NameLivechat = "livechat"

// Livechat allows authentication via Livechat OAuth2.
type Livechat struct {
	*baseProvider
}

// NewLivechatProvider creates new Livechat provider instance with some defaults.
func NewLivechatProvider() *Livechat {
	return &Livechat{&baseProvider{
		ctx:         context.Background(),
		displayName: "LiveChat",
		pkce:        true,
		scopes:      []string{}, // default scopes are specified from the provider dashboard
		authUrl:     "https://accounts.livechat.com/",
		tokenUrl:    "https://accounts.livechat.com/token",
		userApiUrl:  "https://accounts.livechat.com/v2/accounts/me",
	}}
}

// FetchAuthUser returns an AuthUser based on the Livechat accounts API.
//
// API reference: https://developers.livechat.com/docs/authorization
func (p *Livechat) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserData(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Id            string `json:"account_id"`
		Name          string `json:"name"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		AvatarUrl     string `json:"avatar_url"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           extracted.Id,
		Name:         extracted.Name,
		AvatarUrl:    extracted.AvatarUrl,
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
