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

// AdminPasswordResetRequest specifies an admin password reset request form.
type AdminPasswordResetRequest struct {
	config AdminPasswordResetRequestConfig

	Email string `form:"email" json:"email"`
}

// AdminPasswordResetRequestConfig is the [AdminPasswordResetRequest] factory initializer config.
//
// NB! App is required struct member.
type AdminPasswordResetRequestConfig struct {
	App             core.App
	Dao             *daos.Dao
	ResendThreshold float64 // in seconds
}

// NewAdminPasswordResetRequest creates a new [AdminPasswordResetRequest]
// form with initializer config created from the provided [core.App] instance.
//
// If you want to submit the form as part of another transaction, use
// [NewAdminPasswordResetRequestWithConfig] with explicitly set Dao.
func NewAdminPasswordResetRequest(app core.App) *AdminPasswordResetRequest {
	return NewAdminPasswordResetRequestWithConfig(AdminPasswordResetRequestConfig{
		App:             app,
		ResendThreshold: 120, // 2min
	})
}

// NewAdminPasswordResetRequestWithConfig creates a new [AdminPasswordResetRequest]
// form with the provided config or panics on invalid configuration.
func NewAdminPasswordResetRequestWithConfig(config AdminPasswordResetRequestConfig) *AdminPasswordResetRequest {
	form := &AdminPasswordResetRequest{config: config}

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
// This method doesn't verify that admin with `form.Email` exists (this is done on Submit).
func (form *AdminPasswordResetRequest) Validate() error {
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
// On success sends a password reset email to the `form.Email` admin.
func (form *AdminPasswordResetRequest) Submit() error {
	if err := form.Validate(); err != nil {
		return err
	}

	admin, err := form.config.Dao.FindAdminByEmail(form.Email)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	lastResetSentAt := admin.LastResetSentAt.Time()
	if now.Sub(lastResetSentAt).Seconds() < form.config.ResendThreshold {
		return errors.New("You have already requested a password reset.")
	}

	if err := mails.SendAdminPasswordReset(form.config.App, admin); err != nil {
		return err
	}

	// update last sent timestamp
	admin.LastResetSentAt = types.NowDateTime()

	return form.config.Dao.SaveAdmin(admin)
}
