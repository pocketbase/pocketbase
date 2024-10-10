package apis

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/mails"
	"github.com/pocketbase/pocketbase/tools/routine"
)

func recordRequestVerification(e *core.RequestEvent) error {
	collection, err := findAuthCollection(e)
	if err != nil {
		return err
	}

	if collection.Name == core.CollectionNameSuperusers {
		return e.BadRequestError("All superusers are verified by default.", nil)
	}

	form := new(recordRequestVerificationForm)
	if err = e.BindBody(form); err != nil {
		return firstApiError(err, e.BadRequestError("An error occurred while loading the submitted data.", err))
	}
	if err = form.validate(); err != nil {
		return firstApiError(err, e.BadRequestError("An error occurred while validating the submitted data.", err))
	}

	record, err := e.App.FindAuthRecordByEmail(collection, form.Email)
	if err != nil {
		// eagerly write 204 response as a very basic measure against emails enumeration
		e.NoContent(http.StatusNoContent)
		return fmt.Errorf("failed to fetch %s record with email %s: %w", collection.Name, form.Email, err)
	}

	resendKey := getVerificationResendKey(record)
	if !record.Verified() && e.App.Store().Has(resendKey) {
		// eagerly write 204 response as a very basic measure against emails enumeration
		e.NoContent(http.StatusNoContent)
		return errors.New("try again later - you've already requested a verification email")
	}

	event := new(core.RecordRequestVerificationRequestEvent)
	event.RequestEvent = e
	event.Collection = collection
	event.Record = record

	return e.App.OnRecordRequestVerificationRequest().Trigger(event, func(e *core.RecordRequestVerificationRequestEvent) error {
		if e.Record.Verified() {
			return e.NoContent(http.StatusNoContent)
		}

		// run in background because we don't need to show the result to the client
		app := e.App
		routine.FireAndForget(func() {
			if err := mails.SendRecordVerification(app, e.Record); err != nil {
				app.Logger().Error("Failed to send verification email", "error", err)
			}

			app.Store().Set(resendKey, struct{}{})
			time.AfterFunc(2*time.Minute, func() {
				app.Store().Remove(resendKey)
			})
		})

		return e.NoContent(http.StatusNoContent)
	})
}

// -------------------------------------------------------------------

type recordRequestVerificationForm struct {
	Email string `form:"email" json:"email"`
}

func (form *recordRequestVerificationForm) validate() error {
	return validation.ValidateStruct(form,
		validation.Field(&form.Email, validation.Required, validation.Length(1, 255), is.EmailFormat),
	)
}

func getVerificationResendKey(record *core.Record) string {
	return "@limitVerificationEmail_" + record.Collection().Id + record.Id
}
