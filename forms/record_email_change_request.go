package forms

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/mails"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
)

// RecordEmailChangeRequest is an auth record email change request form.
type RecordEmailChangeRequest struct {
	app    core.App
	dao    *daos.Dao
	record *models.Record

	NewEmail string `form:"newEmail" json:"newEmail"`
}

// NewRecordEmailChangeRequest creates a new [RecordEmailChangeRequest] form
// initialized with from the provided [core.App] and [models.Record] instances.
//
// If you want to submit the form as part of a transaction,
// you can change the default Dao via [SetDao()].
func NewRecordEmailChangeRequest(app core.App, record *models.Record) *RecordEmailChangeRequest {
	return &RecordEmailChangeRequest{
		app:    app,
		dao:    app.Dao(),
		record: record,
	}
}

// SetDao replaces the default form Dao instance with the provided one.
func (form *RecordEmailChangeRequest) SetDao(dao *daos.Dao) {
	form.dao = dao
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *RecordEmailChangeRequest) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(
			&form.NewEmail,
			validation.Required,
			validation.Length(1, 255),
			is.EmailFormat,
			validation.By(form.checkUniqueEmail),
		),
	)
}

func (form *RecordEmailChangeRequest) checkUniqueEmail(value any) error {
	v, _ := value.(string)

	if !form.dao.IsRecordValueUnique(form.record.Collection().Id, schema.FieldNameEmail, v) {
		return validation.NewError("validation_record_email_invalid", "User email already exists or it is invalid.")
	}

	return nil
}

// Submit validates and sends the change email request.
//
// You can optionally provide a list of InterceptorFunc to
// further modify the form behavior before persisting it.
func (form *RecordEmailChangeRequest) Submit(interceptors ...InterceptorFunc[*models.Record]) error {
	if err := form.Validate(); err != nil {
		return err
	}

	return runInterceptors(form.record, func(m *models.Record) error {
		return mails.SendRecordChangeEmail(form.app, m, form.NewEmail)
	}, interceptors...)
}
