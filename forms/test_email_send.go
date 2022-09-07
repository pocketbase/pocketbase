package forms

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/mails"
	"github.com/pocketbase/pocketbase/models"
)

const (
	templateVerification  = "verification"
	templatePasswordReset = "password-reset"
	templateEmailChange   = "email-change"
)

// TestEmailSend is a email template test request form.
type TestEmailSend struct {
	app core.App

	Template string `form:"template" json:"template"`
	Email    string `form:"email" json:"email"`
}

// NewTestEmailSend creates and initializes new TestEmailSend form.
func NewTestEmailSend(app core.App) *TestEmailSend {
	return &TestEmailSend{app: app}
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *TestEmailSend) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(
			&form.Email,
			validation.Required,
			validation.Length(1, 255),
			is.EmailFormat,
		),
		validation.Field(
			&form.Template,
			validation.Required,
			validation.In(templateVerification, templateEmailChange, templatePasswordReset),
		),
	)
}

// Submit validates and sends a test email to the form.Email address.
func (form *TestEmailSend) Submit() error {
	if err := form.Validate(); err != nil {
		return err
	}

	// create a test user
	user := &models.User{}
	user.Id = "__pb_test_id__"
	user.Email = form.Email
	user.RefreshTokenKey()

	switch form.Template {
	case templateVerification:
		return mails.SendUserVerification(form.app, user)
	case templatePasswordReset:
		return mails.SendUserPasswordReset(form.app, user)
	case templateEmailChange:
		return mails.SendUserChangeEmail(form.app, user, form.Email)
	}

	return nil
}
