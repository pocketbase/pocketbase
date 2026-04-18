package auth

import (
	"context"
	"encoding/json"
	"io"
	"strconv"

	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

func init() {
	Providers[NameGitee] = wrapFactory(NewGiteeProvider)
}

var _ Provider = (*Gitee)(nil)

// NameGitee is the unique name of the Gitee provider.
const NameGitee string = "gitee"

// Gitee allows authentication via Gitee OAuth2.
type Gitee struct {
	BaseProvider
}

// NewGiteeProvider creates new Gitee provider instance with some defaults.
func NewGiteeProvider() *Gitee {
	return &Gitee{BaseProvider{
		ctx:         context.Background(),
		order:       10,
		logo:        `<svg xmlns="http://www.w3.org/2000/svg" viewBox="120 13 72 72"><g fill="none" fill-rule="evenodd"><path d="M0 0h312v100H0z"/><path fill="#c71d23" d="M156 85a36 36 0 1 1 0-72 36 36 0 0 1 0 72m18.2-40h-20.4q-1.7.1-1.8 1.8v4.4q.2 1.6 1.8 1.8h12.4q1.7.1 1.8 1.8v.9c0 3-2.4 5.3-5.3 5.3h-17q-1.6-.1-1.7-1.8V42.3c0-3 2.4-5.3 5.3-5.3h25q1.5-.1 1.7-1.8v-4.4a2 2 0 0 0-1.8-1.8h-24.9C142 29 136 35 136 42.3v25q.2 1.5 1.8 1.7H164a12 12 0 0 0 12-12V46.8q-.2-1.6-1.8-1.8"/></g></svg>`,
		displayName: "Gitee",
		pkce:        true,
		scopes:      []string{"user_info", "emails"},
		authURL:     "https://gitee.com/oauth/authorize",
		tokenURL:    "https://gitee.com/oauth/token",
		userInfoURL: "https://gitee.com/api/v5/user",
	}}
}

// FetchAuthUser returns an AuthUser instance based the Gitee's user api.
//
// API reference: https://gitee.com/api/v5/swagger#/getV5User
func (p *Gitee) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
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
		AvatarURL:    extracted.AvatarURL,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

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
