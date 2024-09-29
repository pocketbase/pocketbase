package auth

import (
	"context"
	"encoding/json"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

func init() {
	Providers[NameTwitter] = wrapFactory(NewTwitterProvider)
}

var _ Provider = (*Twitter)(nil)

// NameTwitter is the unique name of the Twitter provider.
const NameTwitter string = "twitter"

// Twitter allows authentication via Twitter OAuth2.
type Twitter struct {
	BaseProvider
}

// NewTwitterProvider creates new Twitter provider instance with some defaults.
func NewTwitterProvider() *Twitter {
	return &Twitter{BaseProvider{
		ctx:         context.Background(),
		displayName: "Twitter",
		pkce:        true,
		scopes: []string{
			"users.read",

			// we don't actually use this scope, but for some reason it is required by the `/2/users/me` endpoint
			// (see https://developer.twitter.com/en/docs/twitter-api/users/lookup/api-reference/get-users-me)
			"tweet.read",
		},
		authURL:     "https://twitter.com/i/oauth2/authorize",
		tokenURL:    "https://api.twitter.com/2/oauth2/token",
		userInfoURL: "https://api.twitter.com/2/users/me?user.fields=id,name,username,profile_image_url",
	}}
}

// FetchAuthUser returns an AuthUser instance based on the Twitter's user api.
//
// API reference: https://developer.twitter.com/en/docs/twitter-api/users/lookup/api-reference/get-users-me
func (p *Twitter) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
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
			Id              string `json:"id"`
			Name            string `json:"name"`
			Username        string `json:"username"`
			ProfileImageURL string `json:"profile_image_url"`

			// NB! At the time of writing, Twitter OAuth2 doesn't support returning the user email address
			// (see https://twittercommunity.com/t/which-api-to-get-user-after-oauth2-authorization/162417/33)
			// Email string `json:"email"`
		} `json:"data"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           extracted.Data.Id,
		Name:         extracted.Data.Name,
		Username:     extracted.Data.Username,
		AvatarURL:    extracted.Data.ProfileImageURL,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	return user, nil
}
