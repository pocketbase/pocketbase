package auth

import (
	"context"
	"encoding/json"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

func init() {
	Providers[NameLinear] = wrapFactory(NewLinearProvider)
}

var _ Provider = (*Linear)(nil)

// NameLinear is the unique name of the Linear provider.
const NameLinear string = "linear"

// Linear allows authentication via Linear OAuth2.
type Linear struct {
	BaseProvider
}

// NewLinearProvider creates new Linear provider instance with some defaults.
//
// API reference: https://developers.linear.app/docs/oauth/authentication
func NewLinearProvider() *Linear {
	return &Linear{BaseProvider{
		ctx:         context.Background(),
		displayName: "Linear",
		pkce:        true,
		scopes:      []string{"read"},
		authURL:     "https://linear.app/oauth/authorize",
		tokenURL:    "https://api.linear.app/oauth/token",
		userInfoURL: "https://api.linear.app/graphql",
	}}
}

// FetchAuthUser returns an AuthUser instance based on the Linear's user api.
//
// API reference: https://developers.linear.app/docs/graphql/working-with-the-graphql-api#authentication
func (p *Linear) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
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
			Viewer struct {
				Id       	string `json:"id"`
				DisplayName	string `json:"displayName"`
				Name     	string `json:"name"`
				Email			string `json:"email"`
				AvatarURL 	string `json:"avatarUrl"`
			} `json:"viewer"`
		} `json:"data"`
	}{}

	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	// Email are verified on signup and on change

	user := &AuthUser{
		Id:           extracted.Data.Viewer.Id,
		Name:         extracted.Data.Viewer.DisplayName,
		Username:     extracted.Data.Viewer.Name,
		Email:		  extracted.Data.Viewer.Email,
		AvatarURL:    extracted.Data.Viewer.AvatarURL,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	return user, nil
}