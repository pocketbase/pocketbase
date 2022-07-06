package forms

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

// UserVerificationConfirm defines a user email confirmation form.
type UserVerificationConfirm struct {
	app core.App

	Token string `form:"token" json:"token"`
}

// NewUserVerificationConfirm creates a new user email confirmation form.
func NewUserVerificationConfirm(app core.App) *UserVerificationConfirm {
	return &UserVerificationConfirm{
		app: app,
	}
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

	user, err := form.app.Dao().FindUserByToken(
		v,
		form.app.Settings().UserVerificationToken.Secret,
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

	user, err := form.app.Dao().FindUserByToken(
		form.Token,
		form.app.Settings().UserVerificationToken.Secret,
	)
	if err != nil {
		return nil, err
	}

	if user.Verified {
		return user, nil // already verified
	}

	user.Verified = true

	if err := form.app.Dao().SaveUser(user); err != nil {
		return nil, err
	}

	return user, nil
}
