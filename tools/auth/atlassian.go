package auth

import (
	"context"
	"encoding/json"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

func init() {
	Providers[NameAtlassian] = wrapFactory(NewAtlassianProvider)
}

var _ Provider = (*Atlassian)(nil)

// NameAtlassian is the unique name of the Atlassian provider.
const NameAtlassian string = "atlassian"

// Atlassian allows authentication via Atlassian OAuth2.
type Atlassian struct {
	BaseProvider
}

// NewAtlassianProvider creates a new Atlassian provider instance with some defaults.
func NewAtlassianProvider() *Atlassian {
	// https://developer.atlassian.com/cloud/jira/platform/oauth-2-3lo-apps/
	return &Atlassian{BaseProvider{
		ctx:         context.Background(),
		displayName: "Atlassian",
		pkce:        true,
		scopes:      []string{"read:me"},
		authURL:     "https://auth.atlassian.com/authorize",
		tokenURL:    "https://auth.atlassian.com/oauth/token",
		userInfoURL: "https://api.atlassian.com/me",
	}}
}

// BuildAuthURL implements Provider.BuildAuthURL interface method.
//
// Atlassian requires the "audience" parameter in the auth URL.
func (p *Atlassian) BuildAuthURL(state string, opts ...oauth2.AuthCodeOption) string {
	opts = append(opts, oauth2.SetAuthURLParam("audience", "api.atlassian.com"))
	opts = append(opts, oauth2.SetAuthURLParam("prompt", "consent"))
	return p.BaseProvider.BuildAuthURL(state, opts...)
}

// FetchAuthUser returns an AuthUser instance based on Atlassian's user api.
//
// API reference: https://developer.atlassian.com/cloud/jira/platform/oauth-2-3lo-apps/#how-do-i-retrieve-the-public-profile-of-the-authenticated-user-
func (p *Atlassian) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserInfo(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		AccountId     string `json:"account_id"`
		Email         string `json:"email"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		EmailVerified bool   `json:"email_verified"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           extracted.AccountId,
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
