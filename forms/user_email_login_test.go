package forms_test

import (
	"encoding/json"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/tests"
)

func TestUserEmailLoginPanic(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("The form did not panic")
		}
	}()

	forms.NewUserEmailLogin(nil)
}

func TestUserEmailLoginValidate(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		jsonData       string
		expectedErrors []string
	}{
		// empty payload
		{"{}", []string{"email", "password"}},
		// empty data
		{
			`{"email": "","password": ""}`,
			[]string{"email", "password"},
		},
		// invalid email
		{
			`{"email": "invalid","password": "123"}`,
			[]string{"email"},
		},
		// valid email
		{
			`{"email": "test@example.com","password": "123"}`,
			[]string{},
		},
	}

	for i, s := range scenarios {
		form := forms.NewUserEmailLogin(app)

		// load data
		loadErr := json.Unmarshal([]byte(s.jsonData), form)
		if loadErr != nil {
			t.Errorf("(%d) Failed to load form data: %v", i, loadErr)
			continue
		}

		err := form.Validate()

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
	}
}

func TestUserEmailLoginSubmit(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		email       string
		password    string
		expectError bool
	}{
		// invalid email
		{"invalid", "123456", true},
		// missing user
		{"missing@example.com", "123456", true},
		// invalid password
		{"test@example.com", "123", true},
		// valid email and password
		{"test@example.com", "123456", false},
	}

	for i, s := range scenarios {
		form := forms.NewUserEmailLogin(app)
		form.Email = s.email
		form.Password = s.password

		user, err := form.Submit()

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("(%d) Expected hasErr to be %v, got %v (%v)", i, s.expectError, hasErr, err)
			continue
		}

		if !s.expectError && user.Email != s.email {
			t.Errorf("(%d) Expected user with email %q, got %q", i, s.email, user.Email)
		}
	}
}
