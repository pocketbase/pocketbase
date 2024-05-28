package auth

import (
	"context"
	"encoding/json"
	"golang.org/x/oauth2"
	"net/http"
)

var _ Provider = (*Notion)(nil)

// NameNotion is the unique name of the Patreon provider.
const NameNotion string = "notion"

// Notion allows authentication via Notion OAuth2.
type Notion struct {
	*baseProvider
}

// NewNotionProvider creates new Notion provider instance with some defaults.
func NewNotionProvider() *Notion {
	return &Notion{&baseProvider{
		ctx:        context.Background(),
		authUrl:    "https://api.notion.com/v1/oauth/authorize",
		tokenUrl:   "https://api.notion.com/v1/oauth/token",
		userApiUrl: "https://api.notion.com/v1/users/me",
	}}
}

// FetchAuthUser returns an AuthUser instance based on the Notion's User api.
// API reference: https://developers.notion.com/reference/get-self
func (p *Notion) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	req, err := http.NewRequestWithContext(p.ctx, "GET", p.userApiUrl, nil)
	if err != nil {
		return nil, err
	}

	// See https://developers.notion.com/reference/versioning where header is required for all requests
	req.Header.Set("Notion-Version", "2022-06-28")

	data, err := p.sendRawUserDataRequest(req, token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		AvatarUrl string `json:"avatar_url"`
		Bot       struct {
			Owner struct {
				Type string `json:"type"`
				User struct {
					AvatarUrl string `json:"avatar_url"`
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
		AvatarUrl:    extracted.Bot.Owner.User.AvatarUrl,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	return user, nil
}
