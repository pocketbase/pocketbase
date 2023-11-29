package auth

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

var _ Provider = (*Gitlab)(nil)

// NameGitlab is the unique name of the Gitlab provider.
const NameGitlab string = "gitlab"

// Gitlab allows authentication via Gitlab OAuth2.
type Gitlab struct {
	*baseProvider
}

// NewGitlabProvider creates new Gitlab provider instance with some defaults.
func NewGitlabProvider() *Gitlab {
	return &Gitlab{&baseProvider{
		ctx:         context.Background(),
		displayName: "GitLab",
		pkce:        true,
		scopes:      []string{"read_user"},
		authUrl:     "https://gitlab.com/oauth/authorize",
		tokenUrl:    "https://gitlab.com/oauth/token",
		userApiUrl:  "https://gitlab.com/api/v4/user",
	}}
}

// FetchAuthUser returns an AuthUser instance based the Gitlab's user api.
//
// API reference: https://docs.gitlab.com/ee/api/users.html#for-admin
func (p *Gitlab) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserData(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Id        int    `json:"id"`
		Name      string `json:"name"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		AvatarUrl string `json:"avatar_url"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           strconv.Itoa(extracted.Id),
		Name:         extracted.Name,
		Username:     extracted.Username,
		Email:        extracted.Email,
		AvatarUrl:    extracted.AvatarUrl,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	return user, nil
}
