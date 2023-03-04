package auth

import (
	"encoding/json"
	"errors"
	"strings"

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
			"crm.objects.owners.read",
		},
		authUrl:    "https://app.hubspot.com/oauth/authorize",
		tokenUrl:   "https://api.hubapi.com/oauth/v1/token",
		userApiUrl: "https://api.hubapi.com/crm/v3/owners",
	}}
}

// FetchAuthUser returns an AuthUser instance based on the Hubspot user API.
// https://developers.hubspot.com/docs/api/crm/owners
func (p *Hubspot) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserData(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	var resp struct {
		Results []struct {
			ID        string `json:"id"`
			Email     string `json:"email"`
			FirstName string `json:"firstName"`
			LastName  string `json:"lastName"`
			Archived  bool   `json:"archived"`
		} `json:"results"`
	}

	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	// This case shouldn't happen.
	if len(resp.Results) == 0 {
		return nil, errors.New("no owner found")
	}

	profile := resp.Results[0]
	if profile.Archived {
		return nil, errors.New("owner is archived")
	}

	return &AuthUser{
		Id:           profile.ID,
		Email:        profile.Email,
		Name:         strings.Join([]string{profile.FirstName, profile.LastName}, " "),
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}
