package auth

import (
	"context"
	"encoding/json"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

func init() {
	Providers[NameSpotify] = wrapFactory(NewSpotifyProvider)
}

var _ Provider = (*Spotify)(nil)

// NameSpotify is the unique name of the Spotify provider.
const NameSpotify string = "spotify"

// Spotify allows authentication via Spotify OAuth2.
type Spotify struct {
	BaseProvider
}

// NewSpotifyProvider creates a new Spotify provider instance with some defaults.
func NewSpotifyProvider() *Spotify {
	return &Spotify{BaseProvider{
		ctx:         context.Background(),
		order:       21,
		logo:        `<svg xmlns="http://www.w3.org/2000/svg" width="256" height="256" preserveAspectRatio="xMidYMid"><path fill="#1ed760" d="M128 0C57.3 0 0 57.3 0 128s57.3 128 128 128 128-57.3 128-128C256 57.31 198.7 0 128 0m58.7 184.61a7.97 7.97 0 0 1-10.98 2.65c-30.05-18.36-67.88-22.52-112.44-12.34a7.98 7.98 0 0 1-3.55-15.56c48.76-11.14 90.58-6.34 124.32 14.28a8 8 0 0 1 2.65 10.97m15.67-34.85a10 10 0 0 1-13.73 3.29c-34.4-21.15-86.85-27.27-127.55-14.92a10 10 0 0 1-12.45-6.65 10 10 0 0 1 6.65-12.45c46.49-14.1 104.28-7.27 143.79 17.01a10 10 0 0 1 3.29 13.72m1.34-36.3c-41.25-24.5-109.32-26.75-148.7-14.8a11.97 11.97 0 1 1-6.95-22.9c45.21-13.73 120.37-11.08 167.87 17.12a11.96 11.96 0 1 1-12.21 20.6z"/></svg>`,
		displayName: "Spotify",
		pkce:        true,
		scopes: []string{
			"user-read-private",
			// currently Spotify doesn't return information whether the email is verified or not
			// "user-read-email",
		},
		authURL:     spotify.Endpoint.AuthURL,
		tokenURL:    spotify.Endpoint.TokenURL,
		userInfoURL: "https://api.spotify.com/v1/me",
	}}
}

// FetchAuthUser returns an AuthUser instance based on the Spotify's user api.
//
// API reference: https://developer.spotify.com/documentation/web-api/reference/#/operations/get-current-users-profile
func (p *Spotify) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserInfo(token)
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
			URL string `json:"url"`
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

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	if len(extracted.Images) > 0 {
		user.AvatarURL = extracted.Images[0].URL
	}

	return user, nil
}
