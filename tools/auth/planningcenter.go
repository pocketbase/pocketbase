package auth

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

func init() {
	Providers[NamePlanningcenter] = wrapFactory(NewPlanningcenterProvider)
}

var _ Provider = (*Planningcenter)(nil)

// NamePlanningcenter is the unique name of the Planningcenter provider.
const NamePlanningcenter string = "planningcenter"

// Planningcenter allows authentication via Planningcenter OAuth2.
type Planningcenter struct {
	BaseProvider
}

// NewPlanningcenterProvider creates a new Planningcenter provider instance with some defaults.
func NewPlanningcenterProvider() *Planningcenter {
	return &Planningcenter{BaseProvider{
		ctx:         context.Background(),
		order:       29,
		logo:        `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0.1 26.2 26.4"><path fill="#2565f4" d="M23.2 2.7c1.8.5 3 2.2 3 4v13c0 1.8-1.3 3.5-3 4l-8 2.5a8 8 0 0 1-4.3 0L3 23.8c-1.7-.6-3-2.3-3-4.1v-13c0-1.8 1.3-3.4 3-4l8-2.3a8 8 0 0 1 4.2 0z"/><path fill="#fff" d="m18 7.6-4.7 1.3h-.5L8 7.6q-1.7-.3-1.8 1.3v10.5q0 .3.3.4l1.7.4q.5.1.5-.3v-2.5H9l3.5 1h1.4l5.2-1.6q.9-.2.9-1V9c0-1-1-1.7-2-1.4m-.6 7-.2.1-3.7 1.2h-.6l-3.7-1a.2.2 0 0 1-.2-.3v-3.8q0-.2.3-.2l3 .8a3 3 0 0 0 1.8 0l3-.8q.2 0 .3.2z"/></svg>`,
		displayName: "Planning Center",
		pkce:        true,
		scopes:      []string{"people"},
		authURL:     "https://api.planningcenteronline.com/oauth/authorize",
		tokenURL:    "https://api.planningcenteronline.com/oauth/token",
		userInfoURL: "https://api.planningcenteronline.com/people/v2/me",
	}}
}

// FetchAuthUser returns an AuthUser instance based on the Planningcenter's user api.
//
// API reference: https://developer.planning.center/docs/#/overview/authentication
func (p *Planningcenter) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserInfo(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Data struct {
			Id         string `json:"id"`
			Attributes struct {
				Status    string `json:"status"`
				Name      string `json:"name"`
				AvatarURL string `json:"avatar"`
				// don't map the email because users can have multiple assigned
				// and it's not clear if they are verified
			}
		}
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	if extracted.Data.Attributes.Status != "active" {
		return nil, errors.New("the user is not active")
	}

	user := &AuthUser{
		Id:           extracted.Data.Id,
		Name:         extracted.Data.Attributes.Name,
		AvatarURL:    extracted.Data.Attributes.AvatarURL,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	return user, nil
}
