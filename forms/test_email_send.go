package forms

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/mails"
	"github.com/pocketbase/pocketbase/tools/types"
)

const (
	TestTemplateVerification  = "verification"
	TestTemplatePasswordReset = "password-reset"
	TestTemplateEmailChange   = "email-change"
	TestTemplateOTP           = "otp"
	TestTemplateAuthAlert     = "login-alert"
)

// TestEmailSend is a email template test request form.
type TestEmailSend struct {
	app core.App

	Email      string `form:"email" json:"email"`
	Template   string `form:"template" json:"template"`
	Collection string `form:"collection" json:"collection"` // optional, fallbacks to _superusers
}

// NewTestEmailSend creates and initializes new TestEmailSend form.
func NewTestEmailSend(app core.App) *TestEmailSend {
	return &TestEmailSend{app: app}
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *TestEmailSend) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(
			&form.Collection,
			validation.Length(1, 255),
			validation.By(form.checkAuthCollection),
		),
		validation.Field(
			&form.Email,
			validation.Required,
			validation.Length(1, 255),
			is.EmailFormat,
		),
		validation.Field(
			&form.Template,
			validation.Required,
			validation.In(
				TestTemplateVerification,
				TestTemplatePasswordReset,
				TestTemplateEmailChange,
				TestTemplateOTP,
				TestTemplateAuthAlert,
			),
		),
	)
}

func (form *TestEmailSend) checkAuthCollection(value any) error {
	v, _ := value.(string)
	if v == "" {
		return nil // nothing to check
	}

	c, _ := form.app.FindCollectionByNameOrId(v)
	if c == nil || !c.IsAuth() {
		return validation.NewError("validation_invalid_auth_collection", "Must be a valid auth collection id or name.")
	}

	return nil
}

// Submit validates and sends a test email to the form.Email address.
func (form *TestEmailSend) Submit() error {
	if err := form.Validate(); err != nil {
		return err
	}

	collectionIdOrName := form.Collection
	if collectionIdOrName == "" {
		collectionIdOrName = core.CollectionNameSuperusers
	}

	collection, err := form.app.FindCollectionByNameOrId(collectionIdOrName)
	if err != nil {
		return err
	}

	record := core.NewRecord(collection)
	for _, field := range collection.Fields {
		if field.GetHidden() {
			continue
		}
		record.Set(field.GetName(), "__pb_test_"+field.GetName()+"__")
	}
	record.RefreshTokenKey()
	record.SetEmail(form.Email)

	switch form.Template {
	case TestTemplateVerification:
		return mails.SendRecordVerification(form.app, record)
	case TestTemplatePasswordReset:
		return mails.SendRecordPasswordReset(form.app, record)
	case TestTemplateEmailChange:
		return mails.SendRecordChangeEmail(form.app, record, form.Email)
	case TestTemplateOTP:
		return mails.SendRecordOTP(form.app, record, "_PB_TEST_OTP_ID_", "123456")
	case TestTemplateAuthAlert:
		testEvent := types.NowDateTime().String() + " - TEST_IP TEST_USER_AGENT"
		return mails.SendRecordAuthAlert(form.app, record, testEvent)
	default:
		return errors.New("unknown template " + form.Template)
	}
}
