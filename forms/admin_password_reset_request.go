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

// AdminPasswordResetRequest defines an admin password reset request form.
type AdminPasswordResetRequest struct {
	app             core.App
	resendThreshold float64

	Email string `form:"email" json:"email"`
}

// NewAdminPasswordResetRequest creates new admin password reset request form.
func NewAdminPasswordResetRequest(app core.App) *AdminPasswordResetRequest {
	return &AdminPasswordResetRequest{
		app:             app,
		resendThreshold: 120, // 2 min
	}
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
			is.Email,
		),
	)
}

// Submit validates and submits the form.
// On success sends a password reset email to the `form.Email` admin.
func (form *AdminPasswordResetRequest) Submit() error {
	if err := form.Validate(); err != nil {
		return err
	}

	admin, err := form.app.Dao().FindAdminByEmail(form.Email)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	lastResetSentAt := admin.LastResetSentAt.Time()
	if now.Sub(lastResetSentAt).Seconds() < form.resendThreshold {
		return errors.New("You have already requested a password reset.")
	}

	if err := mails.SendAdminPasswordReset(form.app, admin); err != nil {
		return err
	}

	// update last sent timestamp
	admin.LastResetSentAt = types.NowDateTime()

	return form.app.Dao().SaveAdmin(admin)
}
