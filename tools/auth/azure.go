package auth

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

var _ Provider = (*Azure)(nil)

// NameAzure is the unique name of the Azure provider.
const NameAzure string = "azure"

// Microsoft allows authentication via AzureADEndpoint OAuth2.
type Azure struct {
	*baseProvider
}

// NewMicrosoftProvider creates new Microsoft Azure AD provider instance with some defaults.
func NewAzureProvider() *Azure {
	endpoints := microsoft.AzureADEndpoint("")
	return &Azure{&baseProvider{
		scopes:     []string{"User.Read"},
		authUrl:    endpoints.AuthURL,
		tokenUrl:   endpoints.TokenURL,
		userApiUrl: "https://graph.microsoft.com/v1.0/me",
	}}
}

// FetchAuthUser returns an AuthUser instance based on the Microsoft's user api.
func (p *Azure) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
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
