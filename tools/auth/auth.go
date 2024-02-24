package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

// AuthUser defines a standardized oauth2 user data structure.
type AuthUser struct {
	Id           string         `json:"id"`
	Name         string         `json:"name"`
	Username     string         `json:"username"`
	Email        string         `json:"email"`
	AvatarUrl    string         `json:"avatarUrl"`
	AccessToken  string         `json:"accessToken"`
	RefreshToken string         `json:"refreshToken"`
	Expiry       types.DateTime `json:"expiry"`
	RawUser      map[string]any `json:"rawUser"`
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
	FetchRawUserData(token *oauth2.Token) ([]byte, error)

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
	case NameMicrosoft:
		return NewMicrosoftProvider(), nil
	case NameSpotify:
		return NewSpotifyProvider(), nil
	case NameKakao:
		return NewKakaoProvider(), nil
	case NameTwitch:
		return NewTwitchProvider(), nil
	case NameStrava:
		return NewStravaProvider(), nil
	case NameGitee:
		return NewGiteeProvider(), nil
	case NameLivechat:
		return NewLivechatProvider(), nil
	case NameGitea:
		return NewGiteaProvider(), nil
	case NameOIDC:
		return NewOIDCProvider(), nil
	case NameOIDC + "2":
		return NewOIDCProvider(), nil
	case NameOIDC + "3":
		return NewOIDCProvider(), nil
	case NameApple:
		return NewAppleProvider(), nil
	case NameInstagram:
		return NewInstagramProvider(), nil
	case NameVK:
		return NewVKProvider(), nil
	case NameYandex:
		return NewYandexProvider(), nil
	case NamePatreon:
		return NewPatreonProvider(), nil
	case NameMailcow:
		return NewMailcowProvider(), nil
	case NameBitbucket:
		return NewBitbucketProvider(), nil
	case NamePlanningcenter:
		return NewPlanningcenterProvider(), nil
	default:
		return nil, errors.New("Missing provider " + name)
	}
}
