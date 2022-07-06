package forms

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms/validators"
	"github.com/pocketbase/pocketbase/models"
)

// AdminUpsert defines an admin upsert (create/update) form.
type AdminUpsert struct {
	app      core.App
	admin    *models.Admin
	isCreate bool

	Avatar          int    `form:"avatar" json:"avatar"`
	Email           string `form:"email" json:"email"`
	Password        string `form:"password" json:"password"`
	PasswordConfirm string `form:"passwordConfirm" json:"passwordConfirm"`
}

// NewAdminUpsert creates new upsert form for the provided admin model
// (pass an empty admin model instance (`&models.Admin{}`) for create).
func NewAdminUpsert(app core.App, admin *models.Admin) *AdminUpsert {
	form := &AdminUpsert{
		app:      app,
		admin:    admin,
		isCreate: !admin.HasId(),
	}

	// load defaults
	form.Avatar = admin.Avatar
	form.Email = admin.Email

	return form
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *AdminUpsert) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(
			&form.Avatar,
			validation.Min(0),
			validation.Max(9),
		),
		validation.Field(
			&form.Email,
			validation.Required,
			validation.Length(1, 255),
			is.Email,
			validation.By(form.checkUniqueEmail),
		),
		validation.Field(
			&form.Password,
			validation.When(form.isCreate, validation.Required),
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

	if form.app.Dao().IsAdminEmailUnique(v, form.admin.Id) {
		return nil
	}

	return validation.NewError("validation_admin_email_exists", "Admin email already exists.")
}

// Submit validates the form and upserts the form's admin model.
func (form *AdminUpsert) Submit() error {
	if err := form.Validate(); err != nil {
		return err
	}

	form.admin.Avatar = form.Avatar
	form.admin.Email = form.Email

	if form.Password != "" {
		form.admin.SetPassword(form.Password)
	}

	return form.app.Dao().SaveAdmin(form.admin)
}
