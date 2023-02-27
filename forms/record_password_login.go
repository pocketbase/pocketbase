package forms

import (
	"database/sql"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

// RecordPasswordLogin is record username/email + password login form.
type RecordPasswordLogin struct {
	app        core.App
	dao        *daos.Dao
	collection *models.Collection

	Identity string `form:"identity" json:"identity"`
	Password string `form:"password" json:"password"`
}

// NewRecordPasswordLogin creates a new [RecordPasswordLogin] form initialized
// with from the provided [core.App] and [models.Collection] instance.
//
// If you want to submit the form as part of a transaction,
// you can change the default Dao via [SetDao()].
func NewRecordPasswordLogin(app core.App, collection *models.Collection) *RecordPasswordLogin {
	return &RecordPasswordLogin{
		app:        app,
		dao:        app.Dao(),
		collection: collection,
	}
}

// SetDao replaces the default form Dao instance with the provided one.
func (form *RecordPasswordLogin) SetDao(dao *daos.Dao) {
	form.dao = dao
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *RecordPasswordLogin) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(&form.Identity, validation.Required, validation.Length(1, 255)),
		validation.Field(&form.Password, validation.Required, validation.Length(1, 255)),
	)
}

// Submit validates and submits the form.
// On success returns the authorized record model.
//
// You can optionally provide a list of InterceptorFunc to
// further modify the form behavior before persisting it.
func (form *RecordPasswordLogin) Submit(interceptors ...InterceptorFunc[*models.Record]) (*models.Record, error) {
	if err := form.Validate(); err != nil {
		return nil, err
	}

	authOptions := form.collection.AuthOptions()

	var authRecord *models.Record
	var fetchErr error

	isEmail := is.EmailFormat.Validate(form.Identity) == nil

	if isEmail {
		if authOptions.AllowEmailAuth {
			authRecord, fetchErr = form.dao.FindAuthRecordByEmail(form.collection.Id, form.Identity)
		}
	} else if authOptions.AllowUsernameAuth {
		authRecord, fetchErr = form.dao.FindAuthRecordByUsername(form.collection.Id, form.Identity)
	}

	// ignore not found errors to allow custom fetch implementations
	if fetchErr != nil && !errors.Is(fetchErr, sql.ErrNoRows) {
		return nil, fetchErr
	}

	interceptorsErr := runInterceptors(authRecord, func(m *models.Record) error {
		authRecord = m

		if authRecord == nil || !authRecord.ValidatePassword(form.Password) {
			return errors.New("Invalid login credentials.")
		}

		return nil
	}, interceptors...)

	if interceptorsErr != nil {
		return nil, interceptorsErr
	}

	return authRecord, nil
}
