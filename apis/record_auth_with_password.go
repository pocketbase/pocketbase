package apis

import (
	"database/sql"
	"errors"
	"slices"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/dbutils"
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

	e.Set(core.RequestEventKeyInfoContext, core.RequestInfoContextPasswordAuth)

	var foundRecord *core.Record
	var foundErr error

	if form.IdentityField != "" {
		foundRecord, foundErr = findRecordByIdentityField(e.App, collection, form.IdentityField, form.Identity)
	} else {
		// prioritize email lookup
		isEmail := is.EmailFormat.Validate(form.Identity) == nil
		if isEmail && list.ExistInSlice(core.FieldNameEmail, collection.PasswordAuth.IdentityFields) {
			foundRecord, foundErr = findRecordByIdentityField(e.App, collection, core.FieldNameEmail, form.Identity)
		}

		// search by the other identity fields
		if !isEmail || foundErr != nil {
			for _, name := range collection.PasswordAuth.IdentityFields {
				if !isEmail && name == core.FieldNameEmail {
					continue // no need to search by the email field if it is not an email
				}

				foundRecord, foundErr = findRecordByIdentityField(e.App, collection, name, form.Identity)
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
		validation.Field(
			&form.IdentityField,
			validation.Length(1, 255),
			validation.In(list.ToInterfaceSlice(collection.PasswordAuth.IdentityFields)...),
		),
	)
}

func findRecordByIdentityField(app core.App, collection *core.Collection, field string, value any) (*core.Record, error) {
	if !slices.Contains(collection.PasswordAuth.IdentityFields, field) {
		return nil, errors.New("invalid identity field " + field)
	}

	index, ok := dbutils.FindSingleColumnUniqueIndex(collection.Indexes, field)
	if !ok {
		return nil, errors.New("missing " + field + " unique index constraint")
	}

	var expr dbx.Expression
	if strings.EqualFold(index.Columns[0].Collate, "nocase") {
		// case-insensitive search
		expr = dbx.NewExp("[["+field+"]] = {:identity} COLLATE NOCASE", dbx.Params{"identity": value})
	} else {
		expr = dbx.HashExp{field: value}
	}

	record := &core.Record{}

	err := app.RecordQuery(collection).AndWhere(expr).Limit(1).One(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}
