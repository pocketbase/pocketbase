package forms_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/security"
)

func TestAdminPasswordResetPanic(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("The form did not panic")
		}
	}()

	forms.NewAdminPasswordResetConfirm(nil)
}

func TestAdminPasswordResetConfirmValidate(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	form := forms.NewAdminPasswordResetConfirm(app)

	scenarios := []struct {
		token           string
		password        string
		passwordConfirm string
		expectError     bool
	}{
		{"", "", "", true},
		{"", "123", "", true},
		{"", "", "123", true},
		{"test", "", "", true},
		{"test", "123", "", true},
		{"test", "123", "123", true},
		{
			// expired
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTY0MTAxMzIwMH0.Gp_1b5WVhqjj2o3nJhNUlJmpdiwFLXN72LbMP-26gjA",
			"1234567890",
			"1234567890",
			true,
		},
		{
			// valid
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg5MzQ3NDAwMH0.72IhlL_5CpNGE0ZKM7sV9aAKa3wxQaMZdDiHBo0orpw",
			"1234567890",
			"1234567890",
			false,
		},
	}

	for i, s := range scenarios {
		form.Token = s.token
		form.Password = s.password
		form.PasswordConfirm = s.passwordConfirm

		err := form.Validate()

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("(%d) Expected hasErr to be %v, got %v (%v)", i, s.expectError, hasErr, err)
		}
	}
}

func TestAdminPasswordResetConfirmSubmit(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	form := forms.NewAdminPasswordResetConfirm(app)

	scenarios := []struct {
		token           string
		password        string
		passwordConfirm string
		expectError     bool
	}{
		{"", "", "", true},
		{"", "123", "", true},
		{"", "", "123", true},
		{"test", "", "", true},
		{"test", "123", "", true},
		{"test", "123", "123", true},
		{
			// expired
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTY0MTAxMzIwMH0.Gp_1b5WVhqjj2o3nJhNUlJmpdiwFLXN72LbMP-26gjA",
			"1234567890",
			"1234567890",
			true,
		},
		{
			// valid
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg5MzQ3NDAwMH0.72IhlL_5CpNGE0ZKM7sV9aAKa3wxQaMZdDiHBo0orpw",
			"1234567890",
			"1234567890",
			false,
		},
	}

	for i, s := range scenarios {
		form.Token = s.token
		form.Password = s.password
		form.PasswordConfirm = s.passwordConfirm

		admin, err := form.Submit()

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("(%d) Expected hasErr to be %v, got %v (%v)", i, s.expectError, hasErr, err)
		}

		if s.expectError {
			continue
		}

		claims, _ := security.ParseUnverifiedJWT(s.token)
		tokenAdminId := claims["id"]

		if admin.Id != tokenAdminId {
			t.Errorf("(%d) Expected admin with id %s to be returned, got %v", i, tokenAdminId, admin)
		}

		if !admin.ValidatePassword(form.Password) {
			t.Errorf("(%d) Expected the admin password to have been updated to %q", i, form.Password)
		}
	}
}
