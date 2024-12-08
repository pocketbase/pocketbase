package auth

import (
	"context"
	"encoding/json"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

func init() {
	Providers[NameWebex] = wrapFactory(NewWebexProvider)
}

var _ Provider = (*Webex)(nil)

// NameWebex is the unique name of the Webex provider.
const NameWebex string = "webex"

// Webex allows authentication via Webex OAuth2.
type Webex struct {
	BaseProvider
}

// NewWebexProvider creates a new Webex provider instance with some defaults.
func NewWebexProvider() *Webex {
	return &Webex{BaseProvider{
		ctx:         context.Background(),
		displayName: "Webex",
		pkce:        true,
		scopes:      []string{"spark:people_read"},
		authURL:     "https://webexapis.com/v1/authorize",
		tokenURL:    "https://webexapis.com/v1/access_token",
		userInfoURL: "https://webexapis.com/v1/people/me",
	}}
}

// FetchAuthUser returns an AuthUser instance based on the Webex user api.
//
// API reference: https://developer.webex.com/docs/api/v1/people/get-my-own-details
func (p *Webex) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserInfo(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Id          string   `json:"id"`
		DisplayName string   `json:"displayName"`
		UserName    string   `json:"userName"`
		Avatar      string   `json:"avatar"`
		Emails      []string `json:"emails"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           extracted.Id,
		Name:         extracted.DisplayName,
		Username:     extracted.UserName,
		AvatarURL:    extracted.Avatar,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	if len(extracted.Emails) == 1 {
		user.Email = extracted.Emails[0]
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	return user, nil
}
