package forms_test

import (
	"errors"
	"testing"

	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tests"
)

func TestRecordPasswordLoginValidateAndSubmit(t *testing.T) {
	t.Parallel()

	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	scenarios := []struct {
		testName       string
		collectionName string
		identity       string
		password       string
		expectError    bool
	}{
		{
			"empty data",
			"users",
			"",
			"",
			true,
		},

		// username
		{
			"existing username + wrong password",
			"users",
			"users75657",
			"invalid",
			true,
		},
		{
			"missing username + valid password",
			"users",
			"clients57772", // not in the "users" collection
			"1234567890",
			true,
		},
		{
			"existing username + valid password but in restricted username auth collection",
			"clients",
			"clients57772",
			"1234567890",
			true,
		},
		{
			"existing username + valid password but in restricted username and email auth collection",
			"nologin",
			"test_username",
			"1234567890",
			true,
		},
		{
			"existing username + valid password",
			"users",
			"users75657",
			"1234567890",
			false,
		},

		// email
		{
			"existing email + wrong password",
			"users",
			"test@example.com",
			"invalid",
			true,
		},
		{
			"missing email + valid password",
			"users",
			"test_missing@example.com",
			"1234567890",
			true,
		},
		{
			"existing username + valid password but in restricted username auth collection",
			"clients",
			"test@example.com",
			"1234567890",
			false,
		},
		{
			"existing username + valid password but in restricted username and email auth collection",
			"nologin",
			"test@example.com",
			"1234567890",
			true,
		},
		{
			"existing email + valid password",
			"users",
			"test@example.com",
			"1234567890",
			false,
		},
	}

	for _, s := range scenarios {
		authCollection, err := testApp.Dao().FindCollectionByNameOrId(s.collectionName)
		if err != nil {
			t.Errorf("[%s] Failed to fetch auth collection: %v", s.testName, err)
		}

		form := forms.NewRecordPasswordLogin(testApp, authCollection)
		form.Identity = s.identity
		form.Password = s.password

		record, err := form.Submit()

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("[%s] Expected hasErr to be %v, got %v (%v)", s.testName, s.expectError, hasErr, err)
			continue
		}

		if hasErr {
			continue
		}

		if record.Email() != s.identity && record.Username() != s.identity {
			t.Errorf("[%s] Expected record with identity %q, got \n%v", s.testName, s.identity, record)
		}
	}
}

func TestRecordPasswordLoginInterceptors(t *testing.T) {
	t.Parallel()

	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	authCollection, err := testApp.Dao().FindCollectionByNameOrId("users")
	if err != nil {
		t.Fatal(err)
	}

	form := forms.NewRecordPasswordLogin(testApp, authCollection)
	form.Identity = "test@example.com"
	form.Password = "123456"
	var interceptorRecord *models.Record
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
			interceptorRecord = record
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

	if interceptorRecord == nil || interceptorRecord.Email() != form.Identity {
		t.Fatalf("Expected auth Record model with email %s, got %v", form.Identity, interceptorRecord)
	}
}
