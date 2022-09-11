package forms_test

import (
	"encoding/json"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/security"
)

func TestUserPasswordResetConfirmPanic(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("The form did not panic")
		}
	}()

	forms.NewUserPasswordResetConfirm(nil)
}

func TestUserPasswordResetConfirmValidate(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		jsonData       string
		expectedErrors []string
	}{
		// empty data
		{
			`{}`,
			[]string{"token", "password", "passwordConfirm"},
		},
		// empty fields
		{
			`{"token":"","password":"","passwordConfirm":""}`,
			[]string{"token", "password", "passwordConfirm"},
		},
		// invalid password length
		{
			`{"token":"invalid","password":"1234","passwordConfirm":"1234"}`,
			[]string{"token", "password"},
		},
		// mismatched passwords
		{
			`{"token":"invalid","password":"12345678","passwordConfirm":"87654321"}`,
			[]string{"token", "passwordConfirm"},
		},
		// invalid JWT token
		{
			`{"token":"invalid","password":"12345678","passwordConfirm":"12345678"}`,
			[]string{"token"},
		},
		// expired token
		{
			`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwiZXhwIjoxNjQwOTkxNjYxfQ.cSUFKWLAKEvulWV4fqPD6RRtkZYoyat_Tb8lrA2xqtw",
				"password":"12345678",
				"passwordConfirm":"12345678"
			}`,
			[]string{"token"},
		},
		// valid data
		{
			`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwiZXhwIjoxODkzNDUyNDYxfQ.YfpL4VOdsYh2gS30VIiPShgwwqPgt2CySD8TuuB1XD4",
				"password":"12345678",
				"passwordConfirm":"12345678"
			}`,
			[]string{},
		},
	}

	for i, s := range scenarios {
		form := forms.NewUserPasswordResetConfirm(app)

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

func TestUserPasswordResetConfirmSubmit(t *testing.T) {
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
		// expired token
		{
			`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwiZXhwIjoxNjQwOTkxNjYxfQ.cSUFKWLAKEvulWV4fqPD6RRtkZYoyat_Tb8lrA2xqtw",
				"password":"12345678",
				"passwordConfirm":"12345678"
			}`,
			true,
		},
		// valid data
		{
			`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwiZXhwIjoxODkzNDUyNDYxfQ.YfpL4VOdsYh2gS30VIiPShgwwqPgt2CySD8TuuB1XD4",
				"password":"12345678",
				"passwordConfirm":"12345678"
			}`,
			false,
		},
	}

	for i, s := range scenarios {
		form := forms.NewUserPasswordResetConfirm(app)

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
			t.Errorf("(%d) Expected user with id %s, got %v", i, tokenUserId, user)
		}

		if !user.LastResetSentAt.IsZero() {
			t.Errorf("(%d) Expected user.LastResetSentAt to be empty, got %v", i, user.LastResetSentAt)
		}

		if !user.ValidatePassword(form.Password) {
			t.Errorf("(%d) Expected the user password to have been updated to %q", i, form.Password)
		}
	}
}
