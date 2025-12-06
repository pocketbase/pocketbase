package auth

import (
	"context"
	"encoding/json"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

func init() {
	Providers[NameSlack] = wrapFactory(NewSlackProvider)
}

var _ Provider = (*Slack)(nil)

// NameSlack is the unique name of the Slack provider.
const NameSlack string = "slack"

// Slack allows authentication via Slack OAuth2 (Sign in with Slack / OpenID Connect).
type Slack struct {
	BaseProvider
}

// NewSlackProvider creates a new Slack provider instance with some defaults.
func NewSlackProvider() *Slack {
	// https://api.slack.com/authentication/sign-in-with-slack
	// https://api.slack.com/methods/openid.connect.userInfo
	return &Slack{BaseProvider{
		ctx:         context.Background(),
		displayName: "Slack",
		pkce:        true,
		scopes:      []string{"openid", "email", "profile"},
		authURL:     "https://slack.com/openid/connect/authorize",
		tokenURL:    "https://slack.com/api/openid.connect.token",
		userInfoURL: "https://slack.com/api/openid.connect.userInfo",
	}}
}

// FetchAuthUser returns an AuthUser instance based on Slack's user api.
//
// API reference: https://api.slack.com/methods/openid.connect.userInfo
func (p *Slack) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserInfo(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Sub           string `json:"sub"`
		Name          string `json:"name"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Picture       string `json:"picture"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           extracted.Sub,
		Name:         extracted.Name,
		AvatarURL:    extracted.Picture,
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
