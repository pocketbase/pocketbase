package auth

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

func init() {
	Providers[NameGitlab] = wrapFactory(NewGitlabProvider)
}

var _ Provider = (*Gitlab)(nil)

// NameGitlab is the unique name of the Gitlab provider.
const NameGitlab string = "gitlab"

// Gitlab allows authentication via Gitlab OAuth2.
type Gitlab struct {
	BaseProvider
}

// NewGitlabProvider creates new Gitlab provider instance with some defaults.
func NewGitlabProvider() *Gitlab {
	return &Gitlab{BaseProvider{
		ctx:         context.Background(),
		displayName: "GitLab",
		pkce:        true,
		scopes:      []string{"read_user"},
		authURL:     "https://gitlab.com/oauth/authorize",
		tokenURL:    "https://gitlab.com/oauth/token",
		userInfoURL: "https://gitlab.com/api/v4/user",
	}}
}

// FetchAuthUser returns an AuthUser instance based the Gitlab's user api.
//
// API reference: https://docs.gitlab.com/ee/api/users.html#for-admin
func (p *Gitlab) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserInfo(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Name      string `json:"name"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
		Id        int64  `json:"id"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           strconv.FormatInt(extracted.Id, 10),
		Name:         extracted.Name,
		Username:     extracted.Username,
		Email:        extracted.Email,
		AvatarURL:    extracted.AvatarURL,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	return user, nil
}
