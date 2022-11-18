package auth

import (
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
		scopes:     []string{"email"},
		authUrl:    facebook.Endpoint.AuthURL,
		tokenUrl:   facebook.Endpoint.TokenURL,
		userApiUrl: "https://graph.facebook.com/me?fields=name,email,picture.type(large)",
	}}
}

// FetchAuthUser returns an AuthUser instance based on the Facebook's user api.
func (p *Facebook) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	// https://developers.facebook.com/docs/graph-api/reference/user/
	rawData := struct {
		Id      string
		Name    string
		Email   string
		Picture struct {
			Data struct{ Url string }
		}
	}{}

	if err := p.FetchRawUserData(token, &rawData); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:        rawData.Id,
		Name:      rawData.Name,
		Email:     rawData.Email,
		AvatarUrl: rawData.Picture.Data.Url,
	}

	return user, nil
}
