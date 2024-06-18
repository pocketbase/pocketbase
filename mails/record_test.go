package mails_test

import (
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/mails"
	"github.com/pocketbase/pocketbase/tests"
)

func TestSendRecordPasswordLoginAlert(t *testing.T) {
	t.Parallel()

	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	// ensure that action url normalization will be applied
	testApp.Settings().Meta.AppUrl = "http://localhost:8090////"

	user, _ := testApp.Dao().FindFirstRecordByData("users", "email", "test@example.com")

	err := mails.SendRecordPasswordLoginAlert(testApp, user, "test1", "test2")
	if err != nil {
		t.Fatal(err)
	}

	if testApp.TestMailer.TotalSend != 1 {
		t.Fatalf("Expected one email to be sent, got %d", testApp.TestMailer.TotalSend)
	}

	expectedParts := []string{"using a password", "OAuth2", "test1", "test2", "auth linked"}

	for _, part := range expectedParts {
		if !strings.Contains(testApp.TestMailer.LastMessage.HTML, part) {
			t.Fatalf("Couldn't find %s\n in\n %s", part, testApp.TestMailer.LastMessage.HTML)
		}
	}
}

func TestSendRecordPasswordReset(t *testing.T) {
	t.Parallel()

	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	// ensure that action url normalization will be applied
	testApp.Settings().Meta.AppUrl = "http://localhost:8090////"

	user, _ := testApp.Dao().FindFirstRecordByData("users", "email", "test@example.com")

	err := mails.SendRecordPasswordReset(testApp, user)
	if err != nil {
		t.Fatal(err)
	}

	if testApp.TestMailer.TotalSend != 1 {
		t.Fatalf("Expected one email to be sent, got %d", testApp.TestMailer.TotalSend)
	}

	expectedParts := []string{
		"http://localhost:8090/_/#/auth/confirm-password-reset/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.",
	}
	for _, part := range expectedParts {
		if !strings.Contains(testApp.TestMailer.LastMessage.HTML, part) {
			t.Fatalf("Couldn't find %s \nin\n %s", part, testApp.TestMailer.LastMessage.HTML)
		}
	}
}

func TestSendRecordVerification(t *testing.T) {
	t.Parallel()

	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	user, _ := testApp.Dao().FindFirstRecordByData("users", "email", "test@example.com")

	err := mails.SendRecordVerification(testApp, user)
	if err != nil {
		t.Fatal(err)
	}

	if testApp.TestMailer.TotalSend != 1 {
		t.Fatalf("Expected one email to be sent, got %d", testApp.TestMailer.TotalSend)
	}

	expectedParts := []string{
		"http://localhost:8090/_/#/auth/confirm-verification/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.",
	}
	for _, part := range expectedParts {
		if !strings.Contains(testApp.TestMailer.LastMessage.HTML, part) {
			t.Fatalf("Couldn't find %s \nin\n %s", part, testApp.TestMailer.LastMessage.HTML)
		}
	}
}

func TestSendRecordChangeEmail(t *testing.T) {
	t.Parallel()

	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	user, _ := testApp.Dao().FindFirstRecordByData("users", "email", "test@example.com")

	err := mails.SendRecordChangeEmail(testApp, user, "new_test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	if testApp.TestMailer.TotalSend != 1 {
		t.Fatalf("Expected one email to be sent, got %d", testApp.TestMailer.TotalSend)
	}

	expectedParts := []string{
		"http://localhost:8090/_/#/auth/confirm-email-change/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.",
	}
	for _, part := range expectedParts {
		if !strings.Contains(testApp.TestMailer.LastMessage.HTML, part) {
			t.Fatalf("Couldn't find %s \nin\n %s", part, testApp.TestMailer.LastMessage.HTML)
		}
	}
}
