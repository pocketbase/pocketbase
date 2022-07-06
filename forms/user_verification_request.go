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

// UserVerificationRequest defines a user email verification request form.
type UserVerificationRequest struct {
	app             core.App
	resendThreshold float64

	Email string `form:"email" json:"email"`
}

// NewUserVerificationRequest creates a new user email verification request form.
func NewUserVerificationRequest(app core.App) *UserVerificationRequest {
	return &UserVerificationRequest{
		app:             app,
		resendThreshold: 120, // 2 min
	}
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
//
// // This method doesn't verify that user with `form.Email` exists (this is done on Submit).
func (form *UserVerificationRequest) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(
			&form.Email,
			validation.Required,
			validation.Length(1, 255),
			is.Email,
		),
	)
}

// Submit validates and sends a verification request email
// to the `form.Email` user.
func (form *UserVerificationRequest) Submit() error {
	if err := form.Validate(); err != nil {
		return err
	}

	user, err := form.app.Dao().FindUserByEmail(form.Email)
	if err != nil {
		return err
	}

	if user.Verified {
		return nil // already verified
	}

	now := time.Now().UTC()
	lastVerificationSentAt := user.LastVerificationSentAt.Time()
	if (now.Sub(lastVerificationSentAt)).Seconds() < form.resendThreshold {
		return errors.New("A verification email was already sent.")
	}

	if err := mails.SendUserVerification(form.app, user); err != nil {
		return err
	}

	// update last sent timestamp
	user.LastVerificationSentAt = types.NowDateTime()

	return form.app.Dao().SaveUser(user)
}
