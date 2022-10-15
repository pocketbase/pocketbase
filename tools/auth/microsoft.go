package auth

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

var _ Provider = (*MicrosoftAd)(nil)

// NameMicrosoftAd is the unique name of the MicrosoftAD provider.
const NameMicrosoftAd string = "microsoftAd"

// Microsoft allows authentication via AzureADEndpoint OAuth2.
type MicrosoftAd struct {
	*baseProvider
}

// NewMicrosoftAdProvider creates new Microsoft AD provider instance with some defaults.
func NewMicrosoftAdProvider() *MicrosoftAd {
	endpoints := microsoft.AzureADEndpoint("")
	return &MicrosoftAd{&baseProvider{
		scopes:     []string{"User.Read"},
		authUrl:    endpoints.AuthURL,
		tokenUrl:   endpoints.TokenURL,
		userApiUrl: "https:/graph.microsoft.com/oidc/userinfo",
	}}
}

// FetchAuthUser returns an AuthUser instance based on the Microsoft's user api.
func (p *MicrosoftAd) FetchAuthUser(token *oauth2.Token) (*AuthUser, error) {
	// https://learn.microsoft.com/en-us/azure/active-directory/develop/userinfo
	rawData := struct {
		Id      string `json:"sub"`
		Name    string `json:"name"`
		Email   string `json:"email"`
		Picture string `json:"picture"`
	}{}

	if err := p.FetchRawUserData(token, &rawData); err != nil {
		return nil, err
	}

	user := &AuthUser{
		Id:        rawData.Id,
		Name:      rawData.Name,
		Email:     rawData.Email,
		AvatarUrl: rawData.Picture,
	}

	return user, nil
}
