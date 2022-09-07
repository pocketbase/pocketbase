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
	App core.App
	Dao *daos.Dao
}

// NewUserOauth2Login creates a new [UserOauth2Login] form with
// initializer config created from the provided [core.App] instance.
//
// If you want to submit the form as part of another transaction, use
// [NewUserOauth2LoginWithConfig] with explicitly set Dao.
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

	if form.config.Dao == nil {
		form.config.Dao = form.config.App.Dao()
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

	// load provider configuration
	config := form.config.App.Settings().NamedAuthProviderConfigs()[form.Provider]
	if err := config.SetupProvider(provider); err != nil {
		return nil, nil, err
	}

	provider.SetRedirectUrl(form.RedirectUrl)

	// fetch token
	token, err := provider.FetchToken(
		form.Code,
		oauth2.SetAuthURLParam("code_verifier", form.CodeVerifier),
	)
	if err != nil {
		return nil, nil, err
	}

	// fetch external auth user
	authData, err := provider.FetchAuthUser(token)
	if err != nil {
		return nil, nil, err
	}

	var user *models.User

	// check for existing relation with the external auth user
	rel, _ := form.config.Dao.FindExternalAuthByProvider(form.Provider, authData.Id)
	if rel != nil {
		user, err = form.config.Dao.FindUserById(rel.UserId)
		if err != nil {
			return nil, authData, err
		}
	} else if authData.Email != "" {
		// look for an existing user by the external user's email
		user, _ = form.config.Dao.FindUserByEmail(authData.Email)
	}

	if user == nil && !config.AllowRegistrations {
		return nil, authData, errors.New("New users registration is not allowed for the authorized provider.")
	}

	saveErr := form.config.Dao.RunInTransaction(func(txDao *daos.Dao) error {
		if user == nil {
			user = &models.User{}
			user.Verified = true
			user.Email = authData.Email
			user.SetPassword(security.RandomString(30))

			// create the new user
			if err := txDao.SaveUser(user); err != nil {
				return err
			}
		} else {
			// update the existing user empty email if the authData has one
			// (this in case previously the user was created with
			// an OAuth2 provider that didn't return an email address)
			if user.Email == "" && authData.Email != "" {
				user.Email = authData.Email
				if err := txDao.SaveUser(user); err != nil {
					return err
				}
			}

			// update the existing user verified state
			// (only if the user doesn't have an email or the user email match with the one in authData)
			if !user.Verified && (user.Email == "" || user.Email == authData.Email) {
				user.Verified = true
				if err := txDao.SaveUser(user); err != nil {
					return err
				}
			}
		}

		// create ExternalAuth relation if missing
		if rel == nil {
			rel = &models.ExternalAuth{
				UserId:     user.Id,
				Provider:   form.Provider,
				ProviderId: authData.Id,
			}
			if err := txDao.SaveExternalAuth(rel); err != nil {
				return err
			}
		}

		return nil
	})

	if saveErr != nil {
		return nil, authData, saveErr
	}

	return user, authData, nil
}
