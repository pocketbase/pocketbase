package mails_test

import (
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/mails"
	"github.com/pocketbase/pocketbase/tests"
)

func TestSendAdminPasswordReset(t *testing.T) {
	t.Parallel()

	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	// ensure that action url normalization will be applied
	testApp.Settings().Meta.AppUrl = "http://localhost:8090////"

	admin, _ := testApp.Dao().FindAdminByEmail("test@example.com")

	err := mails.SendAdminPasswordReset(testApp, admin)
	if err != nil {
		t.Fatal(err)
	}

	if testApp.TestMailer.TotalSend != 1 {
		t.Fatalf("Expected one email to be sent, got %d", testApp.TestMailer.TotalSend)
	}

	expectedParts := []string{
		"http://localhost:8090/_/#/confirm-password-reset/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.",
	}
	for _, part := range expectedParts {
		if !strings.Contains(testApp.TestMailer.LastMessage.HTML, part) {
			t.Fatalf("Couldn't find %s \nin\n %s", part, testApp.TestMailer.LastMessage.HTML)
		}
	}
}
