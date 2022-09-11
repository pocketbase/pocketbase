package forms_test

import (
	"encoding/json"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/security"
)

func TestUserEmailChangeConfirmPanic(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("The form did not panic")
		}
	}()

	forms.NewUserEmailChangeConfirm(nil)
}

func TestUserEmailChangeConfirmValidateAndSubmit(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		jsonData       string
		expectedErrors []string
	}{
		// empty payload
		{"{}", []string{"token", "password"}},
		// empty data
		{
			`{"token": "", "password": ""}`,
			[]string{"token", "password"},
		},
		// invalid token payload
		{
			`{
				"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwiZXhwIjoxODYxOTE2NDYxfQ.VjT3wc3IES--1Vye-1KRuk8RpO5mfdhVp2aKGbNluZ0",
				"password": "123456"
			}`,
			[]string{"token", "password"},
		},
		// expired token
		{
			`{
				"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwibmV3RW1haWwiOiJ0ZXN0X25ld0BleGFtcGxlLmNvbSIsImV4cCI6MTY0MDk5MTY2MX0.oPxbpJjcBpdZVBFbIW35FEXTCMkzJ7-RmQdHrz7zP3s",
				"password": "123456"
			}`,
			[]string{"token", "password"},
		},
		// existing new email
		{
			`{
				"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwibmV3RW1haWwiOiJ0ZXN0MkBleGFtcGxlLmNvbSIsImV4cCI6MTg2MTkxNjQ2MX0.RwHRZma5YpCwxHdj3y2obeBNy_GQrG6lT9CQHIUz6Ys",
				"password": "123456"
			}`,
			[]string{"token", "password"},
		},
		// wrong confirmation password
		{
			`{
				"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwibmV3RW1haWwiOiJ0ZXN0X25ld0BleGFtcGxlLmNvbSIsImV4cCI6MTg2MTkxNjQ2MX0.nS2qDonX25tOf9-6bKCwJXOm1CE88z_EVAA2B72NYM0",
				"password": "1234"
			}`,
			[]string{"password"},
		},
		// valid data
		{
			`{
				"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwibmV3RW1haWwiOiJ0ZXN0X25ld0BleGFtcGxlLmNvbSIsImV4cCI6MTg2MTkxNjQ2MX0.nS2qDonX25tOf9-6bKCwJXOm1CE88z_EVAA2B72NYM0",
				"password": "123456"
			}`,
			[]string{},
		},
	}

	for i, s := range scenarios {
		form := forms.NewUserEmailChangeConfirm(app)

		// load data
		loadErr := json.Unmarshal([]byte(s.jsonData), form)
		if loadErr != nil {
			t.Errorf("(%d) Failed to load form data: %v", i, loadErr)
			continue
		}

		user, err := form.Submit()

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

		if len(s.expectedErrors) > 0 {
			continue
		}

		claims, _ := security.ParseUnverifiedJWT(form.Token)
		newEmail, _ := claims["newEmail"].(string)

		// check whether the user was updated
		// ---
		if user.Email != newEmail {
			t.Errorf("(%d) Expected user email %q, got %q", i, newEmail, user.Email)
		}

		if !user.Verified {
			t.Errorf("(%d) Expected user to be verified, got false", i)
		}

		// shouldn't validate second time due to refreshed user token
		if err := form.Validate(); err == nil {
			t.Errorf("(%d) Expected error, got nil", i)
		}
	}
}
