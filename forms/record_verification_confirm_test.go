package forms_test

import (
	"encoding/json"
	"testing"

	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/security"
)

func TestRecordVerificationConfirmValidateAndSubmit(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	authCollection, err := testApp.Dao().FindCollectionByNameOrId("users")
	if err != nil {
		t.Fatal(err)
	}

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
			`{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsImNvbGxlY3Rpb25JZCI6Il9wYl91c2Vyc19hdXRoXyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiZXhwIjoxNjQwOTkxNjYxfQ.Avbt9IP8sBisVz_2AGrlxLDvangVq4PhL2zqQVYLKlE"}`,
			true,
		},
		// valid token (already verified record)
		{
			`{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6Im9hcDY0MGNvdDR5cnUycyIsImVtYWlsIjoidGVzdDJAZXhhbXBsZS5jb20iLCJjb2xsZWN0aW9uSWQiOiJfcGJfdXNlcnNfYXV0aF8iLCJ0eXBlIjoiYXV0aFJlY29yZCIsImV4cCI6MjIwODk4NTI2MX0.PsOABmYUzGbd088g8iIBL4-pf7DUZm0W5Ju6lL5JVRg"}`,
			false,
		},
		// valid token (unverified record)
		{
			`{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsImNvbGxlY3Rpb25JZCI6Il9wYl91c2Vyc19hdXRoXyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiZXhwIjoyMjA4OTg1MjYxfQ.hL16TVmStHFdHLc4a860bRqJ3sFfzjv0_NRNzwsvsrc"}`,
			false,
		},
	}

	for i, s := range scenarios {
		form := forms.NewRecordVerificationConfirm(testApp, authCollection)

		// load data
		loadErr := json.Unmarshal([]byte(s.jsonData), form)
		if loadErr != nil {
			t.Errorf("(%d) Failed to load form data: %v", i, loadErr)
			continue
		}

		record, err := form.Submit()

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("(%d) Expected hasErr to be %v, got %v (%v)", i, s.expectError, hasErr, err)
		}

		if hasErr {
			continue
		}

		claims, _ := security.ParseUnverifiedJWT(form.Token)
		tokenRecordId := claims["id"]

		if record.Id != tokenRecordId {
			t.Errorf("(%d) Expected record.Id %q, got %q", i, tokenRecordId, record.Id)
		}

		if !record.Verified() {
			t.Errorf("(%d) Expected record.Verified() to be true, got false", i)
		}
	}
}
