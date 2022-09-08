package mails_test

import (
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/mails"
	"github.com/pocketbase/pocketbase/tests"
)

func TestSendCustomEmail(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	email := "maxmustermann@example.com"
	title := "Bye Bye Bye"
	data := struct {
		CustomName string
		CustomText string
	}{
		CustomName: "Mr Max Musterman",
		CustomText: "Time to say goodbye!",
	}
	template := `
	<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
	<html xmlns="http://www.w3.org/1999/xhtml">
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
		<meta name="viewport" content="width=device-width,initial-scale=1" />
		<style>
			body {margin:0;}
		</style>
	</head>
	<body>
		<h1>Hi {{.CustomName}}</h1>
		<p>{{.CustomText}}</p>
	</body>
	</html>
	`

	err := mails.SendCustomEmail(testApp, email, title, template, data)
	if err != nil {
		t.Fatal(err)
	}

	if testApp.TestMailer.TotalSend != 1 {
		t.Fatalf("Expected one email to be sent, got %d", testApp.TestMailer.TotalSend)
	}

	if testApp.TestMailer.LastToAddress.Address != email {
		t.Fatalf("Expected email to address to be %s, got %s", email, testApp.TestMailer.LastToAddress.Address)
	}

	if testApp.TestMailer.LastHtmlSubject != title {
		t.Fatalf("Expected email title to be %s, got %s", title, testApp.TestMailer.LastHtmlSubject)
	}

	expectedParts := []string{
		data.CustomName,
		data.CustomText,
		"body {margin:0;}",
	}
	for _, part := range expectedParts {
		if !strings.Contains(testApp.TestMailer.LastHtmlBody, part) {
			t.Fatalf("Couldn't find %s \nin\n %s", part, testApp.TestMailer.LastHtmlBody)
		}
	}
}
