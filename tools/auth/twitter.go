package auth

import (
	"golang.org/x/oauth2"
)

var _ Provider = (*Twitter)(nil)

// NameTwitter is the unique name of the Twitter provider.
const NameTwitter string = "twitter"

// Twitter allows authentication via Twitter OAuth2.
type Twitter struct {
	*baseProvider
}

// NewTwitterProvider creates new Twitter provider instance with some defaults.
func NewTwitterProvider() *Twitter {
	return &Twitter{&baseProvider{
		scopes:     []string{"users.read"},
		authUrl:    "https://twitter.com/i/oauth2/authorize",
		tokenUrl:   "https://api.twitter.com/2/oauth2/token",
		userApiUrl: "https://api.twitter.com/1.1/account/verify_credentials.json?include_email=true",
	}}
}

// FetchAuthUser returns an AuthUser instance based on the Twitter's user api.
func (p *Twitter) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	// https://developer.twitter.com/en/docs/twitter-api/v1/accounts-and-users/manage-account-settings/api-reference/get-account-verify_credentials#example-response
	rawData := struct {
		Id                   string
		Name                 string
		Email                string
		ProfileImageUrlHttps string `json:"profile_image_url_https"`
	}{}

	if err := p.FetchRawUserData(token, &rawData); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:        rawData.Id,
		Name:      rawData.Name,
		Email:     rawData.Email,
		AvatarUrl: rawData.ProfileImageUrlHttps,
	}

	return user, nil
}
