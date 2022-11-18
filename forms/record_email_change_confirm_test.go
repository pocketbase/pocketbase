package forms_test

import (
	"encoding/json"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/security"
)

func TestRecordEmailChangeConfirmValidateAndSubmit(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	authCollection, err := testApp.Dao().FindCollectionByNameOrId("users")
	if err != nil {
		t.Fatal(err)
	}

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
				"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.quDgaCi2rGTRx3qO06CrFvHdeCua_5J7CCVWSaFhkus",
				"password": "123456"
			}`,
			[]string{"token", "password"},
		},
		// expired token
		{
			`{
				"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwibmV3RW1haWwiOiJ0ZXN0X25ld0BleGFtcGxlLmNvbSIsImV4cCI6MTYwOTQ1NTY2MX0.n1OJXJEACMNPT9aMTO48cVJexIiZEtHsz4UNBIfMcf4",
				"password": "123456"
			}`,
			[]string{"token", "password"},
		},
		// existing new email
		{
			`{
				"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwibmV3RW1haWwiOiJ0ZXN0MkBleGFtcGxlLmNvbSIsImV4cCI6MjIwODk4NTI2MX0.Q_o6zpc2URggTU0mWv2CS0rIPbQhFdmrjZ-ASwHh1Ww",
				"password": "1234567890"
			}`,
			[]string{"token", "password"},
		},
		// wrong confirmation password
		{
			`{
				"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwibmV3RW1haWwiOiJ0ZXN0X25ld0BleGFtcGxlLmNvbSIsImV4cCI6MjIwODk4NTI2MX0.hmR7Ye23C68tS1LgHgYgT7NBJczTad34kzcT4sqW3FY",
				"password": "123456"
			}`,
			[]string{"password"},
		},
		// valid data
		{
			`{
				"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwibmV3RW1haWwiOiJ0ZXN0X25ld0BleGFtcGxlLmNvbSIsImV4cCI6MjIwODk4NTI2MX0.hmR7Ye23C68tS1LgHgYgT7NBJczTad34kzcT4sqW3FY",
				"password": "1234567890"
			}`,
			[]string{},
		},
	}

	for i, s := range scenarios {
		form := forms.NewRecordEmailChangeConfirm(testApp, authCollection)

		// load data
		loadErr := json.Unmarshal([]byte(s.jsonData), form)
		if loadErr != nil {
			t.Errorf("(%d) Failed to load form data: %v", i, loadErr)
			continue
		}

		record, err := form.Submit()

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

		if len(errs) > 0 {
			continue
		}

		claims, _ := security.ParseUnverifiedJWT(form.Token)
		newEmail, _ := claims["newEmail"].(string)

		// check whether the user was updated
		// ---
		if record.Email() != newEmail {
			t.Errorf("(%d) Expected record email %q, got %q", i, newEmail, record.Email())
		}

		if !record.Verified() {
			t.Errorf("(%d) Expected record to be verified, got false", i)
		}

		// shouldn't validate second time due to refreshed record token
		if err := form.Validate(); err == nil {
			t.Errorf("(%d) Expected error, got nil", i)
		}
	}
}
