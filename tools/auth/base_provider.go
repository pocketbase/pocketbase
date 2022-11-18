package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
)

// baseProvider defines common fields and methods used by OAuth2 client providers.
type baseProvider struct {
	scopes       []string
	clientId     string
	clientSecret string
	redirectUrl  string
	authUrl      string
	tokenUrl     string
	userApiUrl   string
}

// Scopes implements Provider.Scopes interface.
func (p *baseProvider) Scopes() []string {
	return p.scopes
}

// SetScopes implements Provider.SetScopes interface.
func (p *baseProvider) SetScopes(scopes []string) {
	p.scopes = scopes
}

// ClientId implements Provider.ClientId interface.
func (p *baseProvider) ClientId() string {
	return p.clientId
}

// SetClientId implements Provider.SetClientId interface.
func (p *baseProvider) SetClientId(clientId string) {
	p.clientId = clientId
}

// ClientSecret implements Provider.ClientSecret interface.
func (p *baseProvider) ClientSecret() string {
	return p.clientSecret
}

// SetClientSecret implements Provider.SetClientSecret interface.
func (p *baseProvider) SetClientSecret(secret string) {
	p.clientSecret = secret
}

// RedirectUrl implements Provider.RedirectUrl interface.
func (p *baseProvider) RedirectUrl() string {
	return p.redirectUrl
}

// SetRedirectUrl implements Provider.SetRedirectUrl interface.
func (p *baseProvider) SetRedirectUrl(url string) {
	p.redirectUrl = url
}

// AuthUrl implements Provider.AuthUrl interface.
func (p *baseProvider) AuthUrl() string {
	return p.authUrl
}

// SetAuthUrl implements Provider.SetAuthUrl interface.
func (p *baseProvider) SetAuthUrl(url string) {
	p.authUrl = url
}

// TokenUrl implements Provider.TokenUrl interface.
func (p *baseProvider) TokenUrl() string {
	return p.tokenUrl
}

// SetTokenUrl implements Provider.SetTokenUrl interface.
func (p *baseProvider) SetTokenUrl(url string) {
	p.tokenUrl = url
}

// UserApiUrl implements Provider.UserApiUrl interface.
func (p *baseProvider) UserApiUrl() string {
	return p.userApiUrl
}

// SetUserApiUrl implements Provider.SetUserApiUrl interface.
func (p *baseProvider) SetUserApiUrl(url string) {
	p.userApiUrl = url
}

// BuildAuthUrl implements Provider.BuildAuthUrl interface.
func (p *baseProvider) BuildAuthUrl(state string, opts ...oauth2.AuthCodeOption) string {
	return p.oauth2Config().AuthCodeURL(state, opts...)
}

// FetchToken implements Provider.FetchToken interface.
func (p *baseProvider) FetchToken(code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	return p.oauth2Config().Exchange(context.Background(), code, opts...)
}

// Client implements Provider.Client interface.
func (p *baseProvider) Client(token *oauth2.Token) *http.Client {
	return p.oauth2Config().Client(context.Background(), token)
}

// FetchRawUserData implements Provider.FetchRawUserData interface.
func (p *baseProvider) FetchRawUserData(token *oauth2.Token, result any) error {
	req, err := http.NewRequest("GET", p.userApiUrl, nil)
	if err != nil {
		return err
	}

	return p.sendRawUserDataRequest(req, token, result)
}

// sendRawUserDataRequest sends the specified request and
// unmarshal the response body into result.
func (p *baseProvider) sendRawUserDataRequest(req *http.Request, token *oauth2.Token, result any) error {
	client := p.Client(token)

	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	// http.Client.Get doesn't treat non 2xx responses as error
	if response.StatusCode >= 400 {
		return fmt.Errorf(
			"Failed to fetch OAuth2 user profile via %s (%d):\n%s",
			p.userApiUrl,
			response.StatusCode,
			string(content),
		)
	}

	return json.Unmarshal(content, &result)
}

// oauth2Config constructs a oauth2.Config instance based on the provider settings.
func (p *baseProvider) oauth2Config() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  p.redirectUrl,
		ClientID:     p.clientId,
		ClientSecret: p.clientSecret,
		Scopes:       p.scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  p.authUrl,
			TokenURL: p.tokenUrl,
		},
	}
}
