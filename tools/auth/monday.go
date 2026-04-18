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
		order:       18,
		logo:        `<svg xmlns="http://www.w3.org/2000/svg" xml:space="preserve" viewBox="0 0 127 127"><path fill="#fb275d" d="M24.8 87.2c-3.7 0-7.1-1.9-8.9-5.1a9 9 0 0 1 .3-9.9L34.5 44c1.9-3.1 5.4-4.9 9.1-4.8s7.1 2.1 8.8 5.3 1.5 7-.5 9.9L33.4 82.6a10 10 0 0 1-8.6 4.6"/><path fill="#fc0" d="M56.1 87.2c-3.7 0-7.1-1.9-8.9-5a9 9 0 0 1 .3-9.9l18.4-28.1c1.9-3.1 5.3-5 9.1-4.9s7.1 2.1 8.8 5.3 1.5 7-.7 10l-18.4 28a11 11 0 0 1-8.6 4.6"/><path fill="#00ca72" d="M86.7 87.2c5.6 0 10.2-4.6 10.2-10.2s-4.6-10.2-10.2-10.2S76.5 71.4 76.5 77c0 5.7 4.5 10.2 10.2 10.2"/></svg>`,
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
