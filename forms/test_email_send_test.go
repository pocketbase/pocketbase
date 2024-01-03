package forms_test

import (
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
		expectedErrors []string
	}{
		{"", "", []string{"template", "email"}},
		{"invalid", "test@example.com", []string{"template"}},
		{"verification", "invalid", []string{"email"}},
		{"verification", "test@example.com", nil},
		{"password-reset", "test@example.com", nil},
		{"email-change", "test@example.com", nil},
	}

	for i, s := range scenarios {
		func() {
			app, _ := tests.NewTestApp()
			defer app.Cleanup()

			form := forms.NewTestEmailSend(app)
			form.Email = s.email
			form.Template = s.template

			result := form.Submit()

			// parse errors
			errs, ok := result.(validation.Errors)
			if !ok && result != nil {
				t.Errorf("(%d) Failed to parse errors %v", i, result)
				return
			}

			// check errors
			if len(errs) > len(s.expectedErrors) {
				t.Errorf("(%d) Expected error keys %v, got %v", i, s.expectedErrors, errs)
				return
			}
			for _, k := range s.expectedErrors {
				if _, ok := errs[k]; !ok {
					t.Errorf("(%d) Missing expected error key %q in %v", i, k, errs)
					return
				}
			}

			expectedEmails := 1
			if len(s.expectedErrors) > 0 {
				expectedEmails = 0
			}

			if app.TestMailer.TotalSend != expectedEmails {
				t.Errorf("(%d) Expected %d email(s) to be sent, got %d", i, expectedEmails, app.TestMailer.TotalSend)
			}

			if len(s.expectedErrors) > 0 {
				return
			}

			expectedContent := "Verify"
			if s.template == "password-reset" {
				expectedContent = "Reset password"
			} else if s.template == "email-change" {
				expectedContent = "Confirm new email"
			}

			if !strings.Contains(app.TestMailer.LastMessage.HTML, expectedContent) {
				t.Errorf("(%d) Expected the email to contains %s, got \n%v", i, expectedContent, app.TestMailer.LastMessage.HTML)
			}
		}()
	}
}
