package forms_test

import (
	"encoding/json"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/tests"
)

func TestUserOauth2LoginPanic(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("The form did not panic")
		}
	}()

	forms.NewUserOauth2Login(nil)
}

func TestUserOauth2LoginValidate(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		jsonData       string
		expectedErrors []string
	}{
		// empty payload
		{"{}", []string{"provider", "code", "codeVerifier", "redirectUrl"}},
		// empty data
		{
			`{"provider":"","code":"","codeVerifier":"","redirectUrl":""}`,
			[]string{"provider", "code", "codeVerifier", "redirectUrl"},
		},
		// missing provider
		{
			`{"provider":"missing","code":"123","codeVerifier":"123","redirectUrl":"https://example.com"}`,
			[]string{"provider"},
		},
		// disabled provider
		{
			`{"provider":"github","code":"123","codeVerifier":"123","redirectUrl":"https://example.com"}`,
			[]string{"provider"},
		},
		// enabled provider
		{
			`{"provider":"gitlab","code":"123","codeVerifier":"123","redirectUrl":"https://example.com"}`,
			[]string{},
		},
	}

	for i, s := range scenarios {
		form := forms.NewUserOauth2Login(app)

		// load data
		loadErr := json.Unmarshal([]byte(s.jsonData), form)
		if loadErr != nil {
			t.Errorf("(%d) Failed to load form data: %v", i, loadErr)
			continue
		}

		err := form.Validate()

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
	}
}

// @todo consider mocking a Oauth2 provider to test Submit
