package forms

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

// AdminLogin defines an admin email/pass login form.
type AdminLogin struct {
	app core.App

	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
}

// NewAdminLogin creates new admin login form for the provided app.
func NewAdminLogin(app core.App) *AdminLogin {
	return &AdminLogin{app: app}
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *AdminLogin) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(&form.Email, validation.Required, validation.Length(1, 255), is.Email),
		validation.Field(&form.Password, validation.Required, validation.Length(1, 255)),
	)
}

// Submit validates and submits the admin form.
// On success returns the authorized admin model.
func (form *AdminLogin) Submit() (*models.Admin, error) {
	if err := form.Validate(); err != nil {
		return nil, err
	}

	admin, err := form.app.Dao().FindAdminByEmail(form.Email)
	if err != nil {
		return nil, err
	}

	if admin.ValidatePassword(form.Password) {
		return admin, nil
	}

	return nil, errors.New("Invalid login credentials.")
}
