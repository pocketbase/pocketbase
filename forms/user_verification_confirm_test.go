package forms_test

import (
	"encoding/json"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/security"
)

func TestUserVerificationConfirmPanic(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("The form did not panic")
		}
	}()

	forms.NewUserVerificationConfirm(nil)
}

func TestUserVerificationConfirmValidate(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		jsonData       string
		expectedErrors []string
	}{
		// empty data
		{
			`{}`,
			[]string{"token"},
		},
		// empty fields
		{
			`{"token":""}`,
			[]string{"token"},
		},
		// invalid JWT token
		{
			`{"token":"invalid"}`,
			[]string{"token"},
		},
		// expired token
		{
			`{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwiZXhwIjoxNjQwOTkxNjYxfQ.6KBn19eFa9aFAZ6hvuhQtK7Ovxb6QlBQ97vJtulb_P8"}`,
			[]string{"token"},
		},
		// valid token
		{
			`{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwiZXhwIjoxOTA2MTA2NDIxfQ.yvH96FwtPHGvzhFSKl8Tsi1FnGytKpMrvb7K9F2_zQA"}`,
			[]string{},
		},
	}

	for i, s := range scenarios {
		form := forms.NewUserVerificationConfirm(app)

		// load data
		loadErr := json.Unmarshal([]byte(s.jsonData), form)
		if loadErr != nil {
			t.Errorf("(%d) Failed to load form data: %v", i, loadErr)
			continue
		}

		// parse errors
		result := form.Validate()
		errs, ok := result.(validation.Errors)
		if !ok && result != nil {
			t.Errorf("(%d) Failed to parse errors %v", i, result)
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

func TestUserVerificationConfirmSubmit(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		jsonData    string
		expectError bool
	}{
		// empty data (Validate call check)
		{
			`{}`,
			true,
		},
		// expired token (Validate call check)
		{
			`{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwiZXhwIjoxNjQwOTkxNjYxfQ.6KBn19eFa9aFAZ6hvuhQtK7Ovxb6QlBQ97vJtulb_P8"}`,
			true,
		},
		// valid token (already verified user)
		{
			`{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwiZXhwIjoxOTA2MTA2NDIxfQ.yvH96FwtPHGvzhFSKl8Tsi1FnGytKpMrvb7K9F2_zQA"}`,
			false,
		},
		// valid token (unverified user)
		{
			`{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjdiYzg0ZDI3LTZiYTItYjQyYS0zODNmLTQxOTdjYzNkM2QwYyIsInR5cGUiOiJ1c2VyIiwiZW1haWwiOiJ0ZXN0MkBleGFtcGxlLmNvbSIsImV4cCI6MTkwNjEwNjQyMX0.KbSucLGasQqTkGxUgqaaCjKNOHJ3ZVkL1WTzSApc6oM"}`,
			false,
		},
	}

	for i, s := range scenarios {
		form := forms.NewUserVerificationConfirm(app)

		// load data
		loadErr := json.Unmarshal([]byte(s.jsonData), form)
		if loadErr != nil {
			t.Errorf("(%d) Failed to load form data: %v", i, loadErr)
			continue
		}

		user, err := form.Submit()

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("(%d) Expected hasErr to be %v, got %v (%v)", i, s.expectError, hasErr, err)
		}

		if s.expectError {
			continue
		}

		claims, _ := security.ParseUnverifiedJWT(form.Token)
		tokenUserId := claims["id"]

		if user.Id != tokenUserId {
			t.Errorf("(%d) Expected user.Id %q, got %q", i, tokenUserId, user.Id)
		}

		if !user.Verified {
			t.Errorf("(%d) Expected user.Verified to be true, got false", i)
		}
	}
}
