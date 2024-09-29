package auth

import (
	"context"
	"encoding/json"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

func init() {
	Providers[NameMicrosoft] = wrapFactory(NewMicrosoftProvider)
}

var _ Provider = (*Microsoft)(nil)

// NameMicrosoft is the unique name of the Microsoft provider.
const NameMicrosoft string = "microsoft"

// Microsoft allows authentication via AzureADEndpoint OAuth2.
type Microsoft struct {
	BaseProvider
}

// NewMicrosoftProvider creates new Microsoft AD provider instance with some defaults.
func NewMicrosoftProvider() *Microsoft {
	endpoints := microsoft.AzureADEndpoint("")
	return &Microsoft{BaseProvider{
		ctx:         context.Background(),
		displayName: "Microsoft",
		pkce:        true,
		scopes:      []string{"User.Read"},
		authURL:     endpoints.AuthURL,
		tokenURL:    endpoints.TokenURL,
		userInfoURL: "https://graph.microsoft.com/v1.0/me",
	}}
}

// FetchAuthUser returns an AuthUser instance based on the Microsoft's user api.
//
// API reference:  https://learn.microsoft.com/en-us/azure/active-directory/develop/userinfo
// Graph explorer: https://developer.microsoft.com/en-us/graph/graph-explorer
func (p *Microsoft) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	data, err := p.FetchRawUserInfo(token)
	if err != nil {
		return nil, err
	}

	rawUser := map[string]any{}
	if err := json.Unmarshal(data, &rawUser); err != nil {
		return nil, err
	}

	extracted := struct {
		Id    string `json:"id"`
		Name  string `json:"displayName"`
		Email string `json:"mail"`
	}{}
	if err := json.Unmarshal(data, &extracted); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:           extracted.Id,
		Name:         extracted.Name,
		Email:        extracted.Email,
		RawUser:      rawUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user.Expiry, _ = types.ParseDateTime(token.Expiry)

	return user, nil
}
