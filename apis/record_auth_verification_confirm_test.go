package apis_test

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

func TestRecordConfirmVerification(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:           "empty data",
			Method:         http.MethodPost,
			URL:            "/api/collections/users/confirm-verification",
			Body:           strings.NewReader(``),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"token":{"code":"validation_required"`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
		{
			Name:            "invalid data format",
			Method:          http.MethodPost,
			URL:             "/api/collections/users/confirm-verification",
			Body:            strings.NewReader(`{"password`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "expired token",
			Method: http.MethodPost,
			URL:    "/api/collections/users/confirm-verification",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImV4cCI6MTY0MDk5MTY2MSwidHlwZSI6InZlcmlmaWNhdGlvbiIsImNvbGxlY3Rpb25JZCI6Il9wYl91c2Vyc19hdXRoXyIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSJ9.qqelNNL2Udl6K_TJ282sNHYCpASgA6SIuSVKGfBHMZU"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"token":{"code":"validation_invalid_token"`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
		{
			Name:   "non-verification token",
			Method: http.MethodPost,
			URL:    "/api/collections/users/confirm-verification",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImV4cCI6MjUyNDYwNDQ2MSwidHlwZSI6InBhc3N3b3JkUmVzZXQiLCJjb2xsZWN0aW9uSWQiOiJfcGJfdXNlcnNfYXV0aF8iLCJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20ifQ.xR-xq1oHDy0D8Q4NDOAEyYKGHWd_swzoiSoL8FLFBHY"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"token":{"code":"validation_invalid_token"`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
		{
			Name:   "non auth collection",
			Method: http.MethodPost,
			URL:    "/api/collections/demo1/confirm-verification?expand=rel,missing",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImV4cCI6MjUyNDYwNDQ2MSwidHlwZSI6InZlcmlmaWNhdGlvbiIsImNvbGxlY3Rpb25JZCI6Il9wYl91c2Vyc19hdXRoXyIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSJ9.SetHpu2H-x-q4TIUz-xiQjwi7MNwLCLvSs4O0hUSp0E"
			}`),
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "different auth collection",
			Method: http.MethodPost,
			URL:    "/api/collections/clients/confirm-verification?expand=rel,missing",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImV4cCI6MjUyNDYwNDQ2MSwidHlwZSI6InZlcmlmaWNhdGlvbiIsImNvbGxlY3Rpb25JZCI6Il9wYl91c2Vyc19hdXRoXyIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSJ9.SetHpu2H-x-q4TIUz-xiQjwi7MNwLCLvSs4O0hUSp0E"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{"token":{"code":"validation_token_collection_mismatch"`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
		{
			Name:   "valid token",
			Method: http.MethodPost,
			URL:    "/api/collections/users/confirm-verification",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImV4cCI6MjUyNDYwNDQ2MSwidHlwZSI6InZlcmlmaWNhdGlvbiIsImNvbGxlY3Rpb25JZCI6Il9wYl91c2Vyc19hdXRoXyIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSJ9.SetHpu2H-x-q4TIUz-xiQjwi7MNwLCLvSs4O0hUSp0E"
			}`),
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"*":                                  0,
				"OnRecordConfirmVerificationRequest": 1,
				"OnModelUpdate":                      1,
				"OnModelValidate":                    1,
				"OnModelUpdateExecute":               1,
				"OnModelAfterUpdateSuccess":          1,
				"OnRecordUpdate":                     1,
				"OnRecordValidate":                   1,
				"OnRecordUpdateExecute":              1,
				"OnRecordAfterUpdateSuccess":         1,
			},
		},
		{
			Name:   "valid token (already verified)",
			Method: http.MethodPost,
			URL:    "/api/collections/users/confirm-verification",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6Im9hcDY0MGNvdDR5cnUycyIsImV4cCI6MjUyNDYwNDQ2MSwidHlwZSI6InZlcmlmaWNhdGlvbiIsImNvbGxlY3Rpb25JZCI6Il9wYl91c2Vyc19hdXRoXyIsImVtYWlsIjoidGVzdDJAZXhhbXBsZS5jb20ifQ.QQmM3odNFVk6u4J4-5H8IBM3dfk9YCD7mPW-8PhBAI8"
			}`),
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"*":                                  0,
				"OnRecordConfirmVerificationRequest": 1,
			},
		},
		{
			Name:   "valid verification token from a collection without allowed login",
			Method: http.MethodPost,
			URL:    "/api/collections/nologin/confirm-verification",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImRjNDlrNmpnZWpuNDBoMyIsImV4cCI6MjUyNDYwNDQ2MSwidHlwZSI6InZlcmlmaWNhdGlvbiIsImNvbGxlY3Rpb25JZCI6ImtwdjcwOXNrMmxxYnFrOCIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSJ9.5GmuZr4vmwk3Cb_3ZZWNxwbE75KZC-j71xxIPR9AsVw"
			}`),
			ExpectedStatus:  204,
			ExpectedContent: []string{},
			ExpectedEvents: map[string]int{
				"*":                                  0,
				"OnRecordConfirmVerificationRequest": 1,
				"OnModelUpdate":                      1,
				"OnModelValidate":                    1,
				"OnModelUpdateExecute":               1,
				"OnModelAfterUpdateSuccess":          1,
				"OnRecordUpdate":                     1,
				"OnRecordValidate":                   1,
				"OnRecordUpdateExecute":              1,
				"OnRecordAfterUpdateSuccess":         1,
			},
		},
		{
			Name:   "OnRecordAfterConfirmVerificationRequest error response",
			Method: http.MethodPost,
			URL:    "/api/collections/users/confirm-verification",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsImV4cCI6MjUyNDYwNDQ2MSwidHlwZSI6InZlcmlmaWNhdGlvbiIsImNvbGxlY3Rpb25JZCI6Il9wYl91c2Vyc19hdXRoXyIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSJ9.SetHpu2H-x-q4TIUz-xiQjwi7MNwLCLvSs4O0hUSp0E"
			}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				app.OnRecordConfirmVerificationRequest().BindFunc(func(e *core.RecordConfirmVerificationRequestEvent) error {
					return errors.New("error")
				})
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"*":                                  0,
				"OnRecordConfirmVerificationRequest": 1,
			},
		},

		// rate limit checks
		// -----------------------------------------------------------
		{
			Name:   "RateLimit rule - nologin:confirmVerification",
			Method: http.MethodPost,
			URL:    "/api/collections/nologin/confirm-verification",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImRjNDlrNmpnZWpuNDBoMyIsImV4cCI6MjUyNDYwNDQ2MSwidHlwZSI6InZlcmlmaWNhdGlvbiIsImNvbGxlY3Rpb25JZCI6ImtwdjcwOXNrMmxxYnFrOCIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSJ9.5GmuZr4vmwk3Cb_3ZZWNxwbE75KZC-j71xxIPR9AsVw"
			}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				app.Settings().RateLimits.Enabled = true
				app.Settings().RateLimits.Rules = []core.RateLimitRule{
					{MaxRequests: 100, Label: "abc"},
					{MaxRequests: 100, Label: "*:confirmVerification"},
					{MaxRequests: 0, Label: "nologin:confirmVerification"},
				}
			},
			ExpectedStatus:  429,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "RateLimit rule - *:confirmVerification",
			Method: http.MethodPost,
			URL:    "/api/collections/nologin/confirm-verification",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImRjNDlrNmpnZWpuNDBoMyIsImV4cCI6MjUyNDYwNDQ2MSwidHlwZSI6InZlcmlmaWNhdGlvbiIsImNvbGxlY3Rpb25JZCI6ImtwdjcwOXNrMmxxYnFrOCIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSJ9.5GmuZr4vmwk3Cb_3ZZWNxwbE75KZC-j71xxIPR9AsVw"
			}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				app.Settings().RateLimits.Enabled = true
				app.Settings().RateLimits.Rules = []core.RateLimitRule{
					{MaxRequests: 100, Label: "abc"},
					{MaxRequests: 0, Label: "*:confirmVerification"},
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
