package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	if err := p.fetchRawUserData(token, &rawData); err != nil {
		return nil, err
	}

	var user *AuthUser
	if len(rawData.Data) > 0 {
		user = &AuthUser{
			Id:        rawData.Data[0].Id,
			Name:      rawData.Data[0].DisplayName,
			Username:  rawData.Data[0].Login,
			Email:     rawData.Data[0].Email,
			AvatarUrl: rawData.Data[0].ProfileImageUrl,
		}
	}

	return user, nil
}

// Must be done this way so that we can include the "Client-ID" header.
func (p *Twitch) fetchRawUserData(token *oauth2.Token, result any) error {
	req, err := http.NewRequest("GET", p.userApiUrl, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Client-Id", p.clientId)

	client := p.Client(token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// http.Client.Get doesn't treat non 2xx responses as error
	if resp.StatusCode >= 400 {
		return fmt.Errorf(
			"Failed to fetch OAuth2 user profile via %s (%d):\n%s",
			p.userApiUrl,
			resp.StatusCode,
			string(content),
		)
	}

	return json.Unmarshal(content, &result)
}
