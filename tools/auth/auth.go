package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

// ProviderFactoryFunc defines a function for initializing a new OAuth2 provider.
type ProviderFactoryFunc func() Provider

// Providers defines a map with all of the available OAuth2 providers.
//
// To register a new provider append a new entry in the map.
var Providers = map[string]ProviderFactoryFunc{}

// NewProviderByName returns a new preconfigured provider instance by its name identifier.
func NewProviderByName(name string) (Provider, error) {
	factory, ok := Providers[name]
	if !ok {
		return nil, errors.New("missing provider " + name)
	}

	return factory(), nil
}

// Provider defines a common interface for an OAuth2 client.
type Provider interface {
	// Context returns the context associated with the provider (if any).
	Context() context.Context

	// SetContext assigns the specified context to the current provider.
	SetContext(ctx context.Context)

	// PKCE indicates whether the provider can use the PKCE flow.
	PKCE() bool

	// SetPKCE toggles the state whether the provider can use the PKCE flow or not.
	SetPKCE(enable bool)

	// DisplayName usually returns provider name as it is officially written
	// and it could be used directly in the UI.
	DisplayName() string

	// SetDisplayName sets the provider's display name.
	SetDisplayName(displayName string)

	// Scopes returns the provider access permissions that will be requested.
	Scopes() []string

	// SetScopes sets the provider access permissions that will be requested later.
	SetScopes(scopes []string)

	// ClientId returns the provider client's app ID.
	ClientId() string

	// SetClientId sets the provider client's ID.
	SetClientId(clientId string)

	// ClientSecret returns the provider client's app secret.
	ClientSecret() string

	// SetClientSecret sets the provider client's app secret.
	SetClientSecret(secret string)

	// RedirectURL returns the end address to redirect the user
	// going through the OAuth flow.
	RedirectURL() string

	// SetRedirectURL sets the provider's RedirectURL.
	SetRedirectURL(url string)

	// AuthURL returns the provider's authorization service url.
	AuthURL() string

	// SetAuthURL sets the provider's AuthURL.
	SetAuthURL(url string)

	// TokenURL returns the provider's token exchange service url.
	TokenURL() string

	// SetTokenURL sets the provider's TokenURL.
	SetTokenURL(url string)

	// UserInfoURL returns the provider's user info api url.
	UserInfoURL() string

	// SetUserInfoURL sets the provider's UserInfoURL.
	SetUserInfoURL(url string)

	// Extra returns a shallow copy of any custom config data
	// that the provider may be need.
	Extra() map[string]any

	// SetExtra updates the provider's custom config data.
	SetExtra(data map[string]any)

	// Client returns an http client using the provided token.
	Client(token *oauth2.Token) *http.Client

	// BuildAuthURL returns a URL to the provider's consent page
	// that asks for permissions for the required scopes explicitly.
	BuildAuthURL(state string, opts ...oauth2.AuthCodeOption) string

	// FetchToken converts an authorization code to token.
	FetchToken(code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)

	// FetchRawUserInfo requests and marshalizes into `result` the
	// the OAuth user api response.
	FetchRawUserInfo(token *oauth2.Token) ([]byte, error)

	// FetchAuthUser is similar to FetchRawUserInfo, but normalizes and
	// marshalizes the user api response into a standardized AuthUser struct.
	FetchAuthUser(token *oauth2.Token) (user *AuthUser, err error)
}

// wrapFactory is a helper that wraps a Provider specific factory
// function and returns its result as Provider interface.
func wrapFactory[T Provider](factory func() T) ProviderFactoryFunc {
	return func() Provider {
		return factory()
	}
}

// AuthUser defines a standardized OAuth2 user data structure.
type AuthUser struct {
	Expiry       types.DateTime `json:"expiry"`
	RawUser      map[string]any `json:"rawUser"`
	Id           string         `json:"id"`
	Name         string         `json:"name"`
	Username     string         `json:"username"`
	Email        string         `json:"email"`
	AvatarURL    string         `json:"avatarURL"`
	AccessToken  string         `json:"accessToken"`
	RefreshToken string         `json:"refreshToken"`

	// @todo
	// deprecated: use AvatarURL instead
	// AvatarUrl will be removed after dropping v0.22 support
	AvatarUrl string `json:"avatarUrl"`
}

// MarshalJSON implements the [json.Marshaler] interface.
//
// @todo remove after dropping v0.22 support
func (au AuthUser) MarshalJSON() ([]byte, error) {
	type alias AuthUser // prevent recursion

	au2 := alias(au)
	au2.AvatarUrl = au.AvatarURL // ensure that the legacy field is populated

	return json.Marshal(au2)
}
