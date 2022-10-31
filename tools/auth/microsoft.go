package auth

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

var _ Provider = (*Microsoft)(nil)

// NameMicrosoft is the unique name of the Microsoft provider.
const NameMicrosoft string = "microsoft"

// Microsoft allows authentication via AzureADEndpoint OAuth2.
type Microsoft struct {
	*baseProvider
}

// NewMicrosoftProvider creates new Microsoft AD provider instance with some defaults.
func NewMicrosoftProvider() *Microsoft {
	endpoints := microsoft.AzureADEndpoint("")
	return &Microsoft{&baseProvider{
		scopes:     []string{"User.Read"},
		authUrl:    endpoints.AuthURL,
		tokenUrl:   endpoints.TokenURL,
		userApiUrl: "https://graph.microsoft.com/v1.0/me",
	}}
}

// FetchAuthUser returns an AuthUser instance based on the Microsoft's user api.
func (p *Microsoft) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	// https://learn.microsoft.com/en-us/azure/active-directory/develop/userinfo
	// explore graph: https://developer.microsoft.com/en-us/graph/graph-explorer
	rawData := struct {
		Id    string `json:"id"`
		Name  string `json:"displayName"`
		Email string `json:"mail"`
	}{}

	if err := p.FetchRawUserData(token, &rawData); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:    rawData.Id,
		Name:  rawData.Name,
		Email: rawData.Email,
	}

	return user, nil
}
