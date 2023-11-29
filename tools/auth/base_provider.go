package auth

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
)

// baseProvider defines common fields and methods used by OAuth2 client providers.
type baseProvider struct {
	ctx          context.Context
	clientId     string
	clientSecret string
	displayName  string
	redirectUrl  string
	authUrl      string
	tokenUrl     string
	userApiUrl   string
	scopes       []string
	pkce         bool
}

// Context implements Provider.Context() interface method.
func (p *baseProvider) Context() context.Context {
	return p.ctx
}

// SetContext implements Provider.SetContext() interface method.
func (p *baseProvider) SetContext(ctx context.Context) {
	p.ctx = ctx
}

// PKCE implements Provider.PKCE() interface method.
func (p *baseProvider) PKCE() bool {
	return p.pkce
}

// SetPKCE implements Provider.SetPKCE() interface method.
func (p *baseProvider) SetPKCE(enable bool) {
	p.pkce = enable
}

// DisplayName implements Provider.DisplayName() interface method.
func (p *baseProvider) DisplayName() string {
	return p.displayName
}

// SetDisplayName implements Provider.SetDisplayName() interface method.
func (p *baseProvider) SetDisplayName(displayName string) {
	p.displayName = displayName
}

// Scopes implements Provider.Scopes() interface method.
func (p *baseProvider) Scopes() []string {
	return p.scopes
}

// SetScopes implements Provider.SetScopes() interface method.
func (p *baseProvider) SetScopes(scopes []string) {
	p.scopes = scopes
}

// ClientId implements Provider.ClientId() interface method.
func (p *baseProvider) ClientId() string {
	return p.clientId
}

// SetClientId implements Provider.SetClientId() interface method.
func (p *baseProvider) SetClientId(clientId string) {
	p.clientId = clientId
}

// ClientSecret implements Provider.ClientSecret() interface method.
func (p *baseProvider) ClientSecret() string {
	return p.clientSecret
}

// SetClientSecret implements Provider.SetClientSecret() interface method.
func (p *baseProvider) SetClientSecret(secret string) {
	p.clientSecret = secret
}

// RedirectUrl implements Provider.RedirectUrl() interface method.
func (p *baseProvider) RedirectUrl() string {
	return p.redirectUrl
}

// SetRedirectUrl implements Provider.SetRedirectUrl() interface method.
func (p *baseProvider) SetRedirectUrl(url string) {
	p.redirectUrl = url
}

// AuthUrl implements Provider.AuthUrl() interface method.
func (p *baseProvider) AuthUrl() string {
	return p.authUrl
}

// SetAuthUrl implements Provider.SetAuthUrl() interface method.
func (p *baseProvider) SetAuthUrl(url string) {
	p.authUrl = url
}

// TokenUrl implements Provider.TokenUrl() interface method.
func (p *baseProvider) TokenUrl() string {
	return p.tokenUrl
}

// SetTokenUrl implements Provider.SetTokenUrl() interface method.
func (p *baseProvider) SetTokenUrl(url string) {
	p.tokenUrl = url
}

// UserApiUrl implements Provider.UserApiUrl() interface method.
func (p *baseProvider) UserApiUrl() string {
	return p.userApiUrl
}

// SetUserApiUrl implements Provider.SetUserApiUrl() interface method.
func (p *baseProvider) SetUserApiUrl(url string) {
	p.userApiUrl = url
}

// BuildAuthUrl implements Provider.BuildAuthUrl() interface method.
func (p *baseProvider) BuildAuthUrl(state string, opts ...oauth2.AuthCodeOption) string {
	return p.oauth2Config().AuthCodeURL(state, opts...)
}

// FetchToken implements Provider.FetchToken() interface method.
func (p *baseProvider) FetchToken(code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	return p.oauth2Config().Exchange(p.ctx, code, opts...)
}

// Client implements Provider.Client() interface method.
func (p *baseProvider) Client(token *oauth2.Token) *http.Client {
	return p.oauth2Config().Client(p.ctx, token)
}

// FetchRawUserData implements Provider.FetchRawUserData() interface method.
func (p *baseProvider) FetchRawUserData(token *oauth2.Token) ([]byte, error) {
	req, err := http.NewRequestWithContext(p.ctx, "GET", p.userApiUrl, nil)
	if err != nil {
		return nil, err
	}

	return p.sendRawUserDataRequest(req, token)
}

// sendRawUserDataRequest sends the specified user data request and return its raw response body.
func (p *baseProvider) sendRawUserDataRequest(req *http.Request, token *oauth2.Token) ([]byte, error) {
	client := p.Client(token)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	result, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// http.Client.Get doesn't treat non 2xx responses as error
	if res.StatusCode >= 400 {
		return nil, fmt.Errorf(
			"failed to fetch OAuth2 user profile via %s (%d):\n%s",
			p.userApiUrl,
			res.StatusCode,
			string(result),
		)
	}

	return result, nil
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
