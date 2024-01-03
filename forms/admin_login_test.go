package forms_test

import (
	"errors"
	"testing"

	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tests"
)

func TestAdminLoginValidateAndSubmit(t *testing.T) {
	t.Parallel()

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

func TestAdminLoginInterceptors(t *testing.T) {
	t.Parallel()

	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	form := forms.NewAdminLogin(testApp)
	form.Identity = "test@example.com"
	form.Password = "123456"
	var interceptorAdmin *models.Admin
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
			interceptorAdmin = admin
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

	if interceptorAdmin == nil || interceptorAdmin.Email != form.Identity {
		t.Fatalf("Expected Admin model with email %s, got %v", form.Identity, interceptorAdmin)
	}
}
