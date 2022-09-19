package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/oauth2"
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

	if err := p.fetchUser(token, &rawData); err != nil {
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

// Temporary hack to get the user from the Twitch API.
// Must be done this way so that we can include the "Client-ID" header.
func (p *Twitch) fetchUser(token *oauth2.Token, result any) error {
	req, err := http.NewRequest("GET", p.userApiUrl, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	req.Header.Set("Client-Id", p.clientId)

	client := &http.Client{Timeout: time.Second * 10}

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
