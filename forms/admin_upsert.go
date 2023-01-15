package forms

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/forms/validators"
	"github.com/pocketbase/pocketbase/models"
)

// AdminUpsert is a [models.Admin] upsert (create/update) form.
type AdminUpsert struct {
	app   core.App
	dao   *daos.Dao
	admin *models.Admin

	Id              string `form:"id" json:"id"`
	Avatar          int    `form:"avatar" json:"avatar"`
	Email           string `form:"email" json:"email"`
	Password        string `form:"password" json:"password"`
	PasswordConfirm string `form:"passwordConfirm" json:"passwordConfirm"`
}

// NewAdminUpsert creates a new [AdminUpsert] form with initializer
// config created from the provided [core.App] and [models.Admin] instances
// (for create you could pass a pointer to an empty Admin - `&models.Admin{}`).
//
// If you want to submit the form as part of a transaction,
// you can change the default Dao via [SetDao()].
func NewAdminUpsert(app core.App, admin *models.Admin) *AdminUpsert {
	form := &AdminUpsert{
		app:   app,
		dao:   app.Dao(),
		admin: admin,
	}

	// load defaults
	form.Id = admin.Id
	form.Avatar = admin.Avatar
	form.Email = admin.Email

	return form
}

// SetDao replaces the default form Dao instance with the provided one.
func (form *AdminUpsert) SetDao(dao *daos.Dao) {
	form.dao = dao
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *AdminUpsert) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(
			&form.Id,
			validation.When(
				form.admin.IsNew(),
				validation.Length(models.DefaultIdLength, models.DefaultIdLength),
				validation.Match(idRegex),
				validation.By(validators.UniqueId(form.dao, form.admin.TableName())),
			).Else(validation.In(form.admin.Id)),
		),
		validation.Field(
			&form.Avatar,
			validation.Min(0),
			validation.Max(9),
		),
		validation.Field(
			&form.Email,
			validation.Required,
			validation.Length(1, 255),
			is.EmailFormat,
			validation.By(form.checkUniqueEmail),
		),
		validation.Field(
			&form.Password,
			validation.When(form.admin.IsNew(), validation.Required),
			validation.Length(10, 72),
		),
		validation.Field(
			&form.PasswordConfirm,
			validation.When(form.Password != "", validation.Required),
			validation.By(validators.Compare(form.Password)),
		),
	)
}

func (form *AdminUpsert) checkUniqueEmail(value any) error {
	v, _ := value.(string)

	if form.dao.IsAdminEmailUnique(v, form.admin.Id) {
		return nil
	}

	return validation.NewError("validation_admin_email_exists", "Admin email already exists.")
}

// Submit validates the form and upserts the form admin model.
//
// You can optionally provide a list of InterceptorFunc to further
// modify the form behavior before persisting it.
func (form *AdminUpsert) Submit(interceptors ...InterceptorFunc[*models.Admin]) error {
	if err := form.Validate(); err != nil {
		return err
	}

	// custom insertion id can be set only on create
	if form.admin.IsNew() && form.Id != "" {
		form.admin.MarkAsNew()
		form.admin.SetId(form.Id)
	}

	form.admin.Avatar = form.Avatar
	form.admin.Email = form.Email

	if form.Password != "" {
		form.admin.SetPassword(form.Password)
	}

	return runInterceptors(form.admin, func(admin *models.Admin) error {
		return form.dao.SaveAdmin(admin)
	}, interceptors...)
}
