package forms_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/tests"
)

func TestAdminLoginPanic(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("The form did not panic")
		}
	}()

	forms.NewAdminLogin(nil)
}

func TestAdminLoginValidate(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	form := forms.NewAdminLogin(app)

	scenarios := []struct {
		email       string
		password    string
		expectError bool
	}{
		{"", "", true},
		{"", "123", true},
		{"test@example.com", "", true},
		{"test", "123", true},
		{"test@example.com", "123", false},
	}

	for i, s := range scenarios {
		form.Email = s.email
		form.Password = s.password

		err := form.Validate()

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("(%d) Expected hasErr to be %v, got %v (%v)", i, s.expectError, hasErr, err)
		}
	}
}

func TestAdminLoginSubmit(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	form := forms.NewAdminLogin(app)

	scenarios := []struct {
		email       string
		password    string
		expectError bool
	}{
		{"", "", true},
		{"", "1234567890", true},
		{"test@example.com", "", true},
		{"test", "1234567890", true},
		{"missing@example.com", "1234567890", true},
		{"test@example.com", "123456789", true},
		{"test@example.com", "1234567890", false},
	}

	for i, s := range scenarios {
		form.Email = s.email
		form.Password = s.password

		admin, err := form.Submit()

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("(%d) Expected hasErr to be %v, got %v (%v)", i, s.expectError, hasErr, err)
		}

		if !s.expectError && admin == nil {
			t.Errorf("(%d) Expected admin model to be returned, got nil", i)
		}

		if admin != nil && admin.Email != s.email {
			t.Errorf("(%d) Expected admin with email %s to be returned, got %v", i, s.email, admin)
		}
	}
}
