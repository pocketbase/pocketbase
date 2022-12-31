package auth

import (
	"encoding/json"
	"io"
	"strconv"

	"golang.org/x/oauth2"
)

var _ Provider = (*Gitee)(nil)

// NameGitee is the unique name of the Gitee provider.
const NameGitee string = "gitee"

// Gitee allows authentication via Gitee OAuth2.
type Gitee struct {
	*baseProvider
}

// NewGiteeProvider creates new Gitee provider instance with some defaults.
func NewGiteeProvider() *Gitee {
	return &Gitee{&baseProvider{
		scopes:     []string{"user_info", "emails"},
		authUrl:    "https://gitee.com/oauth/authorize",
		tokenUrl:   "https://gitee.com/oauth/token",
		userApiUrl: "https://gitee.com/api/v5/user",
	}}
}

// FetchAuthUser returns an AuthUser instance based the Gitee's user api.
//
// API reference: https://gitee.com/api/v5/swagger#/getV5User
func (p *Gitee) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserData(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Login     string `json:"login"`
		Id        int    `json:"id"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarUrl string `json:"avatar_url"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:          strconv.Itoa(extracted.Id),
		Name:        extracted.Name,
		Username:    extracted.Login,
		Email:       extracted.Email,
		AvatarUrl:   extracted.AvatarUrl,
		RawUser:     rawUser,
		AccessToken: token.AccessToken,
	}

	// in case user set "Keep my email address private",
	// email should be retrieved via extra API request
	if user.Email == "" {
		client := p.Client(token)

		response, err := client.Get("https://gitee.com/api/v5/emails")
		if err != nil {
			return user, err
		}
		defer response.Body.Close()

		content, err := io.ReadAll(response.Body)
		if err != nil {
			return user, err
		}

		emails := []struct {
			Email string
			State string
			Scope []string
		}{}
		if err := json.Unmarshal(content, &emails); err != nil {
			return user, err
		}

		// extract the verified primary email

		//
		// API reference: https://gitee.com/api/v5/swagger#/getV5Emails
	outer:
		for _, email := range emails {
			for _, scope := range email.Scope {
				if email.State == "confirmed" && scope == "primary" {
					user.Email = email.Email
					break outer
				}
			}
		}
	}

	return user, nil
}
