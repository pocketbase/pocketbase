package apis_test

import (
	"net/http"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

func TestRecordAuthMethodsList(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "missing collection",
			Method:          http.MethodGet,
			URL:             "/api/collections/missing/auth-methods",
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:            "non auth collection",
			Method:          http.MethodGet,
			URL:             "/api/collections/demo1/auth-methods",
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:           "auth collection with none auth methods allowed",
			Method:         http.MethodGet,
			URL:            "/api/collections/nologin/auth-methods",
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"password":{"identityFields":[],"enabled":false}`,
				`"oauth2":{"providers":[],"enabled":false}`,
				`"mfa":{"enabled":false,"duration":0}`,
				`"otp":{"enabled":false,"duration":0}`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
		{
			Name:           "auth collection with all auth methods allowed",
			Method:         http.MethodGet,
			URL:            "/api/collections/users/auth-methods",
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"password":{"identityFields":["email","username"],"enabled":true}`,
				`"mfa":{"enabled":true,"duration":1800}`,
				`"otp":{"enabled":true,"duration":300}`,
				`"oauth2":{`,
				`"providers":[{`,
				`"name":"google"`,
				`"name":"gitlab"`,
				`"state":`,
				`"displayName":`,
				`"codeVerifier":`,
				`"codeChallenge":`,
				`"codeChallengeMethod":`,
				`"authURL":`,
				`redirect_uri="`, // ensures that the redirect_uri is the last url param
			},
			ExpectedEvents: map[string]int{"*": 0},
		},

		// rate limit checks
		// -----------------------------------------------------------
		{
			Name:   "RateLimit rule - nologin:listAuthMethods",
			Method: http.MethodGet,
			URL:    "/api/collections/nologin/auth-methods",
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				app.Settings().RateLimits.Enabled = true
				app.Settings().RateLimits.Rules = []core.RateLimitRule{
					{MaxRequests: 100, Label: "abc"},
					{MaxRequests: 100, Label: "*:listAuthMethods"},
					{MaxRequests: 0, Label: "nologin:listAuthMethods"},
				}
			},
			ExpectedStatus:  429,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "RateLimit rule - *:listAuthMethods",
			Method: http.MethodGet,
			URL:    "/api/collections/nologin/auth-methods",
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				app.Settings().RateLimits.Enabled = true
				app.Settings().RateLimits.Rules = []core.RateLimitRule{
					{MaxRequests: 100, Label: "abc"},
					{MaxRequests: 0, Label: "*:listAuthMethods"},
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
