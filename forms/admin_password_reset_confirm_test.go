package forms_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/security"
)

func TestAdminPasswordResetConfirmValidateAndSubmit(t *testing.T) {
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
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsImV4cCI6MTY0MDk5MTY2MX0.GLwCOsgWTTEKXTK-AyGW838de1OeZGIjfHH0FoRLqZg",
			"1234567890",
			"1234567890",
			true,
		},
		{
			// valid with mismatched passwords
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsImV4cCI6MjIwODk4MTYwMH0.kwFEler6KSMKJNstuaSDvE1QnNdCta5qSnjaIQ0hhhc",
			"1234567890",
			"1234567891",
			true,
		},
		{
			// valid with matching passwords
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsImV4cCI6MjIwODk4MTYwMH0.kwFEler6KSMKJNstuaSDvE1QnNdCta5qSnjaIQ0hhhc",
			"1234567891",
			"1234567891",
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
			continue
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
