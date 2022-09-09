package auth

import (
	"golang.org/x/oauth2"
)

var _ Provider = (*Google)(nil)

// NameGoogle is the unique name of the Google provider.
const NameGoogle string = "google"

// Google allows authentication via Google OAuth2.
type Google struct {
	*baseProvider
}

// NewGoogleProvider creates new Google provider instance with some defaults.
func NewGoogleProvider() *Google {
	return &Google{&baseProvider{
		scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		authUrl:    "https://accounts.google.com/o/oauth2/auth",
		tokenUrl:   "https://accounts.google.com/o/oauth2/token",
		userApiUrl: "https://www.googleapis.com/oauth2/v1/userinfo",
	}}
}

// FetchAuthUser returns an AuthUser instance based the Google's user api.
func (p *Google) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	rawData := struct {
		Id      string `json:"id"`
		Name    string `json:"name"`
		Email   string `json:"email"`
		Picture string `json:"picture"`
	}{}

	if err := p.FetchRawUserData(token, &rawData); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:        rawData.Id,
		Name:      rawData.Name,
		Email:     rawData.Email,
		AvatarUrl: rawData.Picture,
	}

	return user, nil
}
