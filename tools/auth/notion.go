package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

func init() {
	Providers[NameNotion] = wrapFactory(NewNotionProvider)
}

var _ Provider = (*Notion)(nil)

// NameNotion is the unique name of the Notion provider.
const NameNotion string = "notion"

// Notion allows authentication via Notion OAuth2.
type Notion struct {
	BaseProvider
}

// NewNotionProvider creates new Notion provider instance with some defaults.
func NewNotionProvider() *Notion {
	return &Notion{BaseProvider{
		ctx:         context.Background(),
		order:       17,
		logo:        `<svg xmlns="http://www.w3.org/2000/svg" width="100" height="100" fill="none"><path fill="#000" d="M61 0 6 4q-6 2-6 7v61l3 8 13 17q2 4 8 3l65-4q7-1 7-7V21q0-3-4-5L74 3q-4-4-13-3M26 20q-7 1-9-2l-8-6q-1-2 1-2l54-4q5 0 8 2l9 7 1 1zm-6 68V30q0-3 3-4l63-3q3 0 3 3v58q1 5-4 5l-60 3q-5 1-5-4m59-55q1 4-1 4l-3 1v43l-7 2q-4 0-6-4L43 49v29l6 1s0 4-5 4H30l2-3 3-1V41h-5q0-4 4-5l14-1 20 31V39l-5-1q0-3 3-4z"/></svg>`,
		displayName: "Notion",
		pkce:        true,
		authURL:     "https://api.notion.com/v1/oauth/authorize",
		tokenURL:    "https://api.notion.com/v1/oauth/token",
		userInfoURL: "https://api.notion.com/v1/users/me",
	}}
}

// FetchAuthUser returns an AuthUser instance based on the Notion's User api.
// API reference: https://developers.notion.com/reference/get-self
func (p *Notion) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserInfo(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		AvatarURL string `json:"avatar_url"`
		Bot       struct {
			Owner struct {
				Type string `json:"type"`
				User struct {
					AvatarURL string `json:"avatar_url"`
					Id        string `json:"id"`
					Name      string `json:"name"`
					Person    struct {
						Email string `json:"email"`
					} `json:"person"`
				} `json:"user"`
			} `json:"owner"`
			WorkspaceName string `json:"workspace_name"`
		} `json:"bot"`
		Id        string `json:"id"`
		Name      string `json:"name"`
		Object    string `json:"object"`
		RequestId string `json:"request_id"`
		Type      string `json:"type"`
	}{}

	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           extracted.Bot.Owner.User.Id,
		Name:         extracted.Bot.Owner.User.Name,
		Email:        extracted.Bot.Owner.User.Person.Email,
		AvatarURL:    extracted.Bot.Owner.User.AvatarURL,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	return user, nil
}

// FetchRawUserInfo implements Provider.FetchRawUserInfo interface method.
//
// This differ from BaseProvider because Notion requires a version header for all requests
// (https://developers.notion.com/reference/versioning).
func (p *Notion) FetchRawUserInfo(token *oauth2.Token) ([]byte, error) {
	req, err := http.NewRequestWithContext(p.ctx, "GET", p.userInfoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Notion-Version", "2022-06-28")

	return p.sendRawUserInfoRequest(req, token)
}
