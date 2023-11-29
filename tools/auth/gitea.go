package auth

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

var _ Provider = (*Gitea)(nil)

// NameGitea is the unique name of the Gitea provider.
const NameGitea string = "gitea"

// Gitea allows authentication via Gitea OAuth2.
type Gitea struct {
	*baseProvider
}

// NewGiteaProvider creates new Gitea provider instance with some defaults.
func NewGiteaProvider() *Gitea {
	return &Gitea{&baseProvider{
		ctx:         context.Background(),
		displayName: "Gitea",
		pkce:        true,
		scopes:      []string{"read:user", "user:email"},
		authUrl:     "https://gitea.com/login/oauth/authorize",
		tokenUrl:    "https://gitea.com/login/oauth/access_token",
		userApiUrl:  "https://gitea.com/api/v1/user",
	}}
}

// FetchAuthUser returns an AuthUser instance based on Gitea's user api.
//
// API reference: https://try.gitea.io/api/swagger#/user/userGetCurrent
func (p *Gitea) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
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
		Name      string `json:"full_name"`
		Username  string `json:"login"`
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
