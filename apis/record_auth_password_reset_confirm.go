package apis

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/core/validators"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/spf13/cast"
)

func recordConfirmPasswordReset(e *core.RequestEvent) error {
	collection, err := findAuthCollection(e)
	if err != nil {
		return err
	}

	form := new(recordConfirmPasswordResetForm)
	form.app = e.App
	form.collection = collection
	if err = e.BindBody(form); err != nil {
		return e.BadRequestError("An error occurred while loading the submitted data.", err)
	}
	if err = form.validate(); err != nil {
		return firstApiError(err, e.BadRequestError("An error occurred while validating the submitted data.", err))
	}

	authRecord, err := e.App.FindAuthRecordByToken(form.Token, core.TokenTypePasswordReset)
	if err != nil {
		return firstApiError(err, e.BadRequestError("Invalid or expired password reset token.", err))
	}

	event := new(core.RecordConfirmPasswordResetRequestEvent)
	event.RequestEvent = e
	event.Collection = collection
	event.Record = authRecord

	return e.App.OnRecordConfirmPasswordResetRequest().Trigger(event, func(e *core.RecordConfirmPasswordResetRequestEvent) error {
		authRecord.SetPassword(form.Password)

		if !authRecord.Verified() {
			payload, err := security.ParseUnverifiedJWT(form.Token)
			if err == nil && authRecord.Email() == cast.ToString(payload[core.TokenClaimEmail]) {
				// mark as verified if the email hasn't changed
				authRecord.SetVerified(true)
			}
		}

		err = e.App.Save(authRecord)
		if err != nil {
			return firstApiError(err, e.BadRequestError("Failed to set new password.", err))
		}

		e.App.Store().Remove(getPasswordResetResendKey(authRecord))

		return e.NoContent(http.StatusNoContent)
	})
}

// -------------------------------------------------------------------

type recordConfirmPasswordResetForm struct {
	app        core.App
	collection *core.Collection

	Token           string `form:"token" json:"token"`
	Password        string `form:"password" json:"password"`
	PasswordConfirm string `form:"passwordConfirm" json:"passwordConfirm"`
}

func (form *recordConfirmPasswordResetForm) validate() error {
	min := 1
	passField, ok := form.collection.Fields.GetByName(core.FieldNamePassword).(*core.PasswordField)
	if ok && passField != nil && passField.Min > 0 {
		min = passField.Min
	}

	return validation.ValidateStruct(form,
		validation.Field(&form.Token, validation.Required, validation.By(form.checkToken)),
		validation.Field(&form.Password, validation.Required, validation.Length(min, 255)), // the FieldPassword validator will check further the specicic length constraints
		validation.Field(&form.PasswordConfirm, validation.Required, validation.By(validators.Equal(form.Password))),
	)
}

func (form *recordConfirmPasswordResetForm) checkToken(value any) error {
	v, _ := value.(string)
	if v == "" {
		return nil
	}

	record, err := form.app.FindAuthRecordByToken(v, core.TokenTypePasswordReset)
	if err != nil || record == nil {
		return validation.NewError("validation_invalid_token", "Invalid or expired token.")
	}

	if record.Collection().Id != form.collection.Id {
		return validation.NewError("validation_token_collection_mismatch", "The provided token is for different auth collection.")
	}

	return nil
}
