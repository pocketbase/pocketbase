package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

func init() {
	Providers[NameMonday] = wrapFactory(NewMondayProvider)
}

var _ Provider = (*Monday)(nil)

// NameMonday is the unique name of the Monday provider.
const NameMonday = "monday"

// Monday is an auth provider for monday.com.
type Monday struct {
	BaseProvider
}

// NewMondayProvider creates a new Monday provider instance with some defaults.
func NewMondayProvider() *Monday {
	return &Monday{BaseProvider{
		ctx:         context.Background(),
		displayName: "monday.com",
		pkce:        true,
		scopes:      []string{"me:read"},
		authURL:     "https://auth.monday.com/oauth2/authorize",
		tokenURL:    "https://auth.monday.com/oauth2/token",
		userInfoURL: "https://api.monday.com/v2",
	}}
}

// FetchAuthUser returns an AuthUser instance based on the Monday's user api.
//
// API reference: https://developer.monday.com/api-reference/reference/me
func (p *Monday) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserInfo(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Data struct {
			Me struct {
				Id         string `json:"id"`
				Enabled    bool   `json:"enabled"`
				Name       string `json:"name"`
				Email      string `json:"email"`
				IsVerified bool   `json:"is_verified"`
				Avatar     string `json:"photo_small"`
			} `json:"me"`
		} `json:"data"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	if !extracted.Data.Me.Enabled {
		return nil, errors.New("the monday.com user account is not enabled")
	}

	user := &AuthUser{
		Id:           extracted.Data.Me.Id,
		Name:         extracted.Data.Me.Name,
		AvatarURL:    extracted.Data.Me.Avatar,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	if extracted.Data.Me.IsVerified {
		user.Email = extracted.Data.Me.Email
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	return user, nil
}

// FetchRawUserInfo implements Provider.FetchRawUserInfo interface.
//
// monday.com doesn't have a UserInfo endpoint and information on the user
// is retrieved using their GraphQL API (https://developer.monday.com/api-reference/reference/me#queries)
func (p *Monday) FetchRawUserInfo(token *oauth2.Token) ([]byte, error) {
	query := []byte(`{"query": "query { me { id enabled name email is_verified photo_small }}"}`)
	bodyReader := bytes.NewReader(query)

	req, err := http.NewRequestWithContext(p.ctx, "POST", p.userInfoURL, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return p.sendRawUserInfoRequest(req, token)
}
