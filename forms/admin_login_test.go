package forms_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/tests"
)

func TestAdminLoginValidateAndSubmit(t *testing.T) {
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
		{"test", "test", true},
		{"missing@example.com", "1234567890", true},
		{"test@example.com", "123456789", true},
		{"test@example.com", "1234567890", false},
	}

	for i, s := range scenarios {
		form.Identity = s.email
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
