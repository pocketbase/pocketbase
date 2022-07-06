package forms

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms/validators"
	"github.com/pocketbase/pocketbase/models"
)

// UserPasswordResetConfirm defines a user password reset confirmation form.
type UserPasswordResetConfirm struct {
	app core.App

	Token           string `form:"token" json:"token"`
	Password        string `form:"password" json:"password"`
	PasswordConfirm string `form:"passwordConfirm" json:"passwordConfirm"`
}

// NewUserPasswordResetConfirm creates new user password reset confirmation form.
func NewUserPasswordResetConfirm(app core.App) *UserPasswordResetConfirm {
	return &UserPasswordResetConfirm{
		app: app,
	}
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *UserPasswordResetConfirm) Validate() error {
	minPasswordLength := form.app.Settings().EmailAuth.MinPasswordLength

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

	user, err := form.app.Dao().FindUserByToken(
		v,
		form.app.Settings().UserPasswordResetToken.Secret,
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

	user, err := form.app.Dao().FindUserByToken(
		form.Token,
		form.app.Settings().UserPasswordResetToken.Secret,
	)
	if err != nil {
		return nil, err
	}

	if err := user.SetPassword(form.Password); err != nil {
		return nil, err
	}

	if err := form.app.Dao().SaveUser(user); err != nil {
		return nil, err
	}

	return user, nil
}
