package auth

import (
	"encoding/json"
	"strconv"

	"golang.org/x/oauth2"
)

var _ Provider = (*Hubspot)(nil)

// NameHubspot is the unique name of the Hubspot provider.
const NameHubspot string = "hubspot"

// Hubspot allows authentication via Hubspot OAuth2.
type Hubspot struct {
	*baseProvider
}

// NewHubspotProvider creates a new Hubspot provider instance with some defaults.
func NewHubspotProvider() *Hubspot {
	return &Hubspot{&baseProvider{
		scopes: []string{
			"contacts",
		},
		authUrl:    "https://app.hubspot.com/oauth/authorize",
		tokenUrl:   "https://api.hubapi.com/oauth/v1/token",
		userApiUrl: "https://api.hubapi.com/contacts/v1/contact/vid/",
	}}
}

// FetchAuthUser returns an AuthUser instance based on the Hubspot user API.
func (p *Hubspot) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserData(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	var user struct {
		Vid     int `json:"vid"`
		Profile struct {
			Email       string `json:"email"`
			FirstName   string `json:"firstname"`
			LastName    string `json:"lastname"`
			AvatarUrl   string `json:"avatar-url"`
			Company     string `json:"company"`
			JobTitle    string `json:"jobtitle"`
			City        string `json:"city"`
			State       string `json:"state"`
			Zip         string `json:"zip"`
			Country     string `json:"country"`
			PhoneNumber string `json:"phone"`
		} `json:"properties"`
	}

	if err := json.Unmarshal(data, &user); err != nil {
		return nil, err
	}

	authUser := &AuthUser{
		Id:           strconv.Itoa(user.Vid),
		Email:        user.Profile.Email,
		Name:         user.Profile.FirstName + " " + user.Profile.LastName,
		AvatarUrl:    user.Profile.AvatarUrl,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	return authUser, nil
}
