package forms

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/mails"
	"github.com/pocketbase/pocketbase/models"
)

// UserEmailChangeRequest defines a user email change request form.
type UserEmailChangeRequest struct {
	config UserEmailChangeRequestConfig
	user   *models.User

	NewEmail string `form:"newEmail" json:"newEmail"`
}

// UserEmailChangeRequestConfig is the [UserEmailChangeRequest] factory initializer config.
//
// NB! App is required struct member.
type UserEmailChangeRequestConfig struct {
	App core.App
	Dao *daos.Dao
}

// NewUserEmailChangeRequest creates a new [UserEmailChangeRequest]
// form with initializer config created from the provided [core.App] instance.
//
// If you want to submit the form as part of another transaction, use
// [NewUserEmailChangeConfirmWithConfig] with explicitly set Dao.
func NewUserEmailChangeRequest(app core.App, user *models.User) *UserEmailChangeRequest {
	return NewUserEmailChangeRequestWithConfig(UserEmailChangeRequestConfig{
		App: app,
	}, user)
}

// NewUserEmailChangeRequestWithConfig creates a new [UserEmailChangeRequest]
// form with the provided config or panics on invalid configuration.
func NewUserEmailChangeRequestWithConfig(config UserEmailChangeRequestConfig, user *models.User) *UserEmailChangeRequest {
	form := &UserEmailChangeRequest{
		config: config,
		user:   user,
	}

	if form.config.App == nil || form.user == nil {
		panic("Invalid initializer config or nil user model.")
	}

	if form.config.Dao == nil {
		form.config.Dao = form.config.App.Dao()
	}

	return form
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *UserEmailChangeRequest) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(
			&form.NewEmail,
			validation.Required,
			validation.Length(1, 255),
			is.EmailFormat,
			validation.By(form.checkUniqueEmail),
		),
	)
}

func (form *UserEmailChangeRequest) checkUniqueEmail(value any) error {
	v, _ := value.(string)

	if !form.config.Dao.IsUserEmailUnique(v, "") {
		return validation.NewError("validation_user_email_exists", "User email already exists.")
	}

	return nil
}

// Submit validates and sends the change email request.
func (form *UserEmailChangeRequest) Submit() error {
	if err := form.Validate(); err != nil {
		return err
	}

	return mails.SendUserChangeEmail(form.config.App, form.user, form.NewEmail)
}
