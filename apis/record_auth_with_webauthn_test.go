package apis_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

func TestRecordWebAuthnRegisterBegin(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "not an auth collection",
			Method:          http.MethodPost,
			URL:             "/api/collections/demo1/auth-with-webauthn/register-begin",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "auth collection with disabled webauthn",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-webauthn/register-begin",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				usersCol, err := app.FindCollectionByNameOrId("users")
				if err != nil {
					t.Fatal(err)
				}
				usersCol.WebAuthn.Enabled = false
				if err := app.Save(usersCol); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:            "no auth token (unauthenticated)",
			Method:          http.MethodPost,
			URL:             "/api/collections/users/auth-with-webauthn/register-begin",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				usersCol, err := app.FindCollectionByNameOrId("users")
				if err != nil {
					t.Fatal(err)
				}
				usersCol.WebAuthn.Enabled = true
				if err := app.Save(usersCol); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			Name:   "valid auth token with enabled webauthn",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-webauthn/register-begin",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				usersCol, err := app.FindCollectionByNameOrId("users")
				if err != nil {
					t.Fatal(err)
				}
				usersCol.WebAuthn.Enabled = true
				if err := app.Save(usersCol); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"sessionToken":`,
				`"options":{`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRecordWebAuthnRegisterFinish(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "not an auth collection",
			Method:          http.MethodPost,
			URL:             "/api/collections/demo1/auth-with-webauthn/register-finish",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "auth collection with disabled webauthn",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-webauthn/register-finish",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				usersCol, err := app.FindCollectionByNameOrId("users")
				if err != nil {
					t.Fatal(err)
				}
				usersCol.WebAuthn.Enabled = false
				if err := app.Save(usersCol); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "no auth token (unauthenticated)",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-webauthn/register-finish",
			Body:   strings.NewReader(`{"sessionToken":"test123"}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				usersCol, err := app.FindCollectionByNameOrId("users")
				if err != nil {
					t.Fatal(err)
				}
				usersCol.WebAuthn.Enabled = true
				if err := app.Save(usersCol); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "missing session token",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-webauthn/register-finish",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			Body: strings.NewReader(`{}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				usersCol, err := app.FindCollectionByNameOrId("users")
				if err != nil {
					t.Fatal(err)
				}
				usersCol.WebAuthn.Enabled = true
				if err := app.Save(usersCol); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"sessionToken":{"code":"validation_required"`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
		{
			Name:   "invalid session token",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-webauthn/register-finish",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			Body: strings.NewReader(`{"sessionToken":"invalid_token_12345"}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				usersCol, err := app.FindCollectionByNameOrId("users")
				if err != nil {
					t.Fatal(err)
				}
				usersCol.WebAuthn.Enabled = true
				if err := app.Save(usersCol); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRecordWebAuthnLoginBegin(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "not an auth collection",
			Method:          http.MethodPost,
			URL:             "/api/collections/demo1/auth-with-webauthn/login-begin",
			Body:            strings.NewReader(`{"identity":"test@example.com"}`),
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "auth collection with disabled webauthn",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-webauthn/login-begin",
			Body:   strings.NewReader(`{"identity":"test@example.com"}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				usersCol, err := app.FindCollectionByNameOrId("users")
				if err != nil {
					t.Fatal(err)
				}
				usersCol.WebAuthn.Enabled = false
				if err := app.Save(usersCol); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "empty body",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-webauthn/login-begin",
			Body:   strings.NewReader(``),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				usersCol, err := app.FindCollectionByNameOrId("users")
				if err != nil {
					t.Fatal(err)
				}
				usersCol.WebAuthn.Enabled = true
				if err := app.Save(usersCol); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"identity":{"code":"validation_required"`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
		{
			Name:   "non-existing identity",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-webauthn/login-begin",
			Body:   strings.NewReader(`{"identity":"nonexisting@example.com"}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				usersCol, err := app.FindCollectionByNameOrId("users")
				if err != nil {
					t.Fatal(err)
				}
				usersCol.WebAuthn.Enabled = true
				if err := app.Save(usersCol); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "valid identity but user has no passkeys",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-webauthn/login-begin",
			Body:   strings.NewReader(`{"identity":"test@example.com"}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				usersCol, err := app.FindCollectionByNameOrId("users")
				if err != nil {
					t.Fatal(err)
				}
				usersCol.WebAuthn.Enabled = true
				if err := app.Save(usersCol); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRecordWebAuthnLoginFinish(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "not an auth collection",
			Method:          http.MethodPost,
			URL:             "/api/collections/demo1/auth-with-webauthn/login-finish",
			Body:            strings.NewReader(`{"identity":"test@example.com","sessionToken":"test"}`),
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "auth collection with disabled webauthn",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-webauthn/login-finish",
			Body:   strings.NewReader(`{"identity":"test@example.com","sessionToken":"test"}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				usersCol, err := app.FindCollectionByNameOrId("users")
				if err != nil {
					t.Fatal(err)
				}
				usersCol.WebAuthn.Enabled = false
				if err := app.Save(usersCol); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "empty body",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-webauthn/login-finish",
			Body:   strings.NewReader(``),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				usersCol, err := app.FindCollectionByNameOrId("users")
				if err != nil {
					t.Fatal(err)
				}
				usersCol.WebAuthn.Enabled = true
				if err := app.Save(usersCol); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"identity":{"code":"validation_required"`,
				`"sessionToken":{"code":"validation_required"`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
		{
			Name:   "non-existing identity",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-webauthn/login-finish",
			Body:   strings.NewReader(`{"identity":"missing@example.com","sessionToken":"test"}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				usersCol, err := app.FindCollectionByNameOrId("users")
				if err != nil {
					t.Fatal(err)
				}
				usersCol.WebAuthn.Enabled = true
				if err := app.Save(usersCol); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "valid identity but invalid session token",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-webauthn/login-finish",
			Body:   strings.NewReader(`{"identity":"test@example.com","sessionToken":"invalid_token"}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				usersCol, err := app.FindCollectionByNameOrId("users")
				if err != nil {
					t.Fatal(err)
				}
				usersCol.WebAuthn.Enabled = true
				if err := app.Save(usersCol); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRecordAuthMethodsWebAuthn(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:   "webauthn disabled (default)",
			Method: http.MethodGet,
			URL:    "/api/collections/users/auth-methods",
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				usersCol, err := app.FindCollectionByNameOrId("users")
				if err != nil {
					t.Fatal(err)
				}
				usersCol.WebAuthn.Enabled = false
				if err := app.Save(usersCol); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"webauthn":{"enabled":false}`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
		{
			Name:   "webauthn enabled",
			Method: http.MethodGet,
			URL:    "/api/collections/users/auth-methods",
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				usersCol, err := app.FindCollectionByNameOrId("users")
				if err != nil {
					t.Fatal(err)
				}
				usersCol.WebAuthn.Enabled = true
				if err := app.Save(usersCol); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"webauthn":{"enabled":true}`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRecordWebAuthnListCredentials(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthenticated",
			Method:          http.MethodGet,
			URL:             "/api/collections/users/auth-with-webauthn/credentials",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				usersCol, err := app.FindCollectionByNameOrId("users")
				if err != nil {
					t.Fatal(err)
				}
				usersCol.WebAuthn.Enabled = true
				if err := app.Save(usersCol); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			Name:   "authenticated user with no credentials",
			Method: http.MethodGet,
			URL:    "/api/collections/users/auth-with-webauthn/credentials",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				usersCol, err := app.FindCollectionByNameOrId("users")
				if err != nil {
					t.Fatal(err)
				}
				usersCol.WebAuthn.Enabled = true
				if err := app.Save(usersCol); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus:  200,
			ExpectedContent: []string{`[]`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRecordWebAuthnDeleteCredential(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthenticated",
			Method:          http.MethodDelete,
			URL:             "/api/collections/users/auth-with-webauthn/credentials/test123",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				usersCol, err := app.FindCollectionByNameOrId("users")
				if err != nil {
					t.Fatal(err)
				}
				usersCol.WebAuthn.Enabled = true
				if err := app.Save(usersCol); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			Name:   "authenticated but credential not found",
			Method: http.MethodDelete,
			URL:    "/api/collections/users/auth-with-webauthn/credentials/nonexistent",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				usersCol, err := app.FindCollectionByNameOrId("users")
				if err != nil {
					t.Fatal(err)
				}
				usersCol.WebAuthn.Enabled = true
				if err := app.Save(usersCol); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRecordWebAuthnAdminClearCredentials(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized (no auth)",
			Method:          http.MethodDelete,
			URL:             "/api/collections/users/auth-with-webauthn/credentials-by-record/4q1xlclmfloku33",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "unauthorized (regular user)",
			Method: http.MethodDelete,
			URL:    "/api/collections/users/auth-with-webauthn/credentials-by-record/4q1xlclmfloku33",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized superuser - user not found",
			Method: http.MethodDelete,
			URL:    "/api/collections/users/auth-with-webauthn/credentials-by-record/nonexistent",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized superuser - user with no credentials",
			Method: http.MethodDelete,
			URL:    "/api/collections/users/auth-with-webauthn/credentials-by-record/4q1xlclmfloku33",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"deleted":0`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
