package forms

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

// AdminLogin specifies an admin email/pass login form.
type AdminLogin struct {
	config AdminLoginConfig

	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
}

// AdminLoginConfig is the [AdminLogin] factory initializer config.
//
// NB! App is a required struct member.
type AdminLoginConfig struct {
	App core.App
	Dao *daos.Dao
}

// NewAdminLogin creates a new [AdminLogin] form with initializer
// config created from the provided [core.App] instance.
//
// If you want to submit the form as part of another transaction, use
// [NewAdminLoginWithConfig] with explicitly set Dao.
func NewAdminLogin(app core.App) *AdminLogin {
	return NewAdminLoginWithConfig(AdminLoginConfig{
		App: app,
	})
}

// NewAdminLoginWithConfig creates a new [AdminLogin] form
// with the provided config or panics on invalid configuration.
func NewAdminLoginWithConfig(config AdminLoginConfig) *AdminLogin {
	form := &AdminLogin{config: config}

	if form.config.App == nil {
		panic("Missing required config.App instance.")
	}

	if form.config.Dao == nil {
		form.config.Dao = form.config.App.Dao()
	}

	return form
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *AdminLogin) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(&form.Email, validation.Required, validation.Length(1, 255), is.EmailFormat),
		validation.Field(&form.Password, validation.Required, validation.Length(1, 255)),
	)
}

// Submit validates and submits the admin form.
// On success returns the authorized admin model.
func (form *AdminLogin) Submit() (*models.Admin, error) {
	if err := form.Validate(); err != nil {
		return nil, err
	}

	admin, err := form.config.Dao.FindAdminByEmail(form.Email)
	if err != nil {
		return nil, err
	}

	if admin.ValidatePassword(form.Password) {
		return admin, nil
	}

	return nil, errors.New("Invalid login credentials.")
}
