// Strva Authentication documentation: https://developers.strava.com/docs/authentication/
package auth

import (
	"fmt"

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
// Requested scopes, as a comma delimited string, e.g. "activity:read_all,activity:write". Applications should request only the scopes required for the application to function normally. The scope activity:read is required for activity webhooks.
func NewStravaProvider() *Strava {
	return &Strava{&baseProvider{
		scopes: []string{
			"read_all,profile:read_all,activity:read_all",
		},
		authUrl:    "https://www.strava.com/oauth/authorize",
		tokenUrl:   "https://www.strava.com/oauth/token",
		userApiUrl: "https://www.strava.com/api/v3/athlete",
	}}
}

// FetchAuthUser returns an AuthUser instance based the Strava's user api.
// The Strava's userApiUrl does not return an email but the UserOauth2Login.Submit requires a valid email adress to find / save the user
// HACK: The email adress is as follows: "STRAVAID@strava.com"
// TODO: Implement an OAuth strategy which does not require an email adress and sets the user as "active" but email adress to not valid
func (p *Strava) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	/* Stava authenticated athlete json
	{
		"id": 8844168,
		"username": "username",
		"resource_state": 2,
		"firstname": "Firstname",
		"lastname": "latsname",
		"bio": "",
		"city": "CityName",
		"state": "StateName",
		"country": "CountryName",
		"sex": "Sex",
		"premium": true,
		"summit": true,
		"created_at": "2015-04-23T19:49:36Z",
		"updated_at": "2022-07-10T14:08:54Z",
		"badge_type_id": 1,
		"weight": 76.0,
		"profile_medium": "ProfilePictureMedium",
		"profile": "ProfilePictureLarge",
		"friend": null,
		"follower": null
	}
	*/
	rawData := struct {
		LocalId     int    `json:"id"`
		DisplayName string `json:"username"`
		PhotoUrl    string `json:"profile"`
	}{}

	if err := p.FetchRawUserData(token, &rawData); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:        fmt.Sprintf("%d", rawData.LocalId),
		Name:      rawData.DisplayName,
		Email:     fmt.Sprintf("%d@strava.com", rawData.LocalId),
		AvatarUrl: rawData.PhotoUrl,
	}

	return user, nil
}
