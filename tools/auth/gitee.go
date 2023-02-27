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
		Id:           strconv.Itoa(extracted.Id),
		Name:         extracted.Name,
		Username:     extracted.Login,
		AvatarUrl:    extracted.AvatarUrl,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	if extracted.Email != "" && is.EmailFormat.Validate(extracted.Email) == nil {
		// valid public primary email
		user.Email = extracted.Email
	} else {
		// send an additional optional request to retrieve the email
		email, err := p.fetchPrimaryEmail(token)
		if err != nil {
			return nil, err
		}
		user.Email = email
	}

	return user, nil
}

// fetchPrimaryEmail sends an API request to retrieve the verified primary email,
// in case the user hasn't set "Public email address" or has unchecked
// the "Access your emails data" permission during authentication.
//
// NB! This method can succeed and still return an empty email.
// Error responses that are result of insufficient scopes permissions are ignored.
//
// API reference: https://gitee.com/api/v5/swagger#/getV5Emails
func (p *Gitee) fetchPrimaryEmail(token *oauth2.Token) (string, error) {
	client := p.Client(token)

	response, err := client.Get("https://gitee.com/api/v5/emails")
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// ignore common http errors caused by insufficient scope permissions
	if response.StatusCode == 401 || response.StatusCode == 403 || response.StatusCode == 404 {
		return "", nil
	}

	content, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	emails := []struct {
		Email string
		State string
		Scope []string
	}{}
	if err := json.Unmarshal(content, &emails); err != nil {
		// ignore unmarshal error in case "Keep my email address private"
		// was set because response.Body will be something like:
		// {"email":"12285415+test@user.noreply.gitee.com"}
		return "", nil
	}

	// extract the first verified primary email
	for _, email := range emails {
		for _, scope := range email.Scope {
			if email.State == "confirmed" && scope == "primary" && is.EmailFormat.Validate(email.Email) == nil {
				return email.Email, nil
			}
		}
	}

	return "", nil
}
