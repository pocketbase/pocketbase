package forms_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestRecordPasswordResetRequestSubmit(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	authCollection, err := testApp.Dao().FindCollectionByNameOrId("users")
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		jsonData    string
		expectError bool
	}{
		// empty field (Validate call check)
		{
			`{"email":""}`,
			true,
		},
		// invalid email field (Validate call check)
		{
			`{"email":"invalid"}`,
			true,
		},
		// nonexisting user
		{
			`{"email":"missing@example.com"}`,
			true,
		},
		// existing user
		{
			`{"email":"test@example.com"}`,
			false,
		},
		// existing user - reached send threshod
		{
			`{"email":"test@example.com"}`,
			true,
		},
	}

	now := types.NowDateTime()
	time.Sleep(1 * time.Millisecond)

	for i, s := range scenarios {
		testApp.TestMailer.TotalSend = 0 // reset
		form := forms.NewRecordPasswordResetRequest(testApp, authCollection)

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

		expectedMails := 1
		if s.expectError {
			expectedMails = 0
		}
		if testApp.TestMailer.TotalSend != expectedMails {
			t.Errorf("(%d) Expected %d mail(s) to be sent, got %d", i, expectedMails, testApp.TestMailer.TotalSend)
		}

		if s.expectError {
			continue
		}

		// check whether LastResetSentAt was updated
		user, err := testApp.Dao().FindAuthRecordByEmail(authCollection.Id, form.Email)
		if err != nil {
			t.Errorf("(%d) Expected user with email %q to exist, got nil", i, form.Email)
			continue
		}

		if user.LastResetSentAt().Time().Sub(now.Time()) < 0 {
			t.Errorf("(%d) Expected LastResetSentAt to be after %v, got %v", i, now, user.LastResetSentAt())
		}
	}
}
