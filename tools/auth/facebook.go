package auth

import (
	"context"
	"encoding/json"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

var _ Provider = (*Facebook)(nil)

// NameFacebook is the unique name of the Facebook provider.
const NameFacebook string = "facebook"

// Facebook allows authentication via Facebook OAuth2.
type Facebook struct {
	*baseProvider
}

// NewFacebookProvider creates new Facebook provider instance with some defaults.
func NewFacebookProvider() *Facebook {
	return &Facebook{&baseProvider{
		ctx:         context.Background(),
		displayName: "Facebook",
		pkce:        true,
		scopes:      []string{"email"},
		authUrl:     facebook.Endpoint.AuthURL,
		tokenUrl:    facebook.Endpoint.TokenURL,
		userApiUrl:  "https://graph.facebook.com/me?fields=name,email,picture.type(large)",
	}}
}

// FetchAuthUser returns an AuthUser instance based on the Facebook's user api.
//
// API reference: https://developers.facebook.com/docs/graph-api/reference/user/
func (p *Facebook) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserData(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Id      string
		Name    string
		Email   string
		Picture struct {
			Data struct{ Url string }
		}
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           extracted.Id,
		Name:         extracted.Name,
		Email:        extracted.Email,
		AvatarUrl:    extracted.Picture.Data.Url,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	return user, nil
}
