package forms

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms/validators"
	"github.com/pocketbase/pocketbase/models"
)

// AdminPasswordResetConfirm defines an admin password reset confirmation form.
type AdminPasswordResetConfirm struct {
	app core.App

	Token           string `form:"token" json:"token"`
	Password        string `form:"password" json:"password"`
	PasswordConfirm string `form:"passwordConfirm" json:"passwordConfirm"`
}

// NewAdminPasswordResetConfirm creates new admin password reset confirmation form.
func NewAdminPasswordResetConfirm(app core.App) *AdminPasswordResetConfirm {
	return &AdminPasswordResetConfirm{
		app: app,
	}
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *AdminPasswordResetConfirm) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(&form.Token, validation.Required, validation.By(form.checkToken)),
		validation.Field(&form.Password, validation.Required, validation.Length(10, 100)),
		validation.Field(&form.PasswordConfirm, validation.Required, validation.By(validators.Compare(form.Password))),
	)
}

func (form *AdminPasswordResetConfirm) checkToken(value any) error {
	v, _ := value.(string)
	if v == "" {
		return nil // nothing to check
	}

	admin, err := form.app.Dao().FindAdminByToken(
		v,
		form.app.Settings().AdminPasswordResetToken.Secret,
	)
	if err != nil || admin == nil {
		return validation.NewError("validation_invalid_token", "Invalid or expired token.")
	}

	return nil
}

// Submit validates and submits the admin password reset confirmation form.
// On success returns the updated admin model associated to `form.Token`.
func (form *AdminPasswordResetConfirm) Submit() (*models.Admin, error) {
	if err := form.Validate(); err != nil {
		return nil, err
	}

	admin, err := form.app.Dao().FindAdminByToken(
		form.Token,
		form.app.Settings().AdminPasswordResetToken.Secret,
	)
	if err != nil {
		return nil, err
	}

	if err := admin.SetPassword(form.Password); err != nil {
		return nil, err
	}

	if err := form.app.Dao().SaveAdmin(admin); err != nil {
		return nil, err
	}

	return admin, nil
}
