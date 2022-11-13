package auth

import (
	"encoding/json"
	"io/ioutil"
	"strconv"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var _ Provider = (*Github)(nil)

// NameGithub is the unique name of the Github provider.
const NameGithub string = "github"

// Github allows authentication via Github OAuth2.
type Github struct {
	*baseProvider
}

// NewGithubProvider creates new Github provider instance with some defaults.
func NewGithubProvider() *Github {
	return &Github{&baseProvider{
		scopes:     []string{"read:user", "user:email"},
		authUrl:    github.Endpoint.AuthURL,
		tokenUrl:   github.Endpoint.TokenURL,
		userApiUrl: "https://api.github.com/user",
	}}
}

// FetchAuthUser returns an AuthUser instance based the Github's user api.
func (p *Github) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	// https://docs.github.com/en/rest/reference/users#get-the-authenticated-user
	rawData := struct {
		Login     string `json:"login"`
		Id        int    `json:"id"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarUrl string `json:"avatar_url"`
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

	// in case user set "Keep my email address private",
	// email should be retrieved via extra API request
	if user.Email == "" {
		client := p.Client(token)

		response, err := client.Get(p.userApiUrl + "/emails")
		if err != nil {
			return user, err
		}
		defer response.Body.Close()

		content, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return user, err
		}

		emails := []struct {
			Email    string
			Verified bool
			Primary  bool
		}{}
		if err := json.Unmarshal(content, &emails); err != nil {
			return user, err
		}

		// extract the verified primary email
		for _, email := range emails {
			if email.Verified && email.Primary {
				user.Email = email.Email
				break
			}
		}
	}

	return user, nil
}
