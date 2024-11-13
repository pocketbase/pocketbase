package apis_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestRecordRequestOTP(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "not an auth collection",
			Method:          http.MethodPost,
			URL:             "/api/collections/demo1/request-otp",
			Body:            strings.NewReader(`{"email":"test@example.com"}`),
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "auth collection with disabled otp",
			Method: http.MethodPost,
			URL:    "/api/collections/users/request-otp",
			Body:   strings.NewReader(`{"email":"test@example.com"}`),
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
			Name:            "empty body",
			Method:          http.MethodPost,
			URL:             "/api/collections/users/request-otp",
			Body:            strings.NewReader(``),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{"email":{"code":"validation_required","message":"Cannot be blank."}}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:            "invalid body",
			Method:          http.MethodPost,
			URL:             "/api/collections/users/request-otp",
			Body:            strings.NewReader(`{"email`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:           "invalid request data",
			Method:         http.MethodPost,
			URL:            "/api/collections/users/request-otp",
			Body:           strings.NewReader(`{"email":"invalid"}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"email":{"code":"validation_is_email`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
		{
			Name:           "missing auth record",
			Method:         http.MethodPost,
			URL:            "/api/collections/users/request-otp",
			Body:           strings.NewReader(`{"email":"missing@example.com"}`),
			Delay:          100 * time.Millisecond,
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"otpId":"`, // some fake random generated string
			},
			ExpectedEvents: map[string]int{
				"*":                         0,
				"OnRecordRequestOTPRequest": 1,
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				if app.TestMailer.TotalSend() != 0 {
					t.Fatalf("Expected zero emails, got %d", app.TestMailer.TotalSend())
				}
			},
		},
		{
			Name:   "existing auth record (with < 9 non-expired)",
			Method: http.MethodPost,
			URL:    "/api/collections/users/request-otp",
			Body:   strings.NewReader(`{"email":"test@example.com"}`),
			Delay:  100 * time.Millisecond,
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				user, err := app.FindAuthRecordByEmail("users", "test@example.com")
				if err != nil {
					t.Fatal(err)
				}

				// insert 8 non-expired and 2 expired
				for i := 0; i < 10; i++ {
					otp := core.NewOTP(app)
					otp.Id = "otp_" + strconv.Itoa(i)
					otp.SetCollectionRef(user.Collection().Id)
					otp.SetRecordRef(user.Id)
					otp.SetPassword("123456")
					if i >= 8 {
						expiredDate := types.NowDateTime().AddDate(-3, 0, 0)
						otp.SetRaw("created", expiredDate)
						otp.SetRaw("updated", expiredDate)
					}
					if err := app.SaveNoValidate(otp); err != nil {
						t.Fatal(err)
					}
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"otpId":"`,
			},
			NotExpectedContent: []string{
				`"otpId":"otp_`,
			},
			ExpectedEvents: map[string]int{
				"*":                          0,
				"OnRecordRequestOTPRequest":  1,
				"OnMailerSend":               1,
				"OnMailerRecordOTPSend":      1,
				"OnModelCreate":              1,
				"OnModelCreateExecute":       1,
				"OnModelAfterCreateSuccess":  1,
				"OnModelValidate":            2, // + 1 for the OTP update after the email send
				"OnRecordCreate":             1,
				"OnRecordCreateExecute":      1,
				"OnRecordAfterCreateSuccess": 1,
				"OnRecordValidate":           2,
				// OTP update
				"OnModelUpdate":              1,
				"OnModelUpdateExecute":       1,
				"OnModelAfterUpdateSuccess":  1,
				"OnRecordUpdate":             1,
				"OnRecordUpdateExecute":      1,
				"OnRecordAfterUpdateSuccess": 1,
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				if app.TestMailer.TotalSend() != 1 {
					t.Fatalf("Expected 1 email, got %d", app.TestMailer.TotalSend())
				}

				// ensure that sentTo is set
				otps, err := app.FindRecordsByFilter(core.CollectionNameOTPs, "sentTo='test@example.com'", "", 0, 0)
				if err != nil || len(otps) != 1 {
					t.Fatalf("Expected to find 1 OTP with sentTo %q, found %d", "test@example.com", len(otps))
				}
			},
		},
		{
			Name:   "existing auth record with intercepted email (with < 9 non-expired)",
			Method: http.MethodPost,
			URL:    "/api/collections/users/request-otp",
			Body:   strings.NewReader(`{"email":"test@example.com"}`),
			Delay:  100 * time.Millisecond,
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				// prevent email sent
				app.OnMailerRecordOTPSend("users").BindFunc(func(e *core.MailerRecordEvent) error {
					return nil
				})
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"otpId":"`,
			},
			NotExpectedContent: []string{
				`"otpId":"otp_`,
			},
			ExpectedEvents: map[string]int{
				"*":                          0,
				"OnRecordRequestOTPRequest":  1,
				"OnMailerRecordOTPSend":      1,
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
				if app.TestMailer.TotalSend() != 0 {
					t.Fatalf("Expected 0 emails, got %d", app.TestMailer.TotalSend())
				}

				// ensure that there is no OTP with user email as sentTo
				otps, err := app.FindRecordsByFilter(core.CollectionNameOTPs, "sentTo='test@example.com'", "", 0, 0)
				if err != nil || len(otps) != 0 {
					t.Fatalf("Expected to find 0 OTPs with sentTo %q, found %d", "test@example.com", len(otps))
				}
			},
		},
		{
			Name:   "existing auth record (with > 9 non-expired)",
			Method: http.MethodPost,
			URL:    "/api/collections/users/request-otp",
			Body:   strings.NewReader(`{"email":"test@example.com"}`),
			Delay:  100 * time.Millisecond,
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				user, err := app.FindAuthRecordByEmail("users", "test@example.com")
				if err != nil {
					t.Fatal(err)
				}

				// insert 10 non-expired
				for i := 0; i < 10; i++ {
					otp := core.NewOTP(app)
					otp.Id = "otp_" + strconv.Itoa(i)
					otp.SetCollectionRef(user.Collection().Id)
					otp.SetRecordRef(user.Id)
					otp.SetPassword("123456")
					if err := app.SaveNoValidate(otp); err != nil {
						t.Fatal(err)
					}
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"otpId":"otp_9"`,
			},
			ExpectedEvents: map[string]int{
				"*":                         0,
				"OnRecordRequestOTPRequest": 1,
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				if app.TestMailer.TotalSend() != 0 {
					t.Fatalf("Expected 0 sent emails, got %d", app.TestMailer.TotalSend())
				}
			},
		},

		// rate limit checks
		// -----------------------------------------------------------
		{
			Name:   "RateLimit rule - users:requestOTP",
			Method: http.MethodPost,
			URL:    "/api/collections/users/request-otp",
			Body:   strings.NewReader(`{"email":"test@example.com"}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				app.Settings().RateLimits.Enabled = true
				app.Settings().RateLimits.Rules = []core.RateLimitRule{
					{MaxRequests: 100, Label: "abc"},
					{MaxRequests: 100, Label: "*:requestOTP"},
					{MaxRequests: 0, Label: "users:requestOTP"},
				}
			},
			ExpectedStatus:  429,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "RateLimit rule - *:requestOTP",
			Method: http.MethodPost,
			URL:    "/api/collections/users/request-otp",
			Body:   strings.NewReader(`{"email":"test@example.com"}`),
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				app.Settings().RateLimits.Enabled = true
				app.Settings().RateLimits.Rules = []core.RateLimitRule{
					{MaxRequests: 100, Label: "abc"},
					{MaxRequests: 0, Label: "*:requestOTP"},
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
