package forms_test

import (
	"encoding/json"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/tests"
)

func TestUserEmailChangeRequestPanic1(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("The form did not panic")
		}
	}()

	forms.NewUserEmailChangeRequest(nil, nil)
}

func TestUserEmailChangeRequestPanic2(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	defer func() {
		if recover() == nil {
			t.Fatal("The form did not panic")
		}
	}()

	forms.NewUserEmailChangeRequest(testApp, nil)
}

func TestUserEmailChangeRequestValidateAndSubmit(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	user, err := testApp.Dao().FindUserByEmail("test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		jsonData       string
		expectedErrors []string
	}{
		// empty payload
		{"{}", []string{"newEmail"}},
		// empty data
		{
			`{"newEmail": ""}`,
			[]string{"newEmail"},
		},
		// invalid email
		{
			`{"newEmail": "invalid"}`,
			[]string{"newEmail"},
		},
		// existing email token
		{
			`{"newEmail": "test@example.com"}`,
			[]string{"newEmail"},
		},
		// valid new email
		{
			`{"newEmail": "test_new@example.com"}`,
			[]string{},
		},
	}

	for i, s := range scenarios {
		testApp.TestMailer.TotalSend = 0 // reset
		form := forms.NewUserEmailChangeRequest(testApp, user)

		// load data
		loadErr := json.Unmarshal([]byte(s.jsonData), form)
		if loadErr != nil {
			t.Errorf("(%d) Failed to load form data: %v", i, loadErr)
			continue
		}

		err := form.Submit()

		// parse errors
		errs, ok := err.(validation.Errors)
		if !ok && err != nil {
			t.Errorf("(%d) Failed to parse errors %v", i, err)
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

		expectedMails := 1
		if len(s.expectedErrors) > 0 {
			expectedMails = 0
		}
		if testApp.TestMailer.TotalSend != expectedMails {
			t.Errorf("(%d) Expected %d mail(s) to be sent, got %d", i, expectedMails, testApp.TestMailer.TotalSend)
		}
	}
}
