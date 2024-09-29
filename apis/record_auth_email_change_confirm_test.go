package apis_test

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

func TestRecordConfirmEmailChange(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "not an auth collection",
			Method:          http.MethodPost,
			URL:             "/api/collections/demo1/confirm-email-change",
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:           "empty data",
			Method:         http.MethodPost,
			URL:            "/api/collections/users/confirm-email-change",
			Body:           strings.NewReader(``),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":`,
				`"token":{"code":"validation_required"`,
				`"password":{"code":"validation_required"`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
		{
			Name:            "invalid data",
			Method:          http.MethodPost,
			URL:             "/api/collections/users/confirm-email-change",
			Body:            strings.NewReader(`{"token`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "expired token and correct password",
			Method: http.MethodPost,
			URL:    "/api/collections/users/confirm-email-change",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImNvbGxlY3Rpb25JZCI6Il9wYl91c2Vyc19hdXRoXyIsInR5cGUiOiJlbWFpbENoYW5nZSIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsIm5ld0VtYWlsIjoiY2hhbmdlQGV4YW1wbGUuY29tIiwiZXhwIjoxNjQwOTkxNjYxfQ.dff842MO0mgRTHY8dktp0dqG9-7LGQOgRuiAbQpYBls",
				"password":"1234567890"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"token":{`,
				`"code":"validation_invalid_token"`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
		{
			Name:   "non-email change token",
			Method: http.MethodPost,
			URL:    "/api/collections/users/confirm-email-change",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
				"password":"1234567890"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"token":{`,
				`"code":"validation_invalid_token_payload"`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
		{
			Name:   "valid token and incorrect password",
			Method: http.MethodPost,
			URL:    "/api/collections/users/confirm-email-change",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImNvbGxlY3Rpb25JZCI6Il9wYl91c2Vyc19hdXRoXyIsInR5cGUiOiJlbWFpbENoYW5nZSIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsIm5ld0VtYWlsIjoiY2hhbmdlQGV4YW1wbGUuY29tIiwiZXhwIjoyNTI0NjA0NDYxfQ.Y7mVlaEPhJiNPoIvIqbIosZU4c4lEhwysOrRR8c95iU",
				"password":"1234567891"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"password":{`,
				`"code":"validation_invalid_password"`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
		{
			Name:   "valid token and correct password",
			Method: http.MethodPost,
			URL:    "/api/collections/users/confirm-email-change",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImNvbGxlY3Rpb25JZCI6Il9wYl91c2Vyc19hdXRoXyIsInR5cGUiOiJlbWFpbENoYW5nZSIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsIm5ld0VtYWlsIjoiY2hhbmdlQGV4YW1wbGUuY29tIiwiZXhwIjoyNTI0NjA0NDYxfQ.Y7mVlaEPhJiNPoIvIqbIosZU4c4lEhwysOrRR8c95iU",
				"password":"1234567890"
			}`),
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"*":                                 0,
				"OnRecordConfirmEmailChangeRequest": 1,
				"OnModelUpdate":                     1,
				"OnModelUpdateExecute":              1,
				"OnModelAfterUpdateSuccess":         1,
				"OnModelValidate":                   1,
				"OnRecordUpdate":                    1,
				"OnRecordUpdateExecute":             1,
				"OnRecordAfterUpdateSuccess":        1,
				"OnRecordValidate":                  1,
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				_, err := app.FindAuthRecordByEmail("users", "change@example.com")
				if err != nil {
					t.Fatalf("Expected to find user with email %q, got error: %v", "change@example.com", err)
				}
			},
		},
		{
			Name:   "valid token in different auth collection",
			Method: http.MethodPost,
			URL:    "/api/collections/clients/confirm-email-change",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImNvbGxlY3Rpb25JZCI6Il9wYl91c2Vyc19hdXRoXyIsInR5cGUiOiJlbWFpbENoYW5nZSIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsIm5ld0VtYWlsIjoiY2hhbmdlQGV4YW1wbGUuY29tIiwiZXhwIjoyNTI0NjA0NDYxfQ.Y7mVlaEPhJiNPoIvIqbIosZU4c4lEhwysOrRR8c95iU",
				"password":"1234567890"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"token":{"code":"validation_token_collection_mismatch"`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
		{
			Name:   "OnRecordAfterConfirmEmailChangeRequest error response",
			Method: http.MethodPost,
			URL:    "/api/collections/users/confirm-email-change",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImNvbGxlY3Rpb25JZCI6Il9wYl91c2Vyc19hdXRoXyIsInR5cGUiOiJlbWFpbENoYW5nZSIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsIm5ld0VtYWlsIjoiY2hhbmdlQGV4YW1wbGUuY29tIiwiZXhwIjoyNTI0NjA0NDYxfQ.Y7mVlaEPhJiNPoIvIqbIosZU4c4lEhwysOrRR8c95iU",
				"password":"1234567890"
			}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				app.OnRecordConfirmEmailChangeRequest().BindFunc(func(e *core.RecordConfirmEmailChangeRequestEvent) error {
					return errors.New("error")
				})
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"*":                                 0,
				"OnRecordConfirmEmailChangeRequest": 1,
			},
		},

		// rate limit checks
		// -----------------------------------------------------------
		{
			Name:   "RateLimit rule - users:confirmEmailChange",
			Method: http.MethodPost,
			URL:    "/api/collections/users/confirm-email-change",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImNvbGxlY3Rpb25JZCI6Il9wYl91c2Vyc19hdXRoXyIsInR5cGUiOiJlbWFpbENoYW5nZSIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsIm5ld0VtYWlsIjoiY2hhbmdlQGV4YW1wbGUuY29tIiwiZXhwIjoyNTI0NjA0NDYxfQ.Y7mVlaEPhJiNPoIvIqbIosZU4c4lEhwysOrRR8c95iU",
				"password":"1234567890"
			}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				app.Settings().RateLimits.Enabled = true
				app.Settings().RateLimits.Rules = []core.RateLimitRule{
					{MaxRequests: 100, Label: "abc"},
					{MaxRequests: 100, Label: "*:confirmEmailChange"},
					{MaxRequests: 0, Label: "users:confirmEmailChange"},
				}
			},
			ExpectedStatus:  429,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "RateLimit rule - *:confirmEmailChange",
			Method: http.MethodPost,
			URL:    "/api/collections/users/confirm-email-change",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImNvbGxlY3Rpb25JZCI6Il9wYl91c2Vyc19hdXRoXyIsInR5cGUiOiJlbWFpbENoYW5nZSIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsIm5ld0VtYWlsIjoiY2hhbmdlQGV4YW1wbGUuY29tIiwiZXhwIjoyNTI0NjA0NDYxfQ.Y7mVlaEPhJiNPoIvIqbIosZU4c4lEhwysOrRR8c95iU",
				"password":"1234567890"
			}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				app.Settings().RateLimits.Enabled = true
				app.Settings().RateLimits.Rules = []core.RateLimitRule{
					{MaxRequests: 100, Label: "abc"},
					{MaxRequests: 0, Label: "*:confirmEmailChange"},
				}
			},
			ExpectedStatus:  429,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
