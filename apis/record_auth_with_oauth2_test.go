package apis_test

import (
	"bytes"
	"errors"
	"image"
	"image/png"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/auth"
	"github.com/pocketbase/pocketbase/tools/dbutils"
	"golang.org/x/oauth2"
)

var _ auth.Provider = (*oauth2MockProvider)(nil)

type oauth2MockProvider struct {
	auth.BaseProvider

	AuthUser *auth.AuthUser
	Token    *oauth2.Token
}

func (p *oauth2MockProvider) FetchToken(code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	if p.Token == nil {
		return nil, errors.New("failed to fetch OAuth2 token")
	}
	return p.Token, nil
}

func (p *oauth2MockProvider) FetchAuthUser(token *oauth2.Token) (*auth.AuthUser, error) {
	if p.AuthUser == nil {
		return nil, errors.New("failed to fetch OAuth2 user")
	}
	return p.AuthUser, nil
}

func TestRecordAuthWithOAuth2(t *testing.T) {
	t.Parallel()

	// start a test server
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		buf := new(bytes.Buffer)
		png.Encode(buf, image.Rect(0, 0, 1, 1)) // tiny 1x1 png
		http.ServeContent(res, req, "test_avatar.png", time.Now(), bytes.NewReader(buf.Bytes()))
	}))
	defer server.Close()

	scenarios := []tests.ApiScenario{
		{
			Name:   "disabled OAuth2 auth",
			Method: http.MethodPost,
			URL:    "/api/collections/nologin/auth-with-oauth2",
			Body: strings.NewReader(`{
				"provider":     "test",
				"code":         "123",
				"codeVerifier": "456",
				"redirectURL":  "https://example.com"
			}`),
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:            "invalid body",
			Method:          http.MethodPost,
			URL:             "/api/collections/users/auth-with-oauth2",
			Body:            strings.NewReader(`{"provider"`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "trigger form validations (missing provider)",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-oauth2",
			Body: strings.NewReader(`{
				"provider": "missing"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"provider":`,
				`"code":`,
				`"redirectURL":`,
			},
			NotExpectedContent: []string{
				`"codeVerifier":`, // should be optional
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
		{
			Name:   "trigger form validations (existing but disabled provider)",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-oauth2",
			Body: strings.NewReader(`{
				"provider": "apple"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"provider":`,
				`"code":`,
				`"redirectURL":`,
			},
			NotExpectedContent: []string{
				`"codeVerifier":`, // should be optional
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
		{
			Name:   "existing linked OAuth2 (unverified user)",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-oauth2",
			Body: strings.NewReader(`{
				"provider": "test",
				"code":"123",
				"redirectURL": "https://example.com",
				"createData": {
					"name": "test_new"
				}
			}`),
			Headers: map[string]string{
				// users, test2@example.com
				// (auth with some other user from the same collection to ensure that it is ignored)
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6Im9hcDY0MGNvdDR5cnUycyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.GfJo6EHIobgas_AXt-M-tj5IoQendPnrkMSe9ExuSEY",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				user, err := app.FindAuthRecordByEmail("users", "test@example.com")
				if err != nil {
					t.Fatal(err)
				}

				if user.Verified() {
					t.Fatalf("Expected user %q to be unverified", user.Email())
				}

				// ensure that the old password works
				if !user.ValidatePassword("1234567890") {
					t.Fatalf("Expected password %q to be valid", "1234567890")
				}

				// register the test provider
				auth.Providers["test"] = func() auth.Provider {
					return &oauth2MockProvider{
						AuthUser: &auth.AuthUser{Id: "test_id"},
						Token:    &oauth2.Token{AccessToken: "abc"},
					}
				}

				// add the test provider in the collection
				user.Collection().MFA.Enabled = false
				user.Collection().OAuth2.Enabled = true
				user.Collection().OAuth2.Providers = []core.OAuth2ProviderConfig{{
					Name:         "test",
					ClientId:     "123",
					ClientSecret: "456",
				}}
				if err := app.Save(user.Collection()); err != nil {
					t.Fatal(err)
				}

				// stub linked provider
				ea := core.NewExternalAuth(app)
				ea.SetCollectionRef(user.Collection().Id)
				ea.SetRecordRef(user.Id)
				ea.SetProvider("test")
				ea.SetProviderId("test_id")
				if err := app.Save(ea); err != nil {
					t.Fatal(err)
				}

				// test at least once that the correct request info context is properly loaded
				app.OnRecordAuthRequest().BindFunc(func(e *core.RecordAuthRequestEvent) error {
					info, err := e.RequestInfo()
					if err != nil {
						t.Fatal(err)
					}

					if info.Context != core.RequestInfoContextOAuth2 {
						t.Fatalf("Expected request context %q, got %q", core.RequestInfoContextOAuth2, info.Context)
					}

					return e.Next()
				})
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"record":{`,
				`"token":"`,
				`"meta":{`,
				`"email":"test@example.com"`,
				`"id":"4q1xlclmfloku33"`,
				`"id":"test_id"`,
				`"verified":false`, // shouldn't change
			},
			NotExpectedContent: []string{
				// hidden fields
				`"tokenKey"`,
				`"password"`,
			},
			ExpectedEvents: map[string]int{
				"*":                             0,
				"OnRecordAuthWithOAuth2Request": 1,
				"OnRecordAuthRequest":           1,
				"OnRecordEnrich":                1,
				// ---
				"OnModelCreate":              1,
				"OnModelCreateExecute":       1,
				"OnModelAfterCreateSuccess":  1,
				"OnRecordCreate":             1,
				"OnRecordCreateExecute":      1,
				"OnRecordAfterCreateSuccess": 1,
				// ---
				"OnModelUpdate":              1,
				"OnModelUpdateExecute":       1,
				"OnModelAfterUpdateSuccess":  1,
				"OnRecordUpdate":             1,
				"OnRecordUpdateExecute":      1,
				"OnRecordAfterUpdateSuccess": 1,
				// ---
				"OnModelValidate":  2, // create + update
				"OnRecordValidate": 2,
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				user, err := app.FindAuthRecordByEmail("users", "test@example.com")
				if err != nil {
					t.Fatal(err)
				}

				if name := user.GetString("name"); name != "test1" {
					t.Fatalf("Expected name to not change, got %q", name)
				}

				if user.ValidatePassword("1234567890") {
					t.Fatalf("Expected password %q to be changed", "1234567890")
				}

				devices, err := app.FindAllAuthOriginsByRecord(user)
				if len(devices) != 1 {
					t.Fatalf("Expected only 1 auth origin to be created, got %d (%v)", len(devices), err)
				}
			},
		},
		{
			Name:   "existing linked OAuth2 (verified user)",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-oauth2",
			Body: strings.NewReader(`{
				"provider": "test",
				"code":"123",
				"redirectURL": "https://example.com"
			}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				user, err := app.FindAuthRecordByEmail("users", "test2@example.com")
				if err != nil {
					t.Fatal(err)
				}

				if !user.Verified() {
					t.Fatalf("Expected user %q to be verified", user.Email())
				}

				// ensure that the old password works
				if !user.ValidatePassword("1234567890") {
					t.Fatalf("Expected password %q to be valid", "1234567890")
				}

				// register the test provider
				auth.Providers["test"] = func() auth.Provider {
					return &oauth2MockProvider{
						AuthUser: &auth.AuthUser{Id: "test_id"},
						Token:    &oauth2.Token{AccessToken: "abc"},
					}
				}

				// add the test provider in the collection
				user.Collection().MFA.Enabled = false
				user.Collection().OAuth2.Enabled = true
				user.Collection().OAuth2.Providers = []core.OAuth2ProviderConfig{{
					Name:         "test",
					ClientId:     "123",
					ClientSecret: "456",
				}}
				if err := app.Save(user.Collection()); err != nil {
					t.Fatal(err)
				}

				// stub linked provider
				ea := core.NewExternalAuth(app)
				ea.SetCollectionRef(user.Collection().Id)
				ea.SetRecordRef(user.Id)
				ea.SetProvider("test")
				ea.SetProviderId("test_id")
				if err := app.Save(ea); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"record":{`,
				`"token":"`,
				`"meta":{`,
				`"isNew":false`,
				`"email":"test2@example.com"`,
				`"id":"oap640cot4yru2s"`,
				`"id":"test_id"`,
				`"verified":true`,
			},
			NotExpectedContent: []string{
				// hidden fields
				`"tokenKey"`,
				`"password"`,
			},
			ExpectedEvents: map[string]int{
				"*":                             0,
				"OnRecordAuthWithOAuth2Request": 1,
				"OnRecordAuthRequest":           1,
				"OnRecordEnrich":                1,
				// ---
				"OnModelCreate":              1,
				"OnModelCreateExecute":       1,
				"OnModelAfterCreateSuccess":  1,
				"OnRecordCreate":             1,
				"OnRecordCreateExecute":      1,
				"OnRecordAfterCreateSuccess": 1,
				// ---
				"OnModelValidate":  1,
				"OnRecordValidate": 1,
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				user, err := app.FindAuthRecordByEmail("users", "test2@example.com")
				if err != nil {
					t.Fatal(err)
				}

				if !user.ValidatePassword("1234567890") {
					t.Fatalf("Expected old password %q to be valid", "1234567890")
				}

				devices, err := app.FindAllAuthOriginsByRecord(user)
				if len(devices) != 1 {
					t.Fatalf("Expected only 1 auth origin to be created, got %d (%v)", len(devices), err)
				}
			},
		},
		{
			Name:   "link by email",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-oauth2",
			Body: strings.NewReader(`{
				"provider": "test",
				"code":"123",
				"redirectURL": "https://example.com"
			}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				user, err := app.FindAuthRecordByEmail("users", "test@example.com")
				if err != nil {
					t.Fatal(err)
				}

				if user.Verified() {
					t.Fatalf("Expected user %q to be unverified", user.Email())
				}

				// ensure that the old password works
				if !user.ValidatePassword("1234567890") {
					t.Fatalf("Expected password %q to be valid", "1234567890")
				}

				// register the test provider
				auth.Providers["test"] = func() auth.Provider {
					return &oauth2MockProvider{
						AuthUser: &auth.AuthUser{Id: "test_id", Email: "test@example.com"},
						Token:    &oauth2.Token{AccessToken: "abc"},
					}
				}

				// add the test provider in the collection
				user.Collection().MFA.Enabled = false
				user.Collection().OAuth2.Enabled = true
				user.Collection().OAuth2.Providers = []core.OAuth2ProviderConfig{{
					Name:         "test",
					ClientId:     "123",
					ClientSecret: "456",
				}}
				if err := app.Save(user.Collection()); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"record":{`,
				`"token":"`,
				`"meta":{`,
				`"isNew":false`,
				`"email":"test@example.com"`,
				`"id":"4q1xlclmfloku33"`,
				`"id":"test_id"`,
				`"verified":true`, // should be updated
			},
			NotExpectedContent: []string{
				// hidden fields
				`"tokenKey"`,
				`"password"`,
			},
			ExpectedEvents: map[string]int{
				"*":                             0,
				"OnRecordAuthWithOAuth2Request": 1,
				"OnRecordAuthRequest":           1,
				"OnRecordEnrich":                1,
				// ---
				"OnModelCreate":              2, // authOrigins + externalAuths
				"OnModelCreateExecute":       2,
				"OnModelAfterCreateSuccess":  2,
				"OnRecordCreate":             2,
				"OnRecordCreateExecute":      2,
				"OnRecordAfterCreateSuccess": 2,
				// ---
				"OnModelUpdate":              1, // record password and verified states
				"OnModelUpdateExecute":       1,
				"OnModelAfterUpdateSuccess":  1,
				"OnRecordUpdate":             1,
				"OnRecordUpdateExecute":      1,
				"OnRecordAfterUpdateSuccess": 1,
				// ---
				"OnModelValidate":  3, // record + authOrigins + externalAuths
				"OnRecordValidate": 3,
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				user, err := app.FindAuthRecordByEmail("users", "test@example.com")
				if err != nil {
					t.Fatal(err)
				}

				if user.ValidatePassword("1234567890") {
					t.Fatalf("Expected password %q to be changed", "1234567890")
				}

				devices, err := app.FindAllAuthOriginsByRecord(user)
				if len(devices) != 1 {
					t.Fatalf("Expected only 1 auth origin to be created, got %d (%v)", len(devices), err)
				}
			},
		},
		{
			Name:   "link by fallback user (OAuth2 user with different email)",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-oauth2",
			Body: strings.NewReader(`{
				"provider": "test",
				"code":"123",
				"redirectURL": "https://example.com"
			}`),
			Headers: map[string]string{
				// users, test@example.com
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				user, err := app.FindAuthRecordByEmail("users", "test@example.com")
				if err != nil {
					t.Fatal(err)
				}

				if user.Verified() {
					t.Fatalf("Expected user %q to be unverified", user.Email())
				}

				// ensure that the old password works
				if !user.ValidatePassword("1234567890") {
					t.Fatalf("Expected password %q to be valid", "1234567890")
				}

				// register the test provider
				auth.Providers["test"] = func() auth.Provider {
					return &oauth2MockProvider{
						AuthUser: &auth.AuthUser{
							Id:    "test_id",
							Email: "test2@example.com", // different email -> should be ignored
						},
						Token: &oauth2.Token{AccessToken: "abc"},
					}
				}

				// add the test provider in the collection
				user.Collection().MFA.Enabled = false
				user.Collection().OAuth2.Enabled = true
				user.Collection().OAuth2.Providers = []core.OAuth2ProviderConfig{{
					Name:         "test",
					ClientId:     "123",
					ClientSecret: "456",
				}}
				if err := app.Save(user.Collection()); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"record":{`,
				`"token":"`,
				`"meta":{`,
				`"isNew":false`,
				`"email":"test@example.com"`,
				`"id":"4q1xlclmfloku33"`,
				`"id":"test_id"`,
				`"verified":false`, // shouldn't change because the OAuth2 user email is different
			},
			NotExpectedContent: []string{
				// hidden fields
				`"tokenKey"`,
				`"password"`,
			},
			ExpectedEvents: map[string]int{
				"*":                             0,
				"OnRecordAuthWithOAuth2Request": 1,
				"OnRecordAuthRequest":           1,
				"OnRecordEnrich":                1,
				// ---
				"OnModelCreate":              2, // authOrigins + externalAuths
				"OnModelCreateExecute":       2,
				"OnModelAfterCreateSuccess":  2,
				"OnRecordCreate":             2,
				"OnRecordCreateExecute":      2,
				"OnRecordAfterCreateSuccess": 2,
				// ---
				"OnModelValidate":  2,
				"OnRecordValidate": 2,
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				user, err := app.FindAuthRecordByEmail("users", "test@example.com")
				if err != nil {
					t.Fatal(err)
				}

				if !user.ValidatePassword("1234567890") {
					t.Fatalf("Expected password %q not to be changed", "1234567890")
				}

				devices, err := app.FindAllAuthOriginsByRecord(user)
				if len(devices) != 1 {
					t.Fatalf("Expected only 1 auth origin to be created, got %d (%v)", len(devices), err)
				}
			},
		},
		{
			Name:   "link by fallback user (user without email)",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-oauth2",
			Body: strings.NewReader(`{
				"provider": "test",
				"code":"123",
				"redirectURL": "https://example.com"
			}`),
			Headers: map[string]string{
				// users, test@example.com
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				user, err := app.FindAuthRecordByEmail("users", "test@example.com")
				if err != nil {
					t.Fatal(err)
				}

				if user.Verified() {
					t.Fatalf("Expected user %q to be unverified", user.Email())
				}

				// ensure that the old password works
				if !user.ValidatePassword("1234567890") {
					t.Fatalf("Expected password %q to be valid", "1234567890")
				}

				oldTokenKey := user.TokenKey()

				// manually unset the user email
				user.SetEmail("")
				if err = app.Save(user); err != nil {
					t.Fatal(err)
				}

				// resave with the old token key since the email change above
				// would change it and will make the password token invalid
				user.SetTokenKey(oldTokenKey)
				if err = app.Save(user); err != nil {
					t.Fatalf("Failed to restore original user tokenKey: %v", err)
				}

				// register the test provider
				auth.Providers["test"] = func() auth.Provider {
					return &oauth2MockProvider{
						AuthUser: &auth.AuthUser{
							Id:    "test_id",
							Email: "test_oauth2@example.com",
						},
						Token: &oauth2.Token{AccessToken: "abc"},
					}
				}

				// add the test provider in the collection
				user.Collection().MFA.Enabled = false
				user.Collection().OAuth2.Enabled = true
				user.Collection().OAuth2.Providers = []core.OAuth2ProviderConfig{{
					Name:         "test",
					ClientId:     "123",
					ClientSecret: "456",
				}}
				if err := app.Save(user.Collection()); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"record":{`,
				`"token":"`,
				`"meta":{`,
				`"isNew":false`,
				`"email":"test_oauth2@example.com"`,
				`"id":"4q1xlclmfloku33"`,
				`"id":"test_id"`,
				`"verified":true`,
			},
			NotExpectedContent: []string{
				// hidden fields
				`"tokenKey"`,
				`"password"`,
			},
			ExpectedEvents: map[string]int{
				"*":                             0,
				"OnRecordAuthWithOAuth2Request": 1,
				"OnRecordAuthRequest":           1,
				"OnRecordEnrich":                1,
				// ---
				"OnModelCreate":              2, // authOrigins + externalAuths
				"OnModelCreateExecute":       2,
				"OnModelAfterCreateSuccess":  2,
				"OnRecordCreate":             2,
				"OnRecordCreateExecute":      2,
				"OnRecordAfterCreateSuccess": 2,
				// ---
				"OnModelUpdate":              1, // record email set
				"OnModelUpdateExecute":       1,
				"OnModelAfterUpdateSuccess":  1,
				"OnRecordUpdate":             1,
				"OnRecordUpdateExecute":      1,
				"OnRecordAfterUpdateSuccess": 1,
				// ---
				"OnModelValidate":  3, // record + authOrigins + externalAuths
				"OnRecordValidate": 3,
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				user, err := app.FindAuthRecordByEmail("users", "test_oauth2@example.com")
				if err != nil {
					t.Fatal(err)
				}

				if !user.ValidatePassword("1234567890") {
					t.Fatalf("Expected password %q not to be changed", "1234567890")
				}

				devices, err := app.FindAllAuthOriginsByRecord(user)
				if len(devices) != 1 {
					t.Fatalf("Expected only 1 auth origin to be created, got %d (%v)", len(devices), err)
				}
			},
		},
		{
			Name:   "link by fallback user (unverified user with matching email)",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-oauth2",
			Body: strings.NewReader(`{
				"provider": "test",
				"code":"123",
				"redirectURL": "https://example.com"
			}`),
			Headers: map[string]string{
				// users, test@example.com
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				user, err := app.FindAuthRecordByEmail("users", "test@example.com")
				if err != nil {
					t.Fatal(err)
				}

				if user.Verified() {
					t.Fatalf("Expected user %q to be unverified", user.Email())
				}

				// ensure that the old password works
				if !user.ValidatePassword("1234567890") {
					t.Fatalf("Expected password %q to be valid", "1234567890")
				}

				// register the test provider
				auth.Providers["test"] = func() auth.Provider {
					return &oauth2MockProvider{
						AuthUser: &auth.AuthUser{
							Id:    "test_id",
							Email: "test@example.com", // matching email -> should be marked as verified
						},
						Token: &oauth2.Token{AccessToken: "abc"},
					}
				}

				// add the test provider in the collection
				user.Collection().MFA.Enabled = false
				user.Collection().OAuth2.Enabled = true
				user.Collection().OAuth2.Providers = []core.OAuth2ProviderConfig{{
					Name:         "test",
					ClientId:     "123",
					ClientSecret: "456",
				}}
				if err := app.Save(user.Collection()); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"record":{`,
				`"token":"`,
				`"meta":{`,
				`"isNew":false`,
				`"email":"test@example.com"`,
				`"id":"4q1xlclmfloku33"`,
				`"id":"test_id"`,
				`"verified":true`,
			},
			NotExpectedContent: []string{
				// hidden fields
				`"tokenKey"`,
				`"password"`,
			},
			ExpectedEvents: map[string]int{
				"*":                             0,
				"OnRecordAuthWithOAuth2Request": 1,
				"OnRecordAuthRequest":           1,
				"OnRecordEnrich":                1,
				// ---
				"OnModelCreate":              2, // authOrigins + externalAuths
				"OnModelCreateExecute":       2,
				"OnModelAfterCreateSuccess":  2,
				"OnRecordCreate":             2,
				"OnRecordCreateExecute":      2,
				"OnRecordAfterCreateSuccess": 2,
				// ---
				"OnModelUpdate":              1, // record verified update
				"OnModelUpdateExecute":       1,
				"OnModelAfterUpdateSuccess":  1,
				"OnRecordUpdate":             1,
				"OnRecordUpdateExecute":      1,
				"OnRecordAfterUpdateSuccess": 1,
				// ---
				"OnModelValidate":  3, // record + authOrigins + externalAuths
				"OnRecordValidate": 3,
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				user, err := app.FindAuthRecordByEmail("users", "test@example.com")
				if err != nil {
					t.Fatal(err)
				}

				if !user.ValidatePassword("1234567890") {
					t.Fatalf("Expected password %q not to be changed", "1234567890")
				}

				devices, err := app.FindAllAuthOriginsByRecord(user)
				if len(devices) != 1 {
					t.Fatalf("Expected only 1 auth origin to be created, got %d (%v)", len(devices), err)
				}
			},
		},
		{
			Name:   "creating user (no extra create data or custom fields mapping)",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-oauth2",
			Body: strings.NewReader(`{
				"provider": "test",
				"code":"123",
				"redirectURL": "https://example.com"
			}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				usersCol, err := app.FindCollectionByNameOrId("users")
				if err != nil {
					t.Fatal(err)
				}

				// register the test provider
				auth.Providers["test"] = func() auth.Provider {
					return &oauth2MockProvider{
						AuthUser: &auth.AuthUser{Id: "test_id"},
						Token:    &oauth2.Token{AccessToken: "abc"},
					}
				}

				// add the test provider in the collection
				usersCol.MFA.Enabled = false
				usersCol.OAuth2.Enabled = true
				usersCol.OAuth2.Providers = []core.OAuth2ProviderConfig{{
					Name:         "test",
					ClientId:     "123",
					ClientSecret: "456",
				}}
				if err := app.Save(usersCol); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"record":{`,
				`"token":"`,
				`"meta":{`,
				`"isNew":true`,
				`"email":""`,
				`"id":"test_id"`,
				`"verified":true`,
			},
			NotExpectedContent: []string{
				// hidden fields
				`"tokenKey"`,
				`"password"`,
			},
			ExpectedEvents: map[string]int{
				"*":                             0,
				"OnRecordAuthWithOAuth2Request": 1,
				"OnRecordAuthRequest":           1,
				"OnRecordCreateRequest":         1,
				"OnRecordEnrich":                2, // the auth response and from the create request
				// ---
				"OnModelCreate":              3, // record + authOrigins + externalAuths
				"OnModelCreateExecute":       3,
				"OnModelAfterCreateSuccess":  3,
				"OnRecordCreate":             3,
				"OnRecordCreateExecute":      3,
				"OnRecordAfterCreateSuccess": 3,
				// ---
				"OnModelUpdate":              1, // created record verified state change
				"OnModelUpdateExecute":       1,
				"OnModelAfterUpdateSuccess":  1,
				"OnRecordUpdate":             1,
				"OnRecordUpdateExecute":      1,
				"OnRecordAfterUpdateSuccess": 1,
				// ---
				"OnModelValidate":  4,
				"OnRecordValidate": 4,
			},
		},
		{
			Name:   "creating user (submit failure - form auth fields validator)",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-oauth2",
			Body: strings.NewReader(`{
				"provider": "test",
				"code":"123",
				"redirectURL": "https://example.com",
				"createData": {
					"verified": true,
					"email": "invalid",
					"rel": "invalid",
					"file": "invalid"
				}
			}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				usersCol, err := app.FindCollectionByNameOrId("users")
				if err != nil {
					t.Fatal(err)
				}

				// register the test provider
				auth.Providers["test"] = func() auth.Provider {
					return &oauth2MockProvider{
						AuthUser: &auth.AuthUser{Id: "test_id"},
						Token:    &oauth2.Token{AccessToken: "abc"},
					}
				}

				// add the test provider in the collection
				usersCol.MFA.Enabled = false
				usersCol.OAuth2.Enabled = true
				usersCol.OAuth2.Providers = []core.OAuth2ProviderConfig{{
					Name:         "test",
					ClientId:     "123",
					ClientSecret: "456",
				}}
				if err := app.Save(usersCol); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"verified":{"code":"validation_values_mismatch"`,
			},
			NotExpectedContent: []string{
				`"email":`, // ignored because the record validator never ran
				`"rel":`,   // ignored because the record validator never ran
				`"file":`,  // ignored because the record validator never ran
			},
			ExpectedEvents: map[string]int{
				"*":                             0,
				"OnRecordAuthWithOAuth2Request": 1,
				"OnRecordCreateRequest":         1,
			},
		},
		{
			Name:   "creating user (submit failure - record fields validator)",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-oauth2",
			Body: strings.NewReader(`{
				"provider": "test",
				"code":"123",
				"redirectURL": "https://example.com",
				"createData": {
					"email": "invalid",
					"rel": "invalid",
					"file": "invalid"
				}
			}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				usersCol, err := app.FindCollectionByNameOrId("users")
				if err != nil {
					t.Fatal(err)
				}

				// register the test provider
				auth.Providers["test"] = func() auth.Provider {
					return &oauth2MockProvider{
						AuthUser: &auth.AuthUser{Id: "test_id"},
						Token:    &oauth2.Token{AccessToken: "abc"},
					}
				}

				// add the test provider in the collection
				usersCol.MFA.Enabled = false
				usersCol.OAuth2.Enabled = true
				usersCol.OAuth2.Providers = []core.OAuth2ProviderConfig{{
					Name:         "test",
					ClientId:     "123",
					ClientSecret: "456",
				}}
				if err := app.Save(usersCol); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"email":{"code":"validation_is_email"`,
				`"rel":{"code":"validation_missing_rel_records"`,
				`"file":{"code":"validation_invalid_file"`,
			},
			ExpectedEvents: map[string]int{
				"*":                             0,
				"OnRecordAuthWithOAuth2Request": 1,
				"OnRecordCreateRequest":         1,
				"OnModelValidate":               1,
				"OnRecordValidate":              1,
				"OnModelCreate":                 1,
				"OnModelAfterCreateError":       1,
				"OnRecordCreate":                1,
				"OnRecordAfterCreateError":      1,
			},
		},
		{
			Name:   "creating user (valid create data with empty submitted email)",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-oauth2",
			Body: strings.NewReader(`{
				"provider": "test",
				"code":"123",
				"redirectURL": "https://example.com",
				"createData": {
					"email": "",
					"emailVisibility": true,
					"password": "1234567890",
					"passwordConfirm": "1234567890",
					"name": "test_name",
					"username": "test_username",
					"rel": "0yxhwia2amd8gec"
				}
			}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				usersCol, err := app.FindCollectionByNameOrId("users")
				if err != nil {
					t.Fatal(err)
				}

				// register the test provider
				auth.Providers["test"] = func() auth.Provider {
					return &oauth2MockProvider{
						AuthUser: &auth.AuthUser{Id: "test_id"},
						Token:    &oauth2.Token{AccessToken: "abc"},
					}
				}

				// add the test provider in the collection
				usersCol.MFA.Enabled = false
				usersCol.OAuth2.Enabled = true
				usersCol.OAuth2.Providers = []core.OAuth2ProviderConfig{{
					Name:         "test",
					ClientId:     "123",
					ClientSecret: "456",
				}}
				if err := app.Save(usersCol); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"isNew":true`,
				`"email":""`,
				`"emailVisibility":true`,
				`"name":"test_name"`,
				`"username":"test_username"`,
				`"verified":true`,
				`"rel":"0yxhwia2amd8gec"`,
			},
			NotExpectedContent: []string{
				// hidden fields
				`"tokenKey"`,
				`"password"`,
			},
			ExpectedEvents: map[string]int{
				"*":                             0,
				"OnRecordAuthWithOAuth2Request": 1,
				"OnRecordAuthRequest":           1,
				"OnRecordCreateRequest":         1,
				"OnRecordEnrich":                2, // the auth response and from the create request
				// ---
				"OnModelCreate":              3, // record + authOrigins + externalAuths
				"OnModelCreateExecute":       3,
				"OnModelAfterCreateSuccess":  3,
				"OnRecordCreate":             3,
				"OnRecordCreateExecute":      3,
				"OnRecordAfterCreateSuccess": 3,
				// ---
				"OnModelUpdate":              1, // created record verified state change
				"OnModelUpdateExecute":       1,
				"OnModelAfterUpdateSuccess":  1,
				"OnRecordUpdate":             1,
				"OnRecordUpdateExecute":      1,
				"OnRecordAfterUpdateSuccess": 1,
				// ---
				"OnModelValidate":  4,
				"OnRecordValidate": 4,
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				user, err := app.FindFirstRecordByData("users", "username", "test_username")
				if err != nil {
					t.Fatal(err)
				}

				if !user.ValidatePassword("1234567890") {
					t.Fatalf("Expected password %q to be valid", "1234567890")
				}
			},
		},
		{
			Name:   "creating user (valid create data with non-empty valid submitted email)",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-oauth2",
			Body: strings.NewReader(`{
				"provider": "test",
				"code":"123",
				"redirectURL": "https://example.com",
				"createData": {
					"email": "test_create@example.com",
					"emailVisibility": true,
					"password": "1234567890",
					"passwordConfirm": "1234567890",
					"name": "test_name",
					"username": "test_username",
					"rel": "0yxhwia2amd8gec"
				}
			}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				usersCol, err := app.FindCollectionByNameOrId("users")
				if err != nil {
					t.Fatal(err)
				}

				// register the test provider
				auth.Providers["test"] = func() auth.Provider {
					return &oauth2MockProvider{
						AuthUser: &auth.AuthUser{
							Id:    "test_id",
							Email: "oauth2@example.com", // should be ignored because of the explicit submitted email
						},
						Token: &oauth2.Token{AccessToken: "abc"},
					}
				}

				// add the test provider in the collection
				usersCol.MFA.Enabled = false
				usersCol.OAuth2.Enabled = true
				usersCol.OAuth2.Providers = []core.OAuth2ProviderConfig{{
					Name:         "test",
					ClientId:     "123",
					ClientSecret: "456",
				}}
				if err := app.Save(usersCol); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"email":"test_create@example.com"`,
				`"emailVisibility":true`,
				`"name":"test_name"`,
				`"username":"test_username"`,
				`"verified":false`,
				`"rel":"0yxhwia2amd8gec"`,
			},
			NotExpectedContent: []string{
				// hidden fields
				`"tokenKey"`,
				`"password"`,
			},
			ExpectedEvents: map[string]int{
				"*":                             0,
				"OnRecordAuthWithOAuth2Request": 1,
				"OnRecordAuthRequest":           1,
				"OnRecordCreateRequest":         1,
				"OnRecordEnrich":                2, // the auth response and from the create request
				// ---
				"OnModelCreate":              3, // record + authOrigins + externalAuths
				"OnModelCreateExecute":       3,
				"OnModelAfterCreateSuccess":  3,
				"OnRecordCreate":             3,
				"OnRecordCreateExecute":      3,
				"OnRecordAfterCreateSuccess": 3,
				// ---
				"OnModelValidate":  3,
				"OnRecordValidate": 3,
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				user, err := app.FindFirstRecordByData("users", "username", "test_username")
				if err != nil {
					t.Fatal(err)
				}

				if !user.ValidatePassword("1234567890") {
					t.Fatalf("Expected password %q to be valid", "1234567890")
				}
			},
		},
		{
			Name:   "creating user (with mapped OAuth2 fields and avatarURL->file field)",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-oauth2",
			Body: strings.NewReader(`{
				"provider": "test",
				"code":"123",
				"redirectURL": "https://example.com",
				"createData": {
					"name": "test_name",
					"emailVisibility": true,
					"rel": "0yxhwia2amd8gec"
				}
			}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				usersCol, err := app.FindCollectionByNameOrId("users")
				if err != nil {
					t.Fatal(err)
				}

				// register the test provider
				auth.Providers["test"] = func() auth.Provider {
					return &oauth2MockProvider{
						AuthUser: &auth.AuthUser{
							Id:        "oauth2_id",
							Email:     "oauth2@example.com",
							Username:  "oauth2_username",
							AvatarURL: server.URL + "/oauth2_avatar.png",
						},
						Token: &oauth2.Token{AccessToken: "abc"},
					}
				}

				// add the test provider in the collection
				usersCol.MFA.Enabled = false
				usersCol.OAuth2.Enabled = true
				usersCol.OAuth2.Providers = []core.OAuth2ProviderConfig{{
					Name:         "test",
					ClientId:     "123",
					ClientSecret: "456",
				}}
				usersCol.OAuth2.MappedFields = core.OAuth2KnownFields{
					Username:  "name", // should be ignored because of the explicit submitted value
					Id:        "username",
					AvatarURL: "avatar",
				}
				if err := app.Save(usersCol); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"isNew":true`,
				`"email":"oauth2@example.com"`,
				`"emailVisibility":true`,
				`"name":"test_name"`,
				`"username":"oauth2_username"`,
				`"verified":true`,
				`"rel":"0yxhwia2amd8gec"`,
				`"avatar":"oauth2_avatar_`,
			},
			NotExpectedContent: []string{
				// hidden fields
				`"tokenKey"`,
				`"password"`,
			},
			ExpectedEvents: map[string]int{
				"*":                             0,
				"OnRecordAuthWithOAuth2Request": 1,
				"OnRecordAuthRequest":           1,
				"OnRecordCreateRequest":         1,
				"OnRecordEnrich":                2, // the auth response and from the create request
				// ---
				"OnModelCreate":              3, // record + authOrigins + externalAuths
				"OnModelCreateExecute":       3,
				"OnModelAfterCreateSuccess":  3,
				"OnRecordCreate":             3,
				"OnRecordCreateExecute":      3,
				"OnRecordAfterCreateSuccess": 3,
				// ---
				"OnModelUpdate":              1, // created record verified state change
				"OnModelUpdateExecute":       1,
				"OnModelAfterUpdateSuccess":  1,
				"OnRecordUpdate":             1,
				"OnRecordUpdateExecute":      1,
				"OnRecordAfterUpdateSuccess": 1,
				// ---
				"OnModelValidate":  4,
				"OnRecordValidate": 4,
			},
		},
		{
			Name:   "creating user (with mapped OAuth2 avatarURL field but empty OAuth2User.avatarURL value)",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-oauth2",
			Body: strings.NewReader(`{
				"provider": "test",
				"code":"123",
				"redirectURL": "https://example.com"
			}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				usersCol, err := app.FindCollectionByNameOrId("users")
				if err != nil {
					t.Fatal(err)
				}

				// register the test provider
				auth.Providers["test"] = func() auth.Provider {
					return &oauth2MockProvider{
						AuthUser: &auth.AuthUser{
							Id:        "oauth2_id",
							Email:     "oauth2@example.com",
							AvatarURL: "",
						},
						Token: &oauth2.Token{AccessToken: "abc"},
					}
				}

				// add the test provider in the collection
				usersCol.MFA.Enabled = false
				usersCol.OAuth2.Enabled = true
				usersCol.OAuth2.Providers = []core.OAuth2ProviderConfig{{
					Name:         "test",
					ClientId:     "123",
					ClientSecret: "456",
				}}
				usersCol.OAuth2.MappedFields = core.OAuth2KnownFields{
					AvatarURL: "avatar",
				}
				if err := app.Save(usersCol); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"isNew":true`,
				`"email":"oauth2@example.com"`,
				`"emailVisibility":false`,
				`"verified":true`,
				`"avatar":""`,
			},
			NotExpectedContent: []string{
				// hidden fields
				`"tokenKey"`,
				`"password"`,
			},
			ExpectedEvents: map[string]int{
				"*":                             0,
				"OnRecordAuthWithOAuth2Request": 1,
				"OnRecordAuthRequest":           1,
				"OnRecordCreateRequest":         1,
				"OnRecordEnrich":                2, // the auth response and from the create request
				// ---
				"OnModelCreate":              3, // record + authOrigins + externalAuths
				"OnModelCreateExecute":       3,
				"OnModelAfterCreateSuccess":  3,
				"OnRecordCreate":             3,
				"OnRecordCreateExecute":      3,
				"OnRecordAfterCreateSuccess": 3,
				// ---
				"OnModelUpdate":              1, // created record verified state change
				"OnModelUpdateExecute":       1,
				"OnModelAfterUpdateSuccess":  1,
				"OnRecordUpdate":             1,
				"OnRecordUpdateExecute":      1,
				"OnRecordAfterUpdateSuccess": 1,
				// ---
				"OnModelValidate":  4,
				"OnRecordValidate": 4,
			},
		},
		{
			Name:   "creating user (with mapped OAuth2 fields, case-sensitive username and avatarURL->non-file field)",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-oauth2",
			Body: strings.NewReader(`{
				"provider": "test",
				"code":"123",
				"redirectURL": "https://example.com"
			}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				usersCol, err := app.FindCollectionByNameOrId("users")
				if err != nil {
					t.Fatal(err)
				}

				// register the test provider
				auth.Providers["test"] = func() auth.Provider {
					return &oauth2MockProvider{
						AuthUser: &auth.AuthUser{
							Id:        "oauth2_id",
							Email:     "oauth2@example.com",
							Username:  "tESt2_username", // wouldn't match with existing because the related field index is case-sensitive
							Name:      "oauth2_name",
							AvatarURL: server.URL + "/oauth2_avatar.png",
						},
						Token: &oauth2.Token{AccessToken: "abc"},
					}
				}

				// add the test provider in the collection
				usersCol.MFA.Enabled = false
				usersCol.OAuth2.Enabled = true
				usersCol.OAuth2.Providers = []core.OAuth2ProviderConfig{{
					Name:         "test",
					ClientId:     "123",
					ClientSecret: "456",
				}}
				usersCol.OAuth2.MappedFields = core.OAuth2KnownFields{
					Username:  "username",
					AvatarURL: "name",
				}
				if err := app.Save(usersCol); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"isNew":true`,
				`"email":"oauth2@example.com"`,
				`"emailVisibility":false`,
				`"username":"tESt2_username"`,
				`"name":"http://127.`,
				`"verified":true`,
				`"avatar":""`,
			},
			NotExpectedContent: []string{
				// hidden fields
				`"tokenKey"`,
				`"password"`,
			},
			ExpectedEvents: map[string]int{
				"*":                             0,
				"OnRecordAuthWithOAuth2Request": 1,
				"OnRecordAuthRequest":           1,
				"OnRecordCreateRequest":         1,
				"OnRecordEnrich":                2, // the auth response and from the create request
				// ---
				"OnModelCreate":              3, // record + authOrigins + externalAuths
				"OnModelCreateExecute":       3,
				"OnModelAfterCreateSuccess":  3,
				"OnRecordCreate":             3,
				"OnRecordCreateExecute":      3,
				"OnRecordAfterCreateSuccess": 3,
				// ---
				"OnModelUpdate":              1, // created record verified state change
				"OnModelUpdateExecute":       1,
				"OnModelAfterUpdateSuccess":  1,
				"OnRecordUpdate":             1,
				"OnRecordUpdateExecute":      1,
				"OnRecordAfterUpdateSuccess": 1,
				// ---
				"OnModelValidate":  4,
				"OnRecordValidate": 4,
			},
		},
		{
			Name:   "creating user (with mapped OAuth2 fields and duplicated case-insensitive username)",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-oauth2",
			Body: strings.NewReader(`{
				"provider": "test",
				"code":"123",
				"redirectURL": "https://example.com"
			}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				usersCol, err := app.FindCollectionByNameOrId("users")
				if err != nil {
					t.Fatal(err)
				}

				// register the test provider
				auth.Providers["test"] = func() auth.Provider {
					return &oauth2MockProvider{
						AuthUser: &auth.AuthUser{
							Id:       "oauth2_id",
							Email:    "oauth2@example.com",
							Username: "tESt2_username",
							Name:     "oauth2_name",
						},
						Token: &oauth2.Token{AccessToken: "abc"},
					}
				}

				// make the username index case-insensitive to ensure that case-insensitive match is used
				index, ok := dbutils.FindSingleColumnUniqueIndex(usersCol.Indexes, "username")
				if ok {
					index.Columns[0].Collate = "nocase"
					usersCol.RemoveIndex(index.IndexName)
					usersCol.Indexes = append(usersCol.Indexes, index.Build())
				}

				// add the test provider in the collection
				usersCol.MFA.Enabled = false
				usersCol.OAuth2.Enabled = true
				usersCol.OAuth2.Providers = []core.OAuth2ProviderConfig{{
					Name:         "test",
					ClientId:     "123",
					ClientSecret: "456",
				}}
				usersCol.OAuth2.MappedFields = core.OAuth2KnownFields{
					Username: "username",
				}
				if err := app.Save(usersCol); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"isNew":true`,
				`"email":"oauth2@example.com"`,
				`"emailVisibility":false`,
				`"verified":true`,
				`"avatar":""`,
				`"username":"users`,
			},
			NotExpectedContent: []string{
				// hidden fields
				`"tokenKey"`,
				`"password"`,
			},
			ExpectedEvents: map[string]int{
				"*":                             0,
				"OnRecordAuthWithOAuth2Request": 1,
				"OnRecordAuthRequest":           1,
				"OnRecordCreateRequest":         1,
				"OnRecordEnrich":                2, // the auth response and from the create request
				// ---
				"OnModelCreate":              3, // record + authOrigins + externalAuths
				"OnModelCreateExecute":       3,
				"OnModelAfterCreateSuccess":  3,
				"OnRecordCreate":             3,
				"OnRecordCreateExecute":      3,
				"OnRecordAfterCreateSuccess": 3,
				// ---
				"OnModelUpdate":              1, // created record verified state change
				"OnModelUpdateExecute":       1,
				"OnModelAfterUpdateSuccess":  1,
				"OnRecordUpdate":             1,
				"OnRecordUpdateExecute":      1,
				"OnRecordAfterUpdateSuccess": 1,
				// ---
				"OnModelValidate":  4,
				"OnRecordValidate": 4,
			},
		},
		{
			Name:   "creating user (with mapped OAuth2 fields and username that doesn't match the field validations)",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-oauth2",
			Body: strings.NewReader(`{
				"provider": "test",
				"code":"123",
				"redirectURL": "https://example.com"
			}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				usersCol, err := app.FindCollectionByNameOrId("users")
				if err != nil {
					t.Fatal(err)
				}

				// register the test provider
				auth.Providers["test"] = func() auth.Provider {
					return &oauth2MockProvider{
						AuthUser: &auth.AuthUser{
							Id:       "oauth2_id",
							Email:    "oauth2@example.com",
							Username: "!@invalid",
							Name:     "oauth2_name",
						},
						Token: &oauth2.Token{AccessToken: "abc"},
					}
				}

				// add the test provider in the collection
				usersCol.MFA.Enabled = false
				usersCol.OAuth2.Enabled = true
				usersCol.OAuth2.Providers = []core.OAuth2ProviderConfig{{
					Name:         "test",
					ClientId:     "123",
					ClientSecret: "456",
				}}
				usersCol.OAuth2.MappedFields = core.OAuth2KnownFields{
					Username: "username",
				}
				if err := app.Save(usersCol); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"isNew":true`,
				`"email":"oauth2@example.com"`,
				`"emailVisibility":false`,
				`"verified":true`,
				`"avatar":""`,
				`"username":"users`,
			},
			NotExpectedContent: []string{
				// hidden fields
				`"tokenKey"`,
				`"password"`,
			},
			ExpectedEvents: map[string]int{
				"*":                             0,
				"OnRecordAuthWithOAuth2Request": 1,
				"OnRecordAuthRequest":           1,
				"OnRecordCreateRequest":         1,
				"OnRecordEnrich":                2, // the auth response and from the create request
				// ---
				"OnModelCreate":              3, // record + authOrigins + externalAuths
				"OnModelCreateExecute":       3,
				"OnModelAfterCreateSuccess":  3,
				"OnRecordCreate":             3,
				"OnRecordCreateExecute":      3,
				"OnRecordAfterCreateSuccess": 3,
				// ---
				"OnModelUpdate":              1, // created record verified state change
				"OnModelUpdateExecute":       1,
				"OnModelAfterUpdateSuccess":  1,
				"OnRecordUpdate":             1,
				"OnRecordUpdateExecute":      1,
				"OnRecordAfterUpdateSuccess": 1,
				// ---
				"OnModelValidate":  4,
				"OnRecordValidate": 4,
			},
		},

		// rate limit checks
		// -----------------------------------------------------------
		{
			Name:   "RateLimit rule - users:authWithOAuth2",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-oauth2",
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				app.Settings().RateLimits.Enabled = true
				app.Settings().RateLimits.Rules = []core.RateLimitRule{
					{MaxRequests: 100, Label: "abc"},
					{MaxRequests: 100, Label: "*:authWithOAuth2"},
					{MaxRequests: 100, Label: "users:auth"},
					{MaxRequests: 0, Label: "users:authWithOAuth2"},
				}
			},
			ExpectedStatus:  429,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "RateLimit rule - *:authWithOAuth2",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-oauth2",
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				app.Settings().RateLimits.Enabled = true
				app.Settings().RateLimits.Rules = []core.RateLimitRule{
					{MaxRequests: 100, Label: "abc"},
					{MaxRequests: 100, Label: "*:auth"},
					{MaxRequests: 0, Label: "*:authWithOAuth2"},
				}
			},
			ExpectedStatus:  429,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "RateLimit tag - users:auth",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-oauth2",
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				app.Settings().RateLimits.Enabled = true
				app.Settings().RateLimits.Rules = []core.RateLimitRule{
					{MaxRequests: 100, Label: "abc"},
					{MaxRequests: 100, Label: "*:authWithOAuth2"},
					{MaxRequests: 0, Label: "users:auth"},
				}
			},
			ExpectedStatus:  429,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "RateLimit tag - *:auth",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-oauth2",
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
