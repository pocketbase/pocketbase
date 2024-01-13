package forms

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/forms/validators"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/spf13/cast"
)

// RecordPasswordResetConfirm is an auth record password reset confirmation form.
type RecordPasswordResetConfirm struct {
	app        core.App
	collection *models.Collection
	dao        *daos.Dao

	Token           string `form:"token" json:"token"`
	Password        string `form:"password" json:"password"`
	PasswordConfirm string `form:"passwordConfirm" json:"passwordConfirm"`
}

// NewRecordPasswordResetConfirm creates a new [RecordPasswordResetConfirm]
// form initialized with from the provided [core.App] instance.
//
// If you want to submit the form as part of a transaction,
// you can change the default Dao via [SetDao()].
func NewRecordPasswordResetConfirm(app core.App, collection *models.Collection) *RecordPasswordResetConfirm {
	return &RecordPasswordResetConfirm{
		app:        app,
		dao:        app.Dao(),
		collection: collection,
	}
}

// SetDao replaces the default form Dao instance with the provided one.
func (form *RecordPasswordResetConfirm) SetDao(dao *daos.Dao) {
	form.dao = dao
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *RecordPasswordResetConfirm) Validate() error {
	minPasswordLength := form.collection.AuthOptions().MinPasswordLength

	return validation.ValidateStruct(form,
		validation.Field(&form.Token, validation.Required, validation.By(form.checkToken)),
		validation.Field(&form.Password, validation.Required, validation.Length(minPasswordLength, 100)),
		validation.Field(&form.PasswordConfirm, validation.Required, validation.By(validators.Compare(form.Password))),
	)
}

func (form *RecordPasswordResetConfirm) checkToken(value any) error {
	v, _ := value.(string)
	if v == "" {
		return nil // nothing to check
	}

	record, err := form.dao.FindAuthRecordByToken(
		v,
		form.app.Settings().RecordPasswordResetToken.Secret,
	)
	if err != nil || record == nil {
		return validation.NewError("validation_invalid_token", "Invalid or expired token.")
	}

	if record.Collection().Id != form.collection.Id {
		return validation.NewError("validation_token_collection_mismatch", "The provided token is for different auth collection.")
	}

	return nil
}

// Submit validates and submits the form.
// On success returns the updated auth record associated to `form.Token`.
//
// You can optionally provide a list of InterceptorFunc to further
// modify the form behavior before persisting it.
func (form *RecordPasswordResetConfirm) Submit(interceptors ...InterceptorFunc[*models.Record]) (*models.Record, error) {
	if err := form.Validate(); err != nil {
		return nil, err
	}

	authRecord, err := form.dao.FindAuthRecordByToken(
		form.Token,
		form.app.Settings().RecordPasswordResetToken.Secret,
	)
	if err != nil {
		return nil, err
	}

	if err := authRecord.SetPassword(form.Password); err != nil {
		return nil, err
	}

	if !authRecord.Verified() {
		payload, err := security.ParseUnverifiedJWT(form.Token)
		if err != nil {
			return nil, err
		}

		// mark as verified if the email hasn't changed
		if authRecord.Email() == cast.ToString(payload["email"]) {
			authRecord.SetVerified(true)
		}
	}

	interceptorsErr := runInterceptors(authRecord, func(m *models.Record) error {
		authRecord = m
		return form.dao.SaveRecord(authRecord)
	}, interceptors...)

	if interceptorsErr != nil {
		return nil, interceptorsErr
	}

	return authRecord, nil
}
