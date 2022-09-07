package forms

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/forms/validators"
	"github.com/pocketbase/pocketbase/models"
)

// UserPasswordResetConfirm specifies a user password reset confirmation form.
type UserPasswordResetConfirm struct {
	config UserPasswordResetConfirmConfig

	Token           string `form:"token" json:"token"`
	Password        string `form:"password" json:"password"`
	PasswordConfirm string `form:"passwordConfirm" json:"passwordConfirm"`
}

// UserPasswordResetConfirmConfig is the [UserPasswordResetConfirm]
// factory initializer config.
//
// NB! App is required struct member.
type UserPasswordResetConfirmConfig struct {
	App core.App
	Dao *daos.Dao
}

// NewUserPasswordResetConfirm creates a new [UserPasswordResetConfirm]
// form with initializer config created from the provided [core.App] instance.
//
// If you want to submit the form as part of another transaction, use
// [NewUserPasswordResetConfirmWithConfig] with explicitly set Dao.
func NewUserPasswordResetConfirm(app core.App) *UserPasswordResetConfirm {
	return NewUserPasswordResetConfirmWithConfig(UserPasswordResetConfirmConfig{
		App: app,
	})
}

// NewUserPasswordResetConfirmWithConfig creates a new [UserPasswordResetConfirm]
// form with the provided config or panics on invalid configuration.
func NewUserPasswordResetConfirmWithConfig(config UserPasswordResetConfirmConfig) *UserPasswordResetConfirm {
	form := &UserPasswordResetConfirm{config: config}

	if form.config.App == nil {
		panic("Missing required config.App instance.")
	}

	if form.config.Dao == nil {
		form.config.Dao = form.config.App.Dao()
	}

	return form
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *UserPasswordResetConfirm) Validate() error {
	minPasswordLength := form.config.App.Settings().EmailAuth.MinPasswordLength

	return validation.ValidateStruct(form,
		validation.Field(&form.Token, validation.Required, validation.By(form.checkToken)),
		validation.Field(&form.Password, validation.Required, validation.Length(minPasswordLength, 100)),
		validation.Field(&form.PasswordConfirm, validation.Required, validation.By(validators.Compare(form.Password))),
	)
}

func (form *UserPasswordResetConfirm) checkToken(value any) error {
	v, _ := value.(string)
	if v == "" {
		return nil // nothing to check
	}

	user, err := form.config.Dao.FindUserByToken(
		v,
		form.config.App.Settings().UserPasswordResetToken.Secret,
	)
	if err != nil || user == nil {
		return validation.NewError("validation_invalid_token", "Invalid or expired token.")
	}

	return nil
}

// Submit validates and submits the form.
// On success returns the updated user model associated to `form.Token`.
func (form *UserPasswordResetConfirm) Submit() (*models.User, error) {
	if err := form.Validate(); err != nil {
		return nil, err
	}

	user, err := form.config.Dao.FindUserByToken(
		form.Token,
		form.config.App.Settings().UserPasswordResetToken.Secret,
	)
	if err != nil {
		return nil, err
	}

	if err := user.SetPassword(form.Password); err != nil {
		return nil, err
	}

	if err := form.config.Dao.SaveUser(user); err != nil {
		return nil, err
	}

	return user, nil
}
