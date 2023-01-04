package auth

import (
	"encoding/json"
	"fmt"

	"golang.org/x/oauth2"
)

var _ Provider = (*Reddit)(nil)

// NameReddit is the unique name of the Reddit provider.
const NameReddit string = "reddit"

// Reddit allows authentication via Reddit OAuth2.
type Reddit struct {
	*baseProvider
}

// NewRedditProvider creates new Reddit provider instance with some defaults.
func NewRedditProvider() *Reddit {
	return &Reddit{&baseProvider{
		scopes:     []string{"identity"},
		authUrl:    "https://www.reddit.com/api/v1/authorize",
		tokenUrl:   "https://www.reddit.com/api/v1/access_token",
		userApiUrl: "https://oauth.reddit.com/api/v1/me",
	}}
}

// FetchAuthUser returns an AuthUser instance based the Reddit's user api.
func (p *Reddit) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserData(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	fmt.Println(string(data))

	extracted := struct {
		Login            string `json:"login"`
		Id               string `json:"id"`
		Name             string `json:"name"`
		AvatarUrl        string `json:"snoovatar_img"`
		HasVerifiedEmail bool   `json:"has_verified_email"`
		IsSuspended      bool   `json:"is_suspended"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:          extracted.Id,
		Name:        extracted.Name, // Reddit isn't sharing the display name only username
		Username:    extracted.Name,
		AvatarUrl:   extracted.AvatarUrl,
		RawUser:     rawUser,
		AccessToken: token.AccessToken,
	}

	if !extracted.HasVerifiedEmail || extracted.IsSuspended {
		return nil, fmt.Errorf("email not verified or account suspended")
	}

	return user, nil
}
