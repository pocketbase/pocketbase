package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

var _ Provider = (*Monday)(nil)

// NameMonday is the unique name of the Monday provider.
const NameMonday = "monday"

// Monday is an auth provider for monday.com.
type Monday struct {
	*baseProvider
}

// NewMondayProvider creates a new Monday provider instance with some defaults.
func NewMondayProvider() *Monday {
	return &Monday{&baseProvider{
		ctx:         context.Background(),
		displayName: "monday.com",
		pkce:        true,
		scopes:      []string{"me:read"},
		authUrl:     "https://auth.monday.com/oauth2/authorize",
		tokenUrl:    "https://auth.monday.com/oauth2/token",
		userApiUrl:  "https://api.monday.com/v2",
	}}
}

// FetchAuthUser returns an AuthUser instance based on the Monday's user api.
//
// API reference: https://developer.monday.com/api-reference/reference/me
func (p *Monday) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserData(token)
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
				ID         string `json:"id"`
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

	user := &AuthUser{
		Id:           extracted.Data.Me.ID,
		Name:         extracted.Data.Me.Name,
		AvatarUrl:    extracted.Data.Me.Avatar,
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

// FetchRawUserData implements Provider.FetchRawUserData interface.
//
// monday.com doesn't have a UserInfo endpoint and information on the user
// is retrieved using their GraphQL API (https://developer.monday.com/api-reference/reference/me#queries)
func (p *Monday) FetchRawUserData(token *oauth2.Token) ([]byte, error) {
	query := []byte(`{"query": "query { me { id name email is_verified photo_small }}"}`)
	bodyReader := bytes.NewReader(query)

	req, err := http.NewRequestWithContext(p.ctx, "POST", p.userApiUrl, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return p.sendRawUserDataRequest(req, token)
}
