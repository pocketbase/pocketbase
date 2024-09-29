package forms_test

import (
	"fmt"
	"strings"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/tests"
)

func TestEmailSendValidateAndSubmit(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		template       string
		email          string
		collection     string
		expectedErrors []string
	}{
		{"", "", "", []string{"template", "email"}},
		{"invalid", "test@example.com", "", []string{"template"}},
		{forms.TestTemplateVerification, "invalid", "", []string{"email"}},
		{forms.TestTemplateVerification, "test@example.com", "invalid", []string{"collection"}},
		{forms.TestTemplateVerification, "test@example.com", "demo1", []string{"collection"}},
		{forms.TestTemplateVerification, "test@example.com", "users", nil},
		{forms.TestTemplatePasswordReset, "test@example.com", "", nil},
		{forms.TestTemplateEmailChange, "test@example.com", "", nil},
		{forms.TestTemplateOTP, "test@example.com", "", nil},
		{forms.TestTemplateAuthAlert, "test@example.com", "", nil},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%s", i, s.template), func(t *testing.T) {
			app, _ := tests.NewTestApp()
			defer app.Cleanup()

			form := forms.NewTestEmailSend(app)
			form.Email = s.email
			form.Template = s.template
			form.Collection = s.collection

			result := form.Submit()

			// parse errors
			errs, ok := result.(validation.Errors)
			if !ok && result != nil {
				t.Fatalf("Failed to parse errors %v", result)
			}

			// check errors
			if len(errs) > len(s.expectedErrors) {
				t.Fatalf("Expected error keys %v, got %v", s.expectedErrors, errs)
			}
			for _, k := range s.expectedErrors {
				if _, ok := errs[k]; !ok {
					t.Fatalf("Missing expected error key %q in %v", k, errs)
				}
			}

			expectedEmails := 1
			if len(s.expectedErrors) > 0 {
				expectedEmails = 0
			}

			if app.TestMailer.TotalSend() != expectedEmails {
				t.Fatalf("Expected %d email(s) to be sent, got %d", expectedEmails, app.TestMailer.TotalSend())
			}

			if len(s.expectedErrors) > 0 {
				return
			}

			var expectedContent string
			switch s.template {
			case forms.TestTemplatePasswordReset:
				expectedContent = "Reset password"
			case forms.TestTemplateEmailChange:
				expectedContent = "Confirm new email"
			case forms.TestTemplateVerification:
				expectedContent = "Verify"
			case forms.TestTemplateOTP:
				expectedContent = "one-time password"
			case forms.TestTemplateAuthAlert:
				expectedContent = "from a new location"
			default:
				expectedContent = "__UNKNOWN_TEMPLATE__"
			}

			if !strings.Contains(app.TestMailer.LastMessage().HTML, expectedContent) {
				t.Errorf("Expected the email to contains %q, got\n%v", expectedContent, app.TestMailer.LastMessage().HTML)
			}
		})
	}
}
