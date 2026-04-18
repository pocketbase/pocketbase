package auth

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

func init() {
	Providers[NameBitbucket] = wrapFactory(NewBitbucketProvider)
}

var _ Provider = (*Bitbucket)(nil)

// NameBitbucket is the unique name of the Bitbucket provider.
const NameBitbucket = "bitbucket"

// Bitbucket is an auth provider for Bitbucket.
type Bitbucket struct {
	BaseProvider
}

// NewBitbucketProvider creates a new Bitbucket provider instance with some defaults.
func NewBitbucketProvider() *Bitbucket {
	return &Bitbucket{BaseProvider{
		ctx:         context.Background(),
		order:       9,
		logo:        `<svg xmlns="http://www.w3.org/2000/svg" width="2500" height="2256" preserveAspectRatio="xMidYMid" viewBox="-1 -0.6 257.9 230.8"><linearGradient id="a" x1="108.6%" x2="46.9%" y1="13.8%" y2="78.8%"><stop offset=".2" stop-color="#0052cc"/><stop offset="1" stop-color="#2684ff"/></linearGradient><g fill="none"><path d="M101 153h54l13-76H87z"/><path fill="#2684ff" d="M8 0a8 8 0 0 0-8 10l35 211a11 11 0 0 0 11 9h167a8 8 0 0 0 8-7l35-213a8 8 0 0 0-8-10zm147 153h-53L87 77h81z"/><path fill="url(#a)" d="M245 77h-77l-13 76h-53l-63 74 7 3h167a8 8 0 0 0 8-7z"/></g></svg>`,
		displayName: "Bitbucket",
		pkce:        false,
		scopes:      []string{"account"},
		authURL:     "https://bitbucket.org/site/oauth2/authorize",
		tokenURL:    "https://bitbucket.org/site/oauth2/access_token",
		userInfoURL: "https://api.bitbucket.org/2.0/user",
	}}
}

// FetchAuthUser returns an AuthUser instance based on the Bitbucket's user API.
//
// API reference: https://developer.atlassian.com/cloud/bitbucket/rest/api-group-users/#api-user-get
func (p *Bitbucket) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserInfo(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		UUID          string `json:"uuid"`
		Username      string `json:"username"`
		DisplayName   string `json:"display_name"`
		AccountStatus string `json:"account_status"`
		Links         struct {
			Avatar struct {
				Href string `json:"href"`
			} `json:"avatar"`
		} `json:"links"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	if extracted.AccountStatus != "active" {
		return nil, errors.New("the Bitbucket user is not active")
	}

	email, err := p.fetchPrimaryEmail(token)
	if err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           extracted.UUID,
		Name:         extracted.DisplayName,
		Username:     extracted.Username,
		Email:        email,
		AvatarURL:    extracted.Links.Avatar.Href,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	return user, nil
}

// fetchPrimaryEmail sends an API request to retrieve the first
// verified primary email.
//
// NB! This method can succeed and still return an empty email.
// Error responses that are result of insufficient scopes permissions are ignored.
//
// API reference: https://developer.atlassian.com/cloud/bitbucket/rest/api-group-users/#api-user-emails-get
func (p *Bitbucket) fetchPrimaryEmail(token *oauth2.Token) (string, error) {
	response, err := p.Client(token).Get(p.userInfoURL + "/emails")
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// ignore common http errors caused by insufficient scope permissions
	// (the email field is optional, aka. return the auth user without it)
	if response.StatusCode >= 400 {
		return "", nil
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	expected := struct {
		Values []struct {
			Email     string `json:"email"`
			IsPrimary bool   `json:"is_primary"`
		} `json:"values"`
	}{}
	if err := json.Unmarshal(data, &expected); err != nil {
		return "", err
	}

	for _, v := range expected.Values {
		if v.IsPrimary {
			return v.Email, nil
		}
	}

	return "", nil
}
