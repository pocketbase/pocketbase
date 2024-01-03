package forms_test

import (
	"encoding/json"
	"errors"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/security"
)

func TestRecordEmailChangeConfirmValidateAndSubmit(t *testing.T) {
	t.Parallel()

	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	authCollection, err := testApp.Dao().FindCollectionByNameOrId("users")
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		jsonData       string
		expectedErrors []string
	}{
		// empty payload
		{"{}", []string{"token", "password"}},
		// empty data
		{
			`{"token": "", "password": ""}`,
			[]string{"token", "password"},
		},
		// invalid token payload
		{
			`{
				"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.quDgaCi2rGTRx3qO06CrFvHdeCua_5J7CCVWSaFhkus",
				"password": "123456"
			}`,
			[]string{"token", "password"},
		},
		// expired token
		{
			`{
				"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwibmV3RW1haWwiOiJ0ZXN0X25ld0BleGFtcGxlLmNvbSIsImV4cCI6MTYwOTQ1NTY2MX0.n1OJXJEACMNPT9aMTO48cVJexIiZEtHsz4UNBIfMcf4",
				"password": "123456"
			}`,
			[]string{"token", "password"},
		},
		// existing new email
		{
			`{
				"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwibmV3RW1haWwiOiJ0ZXN0MkBleGFtcGxlLmNvbSIsImV4cCI6MjIwODk4NTI2MX0.Q_o6zpc2URggTU0mWv2CS0rIPbQhFdmrjZ-ASwHh1Ww",
				"password": "1234567890"
			}`,
			[]string{"token", "password"},
		},
		// wrong confirmation password
		{
			`{
				"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwibmV3RW1haWwiOiJ0ZXN0X25ld0BleGFtcGxlLmNvbSIsImV4cCI6MjIwODk4NTI2MX0.hmR7Ye23C68tS1LgHgYgT7NBJczTad34kzcT4sqW3FY",
				"password": "123456"
			}`,
			[]string{"password"},
		},
		// valid data
		{
			`{
				"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwibmV3RW1haWwiOiJ0ZXN0X25ld0BleGFtcGxlLmNvbSIsImV4cCI6MjIwODk4NTI2MX0.hmR7Ye23C68tS1LgHgYgT7NBJczTad34kzcT4sqW3FY",
				"password": "1234567890"
			}`,
			[]string{},
		},
	}

	for i, s := range scenarios {
		form := forms.NewRecordEmailChangeConfirm(testApp, authCollection)

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

		record, err := form.Submit(interceptor)

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

		if len(errs) > 0 {
			continue
		}

		claims, _ := security.ParseUnverifiedJWT(form.Token)
		newEmail, _ := claims["newEmail"].(string)

		// check whether the user was updated
		// ---
		if record.Email() != newEmail {
			t.Errorf("(%d) Expected record email %q, got %q", i, newEmail, record.Email())
		}

		if !record.Verified() {
			t.Errorf("(%d) Expected record to be verified, got false", i)
		}

		// shouldn't validate second time due to refreshed record token
		if err := form.Validate(); err == nil {
			t.Errorf("(%d) Expected error, got nil", i)
		}
	}
}

func TestRecordEmailChangeConfirmInterceptors(t *testing.T) {
	t.Parallel()

	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	authCollection, err := testApp.Dao().FindCollectionByNameOrId("users")
	if err != nil {
		t.Fatal(err)
	}

	authRecord, err := testApp.Dao().FindAuthRecordByEmail("users", "test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	form := forms.NewRecordEmailChangeConfirm(testApp, authCollection)
	form.Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwibmV3RW1haWwiOiJ0ZXN0X25ld0BleGFtcGxlLmNvbSIsImV4cCI6MjIwODk4NTI2MX0.hmR7Ye23C68tS1LgHgYgT7NBJczTad34kzcT4sqW3FY"
	form.Password = "1234567890"
	interceptorEmail := authRecord.Email()
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
			interceptorEmail = record.Email()
			interceptor2Called = true
			return testErr
		}
	}

	_, submitErr := form.Submit(interceptor1, interceptor2)
	if submitErr != testErr {
		t.Fatalf("Expected submitError %v, got %v", testErr, submitErr)
	}

	if !interceptor1Called {
		t.Fatalf("Expected interceptor1 to be called")
	}

	if !interceptor2Called {
		t.Fatalf("Expected interceptor2 to be called")
	}

	if interceptorEmail == authRecord.Email() {
		t.Fatalf("Expected the form model to be filled before calling the interceptors")
	}
}
