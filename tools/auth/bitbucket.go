package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

var _ Provider = (*Bitbucket)(nil)

// NameBitbucket is the unique name of the Bitbucket provider.
const NameBitbucket = "bitbucket"

// Bitbucket is an auth provider for Bitbucket.
type Bitbucket struct {
	*baseProvider
}

// NewBitbucketProvider creates a new Bitbucket provider instance with some defaults.
func NewBitbucketProvider() *Bitbucket {
	return &Bitbucket{&baseProvider{
		ctx:         context.Background(),
		displayName: "Bitbucket",
		pkce:        false,
		scopes:      []string{"account"},
		authUrl:     "https://bitbucket.org/site/oauth2/authorize",
		tokenUrl:    "https://bitbucket.org/site/oauth2/access_token",
		userApiUrl:  "https://api.bitbucket.org/2.0/user",
	}}
}

// fetchUserEmail returns the primary email from the oauth2 token.
func (p *Bitbucket) fetchUserEmail(token *oauth2.Token) string {
	const EMAIL_API_URL = "https://api.bitbucket.org/2.0/user/emails"
	req, err := http.NewRequestWithContext(p.ctx, "GET", EMAIL_API_URL, nil)
	if err != nil {
		return ""
	}

	data, err := p.sendRawUserDataRequest(req, token)
	if err != nil {
		return ""
	}

	expected := struct {
		Values []struct {
			Email     string `json:"email"`
			IsPrimary bool   `json:"is_primary"`
		} `json:"values"`
	}{}
	if err := json.Unmarshal(data, &expected); err != nil {
		return ""
	}

	for _, v := range expected.Values {
		if v.IsPrimary {
			return v.Email
		}
	}

	return ""
}

// FetchAuthUser returns an AuthUser instance based on the Bitbucket API.
//
// API reference: https://developer.atlassian.com/cloud/bitbucket/oauth-2/
func (p *Bitbucket) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserData(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		UUID        string `json:"uuid"`
		Username    string `json:"username"`
		DisplayName string `json:"display_name"`
		Links       struct {
			Avatar struct {
				Href string `json:"href"`
			} `json:"avatar"`
		} `json:"links"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           extracted.UUID,
		Name:         extracted.DisplayName,
		Username:     extracted.Username,
		Email:        p.fetchUserEmail(token),
		AvatarUrl:    extracted.Links.Avatar.Href,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	return user, nil
}
