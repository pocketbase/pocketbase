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
