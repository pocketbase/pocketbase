package apis_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

func TestSettingsList(t *testing.T) {
	t.Parallel()

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
				`"backups":{`,
				`"adminAuthToken":{`,
				`"adminPasswordResetToken":{`,
				`"adminFileToken":{`,
				`"recordAuthToken":{`,
				`"recordPasswordResetToken":{`,
				`"recordEmailChangeToken":{`,
				`"recordVerificationToken":{`,
				`"recordFileToken":{`,
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
				`"stravaAuth":{`,
				`"giteeAuth":{`,
				`"livechatAuth":{`,
				`"giteaAuth":{`,
				`"oidcAuth":{`,
				`"oidc2Auth":{`,
				`"oidc3Auth":{`,
				`"appleAuth":{`,
				`"instagramAuth":{`,
				`"vkAuth":{`,
				`"yandexAuth":{`,
				`"patreonAuth":{`,
				`"mailcowAuth":{`,
				`"bitbucketAuth":{`,
				`"planningcenterAuth":{`,
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
	t.Parallel()

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
				`"backups":{`,
				`"adminAuthToken":{`,
				`"adminPasswordResetToken":{`,
				`"adminFileToken":{`,
				`"recordAuthToken":{`,
				`"recordPasswordResetToken":{`,
				`"recordEmailChangeToken":{`,
				`"recordVerificationToken":{`,
				`"recordFileToken":{`,
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
				`"stravaAuth":{`,
				`"giteeAuth":{`,
				`"livechatAuth":{`,
				`"giteaAuth":{`,
				`"oidcAuth":{`,
				`"oidc2Auth":{`,
				`"oidc3Auth":{`,
				`"appleAuth":{`,
				`"instagramAuth":{`,
				`"vkAuth":{`,
				`"yandexAuth":{`,
				`"patreonAuth":{`,
				`"mailcowAuth":{`,
				`"bitbucketAuth":{`,
				`"planningcenterAuth":{`,
				`"secret":"******"`,
				`"clientSecret":"******"`,
				`"appName":"acme_test"`,
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
				`"backups":{`,
				`"adminAuthToken":{`,
				`"adminPasswordResetToken":{`,
				`"adminFileToken":{`,
				`"recordAuthToken":{`,
				`"recordPasswordResetToken":{`,
				`"recordEmailChangeToken":{`,
				`"recordVerificationToken":{`,
				`"recordFileToken":{`,
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
				`"stravaAuth":{`,
				`"giteeAuth":{`,
				`"livechatAuth":{`,
				`"giteaAuth":{`,
				`"oidcAuth":{`,
				`"oidc2Auth":{`,
				`"oidc3Auth":{`,
				`"appleAuth":{`,
				`"instagramAuth":{`,
				`"vkAuth":{`,
				`"yandexAuth":{`,
				`"patreonAuth":{`,
				`"mailcowAuth":{`,
				`"bitbucketAuth":{`,
				`"planningcenterAuth":{`,
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
		{
			Name:   "OnSettingsAfterUpdateRequest error response",
			Method: http.MethodPatch,
			Url:    "/api/settings",
			Body:   strings.NewReader(validData),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				app.OnSettingsAfterUpdateRequest().Add(func(e *core.SettingsUpdateEvent) error {
					return errors.New("error")
				})
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
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
	t.Parallel()

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
			Name:   "authorized as admin (missing body + no s3)",
			Method: http.MethodPost,
			Url:    "/api/settings/test/s3",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"filesystem":{`,
			},
		},
		{
			Name:   "authorized as admin (invalid filesystem)",
			Method: http.MethodPost,
			Url:    "/api/settings/test/s3",
			Body:   strings.NewReader(`{"filesystem":"invalid"}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"filesystem":{`,
			},
		},
		{
			Name:   "authorized as admin (valid filesystem and no s3)",
			Method: http.MethodPost,
			Url:    "/api/settings/test/s3",
			Body:   strings.NewReader(`{"filesystem":"storage"}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{}`,
			},
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
			AfterTestFunc: func(t *testing.T, app *tests.TestApp, res *http.Response) {
				if app.TestMailer.TotalSend != 1 {
					t.Fatalf("[verification] Expected 1 sent email, got %d", app.TestMailer.TotalSend)
				}

				if len(app.TestMailer.LastMessage.To) != 1 {
					t.Fatalf("[verification] Expected 1 recipient, got %v", app.TestMailer.LastMessage.To)
				}

				if app.TestMailer.LastMessage.To[0].Address != "test@example.com" {
					t.Fatalf("[verification] Expected the email to be sent to %s, got %s", "test@example.com", app.TestMailer.LastMessage.To[0].Address)
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
			AfterTestFunc: func(t *testing.T, app *tests.TestApp, res *http.Response) {
				if app.TestMailer.TotalSend != 1 {
					t.Fatalf("[password-reset] Expected 1 sent email, got %d", app.TestMailer.TotalSend)
				}

				if len(app.TestMailer.LastMessage.To) != 1 {
					t.Fatalf("[password-reset] Expected 1 recipient, got %v", app.TestMailer.LastMessage.To)
				}

				if app.TestMailer.LastMessage.To[0].Address != "test@example.com" {
					t.Fatalf("[password-reset] Expected the email to be sent to %s, got %s", "test@example.com", app.TestMailer.LastMessage.To[0].Address)
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
			AfterTestFunc: func(t *testing.T, app *tests.TestApp, res *http.Response) {
				if app.TestMailer.TotalSend != 1 {
					t.Fatalf("[email-change] Expected 1 sent email, got %d", app.TestMailer.TotalSend)
				}

				if len(app.TestMailer.LastMessage.To) != 1 {
					t.Fatalf("[email-change] Expected 1 recipient, got %v", app.TestMailer.LastMessage.To)
				}

				if app.TestMailer.LastMessage.To[0].Address != "test@example.com" {
					t.Fatalf("[email-change] Expected the email to be sent to %s, got %s", "test@example.com", app.TestMailer.LastMessage.To[0].Address)
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
			Url:             "/api/settings/apple/generate-client-secret",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as auth record",
			Method: http.MethodPost,
			Url:    "/api/settings/apple/generate-client-secret",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin (invalid body)",
			Method: http.MethodPost,
			Url:    "/api/settings/apple/generate-client-secret",
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
			Url:    "/api/settings/apple/generate-client-secret",
			Body:   strings.NewReader(`{}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"clientId":{"code":"validation_required"`,
				`"teamId":{"code":"validation_required"`,
				`"keyId":{"code":"validation_required"`,
				`"privateKey":{"code":"validation_required"`,
				`"duration":{"code":"validation_required"`,
			},
		},
		{
			Name:   "authorized as admin (invalid data)",
			Method: http.MethodPost,
			Url:    "/api/settings/apple/generate-client-secret",
			Body: strings.NewReader(`{
				"clientId": "",
				"teamId": "123456789",
				"keyId": "123456789",
				"privateKey": "invalid",
				"duration": -1
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"clientId":{"code":"validation_required"`,
				`"teamId":{"code":"validation_length_invalid"`,
				`"keyId":{"code":"validation_length_invalid"`,
				`"privateKey":{"code":"validation_match_invalid"`,
				`"duration":{"code":"validation_min_greater_equal_than_required"`,
			},
		},
		{
			Name:   "authorized as admin (valid data)",
			Method: http.MethodPost,
			Url:    "/api/settings/apple/generate-client-secret",
			Body: strings.NewReader(fmt.Sprintf(`{
				"clientId": "123",
				"teamId": "1234567890",
				"keyId": "1234567891",
				"privateKey": %q,
				"duration": 1
			}`, privatePem)),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"secret":"`,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
