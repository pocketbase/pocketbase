package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

var _ Provider = (*Gitea)(nil)

// NameGitea is the unique name of the Gitea provider.
const NameGitea string = "gitea"

// Gitea allows authentication via Gitea OAuth2.
type Gitea struct {
	*baseProvider
}

// NewGiteaProvider creates new Gitea provider instance with some defaults.
func NewGiteaProvider() *Gitea {
	return &Gitea{&baseProvider{
		ctx:         context.Background(),
		displayName: "Gitea",
		pkce:        true,
		scopes:      []string{"read:user", "user:email"},
		authUrl:     "https://gitea.com/login/oauth/authorize",
		tokenUrl:    "https://gitea.com/login/oauth/access_token",
		userApiUrl:  "https://gitea.com/api/v1/user",
	}}
}

// FetchAuthUser returns an AuthUser instance based on Gitea's user api.
//
// API reference: https://try.gitea.io/api/swagger#/user/userGetCurrent
func (p *Gitea) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserData(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Name      string `json:"full_name"`
		Username  string `json:"login"`
		AvatarUrl string `json:"avatar_url"`
		Id        int    `json:"id"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           strconv.Itoa(extracted.Id),
		Name:         extracted.Name,
		Username:     extracted.Username,
		AvatarUrl:    extracted.AvatarUrl,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	email, err := p.fetchVerifiedPrimaryEmail(token)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch primary email: %w", err)
	}
	user.Email = email

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	return user, nil
}

// fetchVerifiedPrimaryEmail sends an API request to retrieve the verified
// primary email, in case "Keep my email address private" was set.
//
// NB! This method can succeed and still return an empty email.
// Error responses that are result of insufficient scopes permissions are ignored.
//
// API reference: https://codeberg.org/api/swagger#/user/userListEmails
func (p *Gitea) fetchVerifiedPrimaryEmail(token *oauth2.Token) (string, error) {
	client := p.Client(token)

	response, err := client.Get(p.userApiUrl + "/emails")
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
