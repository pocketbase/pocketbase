package forms_test

import (
	"encoding/json"
	"errors"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tests"
)

func TestRecordEmailChangeRequestValidateAndSubmit(t *testing.T) {
	t.Parallel()

	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	user, err := testApp.Dao().FindAuthRecordByEmail("users", "test@example.com")
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
			`{"newEmail": "test2@example.com"}`,
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
		form := forms.NewRecordEmailChangeRequest(testApp, user)

		// load data
		loadErr := json.Unmarshal([]byte(s.jsonData), form)
		if loadErr != nil {
			t.Errorf("(%d) Failed to load form data: %v", i, loadErr)
			continue
		}

		interceptorCalls := 0
		interceptor := func(next forms.InterceptorNextFunc[*models.Record]) forms.InterceptorNextFunc[*models.Record] {
			return func(r *models.Record) error {
				interceptorCalls++
				return next(r)
			}
		}

		err := form.Submit(interceptor)

		// check interceptor calls
		expectInterceptorCalls := 1
		if len(s.expectedErrors) > 0 {
			expectInterceptorCalls = 0
		}
		if interceptorCalls != expectInterceptorCalls {
			t.Errorf("[%d] Expected interceptor to be called %d, got %d", i, expectInterceptorCalls, interceptorCalls)
		}

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

func TestRecordEmailChangeRequestInterceptors(t *testing.T) {
	t.Parallel()

	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	authRecord, err := testApp.Dao().FindAuthRecordByEmail("users", "test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	form := forms.NewRecordEmailChangeRequest(testApp, authRecord)
	form.NewEmail = "test_new@example.com"
	testErr := errors.New("test_error")

	interceptor1Called := false
	interceptor1 := func(next forms.InterceptorNextFunc[*models.Record]) forms.InterceptorNextFunc[*models.Record] {
		return func(record *models.Record) error {
			interceptor1Called = true
			return next(record)
		}
	}

	interceptor2Called := false
	interceptor2 := func(next forms.InterceptorNextFunc[*models.Record]) forms.InterceptorNextFunc[*models.Record] {
		return func(record *models.Record) error {
			interceptor2Called = true
			return testErr
		}
	}

	submitErr := form.Submit(interceptor1, interceptor2)
	if submitErr != testErr {
		t.Fatalf("Expected submitError %v, got %v", testErr, submitErr)
	}

	if !interceptor1Called {
		t.Fatalf("Expected interceptor1 to be called")
	}

	if !interceptor2Called {
		t.Fatalf("Expected interceptor2 to be called")
	}
}
