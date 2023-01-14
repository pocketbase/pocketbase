package auth

import (
	"encoding/json"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

var _ Provider = (*Spotify)(nil)

// NameSpotify is the unique name of the Spotify provider.
const NameSpotify string = "spotify"

// Spotify allows authentication via Spotify OAuth2.
type Spotify struct {
	*baseProvider
}

// NewSpotifyProvider creates a new Spotify provider instance with some defaults.
func NewSpotifyProvider() *Spotify {
	return &Spotify{&baseProvider{
		scopes: []string{
			"user-read-private",
			// currently Spotify doesn't return information whether the email is verified or not
			// "user-read-email",
		},
		authUrl:    spotify.Endpoint.AuthURL,
		tokenUrl:   spotify.Endpoint.TokenURL,
		userApiUrl: "https://api.spotify.com/v1/me",
	}}
}

// FetchAuthUser returns an AuthUser instance based on the Spotify's user api.
//
// API reference: https://developer.spotify.com/documentation/web-api/reference/#/operations/get-current-users-profile
func (p *Spotify) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserData(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Id     string `json:"id"`
		Name   string `json:"display_name"`
		Images []struct {
			Url string `json:"url"`
		} `json:"images"`
		// don't map the email because per the official docs
		// the email field is "unverified" and there is no proof
		// that it actually belongs to the user
		// Email  string `json:"email"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           extracted.Id,
		Name:         extracted.Name,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
	if len(extracted.Images) > 0 {
		user.AvatarUrl = extracted.Images[0].Url
	}

	return user, nil
}
