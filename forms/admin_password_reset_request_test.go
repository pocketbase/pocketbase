package forms_test

import (
	"errors"
	"testing"

	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tests"
)

func TestAdminPasswordResetRequestValidateAndSubmit(t *testing.T) {
	t.Parallel()

	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	form := forms.NewAdminPasswordResetRequest(testApp)

	scenarios := []struct {
		email       string
		expectError bool
	}{
		{"", true},
		{"", true},
		{"invalid", true},
		{"missing@example.com", true},
		{"test@example.com", false},
		{"test@example.com", true}, // already requested
	}

	for i, s := range scenarios {
		testApp.TestMailer.TotalSend = 0 // reset
		form.Email = s.email

		adminBefore, _ := testApp.Dao().FindAdminByEmail(s.email)

		interceptorCalls := 0
		interceptor := func(next forms.InterceptorNextFunc[*models.Admin]) forms.InterceptorNextFunc[*models.Admin] {
			return func(m *models.Admin) error {
				interceptorCalls++
				return next(m)
			}
		}

		err := form.Submit(interceptor)

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
		}

		adminAfter, _ := testApp.Dao().FindAdminByEmail(s.email)

		if !s.expectError && (adminBefore.LastResetSentAt == adminAfter.LastResetSentAt || adminAfter.LastResetSentAt.IsZero()) {
			t.Errorf("(%d) Expected admin.LastResetSentAt to change, got %q", i, adminAfter.LastResetSentAt)
		}

		expectedMails := 1
		if s.expectError {
			expectedMails = 0
		}
		if testApp.TestMailer.TotalSend != expectedMails {
			t.Errorf("(%d) Expected %d mail(s) to be sent, got %d", i, expectedMails, testApp.TestMailer.TotalSend)
		}
	}
}

func TestAdminPasswordResetRequestInterceptors(t *testing.T) {
	t.Parallel()

	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	admin, err := testApp.Dao().FindAdminByEmail("test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	form := forms.NewAdminPasswordResetRequest(testApp)
	form.Email = admin.Email
	interceptorLastResetSentAt := admin.LastResetSentAt
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
			interceptorLastResetSentAt = admin.LastResetSentAt
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

	if interceptorLastResetSentAt.String() != admin.LastResetSentAt.String() {
		t.Fatalf("Expected the form model to NOT be filled before calling the interceptors")
	}
}
