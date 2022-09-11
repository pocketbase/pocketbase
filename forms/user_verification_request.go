package forms

import (
	"errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/mails"
	"github.com/pocketbase/pocketbase/tools/types"
)

// UserVerificationRequest defines a user email verification request form.
type UserVerificationRequest struct {
	config UserVerificationRequestConfig

	Email string `form:"email" json:"email"`
}

// UserVerificationRequestConfig is the [UserVerificationRequest]
// factory initializer config.
//
// NB! App is required struct member.
type UserVerificationRequestConfig struct {
	App             core.App
	Dao             *daos.Dao
	ResendThreshold float64 // in seconds
}

// NewUserVerificationRequest creates a new [UserVerificationRequest]
// form with initializer config created from the provided [core.App] instance.
//
// If you want to submit the form as part of another transaction, use
// [NewUserVerificationRequestWithConfig] with explicitly set Dao.
func NewUserVerificationRequest(app core.App) *UserVerificationRequest {
	return NewUserVerificationRequestWithConfig(UserVerificationRequestConfig{
		App:             app,
		ResendThreshold: 120, // 2 min
	})
}

// NewUserVerificationRequestWithConfig creates a new [UserVerificationRequest]
// form with the provided config or panics on invalid configuration.
func NewUserVerificationRequestWithConfig(config UserVerificationRequestConfig) *UserVerificationRequest {
	form := &UserVerificationRequest{config: config}

	if form.config.App == nil {
		panic("Missing required config.App instance.")
	}

	if form.config.Dao == nil {
		form.config.Dao = form.config.App.Dao()
	}

	return form
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
			is.EmailFormat,
		),
	)
}

// Submit validates and sends a verification request email
// to the `form.Email` user.
func (form *UserVerificationRequest) Submit() error {
	if err := form.Validate(); err != nil {
		return err
	}

	user, err := form.config.Dao.FindUserByEmail(form.Email)
	if err != nil {
		return err
	}

	if user.Verified {
		return nil // already verified
	}

	now := time.Now().UTC()
	lastVerificationSentAt := user.LastVerificationSentAt.Time()
	if (now.Sub(lastVerificationSentAt)).Seconds() < form.config.ResendThreshold {
		return errors.New("A verification email was already sent.")
	}

	if err := mails.SendUserVerification(form.config.App, user); err != nil {
		return err
	}

	// update last sent timestamp
	user.LastVerificationSentAt = types.NowDateTime()

	return form.config.Dao.SaveUser(user)
}
