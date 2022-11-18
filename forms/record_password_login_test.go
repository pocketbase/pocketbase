package forms_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/tests"
)

func TestRecordEmailLoginValidateAndSubmit(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	scenarios := []struct {
		testName       string
		collectionName string
		identity       string
		password       string
		expectError    bool
	}{
		{
			"empty data",
			"users",
			"",
			"",
			true,
		},

		// username
		{
			"existing username + wrong password",
			"users",
			"users75657",
			"invalid",
			true,
		},
		{
			"missing username + valid password",
			"users",
			"clients57772", // not in the "users" collection
			"1234567890",
			true,
		},
		{
			"existing username + valid password but in restricted username auth collection",
			"clients",
			"clients57772",
			"1234567890",
			true,
		},
		{
			"existing username + valid password but in restricted username and email auth collection",
			"nologin",
			"test_username",
			"1234567890",
			true,
		},
		{
			"existing username + valid password",
			"users",
			"users75657",
			"1234567890",
			false,
		},

		// email
		{
			"existing email + wrong password",
			"users",
			"test@example.com",
			"invalid",
			true,
		},
		{
			"missing email + valid password",
			"users",
			"test_missing@example.com",
			"1234567890",
			true,
		},
		{
			"existing username + valid password but in restricted username auth collection",
			"clients",
			"test@example.com",
			"1234567890",
			false,
		},
		{
			"existing username + valid password but in restricted username and email auth collection",
			"nologin",
			"test@example.com",
			"1234567890",
			true,
		},
		{
			"existing email + valid password",
			"users",
			"test@example.com",
			"1234567890",
			false,
		},
	}

	for _, s := range scenarios {
		authCollection, err := testApp.Dao().FindCollectionByNameOrId(s.collectionName)
		if err != nil {
			t.Errorf("[%s] Failed to fetch auth collection: %v", s.testName, err)
		}

		form := forms.NewRecordPasswordLogin(testApp, authCollection)
		form.Identity = s.identity
		form.Password = s.password

		record, err := form.Submit()

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("[%s] Expected hasErr to be %v, got %v (%v)", s.testName, s.expectError, hasErr, err)
			continue
		}

		if hasErr {
			continue
		}

		if record.Email() != s.identity && record.Username() != s.identity {
			t.Errorf("[%s] Expected record with identity %q, got \n%v", s.testName, s.identity, record)
		}
	}
}
