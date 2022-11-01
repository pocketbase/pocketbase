package auth

import (
	"log"

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
		scopes:     []string{"user-read-private", "user-read-email"},
		authUrl:    spotify.Endpoint.AuthURL,
		tokenUrl:   spotify.Endpoint.TokenURL,
		userApiUrl: "https://api.spotify.com/v1/me",
	}}
}

// FetchAuthUser returns an AuthUser instance based on the Spotify's user api.
func (p *Spotify) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	rawData := struct {
		Id    string `json:"id"`
		Name  string `json:"display_name"`
		Email string `json:"email"`
	}{}

	if err := p.FetchRawUserData(token, &rawData); err != nil {
		return nil, err
	}
	log.Printf("rawData: %+v", rawData)

	user := &AuthUser{
		Id:    rawData.Id,
		Name:  rawData.Name,
		Email: rawData.Email,
	}

	return user, nil
}
