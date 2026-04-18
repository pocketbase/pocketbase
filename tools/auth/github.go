package auth

import (
	"context"
	"encoding/json"
	"io"
	"strconv"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func init() {
	Providers[NameGithub] = wrapFactory(NewGithubProvider)
}

var _ Provider = (*Github)(nil)

// NameGithub is the unique name of the Github provider.
const NameGithub string = "github"

// Github allows authentication via Github OAuth2.
type Github struct {
	BaseProvider
}

// NewGithubProvider creates new Github provider instance with some defaults.
func NewGithubProvider() *Github {
	return &Github{BaseProvider{
		ctx:         context.Background(),
		order:       7,
		logo:        `<svg xmlns="http://www.w3.org/2000/svg" width="256" height="250" preserveAspectRatio="xMidYMid"><path fill="#161614" d="M128 0a128 128 0 0 0-40.5 249.5c6.4 1.1 8.8-2.8 8.8-6.2l-.2-23.8C60.5 227.2 53 204.4 53 204.4c-5.8-14.8-14.2-18.8-14.2-18.8-11.6-7.9.8-7.7.8-7.7 12.9.9 19.7 13.1 19.7 13.1 11.4 19.6 30 14 37.2 10.7 1.2-8.3 4.5-14 8.1-17.1-28.4-3.3-58.3-14.2-58.3-63.3 0-14 5-25.4 13.2-34.3a46 46 0 0 1 1.3-34S71.5 49.7 96 66.3a123 123 0 0 1 64 0c24.5-16.6 35.2-13.1 35.2-13.1a46 46 0 0 1 1.3 33.9c8.2 9 13.2 20.3 13.2 34.3 0 49.2-30 60-58.5 63.2 4.6 4 8.7 11.7 8.7 23.7l-.2 35.1c0 3.4 2.4 7.4 8.8 6.1A128 128 0 0 0 128 0M48 182.3q-.6 1.1-2.3.4c-.9-.4-1.4-1.3-1.1-1.9q.6-1 2.2-.4 1.6.8 1.1 2m6.2 5.7c-.6.5-1.8.3-2.6-.6q-1.2-1.7-.4-2.7 1.2-.8 2.7.6 1.3 1.5.3 2.7m4.4 7.1c-.8.6-2.1 0-2.9-1-.8-1.2-.8-2.6 0-3.1q1.4-.7 2.9 1c.8 1.2.8 2.6 0 3.1m7.3 8.4c-.7.7-2.2.5-3.3-.5s-1.5-2.5-.8-3.3c.8-.8 2.3-.5 3.4.5 1 1 1.4 2.5.7 3.3m9.4 2.8c-.3 1-1.7 1.4-3.2 1-1.4-.4-2.4-1.6-2.1-2.6s1.7-1.5 3.2-1 2.4 1.6 2.1 2.6m10.7 1.2q-.1 1.7-2.7 2-2.5-.2-2.8-2c0-1 1.2-1.9 2.8-1.9s2.7.8 2.7 1.9m10.6-.4c.2 1-.9 2-2.4 2.3s-2.8-.3-3-1.3c-.2-1.1.9-2.2 2.3-2.4 1.6-.3 3 .3 3.1 1.4"/></svg>`,
		displayName: "GitHub",
		pkce:        true, // technically is not supported yet but it is safe as the PKCE params are just ignored
		scopes:      []string{"read:user", "user:email"},
		authURL:     github.Endpoint.AuthURL,
		tokenURL:    github.Endpoint.TokenURL,
		userInfoURL: "https://api.github.com/user",
	}}
}

// FetchAuthUser returns an AuthUser instance based the Github's user api.
//
// API reference: https://docs.github.com/en/rest/reference/users#get-the-authenticated-user
func (p *Github) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserInfo(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Login     string `json:"login"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
		Id        int64  `json:"id"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           strconv.FormatInt(extracted.Id, 10),
		Name:         extracted.Name,
		Username:     extracted.Login,
		Email:        extracted.Email,
		AvatarURL:    extracted.AvatarURL,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	// in case user has set "Keep my email address private", send an
	// **optional** API request to retrieve the verified primary email
	if user.Email == "" {
		email, err := p.fetchPrimaryEmail(token)
		if err != nil {
			return nil, err
		}
		user.Email = email
	}

	return user, nil
}

// fetchPrimaryEmail sends an API request to retrieve the verified
// primary email, in case "Keep my email address private" was set.
//
// NB! This method can succeed and still return an empty email.
// Error responses that are result of insufficient scopes permissions are ignored.
//
// API reference: https://docs.github.com/en/rest/users/emails?apiVersion=2022-11-28
func (p *Github) fetchPrimaryEmail(token *oauth2.Token) (string, error) {
	client := p.Client(token)

	response, err := client.Get(p.userInfoURL + "/emails")
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// ignore common http errors caused by insufficient scope permissions
	// (the email field is optional, aka. return the auth user without it)
	if response.StatusCode == 401 || response.StatusCode == 403 || response.StatusCode == 404 {
		return "", nil
	}

	content, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	emails := []struct {
		Email    string
		Verified bool
		Primary  bool
	}{}
	if err := json.Unmarshal(content, &emails); err != nil {
		return "", err
	}

	// extract the verified primary email
	for _, email := range emails {
		if email.Verified && email.Primary {
			return email.Email, nil
		}
	}

	return "", nil
}
