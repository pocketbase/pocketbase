package forms

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/forms/validators"
	"github.com/pocketbase/pocketbase/models"
)

// AdminUpsert specifies a [models.Admin] upsert (create/update) form.
type AdminUpsert struct {
	config AdminUpsertConfig
	admin  *models.Admin

	Id              string `form:"id" json:"id"`
	Avatar          int    `form:"avatar" json:"avatar"`
	Email           string `form:"email" json:"email"`
	Password        string `form:"password" json:"password"`
	PasswordConfirm string `form:"passwordConfirm" json:"passwordConfirm"`
}

// AdminUpsertConfig is the [AdminUpsert] factory initializer config.
//
// NB! App is a required struct member.
type AdminUpsertConfig struct {
	App core.App
	Dao *daos.Dao
}

// NewAdminUpsert creates a new [AdminUpsert] form with initializer
// config created from the provided [core.App] and [models.Admin] instances
// (for create you could pass a pointer to an empty Admin - `&models.Admin{}`).
//
// If you want to submit the form as part of another transaction, use
// [NewAdminUpsertWithConfig] with explicitly set Dao.
func NewAdminUpsert(app core.App, admin *models.Admin) *AdminUpsert {
	return NewAdminUpsertWithConfig(AdminUpsertConfig{
		App: app,
	}, admin)
}

// NewAdminUpsertWithConfig creates a new [AdminUpsert] form
// with the provided config and [models.Admin] instance or panics on invalid configuration
// (for create you could pass a pointer to an empty Admin - `&models.Admin{}`).
func NewAdminUpsertWithConfig(config AdminUpsertConfig, admin *models.Admin) *AdminUpsert {
	form := &AdminUpsert{
		config: config,
		admin:  admin,
	}

	if form.config.App == nil || form.admin == nil {
		panic("Invalid initializer config or nil upsert model.")
	}

	if form.config.Dao == nil {
		form.config.Dao = form.config.App.Dao()
	}

	// load defaults
	form.Id = admin.Id
	form.Avatar = admin.Avatar
	form.Email = admin.Email

	return form
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
			validation.Length(10, 100),
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

	if form.config.Dao.IsAdminEmailUnique(v, form.admin.Id) {
		return nil
	}

	return validation.NewError("validation_admin_email_exists", "Admin email already exists.")
}

// Submit validates the form and upserts the form admin model.
//
// You can optionally provide a list of InterceptorFunc to further
// modify the form behavior before persisting it.
func (form *AdminUpsert) Submit(interceptors ...InterceptorFunc) error {
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

	return runInterceptors(func() error {
		return form.config.Dao.SaveAdmin(form.admin)
	}, interceptors...)
}
