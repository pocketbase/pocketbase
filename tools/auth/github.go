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
