package forms

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

// UserEmailLogin specifies a user email/pass login form.
type UserEmailLogin struct {
	config UserEmailLoginConfig

	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
}

// UserEmailLoginConfig is the [UserEmailLogin] factory initializer config.
//
// NB! App is required struct member.
type UserEmailLoginConfig struct {
	App core.App
	Dao *daos.Dao
}

// NewUserEmailLogin creates a new [UserEmailLogin] form with
// initializer config created from the provided [core.App] instance.
//
// This factory method is used primarily for convenience (and backward compatibility).
// If you want to submit the form as part of another transaction, use
// [NewUserEmailLoginWithConfig] with explicitly set Dao.
func NewUserEmailLogin(app core.App) *UserEmailLogin {
	return NewUserEmailLoginWithConfig(UserEmailLoginConfig{
		App: app,
	})
}

// NewUserEmailLoginWithConfig creates a new [UserEmailLogin]
// form with the provided config or panics on invalid configuration.
func NewUserEmailLoginWithConfig(config UserEmailLoginConfig) *UserEmailLogin {
	form := &UserEmailLogin{config: config}

	if form.config.App == nil {
		panic("Missing required config.App instance.")
	}

	if form.config.Dao == nil {
		form.config.Dao = form.config.App.Dao()
	}

	return form
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *UserEmailLogin) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(&form.Email, validation.Required, validation.Length(1, 255), is.EmailFormat),
		validation.Field(&form.Password, validation.Required, validation.Length(1, 255)),
	)
}

// Submit validates and submits the form.
// On success returns the authorized user model.
func (form *UserEmailLogin) Submit() (*models.User, error) {
	if err := form.Validate(); err != nil {
		return nil, err
	}

	user, err := form.config.Dao.FindUserByEmail(form.Email)
	if err != nil {
		return nil, err
	}

	if !user.ValidatePassword(form.Password) {
		return nil, validation.NewError("invalid_login", "Invalid login credentials.")
	}

	return user, nil
}
