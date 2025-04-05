package forms

import (
	"context"
	"errors"
	"fmt"
	"slices"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/core/validators"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/spf13/cast"
)

const (
	accessLevelDefault = iota
	accessLevelManager
	accessLevelSuperuser
)

type RecordUpsert struct {
	ctx         context.Context
	app         core.App
	record      *core.Record
	accessLevel int

	// extra password fields
	disablePasswordValidations bool
	password                   string
	passwordConfirm            string
	oldPassword                string
}

// NewRecordUpsert creates a new [RecordUpsert] form from the provided [core.App] and [core.Record] instances
// (for create you could pass a pointer to an empty Record - core.NewRecord(collection)).
func NewRecordUpsert(app core.App, record *core.Record) *RecordUpsert {
	form := &RecordUpsert{
		ctx:    context.Background(),
		app:    app,
		record: record,
	}

	return form
}

// SetContext assigns ctx as context of the current form.
func (form *RecordUpsert) SetContext(ctx context.Context) {
	form.ctx = ctx
}

// SetApp replaces the current form app instance.
//
// This could be used for example if you want to change at later stage
// before submission to change from regular -> transactional app instance.
func (form *RecordUpsert) SetApp(app core.App) {
	form.app = app
}

// SetRecord replaces the current form record instance.
func (form *RecordUpsert) SetRecord(record *core.Record) {
	form.record = record
}

// ResetAccess resets the form access level to the accessLevelDefault.
func (form *RecordUpsert) ResetAccess() {
	form.accessLevel = accessLevelDefault
}

// GrantManagerAccess updates the form access level to "manager" allowing
// directly changing some system record fields (often used with auth collection records).
func (form *RecordUpsert) GrantManagerAccess() {
	form.accessLevel = accessLevelManager
}

// GrantSuperuserAccess updates the form access level to "superuser" allowing
// directly changing all system record fields, including those marked as "Hidden".
func (form *RecordUpsert) GrantSuperuserAccess() {
	form.accessLevel = accessLevelSuperuser
}

// HasManageAccess reports whether the form has "manager" or "superuser" level access.
func (form *RecordUpsert) HasManageAccess() bool {
	return form.accessLevel == accessLevelManager || form.accessLevel == accessLevelSuperuser
}

// Load loads the provided data into the form and the related record.
func (form *RecordUpsert) Load(data map[string]any) {
	excludeFields := []string{core.FieldNameExpand}

	isAuth := form.record.Collection().IsAuth()

	// load the special auth form fields
	if isAuth {
		if v, ok := data["password"]; ok {
			form.password = cast.ToString(v)
		}
		if v, ok := data["passwordConfirm"]; ok {
			form.passwordConfirm = cast.ToString(v)
		}
		if v, ok := data["oldPassword"]; ok {
			form.oldPassword = cast.ToString(v)
		}

		excludeFields = append(excludeFields, "passwordConfirm", "oldPassword") // skip non-schema password fields
	}

	for k, v := range data {
		if slices.Contains(excludeFields, k) {
			continue
		}

		// set only known collection fields
		field := form.record.SetIfFieldExists(k, v)

		// restore original value if hidden field (with exception of the auth "password")
		//
		// note: this is an extra measure against loading hidden fields
		// but usually is not used by the default route handlers since
		// they filter additionally the data before calling Load
		if form.accessLevel != accessLevelSuperuser && field != nil && field.GetHidden() && (!isAuth || field.GetName() != core.FieldNamePassword) {
			form.record.SetRaw(field.GetName(), form.record.Original().GetRaw(field.GetName()))
		}
	}
}

func (form *RecordUpsert) validateFormFields() error {
	isAuth := form.record.Collection().IsAuth()
	if !isAuth {
		return nil
	}

	form.syncPasswordFields()

	isNew := form.record.IsNew()

	original := form.record.Original()

	validateData := map[string]any{
		"email":           form.record.Email(),
		"verified":        form.record.Verified(),
		"password":        form.password,
		"passwordConfirm": form.passwordConfirm,
		"oldPassword":     form.oldPassword,
	}

	return validation.Validate(validateData,
		validation.Map(
			validation.Key(
				"email",
				// don't allow direct email updates if the form doesn't have manage access permissions
				// (aka. allow only admin or authorized auth models to directly update the field)
				validation.When(
					!isNew && !form.HasManageAccess(),
					validation.By(validators.Equal(original.Email())),
				),
			),
			validation.Key(
				"verified",
				// don't allow changing verified if the form doesn't have manage access permissions
				// (aka. allow only admin or authorized auth models to directly change the field)
				validation.When(
					!form.HasManageAccess(),
					validation.By(validators.Equal(original.Verified())),
				),
			),
			validation.Key(
				"password",
				validation.When(
					!form.disablePasswordValidations && (isNew || form.passwordConfirm != "" || form.oldPassword != ""),
					validation.Required,
				),
			),
			validation.Key(
				"passwordConfirm",
				validation.When(
					!form.disablePasswordValidations && (isNew || form.password != "" || form.oldPassword != ""),
					validation.Required,
				),
				validation.When(!form.disablePasswordValidations, validation.By(validators.Equal(form.password))),
			),
			validation.Key(
				"oldPassword",
				// require old password only on update when:
				// - form.HasManageAccess() is not satisfied
				// - changing the existing password
				validation.When(
					!form.disablePasswordValidations && !isNew && !form.HasManageAccess() && (form.password != "" || form.passwordConfirm != ""),
					validation.Required,
					validation.By(form.checkOldPassword),
				),
			),
		),
	)
}

func (form *RecordUpsert) checkOldPassword(value any) error {
	v, ok := value.(string)
	if !ok {
		return validators.ErrUnsupportedValueType
	}

	if !form.record.Original().ValidatePassword(v) {
		return validation.NewError("validation_invalid_old_password", "Missing or invalid old password.")
	}

	return nil
}

// Deprecated: It was previously used as part of the record create action but it is not needed anymore and will be removed in the future.
//
// DrySubmit performs a temp form submit within a transaction and reverts it at the end.
// For actual record persistence, check the [RecordUpsert.Submit()] method.
//
// This method doesn't perform validations, handle file uploads/deletes or trigger app save events!
func (form *RecordUpsert) DrySubmit(callback func(txApp core.App, drySavedRecord *core.Record) error) error {
	isNew := form.record.IsNew()

	clone := form.record.Clone()

	// set an id if it doesn't have already
	// (the value doesn't matter; it is used only during the manual delete/update rollback)
	if clone.IsNew() && clone.Id == "" {
		clone.Id = "_temp_" + security.PseudorandomString(15)
	}

	app := form.app.UnsafeWithoutHooks()

	_, isTransactional := app.DB().(*dbx.Tx)
	if !isTransactional {
		return app.RunInTransaction(func(txApp core.App) error {
			tx, ok := txApp.DB().(*dbx.Tx)
			if !ok {
				return errors.New("failed to get transaction db")
			}
			defer tx.Rollback()

			if err := txApp.SaveNoValidateWithContext(form.ctx, clone); err != nil {
				return validators.NormalizeUniqueIndexError(err, clone.Collection().Name, clone.Collection().Fields.FieldNames())
			}

			if callback != nil {
				return callback(txApp, clone)
			}

			return nil
		})
	}

	// already in a transaction
	// (manual rollback to avoid starting another transaction)
	// ---------------------------------------------------------------
	err := app.SaveNoValidateWithContext(form.ctx, clone)
	if err != nil {
		return validators.NormalizeUniqueIndexError(err, clone.Collection().Name, clone.Collection().Fields.FieldNames())
	}

	manualRollback := func() error {
		if isNew {
			err = app.DeleteWithContext(form.ctx, clone)
			if err != nil {
				return fmt.Errorf("failed to rollback dry submit created record: %w", err)
			}
		} else {
			clone.Load(clone.Original().FieldsData())
			err = app.SaveNoValidateWithContext(form.ctx, clone)
			if err != nil {
				return fmt.Errorf("failed to rollback dry submit updated record: %w", err)
			}
		}

		return nil
	}

	if callback != nil {
		return errors.Join(callback(app, clone), manualRollback())
	}

	return manualRollback()
}

// Submit validates the form specific validations and attempts to save the form record.
func (form *RecordUpsert) Submit() error {
	err := form.validateFormFields()
	if err != nil {
		return err
	}

	// run record validations and persist in db
	return form.app.SaveWithContext(form.ctx, form.record)
}

// syncPasswordFields syncs the form's auth password fields with their
// corresponding record field values.
//
// This could be useful in case the password fields were programmatically set
// directly by modifying the related record model.
func (form *RecordUpsert) syncPasswordFields() {
	if !form.record.Collection().IsAuth() {
		return // not an auth collection
	}

	form.disablePasswordValidations = false

	rawPassword := form.record.GetRaw(core.FieldNamePassword)
	if v, ok := rawPassword.(*core.PasswordFieldValue); ok && v != nil {
		if
		// programmatically set custom plain password value
		(v.Plain != "" && v.Plain != form.password) ||
			// random generated password for new record
			(v.Plain == "" && v.Hash != "" && form.record.IsNew()) {
			form.disablePasswordValidations = true
		}
	}
}
