package forms

import (
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
func (form *RecordPasswordLogin) Submit() (*models.Record, error) {
	if err := form.Validate(); err != nil {
		return nil, err
	}

	authOptions := form.collection.AuthOptions()

	if !authOptions.AllowEmailAuth && !authOptions.AllowUsernameAuth {
		return nil, errors.New("Password authentication is not allowed for the collection.")
	}

	var record *models.Record
	var fetchErr error

	if authOptions.AllowEmailAuth &&
		(!authOptions.AllowUsernameAuth || is.EmailFormat.Validate(form.Identity) == nil) {
		record, fetchErr = form.dao.FindAuthRecordByEmail(form.collection.Id, form.Identity)
	} else {
		record, fetchErr = form.dao.FindAuthRecordByUsername(form.collection.Id, form.Identity)
	}

	if fetchErr != nil || !record.ValidatePassword(form.Password) {
		return nil, errors.New("Invalid login credentials.")
	}

	return record, nil
}
