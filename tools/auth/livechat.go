package auth

import (
	"context"
	"encoding/json"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

func init() {
	Providers[NameLivechat] = wrapFactory(NewLivechatProvider)
}

var _ Provider = (*Livechat)(nil)

// NameLivechat is the unique name of the Livechat provider.
const NameLivechat = "livechat"

// Livechat allows authentication via Livechat OAuth2.
type Livechat struct {
	BaseProvider
}

// NewLivechatProvider creates new Livechat provider instance with some defaults.
func NewLivechatProvider() *Livechat {
	return &Livechat{BaseProvider{
		ctx:         context.Background(),
		order:       27,
		logo:        `<svg xmlns="http://www.w3.org/2000/svg" xml:space="preserve" viewBox="0 0 80 80"><path d="M79.5 49.7a19 19 0 0 1-19 17.3H50L30 80V67l20-13h10.5a6 6 0 0 0 6.1-5.3q1-15-.2-30a5.6 5.6 0 0 0-5.2-5.1Q50.9 13.1 40 13c-10.9-.1-14.4.2-21.2.7-2.8.2-5 2.3-5.2 5.1q-1 15-.2 30c.3 3.1 3 5.3 6.1 5.2H30v13H19.5c-9.9.1-18.2-7.4-19-17.3q-1-16 .2-32C1.5 8.5 8.8 1.3 18 .7a313 313 0 0 1 44.1.1c9.2.6 16.5 7.8 17.3 17q1.1 16 .1 31.9" style="fill:#ff5100"/></svg>`,
		displayName: "LiveChat",
		pkce:        true,
		scopes:      []string{}, // default scopes are specified from the provider dashboard
		authURL:     "https://accounts.livechat.com/",
		tokenURL:    "https://accounts.livechat.com/token",
		userInfoURL: "https://accounts.livechat.com/v2/accounts/me",
	}}
}

// FetchAuthUser returns an AuthUser based on the Livechat accounts API.
//
// API reference: https://developers.livechat.com/docs/authorization
func (p *Livechat) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserInfo(token)
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
		AvatarURL     string `json:"avatar_url"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           extracted.Id,
		Name:         extracted.Name,
		AvatarURL:    extracted.AvatarURL,
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
