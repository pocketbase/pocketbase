package forms

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/forms/validators"
	"github.com/pocketbase/pocketbase/models"
)

// AdminPasswordResetConfirm is an admin password reset confirmation form.
type AdminPasswordResetConfirm struct {
	app core.App
	dao *daos.Dao

	Token           string `form:"token" json:"token"`
	Password        string `form:"password" json:"password"`
	PasswordConfirm string `form:"passwordConfirm" json:"passwordConfirm"`
}

// NewAdminPasswordResetConfirm creates a new [AdminPasswordResetConfirm]
// form initialized with from the provided [core.App] instance.
//
// If you want to submit the form as part of a transaction,
// you can change the default Dao via [SetDao()].
func NewAdminPasswordResetConfirm(app core.App) *AdminPasswordResetConfirm {
	return &AdminPasswordResetConfirm{
		app: app,
		dao: app.Dao(),
	}
}

// SetDao replaces the form Dao instance with the provided one.
//
// This is useful if you want to use a specific transaction Dao instance
// instead of the default app.Dao().
func (form *AdminPasswordResetConfirm) SetDao(dao *daos.Dao) {
	form.dao = dao
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *AdminPasswordResetConfirm) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(&form.Token, validation.Required, validation.By(form.checkToken)),
		validation.Field(&form.Password, validation.Required, validation.Length(10, 72)),
		validation.Field(&form.PasswordConfirm, validation.Required, validation.By(validators.Compare(form.Password))),
	)
}

func (form *AdminPasswordResetConfirm) checkToken(value any) error {
	v, _ := value.(string)
	if v == "" {
		return nil // nothing to check
	}

	admin, err := form.dao.FindAdminByToken(v, form.app.Settings().AdminPasswordResetToken.Secret)
	if err != nil || admin == nil {
		return validation.NewError("validation_invalid_token", "Invalid or expired token.")
	}

	return nil
}

// Submit validates and submits the admin password reset confirmation form.
// On success returns the updated admin model associated to `form.Token`.
//
// You can optionally provide a list of InterceptorFunc to further
// modify the form behavior before persisting it.
func (form *AdminPasswordResetConfirm) Submit(interceptors ...InterceptorFunc[*models.Admin]) (*models.Admin, error) {
	if err := form.Validate(); err != nil {
		return nil, err
	}

	admin, err := form.dao.FindAdminByToken(
		form.Token,
		form.app.Settings().AdminPasswordResetToken.Secret,
	)
	if err != nil {
		return nil, err
	}

	if err := admin.SetPassword(form.Password); err != nil {
		return nil, err
	}

	interceptorsErr := runInterceptors(admin, func(m *models.Admin) error {
		admin = m
		return form.dao.SaveAdmin(m)
	}, interceptors...)

	if interceptorsErr != nil {
		return nil, interceptorsErr
	}

	return admin, nil
}
