package forms_test

import (
	"strings"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/tests"
)

func TestEmailSendValidate(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

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
		form := forms.NewTestEmailSend(app)
		form.Email = s.email
		form.Template = s.template

		result := form.Validate()

		// parse errors
		errs, ok := result.(validation.Errors)
		if !ok && result != nil {
			t.Errorf("(%d) Failed to parse errors %v", i, result)
			continue
		}

		// check errors
		if len(errs) > len(s.expectedErrors) {
			t.Errorf("(%d) Expected error keys %v, got %v", i, s.expectedErrors, errs)
		}
		for _, k := range s.expectedErrors {
			if _, ok := errs[k]; !ok {
				t.Errorf("(%d) Missing expected error key %q in %v", i, k, errs)
			}
		}
	}
}

func TestEmailSendSubmit(t *testing.T) {
	scenarios := []struct {
		template    string
		email       string
		expectError bool
	}{
		{"", "", true},
		{"invalid", "test@example.com", true},
		{"verification", "invalid", true},
		{"verification", "test@example.com", false},
		{"password-reset", "test@example.com", false},
		{"email-change", "test@example.com", false},
	}

	for i, s := range scenarios {
		app, _ := tests.NewTestApp()
		defer app.Cleanup()

		form := forms.NewTestEmailSend(app)
		form.Email = s.email
		form.Template = s.template

		err := form.Submit()

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("(%d) Expected hasErr to be %v, got %v (%v)", i, s.expectError, hasErr, err)
		}

		if hasErr {
			continue
		}

		if app.TestMailer.TotalSend != 1 {
			t.Errorf("(%d) Expected one email to be sent, got %d", i, app.TestMailer.TotalSend)
		}

		expectedContent := "Verify"
		if s.template == "password-reset" {
			expectedContent = "Reset password"
		} else if s.template == "email-change" {
			expectedContent = "Confirm new email"
		}

		if !strings.Contains(app.TestMailer.LastHtmlBody, expectedContent) {
			t.Errorf("(%d) Expected the email to contains %s, got \n%v", i, expectedContent, app.TestMailer.LastHtmlBody)
		}
	}
}
