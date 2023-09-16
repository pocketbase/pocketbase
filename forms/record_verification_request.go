package forms

import (
	"errors"
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

// RecordVerificationRequest is an auth record email verification request form.
type RecordVerificationRequest struct {
	app             core.App
	collection      *models.Collection
	dao             *daos.Dao
	resendThreshold float64 // in seconds

	Email string `form:"email" json:"email"`
}

// NewRecordVerificationRequest creates a new [RecordVerificationRequest]
// form initialized with from the provided [core.App] instance.
//
// If you want to submit the form as part of a transaction,
// you can change the default Dao via [SetDao()].
func NewRecordVerificationRequest(app core.App, collection *models.Collection) *RecordVerificationRequest {
	return &RecordVerificationRequest{
		app:             app,
		dao:             app.Dao(),
		collection:      collection,
		resendThreshold: 120, // 2 min
	}
}

// SetDao replaces the default form Dao instance with the provided one.
func (form *RecordVerificationRequest) SetDao(dao *daos.Dao) {
	form.dao = dao
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
//
// // This method doesn't verify that auth record with `form.Email` exists (this is done on Submit).
func (form *RecordVerificationRequest) Validate() error {
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
// to the `form.Email` auth record.
//
// You can optionally provide a list of InterceptorFunc to further
// modify the form behavior before persisting it.
func (form *RecordVerificationRequest) Submit(interceptors ...InterceptorFunc[*models.Record]) error {
	if err := form.Validate(); err != nil {
		return err
	}

	record, err := form.dao.FindFirstRecordByData(
		form.collection.Id,
		schema.FieldNameEmail,
		form.Email,
	)
	if err != nil {
		return err
	}

	if !record.Verified() {
		now := time.Now().UTC()
		lastVerificationSentAt := record.LastVerificationSentAt().Time()
		if (now.Sub(lastVerificationSentAt)).Seconds() < form.resendThreshold {
			return errors.New("A verification email was already sent.")
		}
	}

	return runInterceptors(record, func(m *models.Record) error {
		if m.Verified() {
			return nil // already verified
		}

		if err := mails.SendRecordVerification(form.app, m); err != nil {
			return err
		}

		// update last sent timestamp
		m.SetLastVerificationSentAt(types.NowDateTime())

		return form.dao.SaveRecord(m)
	}, interceptors...)
}
