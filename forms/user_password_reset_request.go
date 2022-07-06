package forms

import (
	"errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/mails"
	"github.com/pocketbase/pocketbase/tools/types"
)

// UserPasswordResetRequest defines a user password reset request form.
type UserPasswordResetRequest struct {
	app             core.App
	resendThreshold float64

	Email string `form:"email" json:"email"`
}

// NewUserPasswordResetRequest creates new user password reset request form.
func NewUserPasswordResetRequest(app core.App) *UserPasswordResetRequest {
	return &UserPasswordResetRequest{
		app:             app,
		resendThreshold: 120, // 2 min
	}
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
//
// This method doesn't checks whether user with `form.Email` exists (this is done on Submit).
func (form *UserPasswordResetRequest) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(
			&form.Email,
			validation.Required,
			validation.Length(1, 255),
			is.Email,
		),
	)
}

// Submit validates and submits the form.
// On success sends a password reset email to the `form.Email` user.
func (form *UserPasswordResetRequest) Submit() error {
	if err := form.Validate(); err != nil {
		return err
	}

	user, err := form.app.Dao().FindUserByEmail(form.Email)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	lastResetSentAt := user.LastResetSentAt.Time()
	if now.Sub(lastResetSentAt).Seconds() < form.resendThreshold {
		return errors.New("You've already requested a password reset.")
	}

	if err := mails.SendUserPasswordReset(form.app, user); err != nil {
		return err
	}

	// update last sent timestamp
	user.LastResetSentAt = types.NowDateTime()

	return form.app.Dao().SaveUser(user)
}
