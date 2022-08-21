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
			Name:   "authorized as user",
			Method: http.MethodGet,
			Url:    "/api/settings",
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin",
			Method: http.MethodGet,
			Url:    "/api/settings",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"meta":{`,
				`"logs":{`,
				`"smtp":{`,
				`"s3":{`,
				`"adminAuthToken":{`,
				`"adminPasswordResetToken":{`,
				`"userAuthToken":{`,
				`"userPasswordResetToken":{`,
				`"userEmailChangeToken":{`,
				`"userVerificationToken":{`,
				`"emailAuth":{`,
				`"googleAuth":{`,
				`"facebookAuth":{`,
				`"githubAuth":{`,
				`"gitlabAuth":{`,
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
	validData := `{"meta":{"appName":"update_test"},"emailAuth":{"minPasswordLength": 12}}`

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
			Name:   "authorized as user",
			Method: http.MethodPatch,
			Url:    "/api/settings",
			Body:   strings.NewReader(validData),
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
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
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"meta":{`,
				`"logs":{`,
				`"smtp":{`,
				`"s3":{`,
				`"adminAuthToken":{`,
				`"adminPasswordResetToken":{`,
				`"userAuthToken":{`,
				`"userPasswordResetToken":{`,
				`"userEmailChangeToken":{`,
				`"userVerificationToken":{`,
				`"emailAuth":{`,
				`"googleAuth":{`,
				`"facebookAuth":{`,
				`"githubAuth":{`,
				`"gitlabAuth":{`,
				`"secret":"******"`,
				`"clientSecret":"******"`,
				`"appName":"Acme"`,
				`"minPasswordLength":8`,
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
			Body:   strings.NewReader(`{"meta":{"appName":""},"emailAuth":{"minPasswordLength": 3}}`),
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"emailAuth":{"minPasswordLength":{"code":"validation_min_greater_equal_than_required","message":"Must be no less than 5."}}`,
				`"meta":{"appName":{"code":"validation_required","message":"Cannot be blank."}}`,
			},
		},
		{
			Name:   "authorized as admin submitting valid data",
			Method: http.MethodPatch,
			Url:    "/api/settings",
			Body:   strings.NewReader(validData),
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"meta":{`,
				`"logs":{`,
				`"smtp":{`,
				`"s3":{`,
				`"adminAuthToken":{`,
				`"adminPasswordResetToken":{`,
				`"userAuthToken":{`,
				`"userPasswordResetToken":{`,
				`"userEmailChangeToken":{`,
				`"userVerificationToken":{`,
				`"emailAuth":{`,
				`"googleAuth":{`,
				`"facebookAuth":{`,
				`"githubAuth":{`,
				`"gitlabAuth":{`,
				`"secret":"******"`,
				`"clientSecret":"******"`,
				`"appName":"update_test"`,
				`"minPasswordLength":12`,
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
			Name:   "authorized as user",
			Method: http.MethodPost,
			Url:    "/api/settings/test/s3",
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin (no s3)",
			Method: http.MethodPost,
			Url:    "/api/settings/test/s3",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		// @todo consider creating a test S3 filesystem
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
			Name:   "authorized as user",
			Method: http.MethodPost,
			Url:    "/api/settings/test/email",
			Body: strings.NewReader(`{
				"template": "verification",
				"email": "test@example.com"
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
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
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
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
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
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
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			AfterFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				if app.TestMailer.TotalSend != 1 {
					t.Fatalf("[verification] Expected 1 sent email, got %d", app.TestMailer.TotalSend)
				}

				if app.TestMailer.LastToAddress.Address != "test@example.com" {
					t.Fatalf("[verification] Expected the email to be sent to %s, got %s", "test@example.com", app.TestMailer.LastToAddress.Address)
				}

				if !strings.Contains(app.TestMailer.LastHtmlBody, "Verify") {
					t.Fatalf("[verification] Expected to sent a verification email, got \n%v\n%v", app.TestMailer.LastHtmlSubject, app.TestMailer.LastHtmlBody)
				}
			},
			ExpectedStatus:  204,
			ExpectedContent: []string{},
			ExpectedEvents: map[string]int{
				"OnMailerBeforeUserVerificationSend": 1,
				"OnMailerAfterUserVerificationSend":  1,
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
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			AfterFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				if app.TestMailer.TotalSend != 1 {
					t.Fatalf("[password-reset] Expected 1 sent email, got %d", app.TestMailer.TotalSend)
				}

				if app.TestMailer.LastToAddress.Address != "test@example.com" {
					t.Fatalf("[password-reset] Expected the email to be sent to %s, got %s", "test@example.com", app.TestMailer.LastToAddress.Address)
				}

				if !strings.Contains(app.TestMailer.LastHtmlBody, "Reset password") {
					t.Fatalf("[password-reset] Expected to sent a password-reset email, got \n%v\n%v", app.TestMailer.LastHtmlSubject, app.TestMailer.LastHtmlBody)
				}
			},
			ExpectedStatus:  204,
			ExpectedContent: []string{},
			ExpectedEvents: map[string]int{
				"OnMailerBeforeUserResetPasswordSend": 1,
				"OnMailerAfterUserResetPasswordSend":  1,
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
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			AfterFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				if app.TestMailer.TotalSend != 1 {
					t.Fatalf("[email-change] Expected 1 sent email, got %d", app.TestMailer.TotalSend)
				}

				if app.TestMailer.LastToAddress.Address != "test@example.com" {
					t.Fatalf("[email-change] Expected the email to be sent to %s, got %s", "test@example.com", app.TestMailer.LastToAddress.Address)
				}

				if !strings.Contains(app.TestMailer.LastHtmlBody, "Confirm new email") {
					t.Fatalf("[email-change] Expected to sent a confirm new email email, got \n%v\n%v", app.TestMailer.LastHtmlSubject, app.TestMailer.LastHtmlBody)
				}
			},
			ExpectedStatus:  204,
			ExpectedContent: []string{},
			ExpectedEvents: map[string]int{
				"OnMailerBeforeUserChangeEmailSend": 1,
				"OnMailerAfterUserChangeEmailSend":  1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
