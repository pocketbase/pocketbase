package apis_test

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/dbutils"
)

func TestRecordAuthWithPassword(t *testing.T) {
	t.Parallel()

	updateIdentityIndex := func(collectionIdOrName string, fieldCollateMap map[string]string) func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
		return func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
			collection, err := app.FindCollectionByNameOrId("clients")
			if err != nil {
				t.Fatal(err)
			}

			for column, collate := range fieldCollateMap {
				index, ok := dbutils.FindSingleColumnUniqueIndex(collection.Indexes, column)
				if !ok {
					t.Fatalf("Missing unique identityField index for column %q", column)
				}

				index.Columns[0].Collate = collate

				collection.RemoveIndex(index.IndexName)
				collection.Indexes = append(collection.Indexes, index.Build())
			}

			err = app.Save(collection)
			if err != nil {
				t.Fatalf("Failed to update identityField index: %v", err)
			}
		}
	}

	scenarios := []tests.ApiScenario{
		{
			Name:            "disabled password auth",
			Method:          http.MethodPost,
			URL:             "/api/collections/nologin/auth-with-password",
			Body:            strings.NewReader(`{"identity":"test@example.com","password":"1234567890"}`),
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:            "non-auth collection",
			Method:          http.MethodPost,
			URL:             "/api/collections/demo1/auth-with-password",
			Body:            strings.NewReader(`{"identity":"test@example.com","password":"1234567890"}`),
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:            "invalid body format",
			Method:          http.MethodPost,
			URL:             "/api/collections/clients/auth-with-password",
			Body:            strings.NewReader(`{"identity`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:           "empty body params",
			Method:         http.MethodPost,
			URL:            "/api/collections/clients/auth-with-password",
			Body:           strings.NewReader(`{"identity":"","password":""}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"identity":{`,
				`"password":{`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
		{
			Name:   "OnRecordAuthWithPasswordRequest error response",
			Method: http.MethodPost,
			URL:    "/api/collections/clients/auth-with-password",
			Body: strings.NewReader(`{
				"identity":"test@example.com",
				"password":"1234567890"
			}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				app.OnRecordAuthWithPasswordRequest().BindFunc(func(e *core.RecordAuthWithPasswordRequestEvent) error {
					return errors.New("error")
				})
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"*":                               0,
				"OnRecordAuthWithPasswordRequest": 1,
			},
		},
		{
			Name:   "valid identity field and invalid password",
			Method: http.MethodPost,
			URL:    "/api/collections/clients/auth-with-password",
			Body: strings.NewReader(`{
				"identity":"test@example.com",
				"password":"invalid"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{}`,
			},
			ExpectedEvents: map[string]int{
				"*":                               0,
				"OnRecordAuthWithPasswordRequest": 1,
			},
		},
		{
			Name:   "valid identity field (email) and valid password",
			Method: http.MethodPost,
			URL:    "/api/collections/clients/auth-with-password",
			Body: strings.NewReader(`{
				"identity":"test@example.com",
				"password":"1234567890"
			}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				// test at least once that the correct request info context is properly loaded
				app.OnRecordAuthRequest().BindFunc(func(e *core.RecordAuthRequestEvent) error {
					info, err := e.RequestInfo()
					if err != nil {
						t.Fatal(err)
					}

					if info.Context != core.RequestInfoContextPasswordAuth {
						t.Fatalf("Expected request context %q, got %q", core.RequestInfoContextPasswordAuth, info.Context)
					}

					return e.Next()
				})
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"email":"test@example.com"`,
				`"token":`,
			},
			NotExpectedContent: []string{
				// hidden fields
				`"tokenKey"`,
				`"password"`,
			},
			ExpectedEvents: map[string]int{
				"*":                               0,
				"OnRecordAuthWithPasswordRequest": 1,
				"OnRecordAuthRequest":             1,
				"OnRecordEnrich":                  1,
				// authOrigin track
				"OnModelCreate":               1,
				"OnModelCreateExecute":        1,
				"OnModelAfterCreateSuccess":   1,
				"OnModelValidate":             1,
				"OnRecordCreate":              1,
				"OnRecordCreateExecute":       1,
				"OnRecordAfterCreateSuccess":  1,
				"OnRecordValidate":            1,
				"OnMailerSend":                1,
				"OnMailerRecordAuthAlertSend": 1,
			},
		},
		{
			Name:   "valid identity field (username) and valid password",
			Method: http.MethodPost,
			URL:    "/api/collections/clients/auth-with-password",
			Body: strings.NewReader(`{
				"identity":"clients57772",
				"password":"1234567890"
			}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"email":"test@example.com"`,
				`"username":"clients57772"`,
				`"token":`,
			},
			NotExpectedContent: []string{
				// hidden fields
				`"tokenKey"`,
				`"password"`,
			},
			ExpectedEvents: map[string]int{
				"*":                               0,
				"OnRecordAuthWithPasswordRequest": 1,
				"OnRecordAuthRequest":             1,
				"OnRecordEnrich":                  1,
				// authOrigin track
				"OnModelCreate":               1,
				"OnModelCreateExecute":        1,
				"OnModelAfterCreateSuccess":   1,
				"OnModelValidate":             1,
				"OnRecordCreate":              1,
				"OnRecordCreateExecute":       1,
				"OnRecordAfterCreateSuccess":  1,
				"OnRecordValidate":            1,
				"OnMailerSend":                1,
				"OnMailerRecordAuthAlertSend": 1,
			},
		},
		{
			Name:   "unknown explicit identityField",
			Method: http.MethodPost,
			URL:    "/api/collections/clients/auth-with-password",
			Body: strings.NewReader(`{
				"identityField": "created",
				"identity":"test@example.com",
				"password":"1234567890"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"identityField":{"code":"validation_in_invalid"`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
		{
			Name:   "valid identity field and valid password with mismatched explicit identityField",
			Method: http.MethodPost,
			URL:    "/api/collections/clients/auth-with-password",
			Body: strings.NewReader(`{
				"identityField": "username",
				"identity":"test@example.com",
				"password":"1234567890"
			}`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"*":                               0,
				"OnRecordAuthWithPasswordRequest": 1,
			},
		},
		{
			Name:   "valid identity field and valid password with matched explicit identityField",
			Method: http.MethodPost,
			URL:    "/api/collections/clients/auth-with-password",
			Body: strings.NewReader(`{
				"identityField": "username",
				"identity":"clients57772",
				"password":"1234567890"
			}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"email":"test@example.com"`,
				`"username":"clients57772"`,
				`"token":`,
			},
			NotExpectedContent: []string{
				// hidden fields
				`"tokenKey"`,
				`"password"`,
			},
			ExpectedEvents: map[string]int{
				"*":                               0,
				"OnRecordAuthWithPasswordRequest": 1,
				"OnRecordAuthRequest":             1,
				"OnRecordEnrich":                  1,
				// authOrigin track
				"OnModelCreate":               1,
				"OnModelCreateExecute":        1,
				"OnModelAfterCreateSuccess":   1,
				"OnModelValidate":             1,
				"OnRecordCreate":              1,
				"OnRecordCreateExecute":       1,
				"OnRecordAfterCreateSuccess":  1,
				"OnRecordValidate":            1,
				"OnMailerSend":                1,
				"OnMailerRecordAuthAlertSend": 1,
			},
		},
		{
			Name:   "valid identity (unverified) and valid password in onlyVerified collection",
			Method: http.MethodPost,
			URL:    "/api/collections/clients/auth-with-password",
			Body: strings.NewReader(`{
				"identity":"test2@example.com",
				"password":"1234567890"
			}`),
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"*":                               0,
				"OnRecordAuthWithPasswordRequest": 1,
			},
		},
		{
			Name:   "already authenticated record",
			Method: http.MethodPost,
			URL:    "/api/collections/clients/auth-with-password",
			Body: strings.NewReader(`{
				"identity":"test@example.com",
				"password":"1234567890"
			}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"gk390qegs4y47wn"`,
				`"email":"test@example.com"`,
				`"token":`,
			},
			NotExpectedContent: []string{
				// hidden fields
				`"tokenKey"`,
				`"password"`,
			},
			ExpectedEvents: map[string]int{
				"*":                               0,
				"OnRecordAuthWithPasswordRequest": 1,
				"OnRecordAuthRequest":             1,
				"OnRecordEnrich":                  1,
				// authOrigin track
				"OnModelCreate":               1,
				"OnModelCreateExecute":        1,
				"OnModelAfterCreateSuccess":   1,
				"OnModelValidate":             1,
				"OnRecordCreate":              1,
				"OnRecordCreateExecute":       1,
				"OnRecordAfterCreateSuccess":  1,
				"OnRecordValidate":            1,
				"OnMailerSend":                1,
				"OnMailerRecordAuthAlertSend": 1,
			},
		},
		{
			Name:   "with mfa first auth check",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-password",
			Body: strings.NewReader(`{
				"identity":"test@example.com",
				"password":"1234567890"
			}`),
			ExpectedStatus: 401,
			ExpectedContent: []string{
				`"mfaId":"`,
			},
			ExpectedEvents: map[string]int{
				"*":                               0,
				"OnRecordAuthWithPasswordRequest": 1,
				"OnRecordAuthRequest":             1,
				// mfa create
				"OnModelCreate":              1,
				"OnModelCreateExecute":       1,
				"OnModelAfterCreateSuccess":  1,
				"OnModelValidate":            1,
				"OnRecordCreate":             1,
				"OnRecordCreateExecute":      1,
				"OnRecordAfterCreateSuccess": 1,
				"OnRecordValidate":           1,
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				user, err := app.FindAuthRecordByEmail("users", "test@example.com")
				if err != nil {
					t.Fatal(err)
				}

				mfas, err := app.FindAllMFAsByRecord(user)
				if err != nil {
					t.Fatal(err)
				}

				if v := len(mfas); v != 1 {
					t.Fatalf("Expected 1 mfa record to be created, got %d", v)
				}
			},
		},
		{
			Name:   "with mfa second auth check",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-password",
			Body: strings.NewReader(`{
				"mfaId": "` + strings.Repeat("a", 15) + `",
				"identity":"test@example.com",
				"password":"1234567890"
			}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				user, err := app.FindAuthRecordByEmail("users", "test@example.com")
				if err != nil {
					t.Fatal(err)
				}

				// insert a dummy mfa record
				mfa := core.NewMFA(app)
				mfa.Id = strings.Repeat("a", 15)
				mfa.SetCollectionRef(user.Collection().Id)
				mfa.SetRecordRef(user.Id)
				mfa.SetMethod("test")
				if err := app.Save(mfa); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"email":"test@example.com"`,
				`"token":`,
			},
			NotExpectedContent: []string{
				// hidden fields
				`"tokenKey"`,
				`"password"`,
			},
			ExpectedEvents: map[string]int{
				"*":                               0,
				"OnRecordAuthWithPasswordRequest": 1,
				"OnRecordAuthRequest":             1,
				"OnRecordEnrich":                  1,
				// authOrigin track
				"OnModelCreate":               1,
				"OnModelCreateExecute":        1,
				"OnModelAfterCreateSuccess":   1,
				"OnModelValidate":             1,
				"OnRecordCreate":              1,
				"OnRecordCreateExecute":       1,
				"OnRecordAfterCreateSuccess":  1,
				"OnRecordValidate":            1,
				"OnMailerSend":                0, // disabled auth email alerts
				"OnMailerRecordAuthAlertSend": 0,
				// mfa delete
				"OnModelDelete":              1,
				"OnModelDeleteExecute":       1,
				"OnModelAfterDeleteSuccess":  1,
				"OnRecordDelete":             1,
				"OnRecordDeleteExecute":      1,
				"OnRecordAfterDeleteSuccess": 1,
			},
		},
		{
			Name:   "with enabled mfa but unsatisfied mfa rule (aka. skip the mfa check)",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-password",
			Body: strings.NewReader(`{
				"identity":"test@example.com",
				"password":"1234567890"
			}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				users, err := app.FindCollectionByNameOrId("users")
				if err != nil {
					t.Fatal(err)
				}

				users.MFA.Enabled = true
				users.MFA.Rule = "1=2"

				if err := app.Save(users); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"email":"test@example.com"`,
				`"token":`,
			},
			NotExpectedContent: []string{
				// hidden fields
				`"tokenKey"`,
				`"password"`,
			},
			ExpectedEvents: map[string]int{
				"*":                               0,
				"OnRecordAuthWithPasswordRequest": 1,
				"OnRecordAuthRequest":             1,
				"OnRecordEnrich":                  1,
				// authOrigin track
				"OnModelCreate":               1,
				"OnModelCreateExecute":        1,
				"OnModelAfterCreateSuccess":   1,
				"OnModelValidate":             1,
				"OnRecordCreate":              1,
				"OnRecordCreateExecute":       1,
				"OnRecordAfterCreateSuccess":  1,
				"OnRecordValidate":            1,
				"OnMailerSend":                0, // disabled auth email alerts
				"OnMailerRecordAuthAlertSend": 0,
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				user, err := app.FindAuthRecordByEmail("users", "test@example.com")
				if err != nil {
					t.Fatal(err)
				}

				mfas, err := app.FindAllMFAsByRecord(user)
				if err != nil {
					t.Fatal(err)
				}

				if v := len(mfas); v != 0 {
					t.Fatalf("Expected no mfa records to be created, got %d", v)
				}
			},
		},

		// case sensitivity checks
		// -----------------------------------------------------------
		{
			Name:   "with explicit identityField (case-sensitive)",
			Method: http.MethodPost,
			URL:    "/api/collections/clients/auth-with-password",
			Body: strings.NewReader(`{
				"identityField": "username",
				"identity":"Clients57772",
				"password":"1234567890"
			}`),
			BeforeTestFunc:  updateIdentityIndex("clients", map[string]string{"username": ""}),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"*":                               0,
				"OnRecordAuthWithPasswordRequest": 1,
			},
		},
		{
			Name:   "with explicit identityField (case-insensitive)",
			Method: http.MethodPost,
			URL:    "/api/collections/clients/auth-with-password",
			Body: strings.NewReader(`{
				"identityField": "username",
				"identity":"Clients57772",
				"password":"1234567890"
			}`),
			BeforeTestFunc: updateIdentityIndex("clients", map[string]string{"username": "nocase"}),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"email":"test@example.com"`,
				`"username":"clients57772"`,
				`"token":`,
			},
			NotExpectedContent: []string{
				// hidden fields
				`"tokenKey"`,
				`"password"`,
			},
			ExpectedEvents: map[string]int{
				"*":                               0,
				"OnRecordAuthWithPasswordRequest": 1,
				"OnRecordAuthRequest":             1,
				"OnRecordEnrich":                  1,
				// authOrigin track
				"OnModelCreate":               1,
				"OnModelCreateExecute":        1,
				"OnModelAfterCreateSuccess":   1,
				"OnModelValidate":             1,
				"OnRecordCreate":              1,
				"OnRecordCreateExecute":       1,
				"OnRecordAfterCreateSuccess":  1,
				"OnRecordValidate":            1,
				"OnMailerSend":                1,
				"OnMailerRecordAuthAlertSend": 1,
			},
		},
		{
			Name:   "without explicit identityField and non-email field (case-insensitive)",
			Method: http.MethodPost,
			URL:    "/api/collections/clients/auth-with-password",
			Body: strings.NewReader(`{
				"identity":"Clients57772",
				"password":"1234567890"
			}`),
			BeforeTestFunc: updateIdentityIndex("clients", map[string]string{"username": "nocase"}),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"email":"test@example.com"`,
				`"username":"clients57772"`,
				`"token":`,
			},
			NotExpectedContent: []string{
				// hidden fields
				`"tokenKey"`,
				`"password"`,
			},
			ExpectedEvents: map[string]int{
				"*":                               0,
				"OnRecordAuthWithPasswordRequest": 1,
				"OnRecordAuthRequest":             1,
				"OnRecordEnrich":                  1,
				// authOrigin track
				"OnModelCreate":               1,
				"OnModelCreateExecute":        1,
				"OnModelAfterCreateSuccess":   1,
				"OnModelValidate":             1,
				"OnRecordCreate":              1,
				"OnRecordCreateExecute":       1,
				"OnRecordAfterCreateSuccess":  1,
				"OnRecordValidate":            1,
				"OnMailerSend":                1,
				"OnMailerRecordAuthAlertSend": 1,
			},
		},
		{
			Name:   "without explicit identityField and email field (case-insensitive)",
			Method: http.MethodPost,
			URL:    "/api/collections/clients/auth-with-password",
			Body: strings.NewReader(`{
				"identity":"tESt@example.com",
				"password":"1234567890"
			}`),
			BeforeTestFunc: updateIdentityIndex("clients", map[string]string{"email": "nocase"}),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"email":"test@example.com"`,
				`"username":"clients57772"`,
				`"token":`,
			},
			NotExpectedContent: []string{
				// hidden fields
				`"tokenKey"`,
				`"password"`,
			},
			ExpectedEvents: map[string]int{
				"*":                               0,
				"OnRecordAuthWithPasswordRequest": 1,
				"OnRecordAuthRequest":             1,
				"OnRecordEnrich":                  1,
				// authOrigin track
				"OnModelCreate":               1,
				"OnModelCreateExecute":        1,
				"OnModelAfterCreateSuccess":   1,
				"OnModelValidate":             1,
				"OnRecordCreate":              1,
				"OnRecordCreateExecute":       1,
				"OnRecordAfterCreateSuccess":  1,
				"OnRecordValidate":            1,
				"OnMailerSend":                1,
				"OnMailerRecordAuthAlertSend": 1,
			},
		},

		// rate limit checks
		// -----------------------------------------------------------
		{
			Name:   "RateLimit rule - users:authWithPassword",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-password",
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				app.Settings().RateLimits.Enabled = true
				app.Settings().RateLimits.Rules = []core.RateLimitRule{
					{MaxRequests: 100, Label: "abc"},
					{MaxRequests: 100, Label: "*:authWithPassword"},
					{MaxRequests: 100, Label: "users:auth"},
					{MaxRequests: 0, Label: "users:authWithPassword"},
				}
			},
			ExpectedStatus:  429,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "RateLimit rule - *:authWithPassword",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-password",
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				app.Settings().RateLimits.Enabled = true
				app.Settings().RateLimits.Rules = []core.RateLimitRule{
					{MaxRequests: 100, Label: "abc"},
					{MaxRequests: 100, Label: "*:auth"},
					{MaxRequests: 0, Label: "*:authWithPassword"},
				}
			},
			ExpectedStatus:  429,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "RateLimit rule - users:auth",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-password",
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				app.Settings().RateLimits.Enabled = true
				app.Settings().RateLimits.Rules = []core.RateLimitRule{
					{MaxRequests: 100, Label: "abc"},
					{MaxRequests: 100, Label: "*:authWithPassword"},
					{MaxRequests: 0, Label: "users:auth"},
				}
			},
			ExpectedStatus:  429,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "RateLimit rule - *:auth",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-password",
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				app.Settings().RateLimits.Enabled = true
				app.Settings().RateLimits.Rules = []core.RateLimitRule{
					{MaxRequests: 100, Label: "abc"},
					{MaxRequests: 0, Label: "*:auth"},
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
