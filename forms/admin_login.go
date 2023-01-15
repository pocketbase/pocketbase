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
//
// You can optionally provide a list of InterceptorFunc to
// further modify the form behavior before persisting it.
func (form *AdminLogin) Submit(interceptors ...InterceptorFunc[*models.Admin]) (*models.Admin, error) {
	if err := form.Validate(); err != nil {
		return nil, err
	}

	admin, fetchErr := form.dao.FindAdminByEmail(form.Identity)

	// ignore not found errors to allow custom fetch implementations
	if fetchErr != nil && !errors.Is(fetchErr, sql.ErrNoRows) {
		return nil, fetchErr
	}

	interceptorsErr := runInterceptors(admin, func(m *models.Admin) error {
		admin = m

		if admin == nil || !admin.ValidatePassword(form.Password) {
			return errors.New("Invalid login credentials.")
		}

		return nil
	}, interceptors...)

	if interceptorsErr != nil {
		return nil, interceptorsErr
	}

	return admin, nil
}
