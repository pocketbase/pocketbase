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
