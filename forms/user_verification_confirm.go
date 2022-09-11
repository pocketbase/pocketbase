package forms

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

// UserVerificationConfirm specifies a user email verification confirmation form.
type UserVerificationConfirm struct {
	config UserVerificationConfirmConfig

	Token string `form:"token" json:"token"`
}

// UserVerificationConfirmConfig is the [UserVerificationConfirm]
// factory initializer config.
//
// NB! App is required struct member.
type UserVerificationConfirmConfig struct {
	App core.App
	Dao *daos.Dao
}

// NewUserVerificationConfirm creates a new [UserVerificationConfirm]
// form with initializer config created from the provided [core.App] instance.
//
// If you want to submit the form as part of another transaction, use
// [NewUserVerificationConfirmWithConfig] with explicitly set Dao.
func NewUserVerificationConfirm(app core.App) *UserVerificationConfirm {
	return NewUserVerificationConfirmWithConfig(UserVerificationConfirmConfig{
		App: app,
	})
}

// NewUserVerificationConfirmWithConfig creates a new [UserVerificationConfirmConfig]
// form with the provided config or panics on invalid configuration.
func NewUserVerificationConfirmWithConfig(config UserVerificationConfirmConfig) *UserVerificationConfirm {
	form := &UserVerificationConfirm{config: config}

	if form.config.App == nil {
		panic("Missing required config.App instance.")
	}

	if form.config.Dao == nil {
		form.config.Dao = form.config.App.Dao()
	}

	return form
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *UserVerificationConfirm) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(&form.Token, validation.Required, validation.By(form.checkToken)),
	)
}

func (form *UserVerificationConfirm) checkToken(value any) error {
	v, _ := value.(string)
	if v == "" {
		return nil // nothing to check
	}

	user, err := form.config.Dao.FindUserByToken(
		v,
		form.config.App.Settings().UserVerificationToken.Secret,
	)
	if err != nil || user == nil {
		return validation.NewError("validation_invalid_token", "Invalid or expired token.")
	}

	return nil
}

// Submit validates and submits the form.
// On success returns the verified user model associated to `form.Token`.
func (form *UserVerificationConfirm) Submit() (*models.User, error) {
	if err := form.Validate(); err != nil {
		return nil, err
	}

	user, err := form.config.Dao.FindUserByToken(
		form.Token,
		form.config.App.Settings().UserVerificationToken.Secret,
	)
	if err != nil {
		return nil, err
	}

	if user.Verified {
		return user, nil // already verified
	}

	user.Verified = true

	if err := form.config.Dao.SaveUser(user); err != nil {
		return nil, err
	}

	return user, nil
}
