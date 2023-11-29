package auth

import (
	"context"
	"encoding/json"

	"github.com/pocketbase/pocketbase/tools/types"
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
		ctx:         context.Background(),
		displayName: "Google",
		pkce:        true,
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
	data, err := p.FetchRawUserData(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Id            string `json:"id"`
		Name          string `json:"name"`
		Email         string `json:"email"`
		Picture       string `json:"picture"`
		VerifiedEmail bool   `json:"verified_email"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           extracted.Id,
		Name:         extracted.Name,
		AvatarUrl:    extracted.Picture,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	if extracted.VerifiedEmail {
		user.Email = extracted.Email
	}

	return user, nil
}
