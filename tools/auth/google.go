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
		displayName: "Google",
		pkce:        true,
		scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		authURL:     "https://accounts.google.com/o/oauth2/auth",
		tokenURL:    "https://accounts.google.com/o/oauth2/token",
		userInfoURL: "https://www.googleapis.com/userinfo/v2/me",
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
		AvatarURL:    extracted.Picture,
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
