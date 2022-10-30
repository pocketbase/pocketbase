package forms_test

import (
	"encoding/json"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/security"
)

func TestRecordPasswordResetConfirmValidateAndSubmit(t *testing.T) {
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
		// empty data (Validate call check)
		{
			`{}`,
			[]string{"token", "password", "passwordConfirm"},
		},
		// expired token
		{
			`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsImNvbGxlY3Rpb25JZCI6Il9wYl91c2Vyc19hdXRoXyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiZXhwIjoxNjQwOTkxNjYxfQ.TayHoXkOTM0w8InkBEb86npMJEaf6YVUrxrRmMgFjeY",
				"password":"12345678",
				"passwordConfirm":"12345678"
			}`,
			[]string{"token"},
		},
		// valid token but invalid passwords lengths
		{
			`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsImNvbGxlY3Rpb25JZCI6Il9wYl91c2Vyc19hdXRoXyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiZXhwIjoyMjA4OTg1MjYxfQ.R_4FOSUHIuJQ5Crl3PpIPCXMsoHzuTaNlccpXg_3FOg",
				"password":"1234567",
				"passwordConfirm":"1234567"
			}`,
			[]string{"password"},
		},
		// valid token but mismatched passwordConfirm
		{
			`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsImNvbGxlY3Rpb25JZCI6Il9wYl91c2Vyc19hdXRoXyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiZXhwIjoyMjA4OTg1MjYxfQ.R_4FOSUHIuJQ5Crl3PpIPCXMsoHzuTaNlccpXg_3FOg",
				"password":"12345678",
				"passwordConfirm":"12345679"
			}`,
			[]string{"passwordConfirm"},
		},
		// valid token and password
		{
			`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsImNvbGxlY3Rpb25JZCI6Il9wYl91c2Vyc19hdXRoXyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiZXhwIjoyMjA4OTg1MjYxfQ.R_4FOSUHIuJQ5Crl3PpIPCXMsoHzuTaNlccpXg_3FOg",
				"password":"12345678",
				"passwordConfirm":"12345678"
			}`,
			[]string{},
		},
	}

	for i, s := range scenarios {
		form := forms.NewRecordPasswordResetConfirm(testApp, authCollection)

		// load data
		loadErr := json.Unmarshal([]byte(s.jsonData), form)
		if loadErr != nil {
			t.Errorf("(%d) Failed to load form data: %v", i, loadErr)
			continue
		}

		record, submitErr := form.Submit()

		// parse errors
		errs, ok := submitErr.(validation.Errors)
		if !ok && submitErr != nil {
			t.Errorf("(%d) Failed to parse errors %v", i, submitErr)
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

		if len(errs) > 0 || len(s.expectedErrors) > 0 {
			continue
		}

		claims, _ := security.ParseUnverifiedJWT(form.Token)
		tokenRecordId := claims["id"]

		if record.Id != tokenRecordId {
			t.Errorf("(%d) Expected record with id %s, got %v", i, tokenRecordId, record)
		}

		if !record.LastResetSentAt().IsZero() {
			t.Errorf("(%d) Expected record.LastResetSentAt to be empty, got %v", i, record.LastResetSentAt())
		}

		if !record.ValidatePassword(form.Password) {
			t.Errorf("(%d) Expected the record password to have been updated to %q", i, form.Password)
		}
	}
}
