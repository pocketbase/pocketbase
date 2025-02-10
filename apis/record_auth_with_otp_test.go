package apis_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestRecordAuthWithOTP(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "not an auth collection",
			Method:          http.MethodPost,
			URL:             "/api/collections/demo1/auth-with-otp",
			Body:            strings.NewReader(`{"otpId":"test","password":"123456"}`),
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "auth collection with disabled otp",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-otp",
			Body:   strings.NewReader(`{"otpId":"test","password":"123456"}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				usersCol, err := app.FindCollectionByNameOrId("users")
				if err != nil {
					t.Fatal(err)
				}

				usersCol.OTP.Enabled = false

				if err := app.Save(usersCol); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:            "invalid body",
			Method:          http.MethodPost,
			URL:             "/api/collections/users/auth-with-otp",
			Body:            strings.NewReader(`{"email`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:           "empty body",
			Method:         http.MethodPost,
			URL:            "/api/collections/users/auth-with-otp",
			Body:           strings.NewReader(``),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"otpId":{"code":"validation_required"`,
				`"password":{"code":"validation_required"`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
		{
			Name:   "invalid request data",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-otp",
			Body: strings.NewReader(`{
				"otpId":"` + strings.Repeat("a", 256) + `",
				"password":"` + strings.Repeat("a", 72) + `"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"otpId":{"code":"validation_length_out_of_range"`,
				`"password":{"code":"validation_length_out_of_range"`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
		{
			Name:   "missing otp",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-otp",
			Body: strings.NewReader(`{
				"otpId":"missing",
				"password":"123456"
			}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				user, err := app.FindAuthRecordByEmail("users", "test@example.com")
				if err != nil {
					t.Fatal(err)
				}
				otp := core.NewOTP(app)
				otp.Id = strings.Repeat("a", 15)
				otp.SetCollectionRef(user.Collection().Id)
				otp.SetRecordRef(user.Id)
				otp.SetPassword("123456")
				if err := app.Save(otp); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "otp for different collection",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-otp",
			Body: strings.NewReader(`{
				"otpId":"` + strings.Repeat("a", 15) + `",
				"password":"123456"
			}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				client, err := app.FindAuthRecordByEmail("clients", "test@example.com")
				if err != nil {
					t.Fatal(err)
				}
				otp := core.NewOTP(app)
				otp.Id = strings.Repeat("a", 15)
				otp.SetCollectionRef(client.Collection().Id)
				otp.SetRecordRef(client.Id)
				otp.SetPassword("123456")
				if err := app.Save(otp); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "otp with wrong password",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-otp",
			Body: strings.NewReader(`{
				"otpId":"` + strings.Repeat("a", 15) + `",
				"password":"123456"
			}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				user, err := app.FindAuthRecordByEmail("users", "test@example.com")
				if err != nil {
					t.Fatal(err)
				}
				otp := core.NewOTP(app)
				otp.Id = strings.Repeat("a", 15)
				otp.SetCollectionRef(user.Collection().Id)
				otp.SetRecordRef(user.Id)
				otp.SetPassword("1234567890")
				if err := app.Save(otp); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "expired otp with valid password",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-otp",
			Body: strings.NewReader(`{
				"otpId":"` + strings.Repeat("a", 15) + `",
				"password":"123456"
			}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				user, err := app.FindAuthRecordByEmail("users", "test@example.com")
				if err != nil {
					t.Fatal(err)
				}
				otp := core.NewOTP(app)
				otp.Id = strings.Repeat("a", 15)
				otp.SetCollectionRef(user.Collection().Id)
				otp.SetRecordRef(user.Id)
				otp.SetPassword("123456")
				expiredDate := types.NowDateTime().AddDate(-3, 0, 0)
				otp.SetRaw("created", expiredDate)
				otp.SetRaw("updated", expiredDate)
				if err := app.Save(otp); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "valid otp with valid password (enabled MFA)",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-otp",
			Body: strings.NewReader(`{
				"otpId":"` + strings.Repeat("a", 15) + `",
				"password":"123456"
			}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				user, err := app.FindAuthRecordByEmail("users", "test@example.com")
				if err != nil {
					t.Fatal(err)
				}
				otp := core.NewOTP(app)
				otp.Id = strings.Repeat("a", 15)
				otp.SetCollectionRef(user.Collection().Id)
				otp.SetRecordRef(user.Id)
				otp.SetPassword("123456")
				if err := app.Save(otp); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"mfaId":"`},
			ExpectedEvents: map[string]int{
				"*":                          0,
				"OnRecordAuthWithOTPRequest": 1,
				"OnRecordAuthRequest":        1,
				// ---
				"OnModelValidate":           1,
				"OnModelCreate":             1, // mfa record
				"OnModelCreateExecute":      1,
				"OnModelAfterCreateSuccess": 1,
				"OnModelDelete":             1, // otp delete
				"OnModelDeleteExecute":      1,
				"OnModelAfterDeleteSuccess": 1,
				// ---
				"OnRecordValidate":           1,
				"OnRecordCreate":             1,
				"OnRecordCreateExecute":      1,
				"OnRecordAfterCreateSuccess": 1,
				"OnRecordDelete":             1,
				"OnRecordDeleteExecute":      1,
				"OnRecordAfterDeleteSuccess": 1,
			},
		},
		{
			Name:   "valid otp with valid password and empty sentTo (disabled MFA)",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-otp",
			Body: strings.NewReader(`{
				"otpId":"` + strings.Repeat("a", 15) + `",
				"password":"123456"
			}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				user, err := app.FindAuthRecordByEmail("users", "test@example.com")
				if err != nil {
					t.Fatal(err)
				}

				// ensure that the user is unverified
				user.SetVerified(false)
				if err = app.Save(user); err != nil {
					t.Fatal(err)
				}

				// disable MFA
				user.Collection().MFA.Enabled = false
				if err = app.Save(user.Collection()); err != nil {
					t.Fatal(err)
				}

				otp := core.NewOTP(app)
				otp.Id = strings.Repeat("a", 15)
				otp.SetCollectionRef(user.Collection().Id)
				otp.SetRecordRef(user.Id)
				otp.SetPassword("123456")
				if err := app.Save(otp); err != nil {
					t.Fatal(err)
				}

				// test at least once that the correct request info context is properly loaded
				app.OnRecordAuthRequest().BindFunc(func(e *core.RecordAuthRequestEvent) error {
					info, err := e.RequestInfo()
					if err != nil {
						t.Fatal(err)
					}

					if info.Context != core.RequestInfoContextOTP {
						t.Fatalf("Expected request context %q, got %q", core.RequestInfoContextOTP, info.Context)
					}

					return e.Next()
				})
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"token":"`,
				`"record":{`,
				`"email":"test@example.com"`,
			},
			NotExpectedContent: []string{
				`"meta":`,
				// hidden fields
				`"tokenKey"`,
				`"password"`,
			},
			ExpectedEvents: map[string]int{
				"*":                          0,
				"OnRecordAuthWithOTPRequest": 1,
				"OnRecordAuthRequest":        1,
				"OnRecordEnrich":             1,
				// ---
				"OnModelValidate":           1,
				"OnModelCreate":             1, // authOrigin
				"OnModelCreateExecute":      1,
				"OnModelAfterCreateSuccess": 1,
				"OnModelDelete":             1, // otp delete
				"OnModelDeleteExecute":      1,
				"OnModelAfterDeleteSuccess": 1,
				// ---
				"OnRecordValidate":           1,
				"OnRecordCreate":             1,
				"OnRecordCreateExecute":      1,
				"OnRecordAfterCreateSuccess": 1,
				"OnRecordDelete":             1,
				"OnRecordDeleteExecute":      1,
				"OnRecordAfterDeleteSuccess": 1,
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				user, err := app.FindAuthRecordByEmail("users", "test@example.com")
				if err != nil {
					t.Fatal(err)
				}

				if user.Verified() {
					t.Fatal("Expected the user to remain unverified because sentTo != email")
				}
			},
		},
		{
			Name:   "valid otp with valid password and nonempty sentTo=email (disabled MFA)",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-otp",
			Body: strings.NewReader(`{
				"otpId":"` + strings.Repeat("a", 15) + `",
				"password":"123456"
			}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				user, err := app.FindAuthRecordByEmail("users", "test@example.com")
				if err != nil {
					t.Fatal(err)
				}

				// ensure that the user is unverified
				user.SetVerified(false)
				if err = app.Save(user); err != nil {
					t.Fatal(err)
				}

				// disable MFA
				user.Collection().MFA.Enabled = false
				if err = app.Save(user.Collection()); err != nil {
					t.Fatal(err)
				}

				otp := core.NewOTP(app)
				otp.Id = strings.Repeat("a", 15)
				otp.SetCollectionRef(user.Collection().Id)
				otp.SetRecordRef(user.Id)
				otp.SetPassword("123456")
				otp.SetSentTo(user.Email())
				if err := app.Save(otp); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"token":"`,
				`"record":{`,
				`"email":"test@example.com"`,
			},
			NotExpectedContent: []string{
				`"meta":`,
				// hidden fields
				`"tokenKey"`,
				`"password"`,
			},
			ExpectedEvents: map[string]int{
				"*":                          0,
				"OnRecordAuthWithOTPRequest": 1,
				"OnRecordAuthRequest":        1,
				"OnRecordEnrich":             1,
				// ---
				"OnModelValidate": 2, // +1 because of the verified user update
				// authOrigin create
				"OnModelCreate":             1,
				"OnModelCreateExecute":      1,
				"OnModelAfterCreateSuccess": 1,
				// OTP delete
				"OnModelDelete":             1,
				"OnModelDeleteExecute":      1,
				"OnModelAfterDeleteSuccess": 1,
				// user verified update
				"OnModelUpdate":             1,
				"OnModelUpdateExecute":      1,
				"OnModelAfterUpdateSuccess": 1,
				// ---
				"OnRecordValidate":           2,
				"OnRecordCreate":             1,
				"OnRecordCreateExecute":      1,
				"OnRecordAfterCreateSuccess": 1,
				"OnRecordDelete":             1,
				"OnRecordDeleteExecute":      1,
				"OnRecordAfterDeleteSuccess": 1,
				"OnRecordUpdate":             1,
				"OnRecordUpdateExecute":      1,
				"OnRecordAfterUpdateSuccess": 1,
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				user, err := app.FindAuthRecordByEmail("users", "test@example.com")
				if err != nil {
					t.Fatal(err)
				}

				if !user.Verified() {
					t.Fatal("Expected the user to be marked as verified")
				}
			},
		},

		// rate limit checks
		// -----------------------------------------------------------
		{
			Name:   "RateLimit rule - users:authWithOTP",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-otp",
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				app.Settings().RateLimits.Enabled = true
				app.Settings().RateLimits.Rules = []core.RateLimitRule{
					{MaxRequests: 100, Label: "abc"},
					{MaxRequests: 100, Label: "*:authWithOTP"},
					{MaxRequests: 100, Label: "users:auth"},
					{MaxRequests: 0, Label: "users:authWithOTP"},
				}
			},
			ExpectedStatus:  429,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "RateLimit rule - *:authWithOTP",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-otp",
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				app.Settings().RateLimits.Enabled = true
				app.Settings().RateLimits.Rules = []core.RateLimitRule{
					{MaxRequests: 100, Label: "abc"},
					{MaxRequests: 100, Label: "*:auth"},
					{MaxRequests: 0, Label: "*:authWithOTP"},
				}
			},
			ExpectedStatus:  429,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "RateLimit rule - users:auth",
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-otp",
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				app.Settings().RateLimits.Enabled = true
				app.Settings().RateLimits.Rules = []core.RateLimitRule{
					{MaxRequests: 100, Label: "abc"},
					{MaxRequests: 100, Label: "*:authWithOTP"},
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
			URL:    "/api/collections/users/auth-with-otp",
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

func TestRecordAuthWithOTPManualRateLimiterCheck(t *testing.T) {
	t.Parallel()

	var storeCache map[string]any

	otpAId := strings.Repeat("a", 15)
	otpBId := strings.Repeat("b", 15)

	scenarios := []struct {
		otpId          string
		password       string
		expectedStatus int
	}{
		{otpAId, "12345", 400},
		{otpAId, "12345", 400},
		{otpBId, "12345", 400},
		{otpBId, "12345", 400},
		{otpBId, "12345", 400},
		{otpAId, "12345", 429},
		{otpAId, "123456", 429}, // reject even if it is correct
		{otpAId, "123456", 429},
		{otpBId, "123456", 429},
	}

	for _, s := range scenarios {
		(&tests.ApiScenario{
			Method: http.MethodPost,
			URL:    "/api/collections/users/auth-with-otp",
			Body: strings.NewReader(`{
				"otpId":"` + s.otpId + `",
				"password":"` + s.password + `"
			}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				for k, v := range storeCache {
					app.Store().Set(k, v)
				}

				user, err := app.FindAuthRecordByEmail("users", "test@example.com")
				if err != nil {
					t.Fatal(err)
				}

				user.Collection().MFA.Enabled = false
				if err := app.Save(user.Collection()); err != nil {
					t.Fatal(err)
				}

				for _, id := range []string{otpAId, otpBId} {
					otp := core.NewOTP(app)
					otp.Id = id
					otp.SetCollectionRef(user.Collection().Id)
					otp.SetRecordRef(user.Id)
					otp.SetPassword("123456")
					if err := app.Save(otp); err != nil {
						t.Fatal(err)
					}
				}
			},
			ExpectedStatus:  s.expectedStatus,
			ExpectedContent: []string{`"`}, // it doesn't matter anything non-empty
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				storeCache = app.Store().GetAll()
			},
		}).Test(t)
	}
}
