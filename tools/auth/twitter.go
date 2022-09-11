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
		scopes: []string{
			"users.read",

			// we don't actually use this scope, but for some reason it is required by the `/2/users/me` endpoint
			// (see https://developer.twitter.com/en/docs/twitter-api/users/lookup/api-reference/get-users-me)
			"tweet.read",
		},
		authUrl:    "https://twitter.com/i/oauth2/authorize",
		tokenUrl:   "https://api.twitter.com/2/oauth2/token",
		userApiUrl: "https://api.twitter.com/2/users/me?user.fields=id,name,username,profile_image_url",
	}}
}

// FetchAuthUser returns an AuthUser instance based on the Twitter's user api.
func (p *Twitter) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	// https://developer.twitter.com/en/docs/twitter-api/users/lookup/api-reference/get-users-me
	rawData := struct {
		Data struct {
			Id              string `json:"id"`
			Name            string `json:"name"`
			Username        string `json:"username"`
			ProfileImageUrl string `json:"profile_image_url"`

			// NB! At the time of writing, Twitter OAuth2 doesn't support returning the user email address
			// (see https://twittercommunity.com/t/which-api-to-get-user-after-oauth2-authorization/162417/33)
			Email string `json:"email"`
		} `json:"data"`
	}{}

	if err := p.FetchRawUserData(token, &rawData); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:        rawData.Data.Id,
		Name:      rawData.Data.Name,
		Username:  rawData.Data.Username,
		Email:     rawData.Data.Email,
		AvatarUrl: rawData.Data.ProfileImageUrl,
	}

	return user, nil
}
