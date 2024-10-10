package apis

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/spf13/cast"
)

func recordConfirmVerification(e *core.RequestEvent) error {
	collection, err := findAuthCollection(e)
	if err != nil {
		return err
	}

	if collection.Name == core.CollectionNameSuperusers {
		return e.BadRequestError("All superusers are verified by default.", nil)
	}

	form := new(recordConfirmVerificationForm)
	form.app = e.App
	form.collection = collection
	if err = e.BindBody(form); err != nil {
		return firstApiError(err, e.BadRequestError("An error occurred while loading the submitted data.", err))
	}
	if err = form.validate(); err != nil {
		return firstApiError(err, e.BadRequestError("An error occurred while validating the submitted data.", err))
	}

	record, err := form.app.FindAuthRecordByToken(form.Token, core.TokenTypeVerification)
	if err != nil {
		return e.BadRequestError("Invalid or expired verification token.", err)
	}

	wasVerified := record.Verified()

	event := new(core.RecordConfirmVerificationRequestEvent)
	event.RequestEvent = e
	event.Collection = collection
	event.Record = record

	return e.App.OnRecordConfirmVerificationRequest().Trigger(event, func(e *core.RecordConfirmVerificationRequestEvent) error {
		if wasVerified {
			return e.NoContent(http.StatusNoContent)
		}

		e.Record.SetVerified(true)

		if err := e.App.Save(e.Record); err != nil {
			return firstApiError(err, e.BadRequestError("An error occurred while saving the verified state.", err))
		}

		e.App.Store().Remove(getVerificationResendKey(e.Record))

		return e.NoContent(http.StatusNoContent)
	})
}

// -------------------------------------------------------------------

type recordConfirmVerificationForm struct {
	app        core.App
	collection *core.Collection

	Token string `form:"token" json:"token"`
}

func (form *recordConfirmVerificationForm) validate() error {
	return validation.ValidateStruct(form,
		validation.Field(&form.Token, validation.Required, validation.By(form.checkToken)),
	)
}

func (form *recordConfirmVerificationForm) checkToken(value any) error {
	v, _ := value.(string)
	if v == "" {
		return nil // nothing to check
	}

	claims, _ := security.ParseUnverifiedJWT(v)
	email := cast.ToString(claims["email"])
	if email == "" {
		return validation.NewError("validation_invalid_token_claims", "Missing email token claim.")
	}

	record, err := form.app.FindAuthRecordByToken(v, core.TokenTypeVerification)
	if err != nil || record == nil {
		return validation.NewError("validation_invalid_token", "Invalid or expired token.")
	}

	if record.Collection().Id != form.collection.Id {
		return validation.NewError("validation_token_collection_mismatch", "The provided token is for different auth collection.")
	}

	if record.Email() != email {
		return validation.NewError("validation_token_email_mismatch", "The record email doesn't match with the requested token claims.")
	}

	return nil
}
