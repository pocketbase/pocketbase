package forms_test

import (
	"encoding/json"
	"os"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/security"
)

func TestNewSettingsUpsert(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	app.Settings().Meta.AppName = "name_update"

	form := forms.NewSettingsUpsert(app)

	formSettings, _ := json.Marshal(form.Settings)
	appSettings, _ := json.Marshal(app.Settings())

	if string(formSettings) != string(appSettings) {
		t.Errorf("Expected settings \n%s, got \n%s", string(appSettings), string(formSettings))
	}
}

func TestSettingsUpsertValidate(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	form := forms.NewSettingsUpsert(app)

	// check if settings validations are triggered
	// (there are already individual tests for each setting)
	form.Meta.AppName = ""
	form.Logs.MaxDays = -10

	// parse errors
	err := form.Validate()
	jsonResult, _ := json.Marshal(err)

	expected := `{"logs":{"maxDays":"must be no less than 0"},"meta":{"appName":"cannot be blank"}}`

	if string(jsonResult) != expected {
		t.Errorf("Expected %v, got %v", expected, string(jsonResult))
	}
}

func TestSettingsUpsertSubmit(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		jsonData       string
		encryption     bool
		expectedErrors []string
	}{
		// empty (plain)
		{"{}", false, nil},
		// empty (encrypt)
		{"{}", true, nil},
		// failure - invalid data
		{
			`{"emailAuth": {"minPasswordLength": 1}, "logs": {"maxDays": -1}}`,
			false,
			[]string{"emailAuth", "logs"},
		},
		// success - valid data (plain)
		{
			`{"emailAuth": {"minPasswordLength": 6}, "logs": {"maxDays": 0}}`,
			false,
			nil,
		},
		// success - valid data (encrypt)
		{
			`{"emailAuth": {"minPasswordLength": 6}, "logs": {"maxDays": 0}}`,
			true,
			nil,
		},
	}

	for i, s := range scenarios {
		if s.encryption {
			os.Setenv(app.EncryptionEnv(), security.RandomString(32))
		} else {
			os.Unsetenv(app.EncryptionEnv())
		}

		form := forms.NewSettingsUpsert(app)

		// load data
		loadErr := json.Unmarshal([]byte(s.jsonData), form)
		if loadErr != nil {
			t.Errorf("(%d) Failed to load form data: %v", i, loadErr)
			continue
		}

		// parse errors
		result := form.Submit()
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

		if len(s.expectedErrors) > 0 {
			continue
		}

		formSettings, _ := json.Marshal(form.Settings)
		appSettings, _ := json.Marshal(app.Settings())

		if string(formSettings) != string(appSettings) {
			t.Errorf("Expected app settings \n%s, got \n%s", string(appSettings), string(formSettings))
		}
	}
}
