package forms

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms/validators"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/types"
)

// UserUpsert defines a user upsert (create/update) form.
type UserUpsert struct {
	app      core.App
	user     *models.User
	isCreate bool

	Email           string `form:"email" json:"email"`
	Password        string `form:"password" json:"password"`
	PasswordConfirm string `form:"passwordConfirm" json:"passwordConfirm"`
}

// NewUserUpsert creates new upsert form for the provided user model
// (pass an empty user model instance (`&models.User{}`) for create).
func NewUserUpsert(app core.App, user *models.User) *UserUpsert {
	form := &UserUpsert{
		app:      app,
		user:     user,
		isCreate: !user.HasId(),
	}

	// load defaults
	form.Email = user.Email

	return form
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *UserUpsert) Validate() error {
	config := form.app.Settings()

	return validation.ValidateStruct(form,
		validation.Field(
			&form.Email,
			validation.Required,
			validation.Length(1, 255),
			is.Email,
			validation.By(form.checkEmailDomain),
			validation.By(form.checkUniqueEmail),
		),
		validation.Field(
			&form.Password,
			validation.When(form.isCreate, validation.Required),
			validation.Length(config.EmailAuth.MinPasswordLength, 100),
		),
		validation.Field(
			&form.PasswordConfirm,
			validation.When(form.isCreate || form.Password != "", validation.Required),
			validation.By(validators.Compare(form.Password)),
		),
	)
}

func (form *UserUpsert) checkUniqueEmail(value any) error {
	v, _ := value.(string)

	if v == "" || form.app.Dao().IsUserEmailUnique(v, form.user.Id) {
		return nil
	}

	return validation.NewError("validation_user_email_exists", "User email already exists.")
}

func (form *UserUpsert) checkEmailDomain(value any) error {
	val, _ := value.(string)
	if val == "" {
		return nil // nothing to check
	}

	domain := val[strings.LastIndex(val, "@")+1:]
	only := form.app.Settings().EmailAuth.OnlyDomains
	except := form.app.Settings().EmailAuth.ExceptDomains

	// only domains check
	if len(only) > 0 && !list.ExistInSlice(domain, only) {
		return validation.NewError("validation_email_domain_not_allowed", "Email domain is not allowed.")
	}

	// except domains check
	if len(except) > 0 && list.ExistInSlice(domain, except) {
		return validation.NewError("validation_email_domain_not_allowed", "Email domain is not allowed.")
	}

	return nil
}

// Submit validates the form and upserts the form user model.
//
// You can optionally provide a list of InterceptorFunc to further
// modify the form behavior before persisting it.
func (form *UserUpsert) Submit(interceptors ...InterceptorFunc) error {
	if err := form.Validate(); err != nil {
		return err
	}

	if form.Password != "" {
		form.user.SetPassword(form.Password)
	}

	if !form.isCreate && form.Email != form.user.Email {
		form.user.Verified = false
		form.user.LastVerificationSentAt = types.DateTime{} // reset
	}

	form.user.Email = form.Email

	return runInterceptors(func() error {
		return form.app.Dao().SaveUser(form.user)
	}, interceptors...)
}
