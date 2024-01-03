package forms_test

import (
	"encoding/json"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/tests"
)

func TestUserOauth2LoginValidate(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		testName       string
		collectionName string
		jsonData       string
		expectedErrors []string
	}{
		{
			"empty payload",
			"users",
			"{}",
			[]string{"provider", "code", "redirectUrl"},
		},
		{
			"empty data",
			"users",
			`{"provider":"","code":"","codeVerifier":"","redirectUrl":""}`,
			[]string{"provider", "code", "redirectUrl"},
		},
		{
			"missing provider",
			"users",
			`{"provider":"missing","code":"123","codeVerifier":"123","redirectUrl":"https://example.com"}`,
			[]string{"provider"},
		},
		{
			"disabled provider",
			"users",
			`{"provider":"github","code":"123","codeVerifier":"123","redirectUrl":"https://example.com"}`,
			[]string{"provider"},
		},
		{
			"enabled provider",
			"users",
			`{"provider":"gitlab","code":"123","codeVerifier":"123","redirectUrl":"https://example.com"}`,
			[]string{},
		},
		{
			"[#3689] any redirectUrl value",
			"users",
			`{"provider":"gitlab","code":"123","codeVerifier":"123","redirectUrl":"something"}`,
			[]string{},
		},
	}

	for _, s := range scenarios {
		authCollection, _ := app.Dao().FindCollectionByNameOrId(s.collectionName)
		if authCollection == nil {
			t.Errorf("[%s] Failed to fetch auth collection", s.testName)
		}

		form := forms.NewRecordOAuth2Login(app, authCollection, nil)

		// load data
		loadErr := json.Unmarshal([]byte(s.jsonData), form)
		if loadErr != nil {
			t.Errorf("[%s] Failed to load form data: %v", s.testName, loadErr)
			continue
		}

		err := form.Validate()

		// parse errors
		errs, ok := err.(validation.Errors)
		if !ok && err != nil {
			t.Errorf("[%s] Failed to parse errors %v", s.testName, err)
			continue
		}

		// check errors
		if len(errs) > len(s.expectedErrors) {
			t.Errorf("[%s] Expected error keys %v, got %v", s.testName, s.expectedErrors, errs)
		}
		for _, k := range s.expectedErrors {
			if _, ok := errs[k]; !ok {
				t.Errorf("[%s] Missing expected error key %q in %v", s.testName, k, errs)
			}
		}
	}
}

// @todo consider mocking a Oauth2 provider to test Submit
