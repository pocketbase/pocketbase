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

func recordRequestPasswordReset(e *core.RequestEvent) error {
	collection, err := findAuthCollection(e)
	if err != nil {
		return err
	}

	if !collection.PasswordAuth.Enabled {
		return e.BadRequestError("The collection is not configured to allow password authentication.", nil)
	}

	form := new(recordRequestPasswordResetForm)
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

	resendKey := getPasswordResetResendKey(record)
	if e.App.Store().Has(resendKey) {
		// eagerly write 204 response as a very basic measure against emails enumeration
		e.NoContent(http.StatusNoContent)
		return errors.New("try again later - you've already requested a password reset email")
	}

	event := new(core.RecordRequestPasswordResetRequestEvent)
	event.RequestEvent = e
	event.Collection = collection
	event.Record = record

	return e.App.OnRecordRequestPasswordResetRequest().Trigger(event, func(e *core.RecordRequestPasswordResetRequestEvent) error {
		// run in background because we don't need to show the result to the client
		app := e.App
		routine.FireAndForget(func() {
			if err := mails.SendRecordPasswordReset(app, e.Record); err != nil {
				app.Logger().Error("Failed to send password reset email", "error", err)
				return
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

type recordRequestPasswordResetForm struct {
	Email string `form:"email" json:"email"`
}

func (form *recordRequestPasswordResetForm) validate() error {
	return validation.ValidateStruct(form,
		validation.Field(&form.Email, validation.Required, validation.Length(1, 255), is.EmailFormat),
	)
}

func getPasswordResetResendKey(record *core.Record) string {
	return "@limitPasswordResetEmail_" + record.Collection().Id + record.Id
}
