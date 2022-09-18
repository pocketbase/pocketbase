package auth

import (
	"golang.org/x/oauth2"
)

var _ Provider = (*Twitch)(nil)

// NameTwitch is the unique name of the Twitch provider.
const NameTwitch string = "twitch"

// Twitch allows authentication via Twitch OAuth2.
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
		Data []struct {
			Login     string `json:"login"`
			Id        string `json:"id"`
			Name      string `json:"display_name"`
			Email     string `json:"email"`
			AvatarUrl string `json:"profile_image_url"`
		} `json:"data"`
	}{}

	if err := p.FetchRawUserData(token, &rawData); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:        rawData.Data[0].Id,
		Name:      rawData.Data[0].Name,
		Username:  rawData.Data[0].Login,
		Email:     rawData.Data[0].Email,
		AvatarUrl: rawData.Data[0].AvatarUrl,
	}

	return user, nil
}
