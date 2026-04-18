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
		order:       8,
		logo:        `<svg xmlns="http://www.w3.org/2000/svg" width="256" height="247" fill="none"><path fill="#e24329" d="m251.7 97.7-.3-.9-34.7-90.6a9 9 0 0 0-9-5.7q-3 .2-5.2 2a9 9 0 0 0-3.1 4.7L176 78.9H81L57.7 7.2a9.1 9.1 0 0 0-17.2-1L5.6 96.8l-.4.9a64.4 64.4 0 0 0 21.4 74.5h.1l.3.3L80 212l26 19.8 16 12a11 11 0 0 0 13 0l15.9-12 26.2-19.8 53.1-39.9h.2a64.5 64.5 0 0 0 21.3-74.5"/><path fill="#fc6d26" d="m251.7 97.7-.3-.9c-17 3.5-32.9 10.6-46.7 21l-76.2 57.6 48.5 36.7 53.2-39.8.2-.1a64.5 64.5 0 0 0 21.3-74.5"/><path fill="#fca326" d="m80 212.1 26 19.8 16 12a11 11 0 0 0 13 0l15.9-12 26.2-19.8s-22.7-17-48.6-36.7z"/><path fill="#fc6d26" d="M52.2 117.8a117 117 0 0 0-46.6-21l-.4.9a64.4 64.4 0 0 0 21.4 74.5h.1l.3.3L80 212l48.5-36.7z"/></svg>`,
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
