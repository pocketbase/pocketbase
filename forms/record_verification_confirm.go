package forms

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/spf13/cast"
)

// RecordVerificationConfirm is an auth record email verification confirmation form.
type RecordVerificationConfirm struct {
	app        core.App
	collection *models.Collection
	dao        *daos.Dao

	Token string `form:"token" json:"token"`
}

// NewRecordVerificationConfirm creates a new [RecordVerificationConfirm]
// form initialized with from the provided [core.App] instance.
//
// If you want to submit the form as part of a transaction,
// you can change the default Dao via [SetDao()].
func NewRecordVerificationConfirm(app core.App, collection *models.Collection) *RecordVerificationConfirm {
	return &RecordVerificationConfirm{
		app:        app,
		dao:        app.Dao(),
		collection: collection,
	}
}

// SetDao replaces the default form Dao instance with the provided one.
func (form *RecordVerificationConfirm) SetDao(dao *daos.Dao) {
	form.dao = dao
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *RecordVerificationConfirm) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(&form.Token, validation.Required, validation.By(form.checkToken)),
	)
}

func (form *RecordVerificationConfirm) checkToken(value any) error {
	v, _ := value.(string)
	if v == "" {
		return nil // nothing to check
	}

	claims, _ := security.ParseUnverifiedJWT(v)
	email := cast.ToString(claims["email"])
	if email == "" {
		return validation.NewError("validation_invalid_token_claims", "Missing email token claim.")
	}

	record, err := form.dao.FindAuthRecordByToken(
		v,
		form.app.Settings().RecordVerificationToken.Secret,
	)
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

// Submit validates and submits the form.
// On success returns the verified auth record associated to `form.Token`.
//
// You can optionally provide a list of InterceptorFunc to further
// modify the form behavior before persisting it.
func (form *RecordVerificationConfirm) Submit(interceptors ...InterceptorFunc[*models.Record]) (*models.Record, error) {
	if err := form.Validate(); err != nil {
		return nil, err
	}

	record, err := form.dao.FindAuthRecordByToken(
		form.Token,
		form.app.Settings().RecordVerificationToken.Secret,
	)
	if err != nil {
		return nil, err
	}

	wasVerified := record.Verified()

	if !wasVerified {
		record.SetVerified(true)
	}

	interceptorsErr := runInterceptors(record, func(m *models.Record) error {
		record = m

		if wasVerified {
			return nil // already verified
		}

		return form.dao.SaveRecord(m)
	}, interceptors...)

	if interceptorsErr != nil {
		return nil, interceptorsErr
	}

	return record, nil
}
