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
		order:       13,
		logo:        `<svg xmlns="http://www.w3.org/2000/svg" width="300" height="300.3"><path d="M179 127 290 0h-26l-97 110L89 0H0l117 167L0 300h26l103-116 82 116h89M36 20h41l187 262h-41"/></svg>`,
		displayName: "X/Twitter",
		pkce:        true,
		scopes: []string{
			"users.read",
			"users.email",

			// we don't actually use this scope, but for some reason the `/2/users/me` endpoint fails with 403 without it
			// (see https://docs.x.com/fundamentals/authentication/guides/v2-authentication-mapping#x-api-v2-authentication-mapping)
			"tweet.read",
		},
		authURL:     "https://x.com/i/oauth2/authorize",
		tokenURL:    "https://api.x.com/2/oauth2/token",
		userInfoURL: "https://api.x.com/2/users/me?user.fields=id,name,username,profile_image_url,confirmed_email",
	}}
}

// FetchAuthUser returns an AuthUser instance based on the Twitter's user api.
//
// API reference: https://docs.x.com/x-api/users/user-lookup-me
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
			Email           string `json:"confirmed_email"`
			ProfileImageURL string `json:"profile_image_url"`
		} `json:"data"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           extracted.Data.Id,
		Name:         extracted.Data.Name,
		Username:     extracted.Data.Username,
		Email:        extracted.Data.Email,
		AvatarURL:    extracted.Data.ProfileImageURL,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	return user, nil
}
