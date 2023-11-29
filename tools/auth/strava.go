package auth

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

var _ Provider = (*Strava)(nil)

// NameStrava is the unique name of the Strava provider.
const NameStrava string = "strava"

// Strava allows authentication via Strava OAuth2.
type Strava struct {
	*baseProvider
}

// NewStravaProvider creates new Strava provider instance with some defaults.
func NewStravaProvider() *Strava {
	return &Strava{&baseProvider{
		ctx:         context.Background(),
		displayName: "Strava",
		pkce:        true,
		scopes: []string{
			"profile:read_all",
		},
		authUrl:    "https://www.strava.com/oauth/authorize",
		tokenUrl:   "https://www.strava.com/api/v3/oauth/token",
		userApiUrl: "https://www.strava.com/api/v3/athlete",
	}}
}

// FetchAuthUser returns an AuthUser instance based on the Strava's user api.
//
// API reference: https://developers.strava.com/docs/authentication/
func (p *Strava) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserData(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Id              int    `json:"id"`
		FirstName       string `json:"firstname"`
		LastName        string `json:"lastname"`
		Username        string `json:"username"`
		ProfileImageUrl string `json:"profile"`

		// At the time of writing, Strava OAuth2 doesn't support returning the user email address
		// Email string `json:"email"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Name:         extracted.FirstName + " " + extracted.LastName,
		Username:     extracted.Username,
		AvatarUrl:    extracted.ProfileImageUrl,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	if extracted.Id != 0 {
		user.Id = strconv.Itoa(extracted.Id)
	}

	return user, nil
}
