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
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/types"
)

// RecordPasswordResetRequest is an auth record reset password request form.
type RecordPasswordResetRequest struct {
	app             core.App
	dao             *daos.Dao
	collection      *models.Collection
	resendThreshold float64 // in seconds

	Email string `form:"email" json:"email"`
}

// NewRecordPasswordResetRequest creates a new [RecordPasswordResetRequest]
// form initialized with from the provided [core.App] instance.
//
// If you want to submit the form as part of a transaction,
// you can change the default Dao via [SetDao()].
func NewRecordPasswordResetRequest(app core.App, collection *models.Collection) *RecordPasswordResetRequest {
	return &RecordPasswordResetRequest{
		app:             app,
		dao:             app.Dao(),
		collection:      collection,
		resendThreshold: 120, // 2 min
	}
}

// SetDao replaces the default form Dao instance with the provided one.
func (form *RecordPasswordResetRequest) SetDao(dao *daos.Dao) {
	form.dao = dao
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
//
// This method doesn't check whether auth record with `form.Email` exists (this is done on Submit).
func (form *RecordPasswordResetRequest) Validate() error {
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
// On success, sends a password reset email to the `form.Email` auth record.
//
// You can optionally provide a list of InterceptorFunc to further
// modify the form behavior before persisting it.
func (form *RecordPasswordResetRequest) Submit(interceptors ...InterceptorFunc[*models.Record]) error {
	if err := form.Validate(); err != nil {
		return err
	}

	authRecord, err := form.dao.FindAuthRecordByEmail(form.collection.Id, form.Email)
	if err != nil {
		return fmt.Errorf("Failed to fetch %s record with email %s: %w", form.collection.Id, form.Email, err)
	}

	now := time.Now().UTC()
	lastResetSentAt := authRecord.LastResetSentAt().Time()
	if now.Sub(lastResetSentAt).Seconds() < form.resendThreshold {
		return errors.New("You've already requested a password reset.")
	}

	return runInterceptors(authRecord, func(m *models.Record) error {
		if err := mails.SendRecordPasswordReset(form.app, m); err != nil {
			return err
		}

		// update last sent timestamp
		m.Set(schema.FieldNameLastResetSentAt, types.NowDateTime())

		return form.dao.SaveRecord(m)
	}, interceptors...)
}
