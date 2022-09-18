package auth

import (
	"strconv"

	"golang.org/x/oauth2"
)

var _ Provider = (*Twitch)(nil)

// NameTwitch is the unique name of the Twitch provider.
const NameTwitch string = "twitch"

// Github allows authentication via Twitch OAuth2.
type Twitch struct {
	*baseProvider
}

// NewTwitchProvider creates new Github provider instance with some defaults.
func NewTwitchProvider() *Twitch {
	return &Twitch{&baseProvider{
		scopes:     []string{"user:read:email"},
		authUrl:    "https://id.twitch.tv/oauth2/authorize",
		tokenUrl:   "https://id.twitch.tv/oauth2/token",
		userApiUrl: "https://api.twitch.tv/helix/users",
	}}
}

// FetchAuthUser returns an AuthUser instance based the Twitch's user api.
func (p *Twitch) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	// https://dev.twitch.tv/docs/api/reference#get-users
	rawData := struct {
		Login     string `json:"login"`
		Id        int    `json:"id"`
		Name      string `json:"display_name"`
		Email     string `json:"email"`
		AvatarUrl string `json:"profile_image_url"`
	}{}

	if err := p.FetchRawUserData(token, &rawData); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:        strconv.Itoa(rawData.Id),
		Name:      rawData.Name,
		Username:  rawData.Login,
		Email:     rawData.Email,
		AvatarUrl: rawData.AvatarUrl,
	}

	return user, nil
}
