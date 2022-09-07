package forms

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/forms/validators"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/types"
)

// UserUpsert specifies a [models.User] upsert (create/update) form.
type UserUpsert struct {
	config UserUpsertConfig
	user   *models.User

	Id              string `form:"id" json:"id"`
	Email           string `form:"email" json:"email"`
	Password        string `form:"password" json:"password"`
	PasswordConfirm string `form:"passwordConfirm" json:"passwordConfirm"`
}

// UserUpsertConfig is the [UserUpsert] factory initializer config.
//
// NB! App is required struct member.
type UserUpsertConfig struct {
	App core.App
	Dao *daos.Dao
}

// NewUserUpsert creates a new [UserUpsert] form with initializer
// config created from the provided [core.App] instance
// (for create you could pass a pointer to an empty User - `&models.User{}`).
//
// If you want to submit the form as part of another transaction, use
// [NewUserEmailChangeConfirmWithConfig] with explicitly set Dao.
func NewUserUpsert(app core.App, user *models.User) *UserUpsert {
	return NewUserUpsertWithConfig(UserUpsertConfig{
		App: app,
	}, user)
}

// NewUserUpsertWithConfig creates a new [UserUpsert] form with the provided
// config and [models.User] instance or panics on invalid configuration
// (for create you could pass a pointer to an empty User - `&models.User{}`).
func NewUserUpsertWithConfig(config UserUpsertConfig, user *models.User) *UserUpsert {
	form := &UserUpsert{
		config: config,
		user:   user,
	}

	if form.config.App == nil || form.user == nil {
		panic("Invalid initializer config or nil upsert model.")
	}

	if form.config.Dao == nil {
		form.config.Dao = form.config.App.Dao()
	}

	// load defaults
	form.Id = user.Id
	form.Email = user.Email

	return form
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *UserUpsert) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(
			&form.Id,
			validation.When(
				form.user.IsNew(),
				validation.Length(models.DefaultIdLength, models.DefaultIdLength),
				validation.Match(idRegex),
			).Else(validation.In(form.user.Id)),
		),
		validation.Field(
			&form.Email,
			validation.Required,
			validation.Length(1, 255),
			is.EmailFormat,
			validation.By(form.checkEmailDomain),
			validation.By(form.checkUniqueEmail),
		),
		validation.Field(
			&form.Password,
			validation.When(form.user.IsNew(), validation.Required),
			validation.Length(form.config.App.Settings().EmailAuth.MinPasswordLength, 100),
		),
		validation.Field(
			&form.PasswordConfirm,
			validation.When(form.user.IsNew() || form.Password != "", validation.Required),
			validation.By(validators.Compare(form.Password)),
		),
	)
}

func (form *UserUpsert) checkUniqueEmail(value any) error {
	v, _ := value.(string)

	if v == "" || form.config.Dao.IsUserEmailUnique(v, form.user.Id) {
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
	only := form.config.App.Settings().EmailAuth.OnlyDomains
	except := form.config.App.Settings().EmailAuth.ExceptDomains

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

	// custom insertion id can be set only on create
	if form.user.IsNew() && form.Id != "" {
		form.user.MarkAsNew()
		form.user.SetId(form.Id)
	}

	if !form.user.IsNew() && form.Email != form.user.Email {
		form.user.Verified = false
		form.user.LastVerificationSentAt = types.DateTime{} // reset
	}

	form.user.Email = form.Email

	return runInterceptors(func() error {
		return form.config.Dao.SaveUser(form.user)
	}, interceptors...)
}
