package forms

import (
	"errors"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/mails"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

// AdminPasswordResetRequest is an admin password reset request form.
type AdminPasswordResetRequest struct {
	app             core.App
	dao             *daos.Dao
	resendThreshold float64 // in seconds

	Email string `form:"email" json:"email"`
}

// NewAdminPasswordResetRequest creates a new [AdminPasswordResetRequest]
// form initialized with from the provided [core.App] instance.
//
// If you want to submit the form as part of a transaction,
// you can change the default Dao via [SetDao()].
func NewAdminPasswordResetRequest(app core.App) *AdminPasswordResetRequest {
	return &AdminPasswordResetRequest{
		app:             app,
		dao:             app.Dao(),
		resendThreshold: 120, // 2min
	}
}

// SetDao replaces the default form Dao instance with the provided one.
func (form *AdminPasswordResetRequest) SetDao(dao *daos.Dao) {
	form.dao = dao
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
//
// You can optionally provide a list of InterceptorFunc to further
// modify the form behavior before persisting it.
func (form *AdminPasswordResetRequest) Submit(interceptors ...InterceptorFunc[*models.Admin]) error {
	if err := form.Validate(); err != nil {
		return err
	}

	admin, err := form.dao.FindAdminByEmail(form.Email)
	if err != nil {
		return fmt.Errorf("Failed to fetch admin with email %s: %w", form.Email, err)
	}

	now := time.Now().UTC()
	lastResetSentAt := admin.LastResetSentAt.Time()
	if now.Sub(lastResetSentAt).Seconds() < form.resendThreshold {
		return errors.New("You have already requested a password reset.")
	}

	return runInterceptors(admin, func(m *models.Admin) error {
		if err := mails.SendAdminPasswordReset(form.app, m); err != nil {
			return err
		}

		// update last sent timestamp
		m.LastResetSentAt = types.NowDateTime()

		return form.dao.SaveAdmin(m)
	}, interceptors...)
}
