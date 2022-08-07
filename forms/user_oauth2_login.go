package forms

import (
	"errors"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/auth"
	"github.com/pocketbase/pocketbase/tools/security"
	"golang.org/x/oauth2"
)

// UserOauth2Login specifies a user Oauth2 login form.
type UserOauth2Login struct {
	config UserOauth2LoginConfig

	// The name of the OAuth2 client provider (eg. "google")
	Provider string `form:"provider" json:"provider"`

	// The authorization code returned from the initial request.
	Code string `form:"code" json:"code"`

	// The code verifier sent with the initial request as part of the code_challenge.
	CodeVerifier string `form:"codeVerifier" json:"codeVerifier"`

	// The redirect url sent with the initial request.
	RedirectUrl string `form:"redirectUrl" json:"redirectUrl"`
}

// UserOauth2LoginConfig is the [UserOauth2Login] factory initializer config.
//
// NB! App is required struct member.
type UserOauth2LoginConfig struct {
	App   core.App
	TxDao *daos.Dao
}

// NewUserOauth2Login creates a new [UserOauth2Login] form with
// initializer config created from the provided [core.App] instance.
//
// If you want to submit the form as part of another transaction, use
// [NewUserOauth2LoginWithConfig] with explicitly set TxDao.
func NewUserOauth2Login(app core.App) *UserOauth2Login {
	return NewUserOauth2LoginWithConfig(UserOauth2LoginConfig{
		App: app,
	})
}

// NewUserOauth2LoginWithConfig creates a new [UserOauth2Login]
// form with the provided config or panics on invalid configuration.
func NewUserOauth2LoginWithConfig(config UserOauth2LoginConfig) *UserOauth2Login {
	form := &UserOauth2Login{config: config}

	if form.config.App == nil {
		panic("Missing required config.App instance.")
	}

	if form.config.TxDao == nil {
		form.config.TxDao = form.config.App.Dao()
	}

	return form
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

	config, ok := form.config.App.Settings().NamedAuthProviderConfigs()[name]
	if !ok || !config.Enabled {
		return validation.NewError("validation_invalid_provider", fmt.Sprintf("%q is missing or is not enabled.", name))
	}

	return nil
}

// Submit validates and submits the form.
// On success returns the authorized user model and the fetched provider's data.
func (form *UserOauth2Login) Submit() (*models.User, *auth.AuthUser, error) {
	if err := form.Validate(); err != nil {
		return nil, nil, err
	}

	provider, err := auth.NewProviderByName(form.Provider)
	if err != nil {
		return nil, nil, err
	}

	config := form.config.App.Settings().NamedAuthProviderConfigs()[form.Provider]
	config.SetupProvider(provider)

	provider.SetRedirectUrl(form.RedirectUrl)

	// fetch token
	token, err := provider.FetchToken(
		form.Code,
		oauth2.SetAuthURLParam("code_verifier", form.CodeVerifier),
	)
	if err != nil {
		return nil, nil, err
	}

	// fetch auth user
	authData, err := provider.FetchAuthUser(token)
	if err != nil {
		return nil, nil, err
	}

	// login/register the auth user
	user, _ := form.config.TxDao.FindUserByEmail(authData.Email)
	if user != nil {
		// update the existing user's verified state
		if !user.Verified {
			user.Verified = true
			if err := form.config.TxDao.SaveUser(user); err != nil {
				return nil, authData, err
			}
		}
		return user, authData, nil
	}

	if !config.AllowRegistrations {
		// registration of new users is not allowed via the Oauth2 provider
		return nil, authData, errors.New("Cannot find user with the authorized email.")
	}

	// create new user
	user = &models.User{Verified: true}
	upsertForm := NewUserUpsertWithConfig(UserUpsertConfig{
		App:   form.config.App,
		TxDao: form.config.TxDao,
	}, user)
	upsertForm.Email = authData.Email
	upsertForm.Password = security.RandomString(30)
	upsertForm.PasswordConfirm = upsertForm.Password

	event := &core.UserOauth2RegisterEvent{
		User:     user,
		AuthData: authData,
	}

	if err := form.config.App.OnUserBeforeOauth2Register().Trigger(event); err != nil {
		return nil, authData, err
	}

	if err := upsertForm.Submit(); err != nil {
		return nil, authData, err
	}

	if err := form.config.App.OnUserAfterOauth2Register().Trigger(event); err != nil {
		return nil, authData, err
	}

	return user, authData, nil
}
