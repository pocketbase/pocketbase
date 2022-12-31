package auth

import (
	"encoding/json"
	"io"
	"strconv"

	"github.com/go-ozzo/ozzo-validation/v4/is"
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
		AvatarUrl:   extracted.AvatarUrl,
		RawUser:     rawUser,
		AccessToken: token.AccessToken,
	}

	// extract the email when it is available
	if extracted.Email != "" && is.EmailFormat.Validate(extracted.Email) == nil {
		user.Email = extracted.Email
		return user, nil
	}

	// Optional: in case user set "Keep my email address private",
	// or uncheck the "all emails" checkbox, we need to
	// retrieve emails from the emails api.
	// All response errors in this request will be ignored.

	client := p.Client(token)

	response, err := client.Get("https://gitee.com/api/v5/emails")
	if err != nil {
		return user, nil
	}
	defer response.Body.Close()

	content, err := io.ReadAll(response.Body)
	if err != nil {
		return user, nil
	}

	emails := []struct {
		Email string
		State string
		Scope []string
	}{}
	if err := json.Unmarshal(content, &emails); err != nil {
		return user, nil
	}

	// extract the verified primary email

	//
	// API reference: https://gitee.com/api/v5/swagger#/getV5Emails
	for _, email := range emails {
		for _, scope := range email.Scope {
			if email.State == "confirmed" && scope == "primary" {
				user.Email = email.Email
				return user, nil
			}
		}
	}

	return user, nil
}
