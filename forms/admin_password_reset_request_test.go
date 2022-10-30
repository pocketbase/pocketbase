package forms_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/tests"
)

func TestAdminPasswordResetRequestValidateAndSubmit(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	form := forms.NewAdminPasswordResetRequest(testApp)

	scenarios := []struct {
		email       string
		expectError bool
	}{
		{"", true},
		{"", true},
		{"invalid", true},
		{"missing@example.com", true},
		{"test@example.com", false},
		{"test@example.com", true}, // already requested
	}

	for i, s := range scenarios {
		testApp.TestMailer.TotalSend = 0 // reset
		form.Email = s.email

		adminBefore, _ := testApp.Dao().FindAdminByEmail(s.email)

		err := form.Submit()

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("(%d) Expected hasErr to be %v, got %v (%v)", i, s.expectError, hasErr, err)
		}

		adminAfter, _ := testApp.Dao().FindAdminByEmail(s.email)

		if !s.expectError && (adminBefore.LastResetSentAt == adminAfter.LastResetSentAt || adminAfter.LastResetSentAt.IsZero()) {
			t.Errorf("(%d) Expected admin.LastResetSentAt to change, got %q", i, adminAfter.LastResetSentAt)
		}

		expectedMails := 1
		if s.expectError {
			expectedMails = 0
		}
		if testApp.TestMailer.TotalSend != expectedMails {
			t.Errorf("(%d) Expected %d mail(s) to be sent, got %d", i, expectedMails, testApp.TestMailer.TotalSend)
		}
	}
}
