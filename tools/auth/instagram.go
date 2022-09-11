package auth

import (
	"golang.org/x/oauth2"
)

var _ Provider = (*Instagram)(nil)

// NameInstagram is the unique name of the Instagram provider.
const NameInstagram string = "instagram"

// Instagram allows authentication via Instagram OAuth2.
type Instagram struct {
	*baseProvider
}

// NewInstagramProvider creates a new Instagram provider instance with some defaults.
func NewInstagramProvider() *Instagram {
	// https://developers.facebook.com/docs/instagram-basic-display-api/reference/oauth-authorize
	// https://developers.facebook.com/docs/instagram-basic-display-api/reference/me
	return &Instagram{&baseProvider{
		scopes:     []string{"user_profile"},
		authUrl:    "https://api.instagram.com/oauth/authorize",
		tokenUrl:   "https://api.instagram.com/oauth/access_token",
		userApiUrl: "https://graph.instagram.com/v14.0/me",
	}}
}

// FetchAuthUser returns an AuthUser instance from Instagram's user api.
func (p *Instagram) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	// https://developers.facebook.com/docs/instagram-basic-display-api/reference/user
	rawData := struct {
		Id            string `json:"id"`
		Username      string `json:"username"`
		
		// At the time of writing, Instagram OAuth2 doesn't support returning the user email address
		Email         string `json:"email"`
	}{}

	if err := p.FetchRawUserData(token, &rawData); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:        rawData.Id,
		Username:  rawData.Username,
		Email:     rawData.Email,
	}

	return user, nil
}
