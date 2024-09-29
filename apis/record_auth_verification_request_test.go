package apis_test

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

func TestRecordRequestVerification(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "not an auth collection",
			Method:          http.MethodPost,
			URL:             "/api/collections/demo1/request-verification",
			Body:            strings.NewReader(``),
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:            "empty data",
			Method:          http.MethodPost,
			URL:             "/api/collections/users/request-verification",
			Body:            strings.NewReader(``),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{"email":{"code":"validation_required","message":"Cannot be blank."}}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:            "invalid data",
			Method:          http.MethodPost,
			URL:             "/api/collections/users/request-verification",
			Body:            strings.NewReader(`{"email`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:           "missing auth record",
			Method:         http.MethodPost,
			URL:            "/api/collections/users/request-verification",
			Body:           strings.NewReader(`{"email":"missing@example.com"}`),
			Delay:          100 * time.Millisecond,
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{"*": 0},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				if app.TestMailer.TotalSend() != 0 {
					t.Fatalf("Expected zero emails, got %d", app.TestMailer.TotalSend())
				}
			},
		},
		{
			Name:           "already verified auth record",
			Method:         http.MethodPost,
			URL:            "/api/collections/users/request-verification",
			Body:           strings.NewReader(`{"email":"test2@example.com"}`),
			Delay:          100 * time.Millisecond,
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"*":                                  0,
				"OnRecordRequestVerificationRequest": 1,
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				if app.TestMailer.TotalSend() != 0 {
					t.Fatalf("Expected zero emails, got %d", app.TestMailer.TotalSend())
				}
			},
		},
		{
			Name:           "existing auth record",
			Method:         http.MethodPost,
			URL:            "/api/collections/users/request-verification",
			Body:           strings.NewReader(`{"email":"test@example.com"}`),
			Delay:          100 * time.Millisecond,
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"*":                                  0,
				"OnRecordRequestVerificationRequest": 1,
				"OnMailerSend":                       1,
				"OnMailerRecordVerificationSend":     1,
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				if !strings.Contains(app.TestMailer.LastMessage().HTML, "/auth/confirm-verification") {
					t.Fatalf("Expected verification email, got\n%v", app.TestMailer.LastMessage().HTML)
				}
			},
		},
		{
			Name:           "existing auth record (after already sent)",
			Method:         http.MethodPost,
			URL:            "/api/collections/users/request-verification",
			Body:           strings.NewReader(`{"email":"test@example.com"}`),
			Delay:          100 * time.Millisecond,
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"*": 0,
				// terminated before firing the event
				// "OnRecordRequestVerificationRequest": 1,
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				// simulate recent verification sent
				authRecord, err := app.FindFirstRecordByData("users", "email", "test@example.com")
				if err != nil {
					t.Fatal(err)
				}
				resendKey := "@limitVerificationEmail_" + authRecord.Collection().Id + authRecord.Id
				app.Store().Set(resendKey, struct{}{})
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				if app.TestMailer.TotalSend() != 0 {
					t.Fatalf("Expected zero emails, got %d", app.TestMailer.TotalSend())
				}
			},
		},

		// rate limit checks
		// -----------------------------------------------------------
		{
			Name:   "RateLimit rule - users:requestVerification",
			Method: http.MethodPost,
			URL:    "/api/collections/users/request-verification",
			Body:   strings.NewReader(`{"email":"test@example.com"}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				app.Settings().RateLimits.Enabled = true
				app.Settings().RateLimits.Rules = []core.RateLimitRule{
					{MaxRequests: 100, Label: "abc"},
					{MaxRequests: 100, Label: "*:requestVerification"},
					{MaxRequests: 0, Label: "users:requestVerification"},
				}
			},
			ExpectedStatus:  429,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "RateLimit rule - *:requestVerification",
			Method: http.MethodPost,
			URL:    "/api/collections/users/request-verification",
			Body:   strings.NewReader(`{"email":"test@example.com"}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				app.Settings().RateLimits.Enabled = true
				app.Settings().RateLimits.Rules = []core.RateLimitRule{
					{MaxRequests: 100, Label: "abc"},
					{MaxRequests: 0, Label: "*:requestVerification"},
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
