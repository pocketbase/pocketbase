package forms_test

import (
	"errors"
	"testing"

	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/security"
)

func TestAdminPasswordResetConfirmValidateAndSubmit(t *testing.T) {
	t.Parallel()

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

		interceptorCalls := 0
		interceptor := func(next forms.InterceptorNextFunc[*models.Admin]) forms.InterceptorNextFunc[*models.Admin] {
			return func(m *models.Admin) error {
				interceptorCalls++
				return next(m)
			}
		}

		admin, err := form.Submit(interceptor)

		// check interceptor calls
		expectInterceptorCalls := 1
		if s.expectError {
			expectInterceptorCalls = 0
		}
		if interceptorCalls != expectInterceptorCalls {
			t.Errorf("[%d] Expected interceptor to be called %d, got %d", i, expectInterceptorCalls, interceptorCalls)
		}

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

func TestAdminPasswordResetConfirmInterceptors(t *testing.T) {
	t.Parallel()

	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	admin, err := testApp.Dao().FindAdminByEmail("test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	form := forms.NewAdminPasswordResetConfirm(testApp)
	form.Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsImV4cCI6MjIwODk4MTYwMH0.kwFEler6KSMKJNstuaSDvE1QnNdCta5qSnjaIQ0hhhc"
	form.Password = "1234567891"
	form.PasswordConfirm = "1234567891"
	interceptorTokenKey := admin.TokenKey
	testErr := errors.New("test_error")

	interceptor1Called := false
	interceptor1 := func(next forms.InterceptorNextFunc[*models.Admin]) forms.InterceptorNextFunc[*models.Admin] {
		return func(admin *models.Admin) error {
			interceptor1Called = true
			return next(admin)
		}
	}

	interceptor2Called := false
	interceptor2 := func(next forms.InterceptorNextFunc[*models.Admin]) forms.InterceptorNextFunc[*models.Admin] {
		return func(admin *models.Admin) error {
			interceptorTokenKey = admin.TokenKey
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

	if interceptorTokenKey == admin.TokenKey {
		t.Fatalf("Expected the form model to be filled before calling the interceptors")
	}
}
