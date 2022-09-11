package auth

import (
	"errors"
	"net/http"

	"golang.org/x/oauth2"
)

// AuthUser defines a standardized oauth2 user data structure.
type AuthUser struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	AvatarUrl string `json:"avatarUrl"`
}

// Provider defines a common interface for an OAuth2 client.
type Provider interface {
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

	// RedirectUrl returns the end address to redirect the user
	// going through the OAuth flow.
	RedirectUrl() string

	// SetRedirectUrl sets the provider's RedirectUrl.
	SetRedirectUrl(url string)

	// AuthUrl returns the provider's authorization service url.
	AuthUrl() string

	// SetAuthUrl sets the provider's AuthUrl.
	SetAuthUrl(url string)

	// TokenUrl returns the provider's token exchange service url.
	TokenUrl() string

	// SetTokenUrl sets the provider's TokenUrl.
	SetTokenUrl(url string)

	// UserApiUrl returns the provider's user info api url.
	UserApiUrl() string

	// SetUserApiUrl sets the provider's UserApiUrl.
	SetUserApiUrl(url string)

	// Client returns an http client using the provided token.
	Client(token *oauth2.Token) *http.Client

	// BuildAuthUrl returns a URL to the provider's consent page
	// that asks for permissions for the required scopes explicitly.
	BuildAuthUrl(state string, opts ...oauth2.AuthCodeOption) string

	// FetchToken converts an authorization code to token.
	FetchToken(code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)

	// FetchRawUserData requests and marshalizes into `result` the
	// the OAuth user api response.
	FetchRawUserData(token *oauth2.Token, result any) error

	// FetchAuthUser is similar to FetchRawUserData, but normalizes and
	// marshalizes the user api response into a standardized AuthUser struct.
	FetchAuthUser(token *oauth2.Token) (user *AuthUser, err error)
}

// NewProviderByName returns a new preconfigured provider instance by its name identifier.
func NewProviderByName(name string) (Provider, error) {
	switch name {
	case NameGoogle:
		return NewGoogleProvider(), nil
	case NameFacebook:
		return NewFacebookProvider(), nil
	case NameGithub:
		return NewGithubProvider(), nil
	case NameGitlab:
		return NewGitlabProvider(), nil
	case NameDiscord:
		return NewDiscordProvider(), nil
	case NameTwitter:
		return NewTwitterProvider(), nil
	default:
		return nil, errors.New("Missing provider " + name)
	}
}
