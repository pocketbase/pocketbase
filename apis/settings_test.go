package apis_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/tests"
)

func TestSettingsList(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodGet,
			Url:             "/api/settings",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as auth record",
			Method: http.MethodGet,
			Url:    "/api/settings",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin",
			Method: http.MethodGet,
			Url:    "/api/settings",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"meta":{`,
				`"logs":{`,
				`"smtp":{`,
				`"s3":{`,
				`"adminAuthToken":{`,
				`"adminPasswordResetToken":{`,
				`"recordAuthToken":{`,
				`"recordPasswordResetToken":{`,
				`"recordEmailChangeToken":{`,
				`"recordVerificationToken":{`,
				`"emailAuth":{`,
				`"googleAuth":{`,
				`"facebookAuth":{`,
				`"githubAuth":{`,
				`"gitlabAuth":{`,
				`"twitterAuth":{`,
				`"discordAuth":{`,
				`"microsoftAuth":{`,
				`"spotifyAuth":{`,
				`"kakaoAuth":{`,
				`"twitchAuth":{`,
				`"secret":"******"`,
				`"clientSecret":"******"`,
			},
			ExpectedEvents: map[string]int{
				"OnSettingsListRequest": 1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestSettingsSet(t *testing.T) {
	validData := `{"meta":{"appName":"update_test"}}`

	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodPatch,
			Url:             "/api/settings",
			Body:            strings.NewReader(validData),
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as auth record",
			Method: http.MethodPatch,
			Url:    "/api/settings",
			Body:   strings.NewReader(validData),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin submitting empty data",
			Method: http.MethodPatch,
			Url:    "/api/settings",
			Body:   strings.NewReader(``),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"meta":{`,
				`"logs":{`,
				`"smtp":{`,
				`"s3":{`,
				`"adminAuthToken":{`,
				`"adminPasswordResetToken":{`,
				`"recordAuthToken":{`,
				`"recordPasswordResetToken":{`,
				`"recordEmailChangeToken":{`,
				`"recordVerificationToken":{`,
				`"emailAuth":{`,
				`"googleAuth":{`,
				`"facebookAuth":{`,
				`"githubAuth":{`,
				`"gitlabAuth":{`,
				`"discordAuth":{`,
				`"microsoftAuth":{`,
				`"spotifyAuth":{`,
				`"kakaoAuth":{`,
				`"twitchAuth":{`,
				`"secret":"******"`,
				`"clientSecret":"******"`,
				`"appName":"Acme"`,
			},
			ExpectedEvents: map[string]int{
				"OnModelBeforeUpdate":           1,
				"OnModelAfterUpdate":            1,
				"OnSettingsBeforeUpdateRequest": 1,
				"OnSettingsAfterUpdateRequest":  1,
			},
		},
		{
			Name:   "authorized as admin submitting invalid data",
			Method: http.MethodPatch,
			Url:    "/api/settings",
			Body:   strings.NewReader(`{"meta":{"appName":""}}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"meta":{"appName":{"code":"validation_required"`,
			},
		},
		{
			Name:   "authorized as admin submitting valid data",
			Method: http.MethodPatch,
			Url:    "/api/settings",
			Body:   strings.NewReader(validData),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"meta":{`,
				`"logs":{`,
				`"smtp":{`,
				`"s3":{`,
				`"adminAuthToken":{`,
				`"adminPasswordResetToken":{`,
				`"recordAuthToken":{`,
				`"recordPasswordResetToken":{`,
				`"recordEmailChangeToken":{`,
				`"recordVerificationToken":{`,
				`"emailAuth":{`,
				`"googleAuth":{`,
				`"facebookAuth":{`,
				`"githubAuth":{`,
				`"gitlabAuth":{`,
				`"twitterAuth":{`,
				`"discordAuth":{`,
				`"microsoftAuth":{`,
				`"spotifyAuth":{`,
				`"kakaoAuth":{`,
				`"twitchAuth":{`,
				`"secret":"******"`,
				`"clientSecret":"******"`,
				`"appName":"update_test"`,
			},
			ExpectedEvents: map[string]int{
				"OnModelBeforeUpdate":           1,
				"OnModelAfterUpdate":            1,
				"OnSettingsBeforeUpdateRequest": 1,
				"OnSettingsAfterUpdateRequest":  1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestSettingsTestS3(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodPost,
			Url:             "/api/settings/test/s3",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as auth record",
			Method: http.MethodPost,
			Url:    "/api/settings/test/s3",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin (no s3)",
			Method: http.MethodPost,
			Url:    "/api/settings/test/s3",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestSettingsTestEmail(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:   "unauthorized",
			Method: http.MethodPost,
			Url:    "/api/settings/test/email",
			Body: strings.NewReader(`{
				"template": "verification",
				"email": "test@example.com"
			}`),
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as auth record",
			Method: http.MethodPost,
			Url:    "/api/settings/test/email",
			Body: strings.NewReader(`{
				"template": "verification",
				"email": "test@example.com"
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin (invalid body)",
			Method: http.MethodPost,
			Url:    "/api/settings/test/email",
			Body:   strings.NewReader(`{`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin (empty json)",
			Method: http.MethodPost,
			Url:    "/api/settings/test/email",
			Body:   strings.NewReader(`{}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"email":{"code":"validation_required"`,
				`"template":{"code":"validation_required"`,
			},
		},
		{
			Name:   "authorized as admin (verifiation template)",
			Method: http.MethodPost,
			Url:    "/api/settings/test/email",
			Body: strings.NewReader(`{
				"template": "verification",
				"email": "test@example.com"
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			AfterTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				if app.TestMailer.TotalSend != 1 {
					t.Fatalf("[verification] Expected 1 sent email, got %d", app.TestMailer.TotalSend)
				}

				if app.TestMailer.LastMessage.To.Address != "test@example.com" {
					t.Fatalf("[verification] Expected the email to be sent to %s, got %s", "test@example.com", app.TestMailer.LastMessage.To.Address)
				}

				if !strings.Contains(app.TestMailer.LastMessage.HTML, "Verify") {
					t.Fatalf("[verification] Expected to sent a verification email, got \n%v\n%v", app.TestMailer.LastMessage.Subject, app.TestMailer.LastMessage.HTML)
				}
			},
			ExpectedStatus:  204,
			ExpectedContent: []string{},
			ExpectedEvents: map[string]int{
				"OnMailerBeforeRecordVerificationSend": 1,
				"OnMailerAfterRecordVerificationSend":  1,
			},
		},
		{
			Name:   "authorized as admin (password reset template)",
			Method: http.MethodPost,
			Url:    "/api/settings/test/email",
			Body: strings.NewReader(`{
				"template": "password-reset",
				"email": "test@example.com"
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			AfterTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				if app.TestMailer.TotalSend != 1 {
					t.Fatalf("[password-reset] Expected 1 sent email, got %d", app.TestMailer.TotalSend)
				}

				if app.TestMailer.LastMessage.To.Address != "test@example.com" {
					t.Fatalf("[password-reset] Expected the email to be sent to %s, got %s", "test@example.com", app.TestMailer.LastMessage.To.Address)
				}

				if !strings.Contains(app.TestMailer.LastMessage.HTML, "Reset password") {
					t.Fatalf("[password-reset] Expected to sent a password-reset email, got \n%v\n%v", app.TestMailer.LastMessage.Subject, app.TestMailer.LastMessage.HTML)
				}
			},
			ExpectedStatus:  204,
			ExpectedContent: []string{},
			ExpectedEvents: map[string]int{
				"OnMailerBeforeRecordResetPasswordSend": 1,
				"OnMailerAfterRecordResetPasswordSend":  1,
			},
		},
		{
			Name:   "authorized as admin (email change)",
			Method: http.MethodPost,
			Url:    "/api/settings/test/email",
			Body: strings.NewReader(`{
				"template": "email-change",
				"email": "test@example.com"
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			AfterTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				if app.TestMailer.TotalSend != 1 {
					t.Fatalf("[email-change] Expected 1 sent email, got %d", app.TestMailer.TotalSend)
				}

				if app.TestMailer.LastMessage.To.Address != "test@example.com" {
					t.Fatalf("[email-change] Expected the email to be sent to %s, got %s", "test@example.com", app.TestMailer.LastMessage.To.Address)
				}

				if !strings.Contains(app.TestMailer.LastMessage.HTML, "Confirm new email") {
					t.Fatalf("[email-change] Expected to sent a confirm new email email, got \n%v\n%v", app.TestMailer.LastMessage.Subject, app.TestMailer.LastMessage.HTML)
				}
			},
			ExpectedStatus:  204,
			ExpectedContent: []string{},
			ExpectedEvents: map[string]int{
				"OnMailerBeforeRecordChangeEmailSend": 1,
				"OnMailerAfterRecordChangeEmailSend":  1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
