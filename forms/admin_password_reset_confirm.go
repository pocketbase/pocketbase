package forms

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/forms/validators"
	"github.com/pocketbase/pocketbase/models"
)

// AdminPasswordResetConfirm specifies an admin password reset confirmation form.
type AdminPasswordResetConfirm struct {
	config AdminPasswordResetConfirmConfig

	Token           string `form:"token" json:"token"`
	Password        string `form:"password" json:"password"`
	PasswordConfirm string `form:"passwordConfirm" json:"passwordConfirm"`
}

// AdminPasswordResetConfirmConfig is the [AdminPasswordResetConfirm] factory initializer config.
//
// NB! App is required struct member.
type AdminPasswordResetConfirmConfig struct {
	App core.App
	Dao *daos.Dao
}

// NewAdminPasswordResetConfirm creates a new [AdminPasswordResetConfirm]
// form with initializer config created from the provided [core.App] instance.
//
// If you want to submit the form as part of another transaction, use
// [NewAdminPasswordResetConfirmWithConfig] with explicitly set Dao.
func NewAdminPasswordResetConfirm(app core.App) *AdminPasswordResetConfirm {
	return NewAdminPasswordResetConfirmWithConfig(AdminPasswordResetConfirmConfig{
		App: app,
	})
}

// NewAdminPasswordResetConfirmWithConfig creates a new [AdminPasswordResetConfirm]
// form with the provided config or panics on invalid configuration.
func NewAdminPasswordResetConfirmWithConfig(config AdminPasswordResetConfirmConfig) *AdminPasswordResetConfirm {
	form := &AdminPasswordResetConfirm{config: config}

	if form.config.App == nil {
		panic("Missing required config.App instance.")
	}

	if form.config.Dao == nil {
		form.config.Dao = form.config.App.Dao()
	}

	return form
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

	admin, err := form.config.Dao.FindAdminByToken(
		v,
		form.config.App.Settings().AdminPasswordResetToken.Secret,
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

	admin, err := form.config.Dao.FindAdminByToken(
		form.Token,
		form.config.App.Settings().AdminPasswordResetToken.Secret,
	)
	if err != nil {
		return nil, err
	}

	if err := admin.SetPassword(form.Password); err != nil {
		return nil, err
	}

	if err := form.config.Dao.SaveAdmin(admin); err != nil {
		return nil, err
	}

	return admin, nil
}
