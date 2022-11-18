package forms_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestRecordVerificationRequestSubmit(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	authCollection, err := testApp.Dao().FindCollectionByNameOrId("clients")
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		jsonData    string
		expectError bool
		expectMail  bool
	}{
		// empty field (Validate call check)
		{
			`{"email":""}`,
			true,
			false,
		},
		// invalid email field (Validate call check)
		{
			`{"email":"invalid"}`,
			true,
			false,
		},
		// nonexisting user
		{
			`{"email":"missing@example.com"}`,
			true,
			false,
		},
		// existing user (already verified)
		{
			`{"email":"test@example.com"}`,
			false,
			false,
		},
		// existing user (already verified) - repeating request to test threshod skip
		{
			`{"email":"test@example.com"}`,
			false,
			false,
		},
		// existing user (unverified)
		{
			`{"email":"test2@example.com"}`,
			false,
			true,
		},
		// existing user (inverified) - reached send threshod
		{
			`{"email":"test2@example.com"}`,
			true,
			false,
		},
	}

	now := types.NowDateTime()
	time.Sleep(1 * time.Millisecond)

	for i, s := range scenarios {
		testApp.TestMailer.TotalSend = 0 // reset
		form := forms.NewRecordVerificationRequest(testApp, authCollection)

		// load data
		loadErr := json.Unmarshal([]byte(s.jsonData), form)
		if loadErr != nil {
			t.Errorf("(%d) Failed to load form data: %v", i, loadErr)
			continue
		}

		err := form.Submit()

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("(%d) Expected hasErr to be %v, got %v (%v)", i, s.expectError, hasErr, err)
		}

		expectedMails := 0
		if s.expectMail {
			expectedMails = 1
		}
		if testApp.TestMailer.TotalSend != expectedMails {
			t.Errorf("(%d) Expected %d mail(s) to be sent, got %d", i, expectedMails, testApp.TestMailer.TotalSend)
		}

		if s.expectError {
			continue
		}

		user, err := testApp.Dao().FindAuthRecordByEmail(authCollection.Id, form.Email)
		if err != nil {
			t.Errorf("(%d) Expected user with email %q to exist, got nil", i, form.Email)
			continue
		}

		// check whether LastVerificationSentAt was updated
		if !user.Verified() && user.LastVerificationSentAt().Time().Sub(now.Time()) < 0 {
			t.Errorf("(%d) Expected LastVerificationSentAt to be after %v, got %v", i, now, user.LastVerificationSentAt())
		}
	}
}
