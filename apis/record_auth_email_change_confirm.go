package apis

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/security"
)

func recordConfirmEmailChange(e *core.RequestEvent) error {
	collection, err := findAuthCollection(e)
	if err != nil {
		return err
	}

	if collection.Name == core.CollectionNameSuperusers {
		return e.BadRequestError("All superusers can change their emails directly.", nil)
	}

	form := newEmailChangeConfirmForm(e.App, collection)
	if err = e.BindBody(form); err != nil {
		return firstApiError(err, e.BadRequestError("An error occurred while loading the submitted data.", err))
	}
	if err = form.validate(); err != nil {
		return firstApiError(err, e.BadRequestError("An error occurred while validating the submitted data.", err))
	}

	authRecord, newEmail, err := form.parseToken()
	if err != nil {
		return firstApiError(err, e.BadRequestError("Invalid or expired token.", err))
	}

	event := new(core.RecordConfirmEmailChangeRequestEvent)
	event.RequestEvent = e
	event.Collection = collection
	event.Record = authRecord
	event.NewEmail = newEmail

	return e.App.OnRecordConfirmEmailChangeRequest().Trigger(event, func(e *core.RecordConfirmEmailChangeRequestEvent) error {
		e.Record.SetEmail(e.NewEmail)
		e.Record.SetVerified(true)

		if err := e.App.Save(e.Record); err != nil {
			return firstApiError(err, e.BadRequestError("Failed to confirm email change.", err))
		}

		return e.NoContent(http.StatusNoContent)
	})
}

// -------------------------------------------------------------------

func newEmailChangeConfirmForm(app core.App, collection *core.Collection) *EmailChangeConfirmForm {
	return &EmailChangeConfirmForm{
		app:        app,
		collection: collection,
	}
}

type EmailChangeConfirmForm struct {
	app        core.App
	collection *core.Collection

	Token    string `form:"token" json:"token"`
	Password string `form:"password" json:"password"`
}

func (form *EmailChangeConfirmForm) validate() error {
	return validation.ValidateStruct(form,
		validation.Field(&form.Token, validation.Required, validation.By(form.checkToken)),
		validation.Field(&form.Password, validation.Required, validation.Length(1, 100), validation.By(form.checkPassword)),
	)
}

func (form *EmailChangeConfirmForm) checkToken(value any) error {
	_, _, err := form.parseToken()
	return err
}

func (form *EmailChangeConfirmForm) checkPassword(value any) error {
	v, _ := value.(string)
	if v == "" {
		return nil // nothing to check
	}

	authRecord, _, _ := form.parseToken()
	if authRecord == nil || !authRecord.ValidatePassword(v) {
		return validation.NewError("validation_invalid_password", "Missing or invalid auth record password.")
	}

	return nil
}

func (form *EmailChangeConfirmForm) parseToken() (*core.Record, string, error) {
	// check token payload
	claims, _ := security.ParseUnverifiedJWT(form.Token)
	newEmail, _ := claims[core.TokenClaimNewEmail].(string)
	if newEmail == "" {
		return nil, "", validation.NewError("validation_invalid_token_payload", "Invalid token payload - newEmail must be set.")
	}

	// ensure that there aren't other users with the new email
	_, err := form.app.FindAuthRecordByEmail(form.collection, newEmail)
	if err == nil {
		return nil, "", validation.NewError("validation_existing_token_email", "The new email address is already registered: "+newEmail)
	}

	// verify that the token is not expired and its signature is valid
	authRecord, err := form.app.FindAuthRecordByToken(form.Token, core.TokenTypeEmailChange)
	if err != nil {
		return nil, "", validation.NewError("validation_invalid_token", "Invalid or expired token.")
	}

	if authRecord.Collection().Id != form.collection.Id {
		return nil, "", validation.NewError("validation_token_collection_mismatch", "The provided token is for different auth collection.")
	}

	return authRecord, newEmail, nil
}
