package apis

import (
	"database/sql"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/list"
)

func recordAuthWithPassword(e *core.RequestEvent) error {
	collection, err := findAuthCollection(e)
	if err != nil {
		return err
	}

	if !collection.PasswordAuth.Enabled {
		return e.ForbiddenError("The collection is not configured to allow password authentication.", nil)
	}

	form := &authWithPasswordForm{}
	if err = e.BindBody(form); err != nil {
		return firstApiError(err, e.BadRequestError("An error occurred while loading the submitted data.", err))
	}
	if err = form.validate(collection); err != nil {
		return firstApiError(err, e.BadRequestError("An error occurred while validating the submitted data.", err))
	}

	var foundRecord *core.Record
	var foundErr error

	if form.IdentityField != "" {
		foundRecord, foundErr = e.App.FindFirstRecordByData(collection.Id, form.IdentityField, form.Identity)
	} else {
		// prioritize email lookup
		isEmail := is.EmailFormat.Validate(form.Identity) == nil
		if isEmail && list.ExistInSlice(core.FieldNameEmail, collection.PasswordAuth.IdentityFields) {
			foundRecord, foundErr = e.App.FindAuthRecordByEmail(collection.Id, form.Identity)
		}

		// search by the other identity fields
		if !isEmail || foundErr != nil {
			for _, name := range collection.PasswordAuth.IdentityFields {
				if !isEmail && name == core.FieldNameEmail {
					continue // no need to search by the email field if it is not an email
				}

				foundRecord, foundErr = e.App.FindFirstRecordByData(collection.Id, name, form.Identity)
				if foundErr == nil {
					break
				}
			}
		}
	}

	// ignore not found errors to allow custom record find implementations
	if foundErr != nil && !errors.Is(foundErr, sql.ErrNoRows) {
		return e.InternalServerError("", foundErr)
	}

	event := new(core.RecordAuthWithPasswordRequestEvent)
	event.RequestEvent = e
	event.Collection = collection
	event.Record = foundRecord
	event.Identity = form.Identity
	event.Password = form.Password
	event.IdentityField = form.IdentityField

	return e.App.OnRecordAuthWithPasswordRequest().Trigger(event, func(e *core.RecordAuthWithPasswordRequestEvent) error {
		if e.Record == nil || !e.Record.ValidatePassword(e.Password) {
			return e.BadRequestError("Failed to authenticate.", errors.New("invalid login credentials"))
		}

		return RecordAuthResponse(e.RequestEvent, e.Record, core.MFAMethodPassword, nil)
	})
}

// -------------------------------------------------------------------

type authWithPasswordForm struct {
	Identity string `form:"identity" json:"identity"`
	Password string `form:"password" json:"password"`

	// IdentityField specifies the field to use to search for the identity
	// (leave it empty for "auto" detection).
	IdentityField string `form:"identityField" json:"identityField"`
}

func (form *authWithPasswordForm) validate(collection *core.Collection) error {
	return validation.ValidateStruct(form,
		validation.Field(&form.Identity, validation.Required, validation.Length(1, 255)),
		validation.Field(&form.Password, validation.Required, validation.Length(1, 255)),
		validation.Field(&form.IdentityField, validation.In(list.ToInterfaceSlice(collection.PasswordAuth.IdentityFields)...)),
	)
}
