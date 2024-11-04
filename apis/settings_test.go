package apis_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/tests"
)

func TestSettingsList(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodGet,
			URL:             "/api/settings",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as regular user",
			Method: http.MethodGet,
			URL:    "/api/settings",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as superuser",
			Method: http.MethodGet,
			URL:    "/api/settings",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"meta":{`,
				`"logs":{`,
				`"smtp":{`,
				`"s3":{`,
				`"backups":{`,
				`"batch":{`,
			},
			ExpectedEvents: map[string]int{
				"*":                     0,
				"OnSettingsListRequest": 1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestSettingsSet(t *testing.T) {
	t.Parallel()

	validData := `{
		"meta":{"appName":"update_test"},
		"s3":{"secret": "s3_secret"},
		"backups":{"s3":{"secret":"backups_s3_secret"}}
	}`

	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodPatch,
			URL:             "/api/settings",
			Body:            strings.NewReader(validData),
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as regular user",
			Method: http.MethodPatch,
			URL:    "/api/settings",
			Body:   strings.NewReader(validData),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as superuser submitting empty data",
			Method: http.MethodPatch,
			URL:    "/api/settings",
			Body:   strings.NewReader(``),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"meta":{`,
				`"logs":{`,
				`"smtp":{`,
				`"s3":{`,
				`"backups":{`,
				`"batch":{`,
			},
			ExpectedEvents: map[string]int{
				"*":                         0,
				"OnSettingsUpdateRequest":   1,
				"OnModelUpdate":             1,
				"OnModelUpdateExecute":      1,
				"OnModelAfterUpdateSuccess": 1,
				"OnModelValidate":           1,
				"OnSettingsReload":          1,
			},
		},
		{
			Name:   "authorized as superuser submitting invalid data",
			Method: http.MethodPatch,
			URL:    "/api/settings",
			Body:   strings.NewReader(`{"meta":{"appName":""}}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"meta":{"appName":{"code":"validation_required"`,
			},
			ExpectedEvents: map[string]int{
				"*":                       0,
				"OnModelUpdate":           1,
				"OnModelAfterUpdateError": 1,
				"OnModelValidate":         1,
				"OnSettingsUpdateRequest": 1,
			},
		},
		{
			Name:   "authorized as superuser submitting valid data",
			Method: http.MethodPatch,
			URL:    "/api/settings",
			Body:   strings.NewReader(validData),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"meta":{`,
				`"logs":{`,
				`"smtp":{`,
				`"s3":{`,
				`"backups":{`,
				`"batch":{`,
				`"appName":"update_test"`,
			},
			NotExpectedContent: []string{
				"secret",
				"password",
			},
			ExpectedEvents: map[string]int{
				"*":                         0,
				"OnSettingsUpdateRequest":   1,
				"OnModelUpdate":             1,
				"OnModelUpdateExecute":      1,
				"OnModelAfterUpdateSuccess": 1,
				"OnModelValidate":           1,
				"OnSettingsReload":          1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestSettingsTestS3(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodPost,
			URL:             "/api/settings/test/s3",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as regular user",
			Method: http.MethodPost,
			URL:    "/api/settings/test/s3",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as superuser (missing body + no s3)",
			Method: http.MethodPost,
			URL:    "/api/settings/test/s3",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"filesystem":{`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
		{
			Name:   "authorized as superuser (invalid filesystem)",
			Method: http.MethodPost,
			URL:    "/api/settings/test/s3",
			Body:   strings.NewReader(`{"filesystem":"invalid"}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"filesystem":{`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
		{
			Name:   "authorized as superuser (valid filesystem and no s3)",
			Method: http.MethodPost,
			URL:    "/api/settings/test/s3",
			Body:   strings.NewReader(`{"filesystem":"storage"}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{}`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestSettingsTestEmail(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:   "unauthorized",
			Method: http.MethodPost,
			URL:    "/api/settings/test/email",
			Body: strings.NewReader(`{
				"template": "verification",
				"email": "test@example.com"
			}`),
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as regular user",
			Method: http.MethodPost,
			URL:    "/api/settings/test/email",
			Body: strings.NewReader(`{
				"template": "verification",
				"email": "test@example.com"
			}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as superuser (invalid body)",
			Method: http.MethodPost,
			URL:    "/api/settings/test/email",
			Body:   strings.NewReader(`{`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as superuser (empty json)",
			Method: http.MethodPost,
			URL:    "/api/settings/test/email",
			Body:   strings.NewReader(`{}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"email":{"code":"validation_required"`,
				`"template":{"code":"validation_required"`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
		{
			Name:   "authorized as superuser (verifiation template)",
			Method: http.MethodPost,
			URL:    "/api/settings/test/email",
			Body: strings.NewReader(`{
				"template": "verification",
				"email": "test@example.com"
			}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				if app.TestMailer.TotalSend() != 1 {
					t.Fatalf("[verification] Expected 1 sent email, got %d", app.TestMailer.TotalSend())
				}

				if len(app.TestMailer.LastMessage().To) != 1 {
					t.Fatalf("[verification] Expected 1 recipient, got %v", app.TestMailer.LastMessage().To)
				}

				if app.TestMailer.LastMessage().To[0].Address != "test@example.com" {
					t.Fatalf("[verification] Expected the email to be sent to %s, got %s", "test@example.com", app.TestMailer.LastMessage().To[0].Address)
				}

				if !strings.Contains(app.TestMailer.LastMessage().HTML, "Verify") {
					t.Fatalf("[verification] Expected to sent a verification email, got \n%v\n%v", app.TestMailer.LastMessage().Subject, app.TestMailer.LastMessage().HTML)
				}
			},
			ExpectedStatus:  204,
			ExpectedContent: []string{},
			ExpectedEvents: map[string]int{
				"*":                              0,
				"OnMailerSend":                   1,
				"OnMailerRecordVerificationSend": 1,
			},
		},
		{
			Name:   "authorized as superuser (password reset template)",
			Method: http.MethodPost,
			URL:    "/api/settings/test/email",
			Body: strings.NewReader(`{
				"template": "password-reset",
				"email": "test@example.com"
			}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				if app.TestMailer.TotalSend() != 1 {
					t.Fatalf("[password-reset] Expected 1 sent email, got %d", app.TestMailer.TotalSend())
				}

				if len(app.TestMailer.LastMessage().To) != 1 {
					t.Fatalf("[password-reset] Expected 1 recipient, got %v", app.TestMailer.LastMessage().To)
				}

				if app.TestMailer.LastMessage().To[0].Address != "test@example.com" {
					t.Fatalf("[password-reset] Expected the email to be sent to %s, got %s", "test@example.com", app.TestMailer.LastMessage().To[0].Address)
				}

				if !strings.Contains(app.TestMailer.LastMessage().HTML, "Reset password") {
					t.Fatalf("[password-reset] Expected to sent a password-reset email, got \n%v\n%v", app.TestMailer.LastMessage().Subject, app.TestMailer.LastMessage().HTML)
				}
			},
			ExpectedStatus:  204,
			ExpectedContent: []string{},
			ExpectedEvents: map[string]int{
				"*":                               0,
				"OnMailerSend":                    1,
				"OnMailerRecordPasswordResetSend": 1,
			},
		},
		{
			Name:   "authorized as superuser (email change)",
			Method: http.MethodPost,
			URL:    "/api/settings/test/email",
			Body: strings.NewReader(`{
				"template": "email-change",
				"email": "test@example.com"
			}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				if app.TestMailer.TotalSend() != 1 {
					t.Fatalf("[email-change] Expected 1 sent email, got %d", app.TestMailer.TotalSend())
				}

				if len(app.TestMailer.LastMessage().To) != 1 {
					t.Fatalf("[email-change] Expected 1 recipient, got %v", app.TestMailer.LastMessage().To)
				}

				if app.TestMailer.LastMessage().To[0].Address != "test@example.com" {
					t.Fatalf("[email-change] Expected the email to be sent to %s, got %s", "test@example.com", app.TestMailer.LastMessage().To[0].Address)
				}

				if !strings.Contains(app.TestMailer.LastMessage().HTML, "Confirm new email") {
					t.Fatalf("[email-change] Expected to sent a confirm new email email, got \n%v\n%v", app.TestMailer.LastMessage().Subject, app.TestMailer.LastMessage().HTML)
				}
			},
			ExpectedStatus:  204,
			ExpectedContent: []string{},
			ExpectedEvents: map[string]int{
				"*":                             0,
				"OnMailerSend":                  1,
				"OnMailerRecordEmailChangeSend": 1,
			},
		},
		{
			Name:   "authorized as superuser (otp)",
			Method: http.MethodPost,
			URL:    "/api/settings/test/email",
			Body: strings.NewReader(`{
				"template": "otp",
				"email": "test@example.com"
			}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				if app.TestMailer.TotalSend() != 1 {
					t.Fatalf("[otp] Expected 1 sent email, got %d", app.TestMailer.TotalSend())
				}

				if len(app.TestMailer.LastMessage().To) != 1 {
					t.Fatalf("[otp] Expected 1 recipient, got %v", app.TestMailer.LastMessage().To)
				}

				if app.TestMailer.LastMessage().To[0].Address != "test@example.com" {
					t.Fatalf("[otp] Expected the email to be sent to %s, got %s", "test@example.com", app.TestMailer.LastMessage().To[0].Address)
				}

				if !strings.Contains(app.TestMailer.LastMessage().HTML, "one-time password") {
					t.Fatalf("[otp] Expected to sent OTP email, got \n%v\n%v", app.TestMailer.LastMessage().Subject, app.TestMailer.LastMessage().HTML)
				}
			},
			ExpectedStatus:  204,
			ExpectedContent: []string{},
			ExpectedEvents: map[string]int{
				"*":                     0,
				"OnMailerSend":          1,
				"OnMailerRecordOTPSend": 1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestGenerateAppleClientSecret(t *testing.T) {
	t.Parallel()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	encodedKey, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}

	privatePem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: encodedKey,
		},
	)

	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodPost,
			URL:             "/api/settings/apple/generate-client-secret",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as regular user",
			Method: http.MethodPost,
			URL:    "/api/settings/apple/generate-client-secret",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as superuser (invalid body)",
			Method: http.MethodPost,
			URL:    "/api/settings/apple/generate-client-secret",
			Body:   strings.NewReader(`{`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as superuser (empty json)",
			Method: http.MethodPost,
			URL:    "/api/settings/apple/generate-client-secret",
			Body:   strings.NewReader(`{}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"clientId":{"code":"validation_required"`,
				`"teamId":{"code":"validation_required"`,
				`"keyId":{"code":"validation_required"`,
				`"privateKey":{"code":"validation_required"`,
				`"duration":{"code":"validation_required"`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
		{
			Name:   "authorized as superuser (invalid data)",
			Method: http.MethodPost,
			URL:    "/api/settings/apple/generate-client-secret",
			Body: strings.NewReader(`{
				"clientId": "",
				"teamId": "123456789",
				"keyId": "123456789",
				"privateKey": "invalid",
				"duration": -1
			}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"clientId":{"code":"validation_required"`,
				`"teamId":{"code":"validation_length_invalid"`,
				`"keyId":{"code":"validation_length_invalid"`,
				`"privateKey":{"code":"validation_match_invalid"`,
				`"duration":{"code":"validation_min_greater_equal_than_required"`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
		{
			Name:   "authorized as superuser (valid data)",
			Method: http.MethodPost,
			URL:    "/api/settings/apple/generate-client-secret",
			Body: strings.NewReader(fmt.Sprintf(`{
				"clientId": "123",
				"teamId": "1234567890",
				"keyId": "1234567891",
				"privateKey": %q,
				"duration": 1
			}`, privatePem)),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"secret":"`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
