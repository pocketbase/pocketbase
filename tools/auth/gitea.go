package auth

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

func init() {
	Providers[NameGitea] = wrapFactory(NewGiteaProvider)
}

var _ Provider = (*Gitea)(nil)

// NameGitea is the unique name of the Gitea provider.
const NameGitea string = "gitea"

// Gitea allows authentication via Gitea OAuth2.
type Gitea struct {
	BaseProvider
}

// NewGiteaProvider creates new Gitea provider instance with some defaults.
func NewGiteaProvider() *Gitea {
	return &Gitea{BaseProvider{
		ctx:         context.Background(),
		order:       11,
		logo:        `<svg xmlns="http://www.w3.org/2000/svg" xml:space="preserve" viewBox="0 0 640 640"><path d="m396 484-127-61c-12-6-18-21-12-34l61-127c6-12 21-17 34-11l27 13V154h17v118s57 24 83 40q7 3 13 14 3 10-1 19l-61 127c-6 13-22 18-34 12" style="fill:#fff"/><path d="M623 150c-4-4-10-4-10-4l-178 8-39 1v117l-17-8V155l-89-3-157-8q-15-2-39 1c-9 2-34 8-54 27C-5 212 7 276 8 286c2 12 7 44 32 72 46 56 144 55 144 55s12 29 31 56c25 33 50 59 75 62h189s12 0 29-11c14-8 26-23 26-23s13-14 31-45l14-28s55-118 55-232c-1-34-9-40-11-42M126 354c-26-9-37-19-37-19s-19-13-29-40c-16-44-1-71-1-71s8-22 38-30c14-4 31-3 31-3s7 59 16 94c7 30 25 78 25 78s-26-3-43-9m300 108s-6 14-20 15l-10-1-5-2-113-55s-11-6-13-16c-2-8 3-18 3-18l54-112s5-10 12-13l5-1c8-3 18 2 18 2l110 54s13 6 16 16q1 12-2 17c-6 16-55 114-55 114" style="fill:#609926"/><path d="M327 380q-14 1-17 14-2 13 9 20c7 4 17 2 22-5q8-13-1-24l24-49h6l7-4 29 16 5 5c2 6-2 15-2 15-2 7-18 40-18 40q-13 1-18 13-4 14 9 21a18 18 0 0 0 21-28l6-11 13-30c1-2 6-11 3-22-2-11-13-16-13-16-12-8-29-16-29-16l-1-7-4-6 14-29-12-6-14 29q-11 0-16 10t1 20z" style="fill:#609926"/></svg>`,
		displayName: "Gitea",
		pkce:        true,
		scopes:      []string{"read:user", "user:email"},
		authURL:     "https://gitea.com/login/oauth/authorize",
		tokenURL:    "https://gitea.com/login/oauth/access_token",
		userInfoURL: "https://gitea.com/api/v1/user",
	}}
}

// FetchAuthUser returns an AuthUser instance based on Gitea's user api.
//
// API reference: https://try.gitea.io/api/swagger#/user/userGetCurrent
func (p *Gitea) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserInfo(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Name      string `json:"full_name"`
		Username  string `json:"login"`
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
