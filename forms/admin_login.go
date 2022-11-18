package forms

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

// AdminLogin is an admin email/pass login form.
type AdminLogin struct {
	app core.App
	dao *daos.Dao

	Identity string `form:"identity" json:"identity"`
	Password string `form:"password" json:"password"`
}

// NewAdminLogin creates a new [AdminLogin] form initialized with
// the provided [core.App] instance.
//
// If you want to submit the form as part of a transaction,
// you can change the default Dao via [SetDao()].
func NewAdminLogin(app core.App) *AdminLogin {
	return &AdminLogin{
		app: app,
		dao: app.Dao(),
	}
}

// SetDao replaces the default form Dao instance with the provided one.
func (form *AdminLogin) SetDao(dao *daos.Dao) {
	form.dao = dao
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *AdminLogin) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(&form.Identity, validation.Required, validation.Length(1, 255), is.EmailFormat),
		validation.Field(&form.Password, validation.Required, validation.Length(1, 255)),
	)
}

// Submit validates and submits the admin form.
// On success returns the authorized admin model.
func (form *AdminLogin) Submit() (*models.Admin, error) {
	if err := form.Validate(); err != nil {
		return nil, err
	}

	admin, err := form.dao.FindAdminByEmail(form.Identity)
	if err != nil {
		return nil, err
	}

	if admin.ValidatePassword(form.Password) {
		return admin, nil
	}

	return nil, errors.New("Invalid login credentials.")
}
