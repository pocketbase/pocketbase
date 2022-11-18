package auth

import (
	"errors"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/twitch"
)

var _ Provider = (*Twitch)(nil)

// NameTwitch is the unique name of the Twitch provider.
const NameTwitch string = "twitch"

// Twitch allows authentication via Twitch OAuth2.
type Twitch struct {
	*baseProvider
}

// NewTwitchProvider creates new Twitch provider instance with some defaults.
func NewTwitchProvider() *Twitch {
	return &Twitch{&baseProvider{
		scopes:     []string{"user:read:email"},
		authUrl:    twitch.Endpoint.AuthURL,
		tokenUrl:   twitch.Endpoint.TokenURL,
		userApiUrl: "https://api.twitch.tv/helix/users",
	}}
}

// FetchAuthUser returns an AuthUser instance based the Twitch's user api.
func (p *Twitch) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	// https://dev.twitch.tv/docs/api/reference#get-users
	rawData := struct {
		Data []struct {
			Id              string `json:"id"`
			Login           string `json:"login"`
			DisplayName     string `json:"display_name"`
			Email           string `json:"email"`
			ProfileImageUrl string `json:"profile_image_url"`
		} `json:"data"`
	}{}

	if err := p.FetchRawUserData(token, &rawData); err != nil {
		return nil, err
	}

	if len(rawData.Data) == 0 {
		return nil, errors.New("Failed to fetch AuthUser data")
	}

	user := &AuthUser{
		Id:        rawData.Data[0].Id,
		Name:      rawData.Data[0].DisplayName,
		Username:  rawData.Data[0].Login,
		Email:     rawData.Data[0].Email,
		AvatarUrl: rawData.Data[0].ProfileImageUrl,
	}

	return user, nil
}

// FetchRawUserData implements Provider.FetchRawUserData interface.
//
// This differ from baseProvider because Twitch requires the `Client-Id` header.
func (p *Twitch) FetchRawUserData(token *oauth2.Token, result any) error {
	req, err := http.NewRequest("GET", p.userApiUrl, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Client-Id", p.clientId)

	return p.sendRawUserDataRequest(req, token, result)
}
