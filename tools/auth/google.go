package auth

import (
	"context"
	"encoding/json"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

func init() {
	Providers[NameGoogle] = wrapFactory(NewGoogleProvider)
}

var _ Provider = (*Google)(nil)

// NameGoogle is the unique name of the Google provider.
const NameGoogle string = "google"

// Google allows authentication via Google OAuth2.
type Google struct {
	BaseProvider
}

// NewGoogleProvider creates new Google provider instance with some defaults.
func NewGoogleProvider() *Google {
	return &Google{BaseProvider{
		ctx:         context.Background(),
		order:       2,
		logo:        `<svg xmlns="http://www.w3.org/2000/svg" width="256" height="262" preserveAspectRatio="xMidYMid"><path fill="#4285f4" d="M255.9 133.5c0-10.8-.9-18.6-2.8-26.7H130.6v48.4h71.9a64 64 0 0 1-26.7 42.4l-.2 1.6 38.7 30 2.7.3c24.7-22.8 38.9-56.3 38.9-96"/><path fill="#34a853" d="M130.6 261.1c35.2 0 64.8-11.6 86.4-31.6l-41.2-32a76 76 0 0 1-45.2 13.1 79 79 0 0 1-74.3-54.2l-1.5.1-40.3 31.2-.6 1.5A131 131 0 0 0 130.6 261"/><path fill="#fbbc05" d="M56.3 156.4a80 80 0 0 1-.2-51.7V103L15.3 71.3l-1.4.6a131 131 0 0 0 0 117.3z"/><path fill="#eb4335" d="M130.6 50.5c24.5 0 41 10.6 50.4 19.4L218 34c-22.8-21-52.2-34-87.4-34C79.5 0 35.4 29.3 13.9 72l42.2 32.7a79 79 0 0 1 74.5-54.2"/></svg>`,
		displayName: "Google",
		pkce:        true,
		scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		authURL:     "https://accounts.google.com/o/oauth2/v2/auth",
		tokenURL:    "https://oauth2.googleapis.com/token",
		userInfoURL: "https://www.googleapis.com/oauth2/v3/userinfo",
	}}
}

// FetchAuthUser returns an AuthUser instance based the Google's user api.
func (p *Google) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserInfo(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Id            string `json:"sub"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           extracted.Id,
		Name:         extracted.Name,
		AvatarURL:    extracted.Picture,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	if extracted.EmailVerified {
		user.Email = extracted.Email
	}

	return user, nil
}
