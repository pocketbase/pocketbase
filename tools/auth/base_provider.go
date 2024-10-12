package auth

import (
	"context"
	"fmt"
	"io"
	"maps"
	"net/http"

	"golang.org/x/oauth2"
)

// BaseProvider defines common fields and methods used by OAuth2 client providers.
type BaseProvider struct {
	ctx          context.Context
	clientId     string
	clientSecret string
	displayName  string
	redirectURL  string
	authURL      string
	tokenURL     string
	userInfoURL  string
	scopes       []string
	pkce         bool
	extra        map[string]any
}

// Context implements Provider.Context() interface method.
func (p *BaseProvider) Context() context.Context {
	return p.ctx
}

// SetContext implements Provider.SetContext() interface method.
func (p *BaseProvider) SetContext(ctx context.Context) {
	p.ctx = ctx
}

// PKCE implements Provider.PKCE() interface method.
func (p *BaseProvider) PKCE() bool {
	return p.pkce
}

// SetPKCE implements Provider.SetPKCE() interface method.
func (p *BaseProvider) SetPKCE(enable bool) {
	p.pkce = enable
}

// DisplayName implements Provider.DisplayName() interface method.
func (p *BaseProvider) DisplayName() string {
	return p.displayName
}

// SetDisplayName implements Provider.SetDisplayName() interface method.
func (p *BaseProvider) SetDisplayName(displayName string) {
	p.displayName = displayName
}

// Scopes implements Provider.Scopes() interface method.
func (p *BaseProvider) Scopes() []string {
	return p.scopes
}

// SetScopes implements Provider.SetScopes() interface method.
func (p *BaseProvider) SetScopes(scopes []string) {
	p.scopes = scopes
}

// ClientId implements Provider.ClientId() interface method.
func (p *BaseProvider) ClientId() string {
	return p.clientId
}

// SetClientId implements Provider.SetClientId() interface method.
func (p *BaseProvider) SetClientId(clientId string) {
	p.clientId = clientId
}

// ClientSecret implements Provider.ClientSecret() interface method.
func (p *BaseProvider) ClientSecret() string {
	return p.clientSecret
}

// SetClientSecret implements Provider.SetClientSecret() interface method.
func (p *BaseProvider) SetClientSecret(secret string) {
	p.clientSecret = secret
}

// RedirectURL implements Provider.RedirectURL() interface method.
func (p *BaseProvider) RedirectURL() string {
	return p.redirectURL
}

// SetRedirectURL implements Provider.SetRedirectURL() interface method.
func (p *BaseProvider) SetRedirectURL(url string) {
	p.redirectURL = url
}

// AuthURL implements Provider.AuthURL() interface method.
func (p *BaseProvider) AuthURL() string {
	return p.authURL
}

// SetAuthURL implements Provider.SetAuthURL() interface method.
func (p *BaseProvider) SetAuthURL(url string) {
	p.authURL = url
}

// TokenURL implements Provider.TokenURL() interface method.
func (p *BaseProvider) TokenURL() string {
	return p.tokenURL
}

// SetTokenURL implements Provider.SetTokenURL() interface method.
func (p *BaseProvider) SetTokenURL(url string) {
	p.tokenURL = url
}

// UserInfoURL implements Provider.UserInfoURL() interface method.
func (p *BaseProvider) UserInfoURL() string {
	return p.userInfoURL
}

// SetUserInfoURL implements Provider.SetUserInfoURL() interface method.
func (p *BaseProvider) SetUserInfoURL(url string) {
	p.userInfoURL = url
}

// Extra implements Provider.Extra() interface method.
func (p *BaseProvider) Extra() map[string]any {
	return maps.Clone(p.extra)
}

// SetExtra implements Provider.SetExtra() interface method.
func (p *BaseProvider) SetExtra(data map[string]any) {
	p.extra = data
}

// BuildAuthURL implements Provider.BuildAuthURL() interface method.
func (p *BaseProvider) BuildAuthURL(state string, opts ...oauth2.AuthCodeOption) string {
	return p.oauth2Config().AuthCodeURL(state, opts...)
}

// FetchToken implements Provider.FetchToken() interface method.
func (p *BaseProvider) FetchToken(code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	return p.oauth2Config().Exchange(p.ctx, code, opts...)
}

// Client implements Provider.Client() interface method.
func (p *BaseProvider) Client(token *oauth2.Token) *http.Client {
	return p.oauth2Config().Client(p.ctx, token)
}

// FetchRawUserInfo implements Provider.FetchRawUserInfo() interface method.
func (p *BaseProvider) FetchRawUserInfo(token *oauth2.Token) ([]byte, error) {
	req, err := http.NewRequestWithContext(p.ctx, "GET", p.userInfoURL, nil)
	if err != nil {
		return nil, err
	}

	return p.sendRawUserInfoRequest(req, token)
}

// sendRawUserInfoRequest sends the specified user info request and return its raw response body.
func (p *BaseProvider) sendRawUserInfoRequest(req *http.Request, token *oauth2.Token) ([]byte, error) {
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
			p.userInfoURL,
			res.StatusCode,
			string(result),
		)
	}

	return result, nil
}

// oauth2Config constructs a oauth2.Config instance based on the provider settings.
func (p *BaseProvider) oauth2Config() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  p.redirectURL,
		ClientID:     p.clientId,
		ClientSecret: p.clientSecret,
		Scopes:       p.scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  p.authURL,
			TokenURL: p.tokenURL,
		},
	}
}
