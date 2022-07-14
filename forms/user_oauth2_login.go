package forms

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/auth"
	"github.com/pocketbase/pocketbase/tools/security"
	"golang.org/x/oauth2"
)

// UserOauth2Login defines a user Oauth2 login form.
type UserOauth2Login struct {
	app core.App

	// The name of the OAuth2 client provider (eg. "google")
	Provider string `form:"provider" json:"provider"`

	// The authorization code returned from the initial request.
	Code string `form:"code" json:"code"`

	// The code verifier sent with the initial request as part of the code_challenge.
	CodeVerifier string `form:"codeVerifier" json:"codeVerifier"`

	// The redirect url sent with the initial request.
	RedirectUrl string `form:"redirectUrl" json:"redirectUrl"`
}

// NewUserOauth2Login creates a new user Oauth2 login form.
func NewUserOauth2Login(app core.App) *UserOauth2Login {
	return &UserOauth2Login{app: app}
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *UserOauth2Login) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(&form.Provider, validation.Required, validation.By(form.checkProviderName)),
		validation.Field(&form.Code, validation.Required),
		validation.Field(&form.CodeVerifier, validation.Required),
		validation.Field(&form.RedirectUrl, validation.Required, is.URL),
	)
}

func (form *UserOauth2Login) checkProviderName(value any) error {
	name, _ := value.(string)

	config, ok := form.app.Settings().NamedAuthProviderConfigs()[name]
	if !ok || !config.Enabled {
		return validation.NewError("validation_invalid_provider", fmt.Sprintf("%q is missing or is not enabled.", name))
	}

	return nil
}

// Submit validates and submits the form.
// On success returns the authorized user model and the fetched provider's data.
// If exists, the user's give token is used to fetch the user's info rather requesting the access token itself
// Convention:
// - The redirect URL must inlucde the oauth2.Token in the redirect url as query params
// - The query param "access_token" must be given
// - Example redirect URL: http://localhost:3000?access_token=XXX123&refresh_token=ABC123&token_type=Bearer&expires_at=123344
func (form *UserOauth2Login) Submit() (*models.User, *auth.AuthUser, error) {
	if err := form.Validate(); err != nil {
		return nil, nil, err
	}

	provider, err := auth.NewProviderByName(form.Provider)
	if err != nil {
		return nil, nil, err
	}

	config := form.app.Settings().NamedAuthProviderConfigs()[form.Provider]
	config.SetupProvider(provider)

	provider.SetRedirectUrl(form.RedirectUrl)

	// request the access token with the initial access code or if given use the form.Token as the access_token to request the user's profile
	var token oauth2.Token
	// Parse the URL to get the token details (if exists)
	u, err := url.Parse(form.RedirectUrl)
	if err == nil {
		token.AccessToken = u.Query().Get("access_token")
		token.RefreshToken = u.Query().Get("refresh_token")
		token.TokenType = u.Query().Get("token_type")
		i, err := strconv.ParseInt(u.Query().Get("expires_at"), 10, 64)
		if err == nil {
			token.Expiry = time.Unix(i, 0)
		}
	}

	// fetch token
	// The "access_token" does not exist; request the "access_token" from the provider by the returned code
	if len(token.AccessToken) == 0 {
		t, err := provider.FetchToken(
			form.Code,
			oauth2.SetAuthURLParam("code_verifier", form.CodeVerifier),
		)
		if err != nil {
			return nil, nil, err
		}
		token = *t
	}
	if err != nil {
		return nil, nil, err
	}

	// fetch auth user
	authData, err := provider.FetchAuthUser(&token)
	if err != nil {
		return nil, nil, err
	}

	// login/register the auth user
	user, _ := form.app.Dao().FindUserByEmail(authData.Email)
	if user != nil {
		// update the existing user's verified state
		if !user.Verified {
			user.Verified = true
			if err := form.app.Dao().SaveUser(user); err != nil {
				return nil, authData, err
			}
		}
	} else {
		if !config.AllowRegistrations {
			// registration of new users is not allowed via the Oauth2 provider
			return nil, authData, errors.New("Cannot find user with the authorized email.")
		}

		// create new user
		user = &models.User{Verified: true}
		upsertForm := NewUserUpsert(form.app, user)
		upsertForm.Email = authData.Email
		upsertForm.Password = security.RandomString(30)
		upsertForm.PasswordConfirm = upsertForm.Password

		event := &core.UserOauth2RegisterEvent{
			User:     user,
			AuthData: authData,
		}

		if err := form.app.OnUserBeforeOauth2Register().Trigger(event); err != nil {
			return nil, authData, err
		}

		if err := upsertForm.Submit(); err != nil {
			return nil, authData, err
		}

		if err := form.app.OnUserAfterOauth2Register().Trigger(event); err != nil {
			return nil, authData, err
		}
	}

	return user, authData, nil
}
