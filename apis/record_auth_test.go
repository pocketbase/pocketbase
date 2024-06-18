package apis_test

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/subscriptions"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestRecordAuthMethodsList(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "missing collection",
			Method:          http.MethodGet,
			Url:             "/api/collections/missing/auth-methods",
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "non auth collection",
			Method:          http.MethodGet,
			Url:             "/api/collections/demo1/auth-methods",
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:           "auth collection with all auth methods allowed",
			Method:         http.MethodGet,
			Url:            "/api/collections/users/auth-methods",
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"usernamePassword":true`,
				`"emailPassword":true`,
				`"onlyVerified":false`,
				`"authProviders":[{`,
				`"name":"gitlab"`,
				`"state":`,
				`"codeVerifier":`,
				`"codeChallenge":`,
				`"codeChallengeMethod":`,
				`"authUrl":`,
				`redirect_uri="`, // ensures that the redirect_uri is the last url param
			},
		},
		{
			Name:           "auth collection with only email/password auth allowed",
			Method:         http.MethodGet,
			Url:            "/api/collections/clients/auth-methods",
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"usernamePassword":false`,
				`"emailPassword":true`,
				`"onlyVerified":true`,
				`"authProviders":[]`,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRecordAuthWithPassword(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "invalid body format",
			Method:          http.MethodPost,
			Url:             "/api/collections/users/auth-with-password",
			Body:            strings.NewReader(`{"identity`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:           "empty body params",
			Method:         http.MethodPost,
			Url:            "/api/collections/users/auth-with-password",
			Body:           strings.NewReader(`{"identity":"","password":""}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"identity":{`,
				`"password":{`,
			},
		},

		// username
		{
			Name:   "invalid username and valid password",
			Method: http.MethodPost,
			Url:    "/api/collections/users/auth-with-password",
			Body: strings.NewReader(`{
				"identity":"invalid",
				"password":"1234567890"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{}`,
			},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeAuthWithPasswordRequest": 1,
			},
		},
		{
			Name:   "valid username and invalid password",
			Method: http.MethodPost,
			Url:    "/api/collections/users/auth-with-password",
			Body: strings.NewReader(`{
				"identity":"test2_username",
				"password":"invalid"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{}`,
			},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeAuthWithPasswordRequest": 1,
			},
		},
		{
			Name:   "valid username and valid password in restricted collection",
			Method: http.MethodPost,
			Url:    "/api/collections/nologin/auth-with-password",
			Body: strings.NewReader(`{
				"identity":"test_username",
				"password":"1234567890"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{}`,
			},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeAuthWithPasswordRequest": 1,
			},
		},
		{
			Name:   "valid username and valid password in allowed collection",
			Method: http.MethodPost,
			Url:    "/api/collections/users/auth-with-password",
			Body: strings.NewReader(`{
				"identity":"test2_username",
				"password":"1234567890"
			}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"record":{`,
				`"token":"`,
				`"id":"oap640cot4yru2s"`,
				`"email":"test2@example.com"`,
			},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeAuthWithPasswordRequest": 1,
				"OnRecordAfterAuthWithPasswordRequest":  1,
				"OnRecordAuthRequest":                   1,
			},
		},

		// email
		{
			Name:   "invalid email and valid password",
			Method: http.MethodPost,
			Url:    "/api/collections/users/auth-with-password",
			Body: strings.NewReader(`{
				"identity":"missing@example.com",
				"password":"1234567890"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{}`,
			},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeAuthWithPasswordRequest": 1,
			},
		},
		{
			Name:   "valid email and invalid password",
			Method: http.MethodPost,
			Url:    "/api/collections/users/auth-with-password",
			Body: strings.NewReader(`{
				"identity":"test@example.com",
				"password":"invalid"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{}`,
			},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeAuthWithPasswordRequest": 1,
			},
		},
		{
			Name:   "valid email and valid password in restricted collection",
			Method: http.MethodPost,
			Url:    "/api/collections/nologin/auth-with-password",
			Body: strings.NewReader(`{
				"identity":"test@example.com",
				"password":"1234567890"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{}`,
			},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeAuthWithPasswordRequest": 1,
			},
		},
		{
			Name:   "valid email (unverified) and valid password in allowed collection",
			Method: http.MethodPost,
			Url:    "/api/collections/users/auth-with-password",
			Body: strings.NewReader(`{
				"identity":"test@example.com",
				"password":"1234567890"
			}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"record":{`,
				`"token":"`,
				`"id":"4q1xlclmfloku33"`,
				`"email":"test@example.com"`,
				`"verified":false`,
			},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeAuthWithPasswordRequest": 1,
				"OnRecordAfterAuthWithPasswordRequest":  1,
				"OnRecordAuthRequest":                   1,
				// lastLoginAlertSentAt update
				"OnModelAfterUpdate":  1,
				"OnModelBeforeUpdate": 1,
			},
		},

		// onlyVerified collection check
		{
			Name:   "unverified user in onlyVerified collection",
			Method: http.MethodPost,
			Url:    "/api/collections/clients/auth-with-password",
			Body: strings.NewReader(`{
				"identity":"test2@example.com",
				"password":"1234567890"
			}`),
			ExpectedStatus: 403,
			ExpectedContent: []string{
				`"data":{}`,
			},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeAuthWithPasswordRequest": 1,
				"OnRecordAfterAuthWithPasswordRequest":  1,
			},
		},
		{
			Name:   "verified user in onlyVerified collection",
			Method: http.MethodPost,
			Url:    "/api/collections/clients/auth-with-password",
			Body: strings.NewReader(`{
				"identity":"test@example.com",
				"password":"1234567890"
			}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"record":{`,
				`"token":"`,
				`"id":"gk390qegs4y47wn"`,
				`"email":"test@example.com"`,
				`"verified":true`,
			},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeAuthWithPasswordRequest": 1,
				"OnRecordAfterAuthWithPasswordRequest":  1,
				"OnRecordAuthRequest":                   1,
			},
		},

		// with already authenticated record or admin
		{
			Name:   "authenticated record",
			Method: http.MethodPost,
			Url:    "/api/collections/users/auth-with-password",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			Body: strings.NewReader(`{
				"identity":"test@example.com",
				"password":"1234567890"
			}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"record":{`,
				`"token":"`,
				`"id":"4q1xlclmfloku33"`,
				`"email":"test@example.com"`,
			},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeAuthWithPasswordRequest": 1,
				"OnRecordAfterAuthWithPasswordRequest":  1,
				"OnRecordAuthRequest":                   1,
				// lastLoginAlertSentAt update
				"OnModelAfterUpdate":  1,
				"OnModelBeforeUpdate": 1,
			},
		},
		{
			Name:   "authenticated admin",
			Method: http.MethodPost,
			Url:    "/api/collections/users/auth-with-password",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			Body: strings.NewReader(`{
				"identity":"test@example.com",
				"password":"1234567890"
			}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"record":{`,
				`"token":"`,
				`"id":"4q1xlclmfloku33"`,
				`"email":"test@example.com"`,
			},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeAuthWithPasswordRequest": 1,
				"OnRecordAfterAuthWithPasswordRequest":  1,
				"OnRecordAuthRequest":                   1,
				// lastLoginAlertSentAt update
				"OnModelAfterUpdate":  1,
				"OnModelBeforeUpdate": 1,
			},
		},

		// after hooks error checks
		{
			Name:   "OnRecordAfterAuthWithPasswordRequest error response",
			Method: http.MethodPost,
			Url:    "/api/collections/users/auth-with-password",
			Body: strings.NewReader(`{
				"identity":"test2_username",
				"password":"1234567890"
			}`),
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				app.OnRecordAfterAuthWithPasswordRequest().Add(func(e *core.RecordAuthWithPasswordEvent) error {
					return errors.New("error")
				})
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeAuthWithPasswordRequest": 1,
				"OnRecordAfterAuthWithPasswordRequest":  1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRecordAuthRefresh(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodPost,
			Url:             "/api/collections/users/auth-refresh",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "admin",
			Method: http.MethodPost,
			Url:    "/api/collections/users/auth-refresh",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "auth record + not an auth collection",
			Method: http.MethodPost,
			Url:    "/api/collections/demo1/auth-refresh",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "auth record + different auth collection",
			Method: http.MethodPost,
			Url:    "/api/collections/clients/auth-refresh?expand=rel,missing",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "auth record + same auth collection as the token",
			Method: http.MethodPost,
			Url:    "/api/collections/users/auth-refresh?expand=rel,missing",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"token":`,
				`"record":`,
				`"id":"4q1xlclmfloku33"`,
				`"emailVisibility":false`,
				`"email":"test@example.com"`, // the owner can always view their email address
				`"expand":`,
				`"rel":`,
				`"id":"llvuca81nly1qls"`,
			},
			NotExpectedContent: []string{
				`"missing":`,
			},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeAuthRefreshRequest": 1,
				"OnRecordAuthRequest":              1,
				"OnRecordAfterAuthRefreshRequest":  1,
			},
		},
		{
			Name:   "unverified auth record in onlyVerified collection",
			Method: http.MethodPost,
			Url:    "/api/collections/clients/auth-refresh",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6Im8xeTBkZDBzcGQ3ODZtZCIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoidjg1MXE0cjc5MHJoa25sIiwiZXhwIjoyMjA4OTg1MjYxfQ.-JYlrz5DcGzvb0nYx-xqnSFMu9dupyKY7Vg_FUm0OaM",
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeAuthRefreshRequest": 1,
				"OnRecordAfterAuthRefreshRequest":  1,
			},
		},
		{
			Name:   "verified auth record in onlyVerified collection",
			Method: http.MethodPost,
			Url:    "/api/collections/clients/auth-refresh",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6ImdrMzkwcWVnczR5NDd3biIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoidjg1MXE0cjc5MHJoa25sIiwiZXhwIjoyMjA4OTg1MjYxfQ.q34IWXrRWsjLvbbVNRfAs_J4SoTHloNBfdGEiLmy-D8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"token":`,
				`"record":`,
				`"id":"gk390qegs4y47wn"`,
				`"verified":true`,
				`"email":"test@example.com"`,
			},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeAuthRefreshRequest": 1,
				"OnRecordAuthRequest":              1,
				"OnRecordAfterAuthRefreshRequest":  1,
			},
		},
		{
			Name:   "OnRecordAfterAuthRefreshRequest error response",
			Method: http.MethodPost,
			Url:    "/api/collections/users/auth-refresh?expand=rel,missing",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				app.OnRecordAfterAuthRefreshRequest().Add(func(e *core.RecordAuthRefreshEvent) error {
					return errors.New("error")
				})
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"OnRecordBeforeAuthRefreshRequest": 1,
				"OnRecordAfterAuthRefreshRequest":  1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRecordAuthRequestPasswordReset(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "not an auth collection",
			Method:          http.MethodPost,
			Url:             "/api/collections/demo1/request-password-reset",
			Body:            strings.NewReader(``),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "empty data",
			Method:          http.MethodPost,
			Url:             "/api/collections/users/request-password-reset",
			Body:            strings.NewReader(``),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{"email":{"code":"validation_required","message":"Cannot be blank."}}`},
		},
		{
			Name:            "invalid data",
			Method:          http.MethodPost,
			Url:             "/api/collections/users/request-password-reset",
			Body:            strings.NewReader(`{"email`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:           "missing auth record",
			Method:         http.MethodPost,
			Url:            "/api/collections/users/request-password-reset",
			Body:           strings.NewReader(`{"email":"missing@example.com"}`),
			Delay:          100 * time.Millisecond,
			ExpectedStatus: 204,
		},
		{
			Name:           "existing auth record",
			Method:         http.MethodPost,
			Url:            "/api/collections/users/request-password-reset",
			Body:           strings.NewReader(`{"email":"test@example.com"}`),
			Delay:          100 * time.Millisecond,
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnModelBeforeUpdate":                       1,
				"OnModelAfterUpdate":                        1,
				"OnRecordBeforeRequestPasswordResetRequest": 1,
				"OnRecordAfterRequestPasswordResetRequest":  1,
				"OnMailerBeforeRecordResetPasswordSend":     1,
				"OnMailerAfterRecordResetPasswordSend":      1,
			},
		},
		{
			Name:           "existing auth record (after already sent)",
			Method:         http.MethodPost,
			Url:            "/api/collections/clients/request-password-reset",
			Body:           strings.NewReader(`{"email":"test@example.com"}`),
			Delay:          100 * time.Millisecond,
			ExpectedStatus: 204,
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				// simulate recent password request sent
				authRecord, err := app.Dao().FindFirstRecordByData("clients", "email", "test@example.com")
				if err != nil {
					t.Fatal(err)
				}
				authRecord.SetLastResetSentAt(types.NowDateTime())
				dao := daos.New(app.Dao().DB()) // new dao to ignore hooks
				if err := dao.Save(authRecord); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			Name:            "existing auth record in a collection with disabled password login",
			Method:          http.MethodPost,
			Url:             "/api/collections/nologin/request-password-reset",
			Body:            strings.NewReader(`{"email":"test@example.com"}`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRecordAuthConfirmPasswordReset(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:           "empty data",
			Method:         http.MethodPost,
			Url:            "/api/collections/users/confirm-password-reset",
			Body:           strings.NewReader(``),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"password":{"code":"validation_required"`,
				`"passwordConfirm":{"code":"validation_required"`,
				`"token":{"code":"validation_required"`,
			},
		},
		{
			Name:            "invalid data format",
			Method:          http.MethodPost,
			Url:             "/api/collections/users/confirm-password-reset",
			Body:            strings.NewReader(`{"password`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "expired token and invalid password",
			Method: http.MethodPost,
			Url:    "/api/collections/users/confirm-password-reset",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsImNvbGxlY3Rpb25JZCI6Il9wYl91c2Vyc19hdXRoXyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiZXhwIjoxNjQwOTkxNjYxfQ.TayHoXkOTM0w8InkBEb86npMJEaf6YVUrxrRmMgFjeY",
				"password":"1234567",
				"passwordConfirm":"7654321"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"token":{"code":"validation_invalid_token"`,
				`"password":{"code":"validation_length_out_of_range"`,
				`"passwordConfirm":{"code":"validation_values_mismatch"`,
			},
		},
		{
			Name:   "non auth collection",
			Method: http.MethodPost,
			Url:    "/api/collections/demo1/confirm-password-reset?expand=rel,missing",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsImNvbGxlY3Rpb25JZCI6Il9wYl91c2Vyc19hdXRoXyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiZXhwIjoyMjA4OTg1MjYxfQ.R_4FOSUHIuJQ5Crl3PpIPCXMsoHzuTaNlccpXg_3FOg",
				"password":"12345678",
				"passwordConfirm":"12345678"
			}`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "different auth collection",
			Method: http.MethodPost,
			Url:    "/api/collections/clients/confirm-password-reset?expand=rel,missing",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsImNvbGxlY3Rpb25JZCI6Il9wYl91c2Vyc19hdXRoXyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiZXhwIjoyMjA4OTg1MjYxfQ.R_4FOSUHIuJQ5Crl3PpIPCXMsoHzuTaNlccpXg_3FOg",
				"password":"12345678",
				"passwordConfirm":"12345678"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{"token":{"code":"validation_token_collection_mismatch"`,
			},
		},
		{
			Name:   "valid token and data (unverified user)",
			Method: http.MethodPost,
			Url:    "/api/collections/users/confirm-password-reset",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsImNvbGxlY3Rpb25JZCI6Il9wYl91c2Vyc19hdXRoXyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiZXhwIjoyMjA4OTg1MjYxfQ.R_4FOSUHIuJQ5Crl3PpIPCXMsoHzuTaNlccpXg_3FOg",
				"password":"12345678",
				"passwordConfirm":"12345678"
			}`),
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnModelAfterUpdate":                        1,
				"OnModelBeforeUpdate":                       1,
				"OnRecordBeforeConfirmPasswordResetRequest": 1,
				"OnRecordAfterConfirmPasswordResetRequest":  1,
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				user, err := app.Dao().FindAuthRecordByEmail("users", "test@example.com")
				if err != nil {
					t.Fatalf("Failed to fetch confirm password user: %v", err)
				}

				if user.Verified() {
					t.Fatalf("Expected the user to be unverified")
				}
			},
			AfterTestFunc: func(t *testing.T, app *tests.TestApp, res *http.Response) {
				user, err := app.Dao().FindAuthRecordByToken(
					"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsImNvbGxlY3Rpb25JZCI6Il9wYl91c2Vyc19hdXRoXyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiZXhwIjoyMjA4OTg1MjYxfQ.R_4FOSUHIuJQ5Crl3PpIPCXMsoHzuTaNlccpXg_3FOg",
					app.Settings().RecordPasswordResetToken.Secret,
				)
				if err == nil {
					t.Fatalf("Expected the password reset token to be invalidated")
				}

				user, err = app.Dao().FindAuthRecordByEmail("users", "test@example.com")
				if err != nil {
					t.Fatalf("Failed to fetch confirm password user: %v", err)
				}

				if !user.Verified() {
					t.Fatalf("Expected the user to be marked as verified")
				}
			},
		},
		{
			Name:   "valid token and data (unverified user with different email from the one in the token)",
			Method: http.MethodPost,
			Url:    "/api/collections/users/confirm-password-reset",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsImNvbGxlY3Rpb25JZCI6Il9wYl91c2Vyc19hdXRoXyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiZXhwIjoyMjA4OTg1MjYxfQ.R_4FOSUHIuJQ5Crl3PpIPCXMsoHzuTaNlccpXg_3FOg",
				"password":"12345678",
				"passwordConfirm":"12345678"
			}`),
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnModelAfterUpdate":                        1,
				"OnModelBeforeUpdate":                       1,
				"OnRecordBeforeConfirmPasswordResetRequest": 1,
				"OnRecordAfterConfirmPasswordResetRequest":  1,
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				user, err := app.Dao().FindAuthRecordByEmail("users", "test@example.com")
				if err != nil {
					t.Fatalf("Failed to fetch confirm password user: %v", err)
				}

				if user.Verified() {
					t.Fatalf("Expected the user to be unverified")
				}

				// manually change the email to check whether the verified state will be updated
				user.SetEmail("test_update@example.com")
				if err := app.Dao().WithoutHooks().SaveRecord(user); err != nil {
					t.Fatalf("Failed to update user test email")
				}
			},
			AfterTestFunc: func(t *testing.T, app *tests.TestApp, res *http.Response) {
				user, err := app.Dao().FindAuthRecordByToken(
					"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsImNvbGxlY3Rpb25JZCI6Il9wYl91c2Vyc19hdXRoXyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiZXhwIjoyMjA4OTg1MjYxfQ.R_4FOSUHIuJQ5Crl3PpIPCXMsoHzuTaNlccpXg_3FOg",
					app.Settings().RecordPasswordResetToken.Secret,
				)
				if err == nil {
					t.Fatalf("Expected the password reset token to be invalidated")
				}

				user, err = app.Dao().FindAuthRecordByEmail("users", "test_update@example.com")
				if err != nil {
					t.Fatalf("Failed to fetch confirm password user: %v", err)
				}

				if user.Verified() {
					t.Fatalf("Expected the user to remain unverified")
				}
			},
		},
		{
			Name:   "valid token and data (verified user)",
			Method: http.MethodPost,
			Url:    "/api/collections/users/confirm-password-reset",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsImNvbGxlY3Rpb25JZCI6Il9wYl91c2Vyc19hdXRoXyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiZXhwIjoyMjA4OTg1MjYxfQ.R_4FOSUHIuJQ5Crl3PpIPCXMsoHzuTaNlccpXg_3FOg",
				"password":"12345678",
				"passwordConfirm":"12345678"
			}`),
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnModelAfterUpdate":                        1,
				"OnModelBeforeUpdate":                       1,
				"OnRecordBeforeConfirmPasswordResetRequest": 1,
				"OnRecordAfterConfirmPasswordResetRequest":  1,
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				user, err := app.Dao().FindAuthRecordByEmail("users", "test@example.com")
				if err != nil {
					t.Fatalf("Failed to fetch confirm password user: %v", err)
				}

				// ensure that the user is already verified
				user.SetVerified(true)
				if err := app.Dao().WithoutHooks().SaveRecord(user); err != nil {
					t.Fatalf("Failed to update user verified state")
				}
			},
			AfterTestFunc: func(t *testing.T, app *tests.TestApp, res *http.Response) {
				user, err := app.Dao().FindAuthRecordByToken(
					"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsImNvbGxlY3Rpb25JZCI6Il9wYl91c2Vyc19hdXRoXyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiZXhwIjoyMjA4OTg1MjYxfQ.R_4FOSUHIuJQ5Crl3PpIPCXMsoHzuTaNlccpXg_3FOg",
					app.Settings().RecordPasswordResetToken.Secret,
				)
				if err == nil {
					t.Fatalf("Expected the password reset token to be invalidated")
				}

				user, err = app.Dao().FindAuthRecordByEmail("users", "test@example.com")
				if err != nil {
					t.Fatalf("Failed to fetch confirm password user: %v", err)
				}

				if !user.Verified() {
					t.Fatalf("Expected the user to remain verified")
				}
			},
		},
		{
			Name:   "OnRecordAfterConfirmPasswordResetRequest error response",
			Method: http.MethodPost,
			Url:    "/api/collections/users/confirm-password-reset",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsImNvbGxlY3Rpb25JZCI6Il9wYl91c2Vyc19hdXRoXyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiZXhwIjoyMjA4OTg1MjYxfQ.R_4FOSUHIuJQ5Crl3PpIPCXMsoHzuTaNlccpXg_3FOg",
				"password":"12345678",
				"passwordConfirm":"12345678"
			}`),
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				app.OnRecordAfterConfirmPasswordResetRequest().Add(func(e *core.RecordConfirmPasswordResetEvent) error {
					return errors.New("error")
				})
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"OnModelAfterUpdate":                        1,
				"OnModelBeforeUpdate":                       1,
				"OnRecordBeforeConfirmPasswordResetRequest": 1,
				"OnRecordAfterConfirmPasswordResetRequest":  1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRecordAuthRequestVerification(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "not an auth collection",
			Method:          http.MethodPost,
			Url:             "/api/collections/demo1/request-verification",
			Body:            strings.NewReader(``),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "empty data",
			Method:          http.MethodPost,
			Url:             "/api/collections/users/request-verification",
			Body:            strings.NewReader(``),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{"email":{"code":"validation_required","message":"Cannot be blank."}}`},
		},
		{
			Name:            "invalid data",
			Method:          http.MethodPost,
			Url:             "/api/collections/users/request-verification",
			Body:            strings.NewReader(`{"email`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:           "missing auth record",
			Method:         http.MethodPost,
			Url:            "/api/collections/users/request-verification",
			Body:           strings.NewReader(`{"email":"missing@example.com"}`),
			Delay:          100 * time.Millisecond,
			ExpectedStatus: 204,
		},
		{
			Name:           "already verified auth record",
			Method:         http.MethodPost,
			Url:            "/api/collections/users/request-verification",
			Body:           strings.NewReader(`{"email":"test2@example.com"}`),
			Delay:          100 * time.Millisecond,
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnRecordBeforeRequestVerificationRequest": 1,
				"OnRecordAfterRequestVerificationRequest":  1,
			},
		},
		{
			Name:           "existing auth record",
			Method:         http.MethodPost,
			Url:            "/api/collections/users/request-verification",
			Body:           strings.NewReader(`{"email":"test@example.com"}`),
			Delay:          100 * time.Millisecond,
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnModelBeforeUpdate":                      1,
				"OnModelAfterUpdate":                       1,
				"OnRecordBeforeRequestVerificationRequest": 1,
				"OnRecordAfterRequestVerificationRequest":  1,
				"OnMailerBeforeRecordVerificationSend":     1,
				"OnMailerAfterRecordVerificationSend":      1,
			},
		},
		{
			Name:           "existing auth record (after already sent)",
			Method:         http.MethodPost,
			Url:            "/api/collections/users/request-verification",
			Body:           strings.NewReader(`{"email":"test@example.com"}`),
			Delay:          100 * time.Millisecond,
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				// "OnRecordBeforeRequestVerificationRequest": 1,
				// "OnRecordAfterRequestVerificationRequest":  1,
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				// simulate recent verification sent
				authRecord, err := app.Dao().FindFirstRecordByData("users", "email", "test@example.com")
				if err != nil {
					t.Fatal(err)
				}
				authRecord.SetLastVerificationSentAt(types.NowDateTime())
				dao := daos.New(app.Dao().DB()) // new dao to ignore hooks
				if err := dao.Save(authRecord); err != nil {
					t.Fatal(err)
				}
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRecordAuthConfirmVerification(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:           "empty data",
			Method:         http.MethodPost,
			Url:            "/api/collections/users/confirm-verification",
			Body:           strings.NewReader(``),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"token":{"code":"validation_required"`,
			},
		},
		{
			Name:            "invalid data format",
			Method:          http.MethodPost,
			Url:             "/api/collections/users/confirm-verification",
			Body:            strings.NewReader(`{"password`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "expired token",
			Method: http.MethodPost,
			Url:    "/api/collections/users/confirm-verification",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsImNvbGxlY3Rpb25JZCI6Il9wYl91c2Vyc19hdXRoXyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiZXhwIjoxNjQwOTkxNjYxfQ.Avbt9IP8sBisVz_2AGrlxLDvangVq4PhL2zqQVYLKlE"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"token":{"code":"validation_invalid_token"`,
			},
		},
		{
			Name:   "non auth collection",
			Method: http.MethodPost,
			Url:    "/api/collections/demo1/confirm-verification?expand=rel,missing",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsImNvbGxlY3Rpb25JZCI6Il9wYl91c2Vyc19hdXRoXyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiZXhwIjoyMjA4OTg1MjYxfQ.R_4FOSUHIuJQ5Crl3PpIPCXMsoHzuTaNlccpXg_3FOg"
			}`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "different auth collection",
			Method: http.MethodPost,
			Url:    "/api/collections/clients/confirm-verification?expand=rel,missing",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsImNvbGxlY3Rpb25JZCI6Il9wYl91c2Vyc19hdXRoXyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiZXhwIjoyMjA4OTg1MjYxfQ.hL16TVmStHFdHLc4a860bRqJ3sFfzjv0_NRNzwsvsrc"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{"token":{"code":"validation_token_collection_mismatch"`,
			},
		},
		{
			Name:   "valid token",
			Method: http.MethodPost,
			Url:    "/api/collections/users/confirm-verification",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsImNvbGxlY3Rpb25JZCI6Il9wYl91c2Vyc19hdXRoXyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiZXhwIjoyMjA4OTg1MjYxfQ.hL16TVmStHFdHLc4a860bRqJ3sFfzjv0_NRNzwsvsrc"
			}`),
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnModelAfterUpdate":                       1,
				"OnModelBeforeUpdate":                      1,
				"OnRecordBeforeConfirmVerificationRequest": 1,
				"OnRecordAfterConfirmVerificationRequest":  1,
			},
		},
		{
			Name:   "valid token (already verified)",
			Method: http.MethodPost,
			Url:    "/api/collections/users/confirm-verification",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6Im9hcDY0MGNvdDR5cnUycyIsImVtYWlsIjoidGVzdDJAZXhhbXBsZS5jb20iLCJjb2xsZWN0aW9uSWQiOiJfcGJfdXNlcnNfYXV0aF8iLCJ0eXBlIjoiYXV0aFJlY29yZCIsImV4cCI6MjIwODk4NTI2MX0.PsOABmYUzGbd088g8iIBL4-pf7DUZm0W5Ju6lL5JVRg"
			}`),
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnRecordBeforeConfirmVerificationRequest": 1,
				"OnRecordAfterConfirmVerificationRequest":  1,
			},
		},
		{
			Name:   "valid verification token from a collection without allowed login",
			Method: http.MethodPost,
			Url:    "/api/collections/nologin/confirm-verification",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImRjNDlrNmpnZWpuNDBoMyIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsImNvbGxlY3Rpb25JZCI6ImtwdjcwOXNrMmxxYnFrOCIsInR5cGUiOiJhdXRoUmVjb3JkIiwiZXhwIjoyMjA4OTg1MjYxfQ.coREjeTDS3_Go7DP1nxHtevIX5rujwHU-_mRB6oOm3w"
			}`),
			ExpectedStatus:  204,
			ExpectedContent: []string{},
			ExpectedEvents: map[string]int{
				"OnModelAfterUpdate":                       1,
				"OnModelBeforeUpdate":                      1,
				"OnRecordBeforeConfirmVerificationRequest": 1,
				"OnRecordAfterConfirmVerificationRequest":  1,
			},
		},
		{
			Name:   "OnRecordAfterConfirmVerificationRequest error response",
			Method: http.MethodPost,
			Url:    "/api/collections/users/confirm-verification",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsImNvbGxlY3Rpb25JZCI6Il9wYl91c2Vyc19hdXRoXyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiZXhwIjoyMjA4OTg1MjYxfQ.hL16TVmStHFdHLc4a860bRqJ3sFfzjv0_NRNzwsvsrc"
			}`),
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				app.OnRecordAfterConfirmVerificationRequest().Add(func(e *core.RecordConfirmVerificationEvent) error {
					return errors.New("error")
				})
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"OnModelAfterUpdate":                       1,
				"OnModelBeforeUpdate":                      1,
				"OnRecordBeforeConfirmVerificationRequest": 1,
				"OnRecordAfterConfirmVerificationRequest":  1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRecordAuthRequestEmailChange(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodPost,
			Url:             "/api/collections/users/request-email-change",
			Body:            strings.NewReader(`{"newEmail":"change@example.com"}`),
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "not an auth collection",
			Method:          http.MethodPost,
			Url:             "/api/collections/demo1/request-email-change",
			Body:            strings.NewReader(`{"newEmail":"change@example.com"}`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "admin authentication",
			Method: http.MethodPost,
			Url:    "/api/collections/users/request-email-change",
			Body:   strings.NewReader(`{"newEmail":"change@example.com"}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "record authentication but from different auth collection",
			Method: http.MethodPost,
			Url:    "/api/collections/clients/request-email-change",
			Body:   strings.NewReader(`{"newEmail":"change@example.com"}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "invalid data",
			Method: http.MethodPost,
			Url:    "/api/collections/users/request-email-change",
			Body:   strings.NewReader(`{"newEmail`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "empty data",
			Method: http.MethodPost,
			Url:    "/api/collections/users/request-email-change",
			Body:   strings.NewReader(`{}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":`,
				`"newEmail":{"code":"validation_required"`,
			},
		},
		{
			Name:   "valid data (existing email)",
			Method: http.MethodPost,
			Url:    "/api/collections/users/request-email-change",
			Body:   strings.NewReader(`{"newEmail":"test2@example.com"}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":`,
				`"newEmail":{"code":"validation_record_email_invalid"`,
			},
		},
		{
			Name:   "valid data (new email)",
			Method: http.MethodPost,
			Url:    "/api/collections/users/request-email-change",
			Body:   strings.NewReader(`{"newEmail":"change@example.com"}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnMailerBeforeRecordChangeEmailSend":     1,
				"OnMailerAfterRecordChangeEmailSend":      1,
				"OnRecordBeforeRequestEmailChangeRequest": 1,
				"OnRecordAfterRequestEmailChangeRequest":  1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRecordAuthConfirmEmailChange(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:           "not an auth collection",
			Method:         http.MethodPost,
			Url:            "/api/collections/demo1/confirm-email-change",
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{}`,
			},
		},
		{
			Name:           "empty data",
			Method:         http.MethodPost,
			Url:            "/api/collections/users/confirm-email-change",
			Body:           strings.NewReader(``),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":`,
				`"token":{"code":"validation_required"`,
				`"password":{"code":"validation_required"`,
			},
		},
		{
			Name:            "invalid data",
			Method:          http.MethodPost,
			Url:             "/api/collections/users/confirm-email-change",
			Body:            strings.NewReader(`{"token`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "expired token and correct password",
			Method: http.MethodPost,
			Url:    "/api/collections/users/confirm-email-change",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImNvbGxlY3Rpb25JZCI6Il9wYl91c2Vyc19hdXRoXyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwibmV3RW1haWwiOiJjaGFuZ2VAZXhhbXBsZS5jb20iLCJleHAiOjE2NDA5OTE2NjF9.D20jh5Ss7SZyXRUXjjEyLCYo9Ky0N5cE5dKB_MGJ8G8",
				"password":"1234567890"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"token":{`,
				`"code":"validation_invalid_token"`,
			},
		},
		{
			Name:   "valid token and incorrect password",
			Method: http.MethodPost,
			Url:    "/api/collections/users/confirm-email-change",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImNvbGxlY3Rpb25JZCI6Il9wYl91c2Vyc19hdXRoXyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwibmV3RW1haWwiOiJjaGFuZ2VAZXhhbXBsZS5jb20iLCJleHAiOjIyMDg5ODUyNjF9.1sG6cL708pRXXjiHRZhG-in0X5fnttSf5nNcadKoYRs",
				"password":"1234567891"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"password":{`,
				`"code":"validation_invalid_password"`,
			},
		},
		{
			Name:   "valid token and correct password",
			Method: http.MethodPost,
			Url:    "/api/collections/users/confirm-email-change",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImNvbGxlY3Rpb25JZCI6Il9wYl91c2Vyc19hdXRoXyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwibmV3RW1haWwiOiJjaGFuZ2VAZXhhbXBsZS5jb20iLCJleHAiOjIyMDg5ODUyNjF9.1sG6cL708pRXXjiHRZhG-in0X5fnttSf5nNcadKoYRs",
				"password":"1234567890"
			}`),
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnModelAfterUpdate":                      1,
				"OnModelBeforeUpdate":                     1,
				"OnRecordBeforeConfirmEmailChangeRequest": 1,
				"OnRecordAfterConfirmEmailChangeRequest":  1,
			},
		},
		{
			Name:   "valid token and correct password in different auth collection",
			Method: http.MethodPost,
			Url:    "/api/collections/clients/confirm-email-change",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImNvbGxlY3Rpb25JZCI6Il9wYl91c2Vyc19hdXRoXyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwibmV3RW1haWwiOiJjaGFuZ2VAZXhhbXBsZS5jb20iLCJleHAiOjIyMDg5ODUyNjF9.1sG6cL708pRXXjiHRZhG-in0X5fnttSf5nNcadKoYRs",
				"password":"1234567890"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"token":{"code":"validation_token_collection_mismatch"`,
			},
		},
		{
			Name:   "OnRecordAfterConfirmEmailChangeRequest error response",
			Method: http.MethodPost,
			Url:    "/api/collections/users/confirm-email-change",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImNvbGxlY3Rpb25JZCI6Il9wYl91c2Vyc19hdXRoXyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwibmV3RW1haWwiOiJjaGFuZ2VAZXhhbXBsZS5jb20iLCJleHAiOjIyMDg5ODUyNjF9.1sG6cL708pRXXjiHRZhG-in0X5fnttSf5nNcadKoYRs",
				"password":"1234567890"
			}`),
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				app.OnRecordAfterConfirmEmailChangeRequest().Add(func(e *core.RecordConfirmEmailChangeEvent) error {
					return errors.New("error")
				})
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"OnModelAfterUpdate":                      1,
				"OnModelBeforeUpdate":                     1,
				"OnRecordBeforeConfirmEmailChangeRequest": 1,
				"OnRecordAfterConfirmEmailChangeRequest":  1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRecordAuthListExternalsAuths(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodGet,
			Url:             "/api/collections/users/records/4q1xlclmfloku33/external-auths",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "admin + nonexisting record id",
			Method: http.MethodGet,
			Url:    "/api/collections/users/records/missing/external-auths",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "admin + existing record id and no external auths",
			Method: http.MethodGet,
			Url:    "/api/collections/users/records/oap640cot4yru2s/external-auths",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus:  200,
			ExpectedContent: []string{`[]`},
			ExpectedEvents:  map[string]int{"OnRecordListExternalAuthsRequest": 1},
		},
		{
			Name:   "admin + existing user id and 2 external auths",
			Method: http.MethodGet,
			Url:    "/api/collections/users/records/4q1xlclmfloku33/external-auths",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"clmflokuq1xl341"`,
				`"id":"dlmflokuq1xl342"`,
				`"recordId":"4q1xlclmfloku33"`,
				`"collectionId":"_pb_users_auth_"`,
			},
			ExpectedEvents: map[string]int{"OnRecordListExternalAuthsRequest": 1},
		},
		{
			Name:   "auth record + trying to list another user external auths",
			Method: http.MethodGet,
			Url:    "/api/collections/users/records/4q1xlclmfloku33/external-auths",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6Im9hcDY0MGNvdDR5cnUycyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.uatnTBFqMnF0p4FkmwEpA9R-uGFu0Putwyk6NJCKBno",
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "auth record + trying to list another user external auths from different collection",
			Method: http.MethodGet,
			Url:    "/api/collections/clients/records/o1y0dd0spd786md/external-auths",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6Im9hcDY0MGNvdDR5cnUycyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.uatnTBFqMnF0p4FkmwEpA9R-uGFu0Putwyk6NJCKBno",
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "auth record + owner without external auths",
			Method: http.MethodGet,
			Url:    "/api/collections/users/records/oap640cot4yru2s/external-auths",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6Im9hcDY0MGNvdDR5cnUycyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.uatnTBFqMnF0p4FkmwEpA9R-uGFu0Putwyk6NJCKBno",
			},
			ExpectedStatus:  200,
			ExpectedContent: []string{`[]`},
			ExpectedEvents:  map[string]int{"OnRecordListExternalAuthsRequest": 1},
		},
		{
			Name:   "authorized as user - owner with 2 external auths",
			Method: http.MethodGet,
			Url:    "/api/collections/users/records/4q1xlclmfloku33/external-auths",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"clmflokuq1xl341"`,
				`"id":"dlmflokuq1xl342"`,
				`"recordId":"4q1xlclmfloku33"`,
				`"collectionId":"_pb_users_auth_"`,
			},
			ExpectedEvents: map[string]int{"OnRecordListExternalAuthsRequest": 1},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRecordAuthUnlinkExternalsAuth(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodDelete,
			Url:             "/api/collections/users/records/4q1xlclmfloku33/external-auths/google",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "admin - nonexisting recod id",
			Method: http.MethodDelete,
			Url:    "/api/collections/users/records/missing/external-auths/google",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "admin - nonlinked provider",
			Method: http.MethodDelete,
			Url:    "/api/collections/users/records/4q1xlclmfloku33/external-auths/facebook",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "admin - linked provider",
			Method: http.MethodDelete,
			Url:    "/api/collections/users/records/4q1xlclmfloku33/external-auths/google",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus:  204,
			ExpectedContent: []string{},
			ExpectedEvents: map[string]int{
				"OnModelAfterDelete":                      1,
				"OnModelBeforeDelete":                     1,
				"OnRecordAfterUnlinkExternalAuthRequest":  1,
				"OnRecordBeforeUnlinkExternalAuthRequest": 1,
			},
			AfterTestFunc: func(t *testing.T, app *tests.TestApp, res *http.Response) {
				record, err := app.Dao().FindRecordById("users", "4q1xlclmfloku33")
				if err != nil {
					t.Fatal(err)
				}
				auth, _ := app.Dao().FindExternalAuthByRecordAndProvider(record, "google")
				if auth != nil {
					t.Fatalf("Expected the google ExternalAuth to be deleted, got got \n%v", auth)
				}
			},
		},
		{
			Name:   "auth record - trying to unlink another user external auth",
			Method: http.MethodDelete,
			Url:    "/api/collections/users/records/4q1xlclmfloku33/external-auths/google",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6Im9hcDY0MGNvdDR5cnUycyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.uatnTBFqMnF0p4FkmwEpA9R-uGFu0Putwyk6NJCKBno",
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "auth record - trying to unlink another user external auth from different collection",
			Method: http.MethodDelete,
			Url:    "/api/collections/clients/records/o1y0dd0spd786md/external-auths/google",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "auth record - owner with existing external auth",
			Method: http.MethodDelete,
			Url:    "/api/collections/users/records/4q1xlclmfloku33/external-auths/google",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus:  204,
			ExpectedContent: []string{},
			ExpectedEvents: map[string]int{
				"OnModelAfterDelete":                      1,
				"OnModelBeforeDelete":                     1,
				"OnRecordAfterUnlinkExternalAuthRequest":  1,
				"OnRecordBeforeUnlinkExternalAuthRequest": 1,
			},
			AfterTestFunc: func(t *testing.T, app *tests.TestApp, res *http.Response) {
				record, err := app.Dao().FindRecordById("users", "4q1xlclmfloku33")
				if err != nil {
					t.Fatal(err)
				}
				auth, _ := app.Dao().FindExternalAuthByRecordAndProvider(record, "google")
				if auth != nil {
					t.Fatalf("Expected the google ExternalAuth to be deleted, got got \n%v", auth)
				}
			},
		},
		{
			Name:   "OnRecordBeforeUnlinkExternalAuthRequest error response",
			Method: http.MethodDelete,
			Url:    "/api/collections/users/records/4q1xlclmfloku33/external-auths/google",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				app.OnRecordAfterUnlinkExternalAuthRequest().Add(func(e *core.RecordUnlinkExternalAuthEvent) error {
					return errors.New("error")
				})
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"OnModelAfterDelete":                      1,
				"OnModelBeforeDelete":                     1,
				"OnRecordAfterUnlinkExternalAuthRequest":  1,
				"OnRecordBeforeUnlinkExternalAuthRequest": 1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRecordAuthOAuth2Redirect(t *testing.T) {
	t.Parallel()

	clientStubs := make([]map[string]subscriptions.Client, 0, 10)

	for i := 0; i < 10; i++ {
		c1 := subscriptions.NewDefaultClient()

		c2 := subscriptions.NewDefaultClient()
		c2.Subscribe("@oauth2")

		c3 := subscriptions.NewDefaultClient()
		c3.Subscribe("test1", "@oauth2")

		c4 := subscriptions.NewDefaultClient()
		c4.Subscribe("test1", "test2")

		c5 := subscriptions.NewDefaultClient()
		c5.Subscribe("@oauth2")
		c5.Discard()

		clientStubs = append(clientStubs, map[string]subscriptions.Client{
			"c1": c1,
			"c2": c2,
			"c3": c3,
			"c4": c4,
			"c5": c5,
		})
	}

	checkFailureRedirect := func(t *testing.T, app *tests.TestApp, res *http.Response) {
		loc := res.Header.Get("Location")
		if !strings.Contains(loc, "/oauth2-redirect-failure") {
			t.Fatalf("Expected failure redirect, got %q", loc)
		}
	}

	checkSuccessRedirect := func(t *testing.T, app *tests.TestApp, res *http.Response) {
		loc := res.Header.Get("Location")
		if !strings.Contains(loc, "/oauth2-redirect-success") {
			t.Fatalf("Expected success redirect, got %q", loc)
		}
	}

	checkClientMessages := func(t *testing.T, clientId string, msg subscriptions.Message, expectedMessages map[string][]string) {
		if len(expectedMessages[clientId]) == 0 {
			t.Fatalf("Unexpected client %q message, got %s:\n%s", clientId, msg.Name, msg.Data)
		}

		if msg.Name != "@oauth2" {
			t.Fatalf("Expected @oauth2 msg.Name, got %q", msg.Name)
		}

		for _, txt := range expectedMessages[clientId] {
			if !strings.Contains(string(msg.Data), txt) {
				t.Fatalf("Failed to find %q in \n%s", txt, msg.Data)
			}
		}
	}

	beforeTestFunc := func(
		clients map[string]subscriptions.Client,
		expectedMessages map[string][]string,
	) func(*testing.T, *tests.TestApp, *echo.Echo) {
		return func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
			for _, client := range clients {
				app.SubscriptionsBroker().Register(client)
			}

			ctx, cancelFunc := context.WithTimeout(context.Background(), 100*time.Millisecond)

			// add to the app store so that it can be cancelled manually after test completion
			app.Store().Set("cancelFunc", cancelFunc)

			go func() {
				defer cancelFunc()

				for {
					select {
					case msg := <-clients["c1"].Channel():
						checkClientMessages(t, "c1", msg, expectedMessages)
					case msg := <-clients["c2"].Channel():
						checkClientMessages(t, "c2", msg, expectedMessages)
					case msg := <-clients["c3"].Channel():
						checkClientMessages(t, "c3", msg, expectedMessages)
					case msg := <-clients["c4"].Channel():
						checkClientMessages(t, "c4", msg, expectedMessages)
					case msg := <-clients["c5"].Channel():
						checkClientMessages(t, "c5", msg, expectedMessages)
					case <-ctx.Done():
						for _, c := range clients {
							close(c.Channel())
						}
						return
					}
				}
			}()
		}
	}

	scenarios := []tests.ApiScenario{
		{
			Name:           "no state query param",
			Method:         http.MethodGet,
			Url:            "/api/oauth2-redirect?code=123",
			BeforeTestFunc: beforeTestFunc(clientStubs[0], nil),
			ExpectedStatus: http.StatusTemporaryRedirect,
			AfterTestFunc: func(t *testing.T, app *tests.TestApp, res *http.Response) {
				app.Store().Get("cancelFunc").(context.CancelFunc)()

				checkFailureRedirect(t, app, res)
			},
		},
		{
			Name:           "invalid or missing client",
			Method:         http.MethodGet,
			Url:            "/api/oauth2-redirect?code=123&state=missing",
			BeforeTestFunc: beforeTestFunc(clientStubs[1], nil),
			ExpectedStatus: http.StatusTemporaryRedirect,
			AfterTestFunc: func(t *testing.T, app *tests.TestApp, res *http.Response) {
				app.Store().Get("cancelFunc").(context.CancelFunc)()

				checkFailureRedirect(t, app, res)
			},
		},
		{
			Name:   "no code query param",
			Method: http.MethodGet,
			Url:    "/api/oauth2-redirect?state=" + clientStubs[2]["c3"].Id(),
			BeforeTestFunc: beforeTestFunc(clientStubs[2], map[string][]string{
				"c3": {`"state":"` + clientStubs[2]["c3"].Id(), `"code":""`},
			}),
			ExpectedStatus: http.StatusTemporaryRedirect,
			AfterTestFunc: func(t *testing.T, app *tests.TestApp, res *http.Response) {
				app.Store().Get("cancelFunc").(context.CancelFunc)()

				checkFailureRedirect(t, app, res)

				if clientStubs[2]["c3"].HasSubscription("@oauth2") {
					t.Fatalf("Expected oauth2 subscription to be removed")
				}
			},
		},
		{
			Name:   "error query param",
			Method: http.MethodGet,
			Url:    "/api/oauth2-redirect?error=example&code=123&state=" + clientStubs[3]["c3"].Id(),
			BeforeTestFunc: beforeTestFunc(clientStubs[3], map[string][]string{
				"c3": {`"state":"` + clientStubs[3]["c3"].Id(), `"code":"123"`, `"error":"example"`},
			}),
			ExpectedStatus: http.StatusTemporaryRedirect,
			AfterTestFunc: func(t *testing.T, app *tests.TestApp, res *http.Response) {
				app.Store().Get("cancelFunc").(context.CancelFunc)()

				checkFailureRedirect(t, app, res)

				if clientStubs[3]["c3"].HasSubscription("@oauth2") {
					t.Fatalf("Expected oauth2 subscription to be removed")
				}
			},
		},
		{
			Name:           "discarded client with @oauth2 subscription",
			Method:         http.MethodGet,
			Url:            "/api/oauth2-redirect?code=123&state=" + clientStubs[4]["c5"].Id(),
			BeforeTestFunc: beforeTestFunc(clientStubs[4], nil),
			ExpectedStatus: http.StatusTemporaryRedirect,
			AfterTestFunc: func(t *testing.T, app *tests.TestApp, res *http.Response) {
				app.Store().Get("cancelFunc").(context.CancelFunc)()

				checkFailureRedirect(t, app, res)
			},
		},
		{
			Name:           "client without @oauth2 subscription",
			Method:         http.MethodGet,
			Url:            "/api/oauth2-redirect?code=123&state=" + clientStubs[4]["c4"].Id(),
			BeforeTestFunc: beforeTestFunc(clientStubs[5], nil),
			ExpectedStatus: http.StatusTemporaryRedirect,
			AfterTestFunc: func(t *testing.T, app *tests.TestApp, res *http.Response) {
				app.Store().Get("cancelFunc").(context.CancelFunc)()

				checkFailureRedirect(t, app, res)
			},
		},
		{
			Name:   "client with @oauth2 subscription",
			Method: http.MethodGet,
			Url:    "/api/oauth2-redirect?code=123&state=" + clientStubs[6]["c3"].Id(),
			BeforeTestFunc: beforeTestFunc(clientStubs[6], map[string][]string{
				"c3": {`"state":"` + clientStubs[6]["c3"].Id(), `"code":"123"`},
			}),
			ExpectedStatus: http.StatusTemporaryRedirect,
			AfterTestFunc: func(t *testing.T, app *tests.TestApp, res *http.Response) {
				app.Store().Get("cancelFunc").(context.CancelFunc)()

				checkSuccessRedirect(t, app, res)

				if clientStubs[6]["c3"].HasSubscription("@oauth2") {
					t.Fatalf("Expected oauth2 subscription to be removed")
				}
			},
		},
		{
			Name:   "(POST) client with @oauth2 subscription",
			Method: http.MethodPost,
			Url:    "/api/oauth2-redirect",
			Body:   strings.NewReader("code=123&state=" + clientStubs[7]["c3"].Id()),
			RequestHeaders: map[string]string{
				"content-type": "application/x-www-form-urlencoded",
			},
			BeforeTestFunc: beforeTestFunc(clientStubs[7], map[string][]string{
				"c3": {`"state":"` + clientStubs[7]["c3"].Id(), `"code":"123"`},
			}),
			ExpectedStatus: http.StatusSeeOther,
			AfterTestFunc: func(t *testing.T, app *tests.TestApp, res *http.Response) {
				app.Store().Get("cancelFunc").(context.CancelFunc)()

				checkSuccessRedirect(t, app, res)

				if clientStubs[7]["c3"].HasSubscription("@oauth2") {
					t.Fatalf("Expected oauth2 subscription to be removed")
				}
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
