package apis_test

import (
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestAdminAuthWithPassword(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "empty data",
			Method:          http.MethodPost,
			Url:             "/api/admins/auth-with-password",
			Body:            strings.NewReader(``),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{"identity":{"code":"validation_required","message":"Cannot be blank."},"password":{"code":"validation_required","message":"Cannot be blank."}}`},
		},
		{
			Name:            "invalid data",
			Method:          http.MethodPost,
			Url:             "/api/admins/auth-with-password",
			Body:            strings.NewReader(`{`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "wrong email",
			Method:          http.MethodPost,
			Url:             "/api/admins/auth-with-password",
			Body:            strings.NewReader(`{"identity":"missing@example.com","password":"1234567890"}`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"OnAdminBeforeAuthWithPasswordRequest": 1,
			},
		},
		{
			Name:            "wrong password",
			Method:          http.MethodPost,
			Url:             "/api/admins/auth-with-password",
			Body:            strings.NewReader(`{"identity":"test@example.com","password":"invalid"}`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"OnAdminBeforeAuthWithPasswordRequest": 1,
			},
		},
		{
			Name:           "valid email/password (guest)",
			Method:         http.MethodPost,
			Url:            "/api/admins/auth-with-password",
			Body:           strings.NewReader(`{"identity":"test@example.com","password":"1234567890"}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"admin":{"id":"sywbhecnh46rhm0"`,
				`"token":`,
			},
			ExpectedEvents: map[string]int{
				"OnAdminBeforeAuthWithPasswordRequest": 1,
				"OnAdminAfterAuthWithPasswordRequest":  1,
				"OnAdminAuthRequest":                   1,
			},
		},
		{
			Name:   "valid email/password (already authorized)",
			Method: http.MethodPost,
			Url:    "/api/admins/auth-with-password",
			Body:   strings.NewReader(`{"identity":"test@example.com","password":"1234567890"}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4MTYwMH0.han3_sG65zLddpcX2ic78qgy7FKecuPfOpFa8Dvi5Bg",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"admin":{"id":"sywbhecnh46rhm0"`,
				`"token":`,
			},
			ExpectedEvents: map[string]int{
				"OnAdminBeforeAuthWithPasswordRequest": 1,
				"OnAdminAfterAuthWithPasswordRequest":  1,
				"OnAdminAuthRequest":                   1,
			},
		},
		{
			Name:   "OnAdminAfterAuthWithPasswordRequest error response",
			Method: http.MethodPost,
			Url:    "/api/admins/auth-with-password",
			Body:   strings.NewReader(`{"identity":"test@example.com","password":"1234567890"}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4MTYwMH0.han3_sG65zLddpcX2ic78qgy7FKecuPfOpFa8Dvi5Bg",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				app.OnAdminAfterAuthWithPasswordRequest().Add(func(e *core.AdminAuthWithPasswordEvent) error {
					return errors.New("error")
				})
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"OnAdminBeforeAuthWithPasswordRequest": 1,
				"OnAdminAfterAuthWithPasswordRequest":  1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestAdminRequestPasswordReset(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "empty data",
			Method:          http.MethodPost,
			Url:             "/api/admins/request-password-reset",
			Body:            strings.NewReader(``),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{"email":{"code":"validation_required","message":"Cannot be blank."}}`},
		},
		{
			Name:            "invalid data",
			Method:          http.MethodPost,
			Url:             "/api/admins/request-password-reset",
			Body:            strings.NewReader(`{"email`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:           "missing admin",
			Method:         http.MethodPost,
			Url:            "/api/admins/request-password-reset",
			Body:           strings.NewReader(`{"email":"missing@example.com"}`),
			Delay:          100 * time.Millisecond,
			ExpectedStatus: 204,
		},
		{
			Name:           "existing admin",
			Method:         http.MethodPost,
			Url:            "/api/admins/request-password-reset",
			Body:           strings.NewReader(`{"email":"test@example.com"}`),
			Delay:          100 * time.Millisecond,
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnModelBeforeUpdate":                      1,
				"OnModelAfterUpdate":                       1,
				"OnMailerBeforeAdminResetPasswordSend":     1,
				"OnMailerAfterAdminResetPasswordSend":      1,
				"OnAdminBeforeRequestPasswordResetRequest": 1,
				"OnAdminAfterRequestPasswordResetRequest":  1,
			},
		},
		{
			Name:           "existing admin (after already sent)",
			Method:         http.MethodPost,
			Url:            "/api/admins/request-password-reset",
			Body:           strings.NewReader(`{"email":"test@example.com"}`),
			Delay:          100 * time.Millisecond,
			ExpectedStatus: 204,
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				// simulate recent password request
				admin, err := app.Dao().FindAdminByEmail("test@example.com")
				if err != nil {
					t.Fatal(err)
				}
				admin.LastResetSentAt = types.NowDateTime()
				dao := daos.New(app.Dao().DB()) // new dao to ignore hooks
				if err := dao.Save(admin); err != nil {
					t.Fatal(err)
				}
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestAdminConfirmPasswordReset(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "empty data",
			Method:          http.MethodPost,
			Url:             "/api/admins/confirm-password-reset",
			Body:            strings.NewReader(``),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{"password":{"code":"validation_required","message":"Cannot be blank."},"passwordConfirm":{"code":"validation_required","message":"Cannot be blank."},"token":{"code":"validation_required","message":"Cannot be blank."}}`},
		},
		{
			Name:            "invalid data",
			Method:          http.MethodPost,
			Url:             "/api/admins/confirm-password-reset",
			Body:            strings.NewReader(`{"password`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "expired token",
			Method: http.MethodPost,
			Url:    "/api/admins/confirm-password-reset",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsImV4cCI6MTY0MDk5MTY2MX0.GLwCOsgWTTEKXTK-AyGW838de1OeZGIjfHH0FoRLqZg",
				"password":"1234567890",
				"passwordConfirm":"1234567890"
			}`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{"token":{"code":"validation_invalid_token","message":"Invalid or expired token."}}}`},
		},
		{
			Name:   "valid token + invalid password",
			Method: http.MethodPost,
			Url:    "/api/admins/confirm-password-reset",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsImV4cCI6MjIwODk4MTYwMH0.kwFEler6KSMKJNstuaSDvE1QnNdCta5qSnjaIQ0hhhc",
				"password":"123456",
				"passwordConfirm":"123456"
			}`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{"password":{"code":"validation_length_out_of_range"`},
		},
		{
			Name:   "valid token + valid password",
			Method: http.MethodPost,
			Url:    "/api/admins/confirm-password-reset",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsImV4cCI6MjIwODk4MTYwMH0.kwFEler6KSMKJNstuaSDvE1QnNdCta5qSnjaIQ0hhhc",
				"password":"1234567891",
				"passwordConfirm":"1234567891"
			}`),
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnModelBeforeUpdate":                      1,
				"OnModelAfterUpdate":                       1,
				"OnAdminBeforeConfirmPasswordResetRequest": 1,
				"OnAdminAfterConfirmPasswordResetRequest":  1,
			},
		},
		{
			Name:   "OnAdminAfterConfirmPasswordResetRequest error response",
			Method: http.MethodPost,
			Url:    "/api/admins/confirm-password-reset",
			Body: strings.NewReader(`{
				"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsImV4cCI6MjIwODk4MTYwMH0.kwFEler6KSMKJNstuaSDvE1QnNdCta5qSnjaIQ0hhhc",
				"password":"1234567891",
				"passwordConfirm":"1234567891"
			}`),
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				app.OnAdminAfterConfirmPasswordResetRequest().Add(func(e *core.AdminConfirmPasswordResetEvent) error {
					return errors.New("error")
				})
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"OnModelBeforeUpdate":                      1,
				"OnModelAfterUpdate":                       1,
				"OnAdminBeforeConfirmPasswordResetRequest": 1,
				"OnAdminAfterConfirmPasswordResetRequest":  1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestAdminRefresh(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodPost,
			Url:             "/api/admins/auth-refresh",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as user",
			Method: http.MethodPost,
			Url:    "/api/admins/auth-refresh",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin (expired token)",
			Method: http.MethodPost,
			Url:    "/api/admins/auth-refresh",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTY0MDk5MTY2MX0.I7w8iktkleQvC7_UIRpD7rNzcU4OnF7i7SFIUu6lD_4",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin (valid token)",
			Method: http.MethodPost,
			Url:    "/api/admins/auth-refresh",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"admin":{"id":"sywbhecnh46rhm0"`,
				`"token":`,
			},
			ExpectedEvents: map[string]int{
				"OnAdminAuthRequest":              1,
				"OnAdminBeforeAuthRefreshRequest": 1,
				"OnAdminAfterAuthRefreshRequest":  1,
			},
		},
		{
			Name:   "OnAdminAfterAuthRefreshRequest error response",
			Method: http.MethodPost,
			Url:    "/api/admins/auth-refresh",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				app.OnAdminAfterAuthRefreshRequest().Add(func(e *core.AdminAuthRefreshEvent) error {
					return errors.New("error")
				})
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"OnAdminBeforeAuthRefreshRequest": 1,
				"OnAdminAfterAuthRefreshRequest":  1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestAdminsList(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodGet,
			Url:             "/api/admins",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as user",
			Method: http.MethodGet,
			Url:    "/api/admins",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin",
			Method: http.MethodGet,
			Url:    "/api/admins",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":1`,
				`"perPage":30`,
				`"totalItems":3`,
				`"items":[{`,
				`"id":"sywbhecnh46rhm0"`,
				`"id":"sbmbsdb40jyxf7h"`,
				`"id":"9q2trqumvlyr3bd"`,
			},
			ExpectedEvents: map[string]int{
				"OnAdminsListRequest": 1,
			},
		},
		{
			Name:   "authorized as admin + paging and sorting",
			Method: http.MethodGet,
			Url:    "/api/admins?page=2&perPage=1&sort=-created",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":2`,
				`"perPage":1`,
				`"totalItems":3`,
				`"items":[{`,
				`"id":"sbmbsdb40jyxf7h"`,
			},
			NotExpectedContent: []string{
				`"tokenKey"`,
				`"passwordHash"`,
			},
			ExpectedEvents: map[string]int{
				"OnAdminsListRequest": 1,
			},
		},
		{
			Name:   "authorized as admin + invalid filter",
			Method: http.MethodGet,
			Url:    "/api/admins?filter=invalidfield~'test2'",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + valid filter",
			Method: http.MethodGet,
			Url:    "/api/admins?filter=email~'test3'",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":1`,
				`"perPage":30`,
				`"totalItems":1`,
				`"items":[{`,
				`"id":"9q2trqumvlyr3bd"`,
			},
			NotExpectedContent: []string{
				`"tokenKey"`,
				`"passwordHash"`,
			},
			ExpectedEvents: map[string]int{
				"OnAdminsListRequest": 1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestAdminView(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodGet,
			Url:             "/api/admins/sbmbsdb40jyxf7h",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as user",
			Method: http.MethodGet,
			Url:    "/api/admins/sbmbsdb40jyxf7h",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + nonexisting admin id",
			Method: http.MethodGet,
			Url:    "/api/admins/nonexisting",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + existing admin id",
			Method: http.MethodGet,
			Url:    "/api/admins/sbmbsdb40jyxf7h",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"sbmbsdb40jyxf7h"`,
			},
			NotExpectedContent: []string{
				`"tokenKey"`,
				`"passwordHash"`,
			},
			ExpectedEvents: map[string]int{
				"OnAdminViewRequest": 1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestAdminDelete(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodDelete,
			Url:             "/api/admins/sbmbsdb40jyxf7h",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as user",
			Method: http.MethodDelete,
			Url:    "/api/admins/sbmbsdb40jyxf7h",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + missing admin id",
			Method: http.MethodDelete,
			Url:    "/api/admins/missing",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + existing admin id",
			Method: http.MethodDelete,
			Url:    "/api/admins/sbmbsdb40jyxf7h",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnModelBeforeDelete":        1,
				"OnModelAfterDelete":         1,
				"OnAdminBeforeDeleteRequest": 1,
				"OnAdminAfterDeleteRequest":  1,
			},
		},
		{
			Name:   "authorized as admin - try to delete the only remaining admin",
			Method: http.MethodDelete,
			Url:    "/api/admins/sywbhecnh46rhm0",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				// delete all admins except the authorized one
				adminModel := &models.Admin{}
				_, err := app.Dao().DB().Delete(adminModel.TableName(), dbx.Not(dbx.HashExp{
					"id": "sywbhecnh46rhm0",
				})).Execute()
				if err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"OnAdminBeforeDeleteRequest": 1,
			},
		},
		{
			Name:   "OnAdminAfterDeleteRequest error response",
			Method: http.MethodDelete,
			Url:    "/api/admins/sbmbsdb40jyxf7h",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				app.OnAdminAfterDeleteRequest().Add(func(e *core.AdminDeleteEvent) error {
					return errors.New("error")
				})
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"OnModelBeforeDelete":        1,
				"OnModelAfterDelete":         1,
				"OnAdminBeforeDeleteRequest": 1,
				"OnAdminAfterDeleteRequest":  1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestAdminCreate(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized (while having at least 1 existing admin)",
			Method:          http.MethodPost,
			Url:             "/api/admins",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "unauthorized (while having 0 existing admins)",
			Method: http.MethodPost,
			Url:    "/api/admins",
			Body:   strings.NewReader(`{"email":"testnew@example.com","password":"1234567890","passwordConfirm":"1234567890","avatar":3}`),
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				// delete all admins
				_, err := app.Dao().DB().NewQuery("DELETE FROM {{_admins}}").Execute()
				if err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":`,
				`"email":"testnew@example.com"`,
				`"avatar":3`,
			},
			ExpectedEvents: map[string]int{
				"OnModelBeforeCreate":        1,
				"OnModelAfterCreate":         1,
				"OnAdminBeforeCreateRequest": 1,
				"OnAdminAfterCreateRequest":  1,
			},
		},
		{
			Name:   "authorized as user",
			Method: http.MethodPost,
			Url:    "/api/admins",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + empty data",
			Method: http.MethodPost,
			Url:    "/api/admins",
			Body:   strings.NewReader(``),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{"email":{"code":"validation_required","message":"Cannot be blank."},"password":{"code":"validation_required","message":"Cannot be blank."}}`},
		},
		{
			Name:   "authorized as admin + invalid data format",
			Method: http.MethodPost,
			Url:    "/api/admins",
			Body:   strings.NewReader(`{`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + invalid data",
			Method: http.MethodPost,
			Url:    "/api/admins",
			Body: strings.NewReader(`{
				"email":"test@example.com",
				"password":"1234",
				"passwordConfirm":"4321",
				"avatar":99
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"avatar":{"code":"validation_max_less_equal_than_required"`,
				`"email":{"code":"validation_admin_email_exists"`,
				`"password":{"code":"validation_length_out_of_range"`,
				`"passwordConfirm":{"code":"validation_values_mismatch"`,
			},
		},
		{
			Name:   "authorized as admin + valid data",
			Method: http.MethodPost,
			Url:    "/api/admins",
			Body: strings.NewReader(`{
				"email":"testnew@example.com",
				"password":"1234567890",
				"passwordConfirm":"1234567890",
				"avatar":3
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":`,
				`"email":"testnew@example.com"`,
				`"avatar":3`,
			},
			NotExpectedContent: []string{
				`"password"`,
				`"passwordConfirm"`,
				`"tokenKey"`,
				`"passwordHash"`,
			},
			ExpectedEvents: map[string]int{
				"OnModelBeforeCreate":        1,
				"OnModelAfterCreate":         1,
				"OnAdminBeforeCreateRequest": 1,
				"OnAdminAfterCreateRequest":  1,
			},
		},
		{
			Name:   "OnAdminAfterCreateRequest error response",
			Method: http.MethodPost,
			Url:    "/api/admins",
			Body: strings.NewReader(`{
				"email":"testnew@example.com",
				"password":"1234567890",
				"passwordConfirm":"1234567890",
				"avatar":3
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				app.OnAdminAfterCreateRequest().Add(func(e *core.AdminCreateEvent) error {
					return errors.New("error")
				})
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"OnModelBeforeCreate":        1,
				"OnModelAfterCreate":         1,
				"OnAdminBeforeCreateRequest": 1,
				"OnAdminAfterCreateRequest":  1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestAdminUpdate(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodPatch,
			Url:             "/api/admins/sbmbsdb40jyxf7h",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as user",
			Method: http.MethodPatch,
			Url:    "/api/admins/sbmbsdb40jyxf7h",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + missing admin",
			Method: http.MethodPatch,
			Url:    "/api/admins/missing",
			Body:   strings.NewReader(``),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + empty data",
			Method: http.MethodPatch,
			Url:    "/api/admins/sbmbsdb40jyxf7h",
			Body:   strings.NewReader(``),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"sbmbsdb40jyxf7h"`,
				`"email":"test2@example.com"`,
				`"avatar":2`,
			},
			ExpectedEvents: map[string]int{
				"OnModelBeforeUpdate":        1,
				"OnModelAfterUpdate":         1,
				"OnAdminBeforeUpdateRequest": 1,
				"OnAdminAfterUpdateRequest":  1,
			},
		},
		{
			Name:   "authorized as admin + invalid formatted data",
			Method: http.MethodPatch,
			Url:    "/api/admins/sbmbsdb40jyxf7h",
			Body:   strings.NewReader(`{`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + invalid data",
			Method: http.MethodPatch,
			Url:    "/api/admins/sbmbsdb40jyxf7h",
			Body: strings.NewReader(`{
				"email":"test@example.com",
				"password":"1234",
				"passwordConfirm":"4321",
				"avatar":99
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"avatar":{"code":"validation_max_less_equal_than_required"`,
				`"email":{"code":"validation_admin_email_exists"`,
				`"password":{"code":"validation_length_out_of_range"`,
				`"passwordConfirm":{"code":"validation_values_mismatch"`,
			},
		},
		{
			Name:   "authorized as admin + valid data",
			Method: http.MethodPatch,
			Url:    "/api/admins/sbmbsdb40jyxf7h",
			Body: strings.NewReader(`{
				"email":"testnew@example.com",
				"password":"1234567891",
				"passwordConfirm":"1234567891",
				"avatar":5
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"sbmbsdb40jyxf7h"`,
				`"email":"testnew@example.com"`,
				`"avatar":5`,
			},
			NotExpectedContent: []string{
				`"password"`,
				`"passwordConfirm"`,
				`"tokenKey"`,
				`"passwordHash"`,
			},
			ExpectedEvents: map[string]int{
				"OnModelBeforeUpdate":        1,
				"OnModelAfterUpdate":         1,
				"OnAdminBeforeUpdateRequest": 1,
				"OnAdminAfterUpdateRequest":  1,
			},
		},
		{
			Name:   "OnAdminAfterUpdateRequest error response",
			Method: http.MethodPatch,
			Url:    "/api/admins/sbmbsdb40jyxf7h",
			Body: strings.NewReader(`{
				"email":"testnew@example.com",
				"password":"1234567891",
				"passwordConfirm":"1234567891",
				"avatar":5
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				app.OnAdminAfterUpdateRequest().Add(func(e *core.AdminUpdateEvent) error {
					return errors.New("error")
				})
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"OnModelBeforeUpdate":        1,
				"OnModelAfterUpdate":         1,
				"OnAdminBeforeUpdateRequest": 1,
				"OnAdminAfterUpdateRequest":  1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
