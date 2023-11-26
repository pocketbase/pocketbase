package forms

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/security"
)

// RecordEmailChangeConfirm is an auth record email change confirmation form.
type RecordEmailChangeConfirm struct {
	app        core.App
	dao        *daos.Dao
	collection *models.Collection

	Token    string `form:"token" json:"token"`
	Password string `form:"password" json:"password"`
}

// NewRecordEmailChangeConfirm creates a new [RecordEmailChangeConfirm] form
// initialized with from the provided [core.App] and [models.Collection] instances.
//
// If you want to submit the form as part of a transaction,
// you can change the default Dao via [SetDao()].
func NewRecordEmailChangeConfirm(app core.App, collection *models.Collection) *RecordEmailChangeConfirm {
	return &RecordEmailChangeConfirm{
		app:        app,
		dao:        app.Dao(),
		collection: collection,
	}
}

// SetDao replaces the default form Dao instance with the provided one.
func (form *RecordEmailChangeConfirm) SetDao(dao *daos.Dao) {
	form.dao = dao
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *RecordEmailChangeConfirm) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(
			&form.Token,
			validation.Required,
			validation.By(form.checkToken),
		),
		validation.Field(
			&form.Password,
			validation.Required,
			validation.Length(1, 100),
			validation.By(form.checkPassword),
		),
	)
}

func (form *RecordEmailChangeConfirm) checkToken(value any) error {
	v, _ := value.(string)
	if v == "" {
		return nil // nothing to check
	}

	authRecord, _, err := form.parseToken(v)
	if err != nil {
		return err
	}

	if authRecord.Collection().Id != form.collection.Id {
		return validation.NewError("validation_token_collection_mismatch", "The provided token is for different auth collection.")
	}

	return nil
}

func (form *RecordEmailChangeConfirm) checkPassword(value any) error {
	v, _ := value.(string)
	if v == "" {
		return nil // nothing to check
	}

	authRecord, _, _ := form.parseToken(form.Token)
	if authRecord == nil || !authRecord.ValidatePassword(v) {
		return validation.NewError("validation_invalid_password", "Missing or invalid auth record password.")
	}

	return nil
}

func (form *RecordEmailChangeConfirm) parseToken(token string) (*models.Record, string, error) {
	// check token payload
	claims, _ := security.ParseUnverifiedJWT(token)
	newEmail, _ := claims["newEmail"].(string)
	if newEmail == "" {
		return nil, "", validation.NewError("validation_invalid_token_payload", "Invalid token payload - newEmail must be set.")
	}

	// ensure that there aren't other users with the new email
	if !form.dao.IsRecordValueUnique(form.collection.Id, schema.FieldNameEmail, newEmail) {
		return nil, "", validation.NewError("validation_existing_token_email", "The new email address is already registered: "+newEmail)
	}

	// verify that the token is not expired and its signature is valid
	authRecord, err := form.dao.FindAuthRecordByToken(
		token,
		form.app.Settings().RecordEmailChangeToken.Secret,
	)
	if err != nil || authRecord == nil {
		return nil, "", validation.NewError("validation_invalid_token", "Invalid or expired token.")
	}

	return authRecord, newEmail, nil
}

// Submit validates and submits the auth record email change confirmation form.
// On success returns the updated auth record associated to `form.Token`.
//
// You can optionally provide a list of InterceptorFunc to
// further modify the form behavior before persisting it.
func (form *RecordEmailChangeConfirm) Submit(interceptors ...InterceptorFunc[*models.Record]) (*models.Record, error) {
	if err := form.Validate(); err != nil {
		return nil, err
	}

	authRecord, newEmail, err := form.parseToken(form.Token)
	if err != nil {
		return nil, err
	}

	authRecord.SetEmail(newEmail)
	authRecord.SetVerified(true)

	// @todo consider removing if not necessary anymore
	authRecord.RefreshTokenKey() // invalidate old tokens

	interceptorsErr := runInterceptors(authRecord, func(m *models.Record) error {
		authRecord = m
		return form.dao.SaveRecord(m)
	}, interceptors...)

	if interceptorsErr != nil {
		return nil, interceptorsErr
	}

	return authRecord, nil
}
