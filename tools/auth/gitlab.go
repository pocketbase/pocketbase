package auth

import (
	"strconv"

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
		scopes:     []string{"read_user"},
		authUrl:    "https://gitlab.com/oauth/authorize",
		tokenUrl:   "https://gitlab.com/oauth/token",
		userApiUrl: "https://gitlab.com/api/v4/user",
	}}
}

// FetchAuthUser returns an AuthUser instance based the Gitlab's user api.
func (p *Gitlab) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	// https://docs.gitlab.com/ee/api/users.html#for-admin
	rawData := struct {
		Id        int    `json:"id"`
		Name      string `json:"name"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		AvatarUrl string `json:"avatar_url"`
	}{}

	if err := p.FetchRawUserData(token, &rawData); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:        strconv.Itoa(rawData.Id),
		Name:      rawData.Name,
		Username:  rawData.Username,
		Email:     rawData.Email,
		AvatarUrl: rawData.AvatarUrl,
	}

	return user, nil
}
