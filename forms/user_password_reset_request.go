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

// UserPasswordResetRequest specifies a user password reset request form.
type UserPasswordResetRequest struct {
	config UserPasswordResetRequestConfig

	Email string `form:"email" json:"email"`
}

// UserPasswordResetRequestConfig is the [UserPasswordResetRequest]
// factory initializer config.
//
// NB! App is required struct member.
type UserPasswordResetRequestConfig struct {
	App             core.App
	Dao             *daos.Dao
	ResendThreshold float64 // in seconds
}

// NewUserPasswordResetRequest creates a new [UserPasswordResetRequest]
// form with initializer config created from the provided [core.App] instance.
//
// If you want to submit the form as part of another transaction, use
// [NewUserPasswordResetRequestWithConfig] with explicitly set Dao.
func NewUserPasswordResetRequest(app core.App) *UserPasswordResetRequest {
	return NewUserPasswordResetRequestWithConfig(UserPasswordResetRequestConfig{
		App:             app,
		ResendThreshold: 120, // 2 min
	})
}

// NewUserPasswordResetRequestWithConfig creates a new [UserPasswordResetRequest]
// form with the provided config or panics on invalid configuration.
func NewUserPasswordResetRequestWithConfig(config UserPasswordResetRequestConfig) *UserPasswordResetRequest {
	form := &UserPasswordResetRequest{config: config}

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
// This method doesn't checks whether user with `form.Email` exists (this is done on Submit).
func (form *UserPasswordResetRequest) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(
			&form.Email,
			validation.Required,
			validation.Length(1, 255),
			is.EmailFormat,
		),
	)
}

// Submit validates and submits the form.
// On success sends a password reset email to the `form.Email` user.
func (form *UserPasswordResetRequest) Submit() error {
	if err := form.Validate(); err != nil {
		return err
	}

	user, err := form.config.Dao.FindUserByEmail(form.Email)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	lastResetSentAt := user.LastResetSentAt.Time()
	if now.Sub(lastResetSentAt).Seconds() < form.config.ResendThreshold {
		return errors.New("You've already requested a password reset.")
	}

	if err := mails.SendUserPasswordReset(form.config.App, user); err != nil {
		return err
	}

	// update last sent timestamp
	user.LastResetSentAt = types.NowDateTime()

	return form.config.Dao.SaveUser(user)
}
