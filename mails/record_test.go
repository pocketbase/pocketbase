package mails_test

import (
	"html"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/mails"
	"github.com/pocketbase/pocketbase/tests"
)

func TestSendRecordAuthAlert(t *testing.T) {
	t.Parallel()

	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	user, _ := testApp.FindFirstRecordByData("users", "email", "test@example.com")

	// to test that it is escaped
	user.Set("name", "<p>"+user.GetString("name")+"</p>")

	err := mails.SendRecordAuthAlert(testApp, user)
	if err != nil {
		t.Fatal(err)
	}

	if testApp.TestMailer.TotalSend() != 1 {
		t.Fatalf("Expected one email to be sent, got %d", testApp.TestMailer.TotalSend())
	}

	expectedParts := []string{
		html.EscapeString(user.GetString("name")) + "{RECORD:tokenKey}", // public and private record placeholder checks
		"login to your " + testApp.Settings().Meta.AppName + " account from a new location",
		"If this was you",
		"If this wasn't you",
	}
	for _, part := range expectedParts {
		if !strings.Contains(testApp.TestMailer.LastMessage().HTML, part) {
			t.Fatalf("Couldn't find %s \nin\n %s", part, testApp.TestMailer.LastMessage().HTML)
		}
	}
}

func TestSendRecordPasswordReset(t *testing.T) {
	t.Parallel()

	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	user, _ := testApp.FindFirstRecordByData("users", "email", "test@example.com")

	// to test that it is escaped
	user.Set("name", "<p>"+user.GetString("name")+"</p>")

	err := mails.SendRecordPasswordReset(testApp, user)
	if err != nil {
		t.Fatal(err)
	}

	if testApp.TestMailer.TotalSend() != 1 {
		t.Fatalf("Expected one email to be sent, got %d", testApp.TestMailer.TotalSend())
	}

	expectedParts := []string{
		html.EscapeString(user.GetString("name")) + "{RECORD:tokenKey}", // the record name as {RECORD:name}
		"http://localhost:8090/_/#/auth/confirm-password-reset/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.",
	}
	for _, part := range expectedParts {
		if !strings.Contains(testApp.TestMailer.LastMessage().HTML, part) {
			t.Fatalf("Couldn't find %s \nin\n %s", part, testApp.TestMailer.LastMessage().HTML)
		}
	}
}

func TestSendRecordVerification(t *testing.T) {
	t.Parallel()

	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	user, _ := testApp.FindFirstRecordByData("users", "email", "test@example.com")

	// to test that it is escaped
	user.Set("name", "<p>"+user.GetString("name")+"</p>")

	err := mails.SendRecordVerification(testApp, user)
	if err != nil {
		t.Fatal(err)
	}

	if testApp.TestMailer.TotalSend() != 1 {
		t.Fatalf("Expected one email to be sent, got %d", testApp.TestMailer.TotalSend())
	}

	expectedParts := []string{
		html.EscapeString(user.GetString("name")) + "{RECORD:tokenKey}", // the record name as {RECORD:name}
		"http://localhost:8090/_/#/auth/confirm-verification/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.",
	}
	for _, part := range expectedParts {
		if !strings.Contains(testApp.TestMailer.LastMessage().HTML, part) {
			t.Fatalf("Couldn't find %s \nin\n %s", part, testApp.TestMailer.LastMessage().HTML)
		}
	}
}

func TestSendRecordChangeEmail(t *testing.T) {
	t.Parallel()

	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	user, _ := testApp.FindFirstRecordByData("users", "email", "test@example.com")

	// to test that it is escaped
	user.Set("name", "<p>"+user.GetString("name")+"</p>")

	err := mails.SendRecordChangeEmail(testApp, user, "new_test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	if testApp.TestMailer.TotalSend() != 1 {
		t.Fatalf("Expected one email to be sent, got %d", testApp.TestMailer.TotalSend())
	}

	expectedParts := []string{
		html.EscapeString(user.GetString("name")) + "{RECORD:tokenKey}", // the record name as {RECORD:name}
		"http://localhost:8090/_/#/auth/confirm-email-change/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.",
	}
	for _, part := range expectedParts {
		if !strings.Contains(testApp.TestMailer.LastMessage().HTML, part) {
			t.Fatalf("Couldn't find %s \nin\n %s", part, testApp.TestMailer.LastMessage().HTML)
		}
	}
}

func TestSendRecordOTP(t *testing.T) {
	t.Parallel()

	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	user, _ := testApp.FindFirstRecordByData("users", "email", "test@example.com")

	// to test that it is escaped
	user.Set("name", "<p>"+user.GetString("name")+"</p>")

	err := mails.SendRecordOTP(testApp, user, "test_otp_id", "test_otp_code")
	if err != nil {
		t.Fatal(err)
	}

	if testApp.TestMailer.TotalSend() != 1 {
		t.Fatalf("Expected one email to be sent, got %d", testApp.TestMailer.TotalSend())
	}

	expectedParts := []string{
		html.EscapeString(user.GetString("name")) + "{RECORD:tokenKey}", // the record name as {RECORD:name}
		"one-time password",
		"test_otp_code",
	}
	for _, part := range expectedParts {
		if !strings.Contains(testApp.TestMailer.LastMessage().HTML, part) {
			t.Fatalf("Couldn't find %s \nin\n %s", part, testApp.TestMailer.LastMessage().HTML)
		}
	}
}
