package forms_test

import (
	"encoding/json"
	"errors"
	"os"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models/settings"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/security"
)

func TestNewSettingsUpsert(t *testing.T) {
	t.Parallel()

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

func TestSettingsUpsertValidateAndSubmit(t *testing.T) {
	t.Parallel()

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
			`{"meta": {"appName": ""}, "logs": {"maxDays": -1}}`,
			false,
			[]string{"meta", "logs"},
		},
		// success - valid data (plain)
		{
			`{"meta": {"appName": "test"}, "logs": {"maxDays": 0}}`,
			false,
			nil,
		},
		// success - valid data (encrypt)
		{
			`{"meta": {"appName": "test"}, "logs": {"maxDays": 7}}`,
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

		interceptorCalls := 0
		interceptor := func(next forms.InterceptorNextFunc[*settings.Settings]) forms.InterceptorNextFunc[*settings.Settings] {
			return func(s *settings.Settings) error {
				interceptorCalls++
				return next(s)
			}
		}

		// parse errors
		result := form.Submit(interceptor)
		errs, ok := result.(validation.Errors)
		if !ok && result != nil {
			t.Errorf("(%d) Failed to parse errors %v", i, result)
			continue
		}

		// check interceptor calls
		expectInterceptorCall := 1
		if len(s.expectedErrors) > 0 {
			expectInterceptorCall = 0
		}
		if interceptorCalls != expectInterceptorCall {
			t.Errorf("(%d) Expected interceptor to be called %d, got %d", i, expectInterceptorCall, interceptorCalls)
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

func TestSettingsUpsertSubmitInterceptors(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	form := forms.NewSettingsUpsert(app)
	form.Meta.AppName = "test_new"

	testErr := errors.New("test_error")

	interceptor1Called := false
	interceptor1 := func(next forms.InterceptorNextFunc[*settings.Settings]) forms.InterceptorNextFunc[*settings.Settings] {
		return func(s *settings.Settings) error {
			interceptor1Called = true
			return next(s)
		}
	}

	interceptor2Called := false
	interceptor2 := func(next forms.InterceptorNextFunc[*settings.Settings]) forms.InterceptorNextFunc[*settings.Settings] {
		return func(s *settings.Settings) error {
			interceptor2Called = true
			return testErr
		}
	}

	submitErr := form.Submit(interceptor1, interceptor2)
	if submitErr != testErr {
		t.Fatalf("Expected submitError %v, got %v", testErr, submitErr)
	}

	if !interceptor1Called {
		t.Fatalf("Expected interceptor1 to be called")
	}

	if !interceptor2Called {
		t.Fatalf("Expected interceptor2 to be called")
	}
}
