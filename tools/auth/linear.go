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
	Providers[NameLinear] = wrapFactory(NewLinearProvider)
}

var _ Provider = (*Linear)(nil)

// NameLinear is the unique name of the Linear provider.
const NameLinear string = "linear"

// Linear allows authentication via Linear OAuth2.
type Linear struct {
	BaseProvider
}

// NewLinearProvider creates new Linear provider instance with some defaults.
//
// API reference: https://developers.linear.app/docs/oauth/authentication
func NewLinearProvider() *Linear {
	return &Linear{BaseProvider{
		ctx:         context.Background(),
		order:       16,
		logo:        `<svg xmlns="http://www.w3.org/2000/svg" width="200" height="200" fill="#222326" viewBox="0 0 100 100"><path d="M1.2 61.5c-.2-1 1-1.5 1.6-.8l36.5 36.5c.7.7.1 1.8-.8 1.6A50 50 0 0 1 1.2 61.5M0 47q0 .4.3.7l52 52.1q.4.3.8.3a50 50 0 0 0 7-1 1 1 0 0 0 .5-1.6l-58-58a1 1 0 0 0-1.7.5 50 50 0 0 0-.9 7m4.2-17.2q-.2.6.2 1.1l64.8 64.8q.5.5 1 .2l5.3-2.7q.8-.6.2-1.5L8.4 24.3a1 1 0 0 0-1.5.2Q5.4 27 4.2 29.7m8.5-11.6a1 1 0 0 1 0-1.4A50 50 0 0 1 100 50a50 50 0 0 1-16.7 37.4 1 1 0 0 1-1.4 0z"/></svg>`,
		displayName: "Linear",
		pkce:        false, // Linear doesn't support PKCE at the moment and returns an error if enabled
		scopes:      []string{"read"},
		authURL:     "https://linear.app/oauth/authorize",
		tokenURL:    "https://api.linear.app/oauth/token",
		userInfoURL: "https://api.linear.app/graphql",
	}}
}

// FetchAuthUser returns an AuthUser instance based on the Linear's user api.
//
// API reference: https://developers.linear.app/docs/graphql/working-with-the-graphql-api#authentication
func (p *Linear) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
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
			Viewer struct {
				Id          string `json:"id"`
				DisplayName string `json:"displayName"`
				Name        string `json:"name"`
				Email       string `json:"email"`
				AvatarURL   string `json:"avatarUrl"`
				Active      bool   `json:"active"`
			} `json:"viewer"`
		} `json:"data"`
	}{}

	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	if !extracted.Data.Viewer.Active {
		return nil, errors.New("the Linear user account is not active")
	}

	user := &AuthUser{
		Id:           extracted.Data.Viewer.Id,
		Name:         extracted.Data.Viewer.Name,
		Username:     extracted.Data.Viewer.DisplayName,
		Email:        extracted.Data.Viewer.Email,
		AvatarURL:    extracted.Data.Viewer.AvatarURL,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	return user, nil
}

// FetchRawUserInfo implements Provider.FetchRawUserInfo interface method.
//
// Linear doesn't have a UserInfo endpoint and information on the user
// is retrieved using their GraphQL API (https://developers.linear.app/docs/graphql/working-with-the-graphql-api#queries-and-mutations)
func (p *Linear) FetchRawUserInfo(token *oauth2.Token) ([]byte, error) {
	query := []byte(`{"query": "query Me { viewer { id displayName name email avatarUrl active } }"}`)
	bodyReader := bytes.NewReader(query)

	req, err := http.NewRequestWithContext(p.ctx, "POST", p.userInfoURL, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return p.sendRawUserInfoRequest(req, token)
}
