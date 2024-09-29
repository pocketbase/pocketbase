package apis

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/mails"
)

func recordRequestEmailChange(e *core.RequestEvent) error {
	collection, err := findAuthCollection(e)
	if err != nil {
		return err
	}

	if collection.Name == core.CollectionNameSuperusers {
		return e.BadRequestError("All superusers can change their emails directly.", nil)
	}

	record := e.Auth
	if record == nil {
		return e.UnauthorizedError("The request requires valid auth record.", nil)
	}

	form := newEmailChangeRequestForm(e.App, record)
	if err = e.BindBody(form); err != nil {
		return firstApiError(err, e.BadRequestError("An error occurred while loading the submitted data.", err))
	}
	if err = form.validate(); err != nil {
		return firstApiError(err, e.BadRequestError("An error occurred while validating the submitted data.", err))
	}

	event := new(core.RecordRequestEmailChangeRequestEvent)
	event.RequestEvent = e
	event.Collection = collection
	event.Record = record
	event.NewEmail = form.NewEmail

	return e.App.OnRecordRequestEmailChangeRequest().Trigger(event, func(e *core.RecordRequestEmailChangeRequestEvent) error {
		if err := mails.SendRecordChangeEmail(e.App, e.Record, e.NewEmail); err != nil {
			return firstApiError(err, e.BadRequestError("Failed to request email change.", err))
		}

		return e.NoContent(http.StatusNoContent)
	})
}

// -------------------------------------------------------------------

func newEmailChangeRequestForm(app core.App, record *core.Record) *emailChangeRequestForm {
	return &emailChangeRequestForm{
		app:    app,
		record: record,
	}
}

type emailChangeRequestForm struct {
	app    core.App
	record *core.Record

	NewEmail string `form:"newEmail" json:"newEmail"`
}

func (form *emailChangeRequestForm) validate() error {
	return validation.ValidateStruct(form,
		validation.Field(&form.NewEmail,
			validation.Required,
			validation.Length(1, 255),
			is.EmailFormat,
			validation.NotIn(form.record.Email()),
			validation.By(form.checkUniqueEmail),
		),
	)
}

func (form *emailChangeRequestForm) checkUniqueEmail(value any) error {
	v, _ := value.(string)
	if v == "" {
		return nil
	}

	found, _ := form.app.FindAuthRecordByEmail(form.record.Collection(), v)
	if found != nil && found.Id != form.record.Id {
		return validation.NewError("validation_invalid_new_email", "Invalid new email address.")
	}

	return nil
}
