package apis_test

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestUsersAuthMethods(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Method:         http.MethodGet,
			Url:            "/api/users/auth-methods",
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"emailPassword":true`,
				`"authProviders":[{`,
				`"authProviders":[{`,
				`"name":"gitlab"`,
				`"state":`,
				`"codeVerifier":`,
				`"codeChallenge":`,
				`"codeChallengeMethod":`,
				`"authUrl":`,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestUserEmailAuth(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:   "authorized as user",
			Method: http.MethodPost,
			Url:    "/api/users/auth-via-email",
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin",
			Method: http.MethodPost,
			Url:    "/api/users/auth-via-email",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "invalid body format",
			Method:          http.MethodPost,
			Url:             "/api/users/auth-via-email",
			Body:            strings.NewReader(`{"email`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:           "invalid data",
			Method:         http.MethodPost,
			Url:            "/api/users/auth-via-email",
			Body:           strings.NewReader(`{"email":"","password":""}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"email":{`,
				`"password":{`,
			},
		},
		{
			Name:   "disabled email/pass auth with valid data",
			Method: http.MethodPost,
			Url:    "/api/users/auth-via-email",
			Body:   strings.NewReader(`{"email":"test@example.com","password":"123456"}`),
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				app.Settings().EmailAuth.Enabled = false
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:           "valid data",
			Method:         http.MethodPost,
			Url:            "/api/users/auth-via-email",
			Body:           strings.NewReader(`{"email":"test2@example.com","password":"123456"}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"token"`,
				`"user"`,
				`"id":"7bc84d27-6ba2-b42a-383f-4197cc3d3d0c"`,
				`"email":"test2@example.com"`,
				`"verified":false`, // unverified user should be able to authenticate
			},
			ExpectedEvents: map[string]int{"OnUserAuthRequest": 1},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestUserRequestPasswordReset(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:            "empty data",
			Method:          http.MethodPost,
			Url:             "/api/users/request-password-reset",
			Body:            strings.NewReader(``),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{"email":{"code":"validation_required","message":"Cannot be blank."}}`},
		},
		{
			Name:            "invalid data",
			Method:          http.MethodPost,
			Url:             "/api/users/request-password-reset",
			Body:            strings.NewReader(`{"email`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:           "missing user",
			Method:         http.MethodPost,
			Url:            "/api/users/request-password-reset",
			Body:           strings.NewReader(`{"email":"missing@example.com"}`),
			Delay:          100 * time.Millisecond,
			ExpectedStatus: 204,
		},
		{
			Name:           "existing user",
			Method:         http.MethodPost,
			Url:            "/api/users/request-password-reset",
			Body:           strings.NewReader(`{"email":"test@example.com"}`),
			Delay:          100 * time.Millisecond,
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnModelBeforeUpdate":                 1,
				"OnModelAfterUpdate":                  1,
				"OnMailerBeforeUserResetPasswordSend": 1,
				"OnMailerAfterUserResetPasswordSend":  1,
			},
		},
		{
			Name:           "existing user (after already sent)",
			Method:         http.MethodPost,
			Url:            "/api/users/request-password-reset",
			Body:           strings.NewReader(`{"email":"test@example.com"}`),
			Delay:          100 * time.Millisecond,
			ExpectedStatus: 204,
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				// simulate recent password request
				user, err := app.Dao().FindUserByEmail("test@example.com")
				if err != nil {
					t.Fatal(err)
				}
				user.LastResetSentAt = types.NowDateTime()
				dao := daos.New(app.Dao().DB()) // new dao to ignore hooks
				if err := dao.Save(user); err != nil {
					t.Fatal(err)
				}
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestUserConfirmPasswordReset(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:            "empty data",
			Method:          http.MethodPost,
			Url:             "/api/users/confirm-password-reset",
			Body:            strings.NewReader(``),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{"password":{"code":"validation_required","message":"Cannot be blank."},"passwordConfirm":{"code":"validation_required","message":"Cannot be blank."},"token":{"code":"validation_required","message":"Cannot be blank."}}`},
		},
		{
			Name:            "invalid data format",
			Method:          http.MethodPost,
			Url:             "/api/users/confirm-password-reset",
			Body:            strings.NewReader(`{"password`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:           "expired token",
			Method:         http.MethodPost,
			Url:            "/api/users/confirm-password-reset",
			Body:           strings.NewReader(`{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0eXBlIjoidXNlciIsImlkIjoiNGQwMTk3Y2MtMmI0YS0zZjgzLWEyNmItZDc3YmM4NDIzZDNjIiwiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwiZXhwIjoxNjQxMDMxMjAwfQ.t2lVe0ny9XruQsSFQdXqBi0I85i6vIUAQjFXZY5HPxc","password":"123456789","passwordConfirm":"123456789"}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"token":{`,
				`"code":"validation_invalid_token"`,
			},
		},
		{
			Name:           "valid token and data",
			Method:         http.MethodPost,
			Url:            "/api/users/confirm-password-reset",
			Body:           strings.NewReader(`{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0eXBlIjoidXNlciIsImlkIjoiNGQwMTk3Y2MtMmI0YS0zZjgzLWEyNmItZDc3YmM4NDIzZDNjIiwiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwiZXhwIjoxODYxOTU2MDAwfQ.V1gEbY4caEIF6IhQAJ8KZD4RvOGvTCFuYg1fTRSvhe0","password":"123456789","passwordConfirm":"123456789"}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"token":`,
				`"user":`,
				`"id":"4d0197cc-2b4a-3f83-a26b-d77bc8423d3c"`,
				`"email":"test@example.com"`,
			},
			ExpectedEvents: map[string]int{"OnUserAuthRequest": 1, "OnModelAfterUpdate": 1, "OnModelBeforeUpdate": 1},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestUserRequestVerification(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:            "empty data",
			Method:          http.MethodPost,
			Url:             "/api/users/request-verification",
			Body:            strings.NewReader(``),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{"email":{"code":"validation_required","message":"Cannot be blank."}}`},
		},
		{
			Name:            "invalid data",
			Method:          http.MethodPost,
			Url:             "/api/users/request-verification",
			Body:            strings.NewReader(`{"email`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:           "missing user",
			Method:         http.MethodPost,
			Url:            "/api/users/request-verification",
			Body:           strings.NewReader(`{"email":"missing@example.com"}`),
			Delay:          100 * time.Millisecond,
			ExpectedStatus: 204,
		},
		{
			Name:           "existing already verified user",
			Method:         http.MethodPost,
			Url:            "/api/users/request-verification",
			Body:           strings.NewReader(`{"email":"test@example.com"}`),
			Delay:          100 * time.Millisecond,
			ExpectedStatus: 204,
		},
		{
			Name:           "existing unverified user",
			Method:         http.MethodPost,
			Url:            "/api/users/request-verification",
			Body:           strings.NewReader(`{"email":"test2@example.com"}`),
			Delay:          100 * time.Millisecond,
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnModelBeforeUpdate":                1,
				"OnModelAfterUpdate":                 1,
				"OnMailerBeforeUserVerificationSend": 1,
				"OnMailerAfterUserVerificationSend":  1,
			},
		},
		{
			Name:           "existing unverified user (after already sent)",
			Method:         http.MethodPost,
			Url:            "/api/users/request-verification",
			Body:           strings.NewReader(`{"email":"test2@example.com"}`),
			Delay:          100 * time.Millisecond,
			ExpectedStatus: 204,
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				// simulate recent verification sent
				user, err := app.Dao().FindUserByEmail("test2@example.com")
				if err != nil {
					t.Fatal(err)
				}
				user.LastVerificationSentAt = types.NowDateTime()
				dao := daos.New(app.Dao().DB()) // new dao to ignore hooks
				if err := dao.Save(user); err != nil {
					t.Fatal(err)
				}
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestUserConfirmVerification(t *testing.T) {
	scenarios := []tests.ApiScenario{
		// empty data
		{
			Method:         http.MethodPost,
			Url:            "/api/users/confirm-verification",
			Body:           strings.NewReader(``),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":`,
				`"token":{"code":"validation_required"`,
			},
		},
		// invalid data
		{
			Method:          http.MethodPost,
			Url:             "/api/users/confirm-verification",
			Body:            strings.NewReader(`{"token`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		// expired token
		{
			Method:         http.MethodPost,
			Url:            "/api/users/confirm-verification",
			Body:           strings.NewReader(`{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0eXBlIjoidXNlciIsImlkIjoiN2JjODRkMjctNmJhMi1iNDJhLTM4M2YtNDE5N2NjM2QzZDBjIiwiZW1haWwiOiJ0ZXN0MkBleGFtcGxlLmNvbSIsImV4cCI6MTY0MTAzMTIwMH0.YCqyREksfqn7cWu-innNNTbWQCr9DgYr7dduM2wxrtQ"}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"token":{`,
				`"code":"validation_invalid_token"`,
			},
		},
		// valid token
		{
			Method:         http.MethodPost,
			Url:            "/api/users/confirm-verification",
			Body:           strings.NewReader(`{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0eXBlIjoidXNlciIsImlkIjoiN2JjODRkMjctNmJhMi1iNDJhLTM4M2YtNDE5N2NjM2QzZDBjIiwiZW1haWwiOiJ0ZXN0MkBleGFtcGxlLmNvbSIsImV4cCI6MTg2MTk1NjAwMH0.OsxRKuZrNTnwyVjvCwB4jY8TbT-NPZ-UFCpRhCvuv2U"}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"token":`,
				`"user":`,
				`"id":"7bc84d27-6ba2-b42a-383f-4197cc3d3d0c"`,
				`"email":"test2@example.com"`,
				`"verified":true`,
			},
			ExpectedEvents: map[string]int{
				"OnUserAuthRequest":   1,
				"OnModelAfterUpdate":  1,
				"OnModelBeforeUpdate": 1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestUserRequestEmailChange(t *testing.T) {
	scenarios := []tests.ApiScenario{
		// unauthorized
		{
			Method:          http.MethodPost,
			Url:             "/api/users/request-email-change",
			Body:            strings.NewReader(`{"newEmail":"change@example.com"}`),
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		// authorized as admin
		{
			Method: http.MethodPost,
			Url:    "/api/users/request-email-change",
			Body:   strings.NewReader(`{"newEmail":"change@example.com"}`),
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		// invalid data
		{
			Method: http.MethodPost,
			Url:    "/api/users/request-email-change",
			Body:   strings.NewReader(`{"newEmail`),
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		// empty data
		{
			Method: http.MethodPost,
			Url:    "/api/users/request-email-change",
			Body:   strings.NewReader(``),
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":`,
				`"newEmail":{"code":"validation_required"`,
			},
		},
		// valid data (existing email)
		{
			Method: http.MethodPost,
			Url:    "/api/users/request-email-change",
			Body:   strings.NewReader(`{"newEmail":"test2@example.com"}`),
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":`,
				`"newEmail":{"code":"validation_user_email_exists"`,
			},
		},
		// valid data (new email)
		{
			Method: http.MethodPost,
			Url:    "/api/users/request-email-change",
			Body:   strings.NewReader(`{"newEmail":"change@example.com"}`),
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnMailerBeforeUserChangeEmailSend": 1,
				"OnMailerAfterUserChangeEmailSend":  1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestUserConfirmEmailChange(t *testing.T) {
	scenarios := []tests.ApiScenario{
		// empty data
		{
			Method:         http.MethodPost,
			Url:            "/api/users/confirm-email-change",
			Body:           strings.NewReader(``),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":`,
				`"token":{"code":"validation_required"`,
				`"password":{"code":"validation_required"`,
			},
		},
		// invalid data
		{
			Method:          http.MethodPost,
			Url:             "/api/users/confirm-email-change",
			Body:            strings.NewReader(`{"token`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		// expired token and correct password
		{
			Method:         http.MethodPost,
			Url:            "/api/users/confirm-email-change",
			Body:           strings.NewReader(`{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjdiYzg0ZDI3LTZiYTItYjQyYS0zODNmLTQxOTdjYzNkM2QwYyIsInR5cGUiOiJ1c2VyIiwiZW1haWwiOiJ0ZXN0MkBleGFtcGxlLmNvbSIsIm5ld0VtYWlsIjoiY2hhbmdlQGV4YW1wbGUuY29tIiwiZXhwIjoxNjQwOTkxNjAwfQ.DOqNtSDcXbWix8OsK13X-tjfWi6jZNlAzIZiwG_YDOs","password":"123456"}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"token":{`,
				`"code":"validation_invalid_token"`,
			},
		},
		// valid token and incorrect password
		{
			Method:         http.MethodPost,
			Url:            "/api/users/confirm-email-change",
			Body:           strings.NewReader(`{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjdiYzg0ZDI3LTZiYTItYjQyYS0zODNmLTQxOTdjYzNkM2QwYyIsInR5cGUiOiJ1c2VyIiwiZW1haWwiOiJ0ZXN0MkBleGFtcGxlLmNvbSIsIm5ld0VtYWlsIjoiY2hhbmdlQGV4YW1wbGUuY29tIiwiZXhwIjoxODkzNDUyNDAwfQ.aWMQJ_c49yFbzHO5TNhlkbKRokQ_isc2RbLGuSJx44c","password":"654321"}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"password":{`,
				`"code":"validation_invalid_password"`,
			},
		},
		// valid token and correct password
		{
			Method:         http.MethodPost,
			Url:            "/api/users/confirm-email-change",
			Body:           strings.NewReader(`{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjdiYzg0ZDI3LTZiYTItYjQyYS0zODNmLTQxOTdjYzNkM2QwYyIsInR5cGUiOiJ1c2VyIiwiZW1haWwiOiJ0ZXN0MkBleGFtcGxlLmNvbSIsIm5ld0VtYWlsIjoiY2hhbmdlQGV4YW1wbGUuY29tIiwiZXhwIjoxODkzNDUyNDAwfQ.aWMQJ_c49yFbzHO5TNhlkbKRokQ_isc2RbLGuSJx44c","password":"123456"}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"token":`,
				`"user":`,
				`"id":"7bc84d27-6ba2-b42a-383f-4197cc3d3d0c"`,
				`"email":"change@example.com"`,
				`"verified":true`,
			},
			ExpectedEvents: map[string]int{"OnUserAuthRequest": 1, "OnModelAfterUpdate": 1, "OnModelBeforeUpdate": 1},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestUserRefresh(t *testing.T) {
	scenarios := []tests.ApiScenario{
		// unauthorized
		{
			Method:          http.MethodPost,
			Url:             "/api/users/refresh",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		// authorized as admin
		{
			Method: http.MethodPost,
			Url:    "/api/users/refresh",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		// authorized as user
		{
			Method: http.MethodPost,
			Url:    "/api/users/refresh",
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"token":`,
				`"user":`,
				`"id":"4d0197cc-2b4a-3f83-a26b-d77bc8423d3c"`,
			},
			ExpectedEvents: map[string]int{"OnUserAuthRequest": 1},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestUsersList(t *testing.T) {
	scenarios := []tests.ApiScenario{
		// unauthorized
		{
			Method:          http.MethodGet,
			Url:             "/api/users",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		// authorized as user
		{
			Method: http.MethodGet,
			Url:    "/api/users",
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		// authorized as admin
		{
			Method: http.MethodGet,
			Url:    "/api/users",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":1`,
				`"perPage":30`,
				`"totalItems":3`,
				`"items":[{`,
				`"id":"4d0197cc-2b4a-3f83-a26b-d77bc8423d3c"`,
				`"id":"7bc84d27-6ba2-b42a-383f-4197cc3d3d0c"`,
				`"id":"97cc3d3d-6ba2-383f-b42a-7bc84d27410c"`,
			},
			ExpectedEvents: map[string]int{"OnUsersListRequest": 1},
		},
		// authorized as admin + paging and sorting
		{
			Method: http.MethodGet,
			Url:    "/api/users?page=2&perPage=2&sort=-created",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":2`,
				`"perPage":2`,
				`"totalItems":3`,
				`"items":[{`,
				`"id":"4d0197cc-2b4a-3f83-a26b-d77bc8423d3c"`,
			},
			ExpectedEvents: map[string]int{"OnUsersListRequest": 1},
		},
		// authorized as admin + invalid filter
		{
			Method: http.MethodGet,
			Url:    "/api/users?filter=invalidfield~'test2'",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		// authorized as admin + valid filter
		{
			Method: http.MethodGet,
			Url:    "/api/users?filter=verified=true",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":1`,
				`"perPage":30`,
				`"totalItems":2`,
				`"items":[{`,
				`"id":"4d0197cc-2b4a-3f83-a26b-d77bc8423d3c"`,
				`"id":"97cc3d3d-6ba2-383f-b42a-7bc84d27410c"`,
			},
			ExpectedEvents: map[string]int{"OnUsersListRequest": 1},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestUserView(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodGet,
			Url:             "/api/users/4d0197cc-2b4a-3f83-a26b-d77bc8423d3c",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + nonexisting user id",
			Method: http.MethodGet,
			Url:    "/api/users/00000000-0000-0000-0000-d77bc8423d3c",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + existing user id",
			Method: http.MethodGet,
			Url:    "/api/users/4d0197cc-2b4a-3f83-a26b-d77bc8423d3c",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"4d0197cc-2b4a-3f83-a26b-d77bc8423d3c"`,
			},
			ExpectedEvents: map[string]int{"OnUserViewRequest": 1},
		},
		{
			Name:   "authorized as user - trying to view another user",
			Method: http.MethodGet,
			Url:    "/api/users/7bc84d27-6ba2-b42a-383f-4197cc3d3d0c",
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as user - owner",
			Method: http.MethodGet,
			Url:    "/api/users/4d0197cc-2b4a-3f83-a26b-d77bc8423d3c",
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"4d0197cc-2b4a-3f83-a26b-d77bc8423d3c"`,
			},
			ExpectedEvents: map[string]int{"OnUserViewRequest": 1},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestUserDelete(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodDelete,
			Url:             "/api/users/4d0197cc-2b4a-3f83-a26b-d77bc8423d3c",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + nonexisting user id",
			Method: http.MethodDelete,
			Url:    "/api/users/00000000-0000-0000-0000-d77bc8423d3c",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + existing user id",
			Method: http.MethodDelete,
			Url:    "/api/users/4d0197cc-2b4a-3f83-a26b-d77bc8423d3c",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnUserBeforeDeleteRequest": 1,
				"OnUserAfterDeleteRequest":  1,
				"OnModelBeforeDelete":       2, // cascade delete to related Record model
				"OnModelAfterDelete":        2, // cascade delete to related Record model
			},
		},
		{
			Name:   "authorized as user - trying to delete another user",
			Method: http.MethodDelete,
			Url:    "/api/users/7bc84d27-6ba2-b42a-383f-4197cc3d3d0c",
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as user - owner",
			Method: http.MethodDelete,
			Url:    "/api/users/4d0197cc-2b4a-3f83-a26b-d77bc8423d3c",
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnUserBeforeDeleteRequest": 1,
				"OnUserAfterDeleteRequest":  1,
				"OnModelBeforeDelete":       2, // cascade delete to related Record model
				"OnModelAfterDelete":        2, // cascade delete to related Record model
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestUserCreate(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:           "empty data",
			Method:         http.MethodPost,
			Url:            "/api/users",
			Body:           strings.NewReader(``),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"email":{"code":"validation_required"`,
				`"password":{"code":"validation_required"`,
			},
		},
		{
			Name:           "invalid data",
			Method:         http.MethodPost,
			Url:            "/api/users",
			Body:           strings.NewReader(`{"email":"test@example.com","password":"1234","passwordConfirm":"4321"}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"email":{"code":"validation_user_email_exists"`,
				`"password":{"code":"validation_length_out_of_range"`,
				`"passwordConfirm":{"code":"validation_values_mismatch"`,
			},
		},
		{
			Name:   "valid data but with disabled email/pass auth",
			Method: http.MethodPost,
			Url:    "/api/users",
			Body:   strings.NewReader(`{"email":"newuser@example.com","password":"123456789","passwordConfirm":"123456789"}`),
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				app.Settings().EmailAuth.Enabled = false
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:           "valid data",
			Method:         http.MethodPost,
			Url:            "/api/users",
			Body:           strings.NewReader(`{"email":"newuser@example.com","password":"123456789","passwordConfirm":"123456789"}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":`,
				`"email":"newuser@example.com"`,
			},
			ExpectedEvents: map[string]int{
				"OnUserBeforeCreateRequest": 1,
				"OnUserAfterCreateRequest":  1,
				"OnModelBeforeCreate":       2, // +1 for the created profile record
				"OnModelAfterCreate":        2, // +1 for the created profile record
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestUserUpdate(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodPatch,
			Url:             "/api/users/4d0197cc-2b4a-3f83-a26b-d77bc8423d3c",
			Body:            strings.NewReader(`{"email":"new@example.com"}`),
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as user (owner)",
			Method: http.MethodPatch,
			Url:    "/api/users/4d0197cc-2b4a-3f83-a26b-d77bc8423d3c",
			Body:   strings.NewReader(`{"email":"new@example.com"}`),
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin - invalid/missing user id",
			Method: http.MethodPatch,
			Url:    "/api/users/invalid",
			Body:   strings.NewReader(``),
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin - empty data",
			Method: http.MethodPatch,
			Url:    "/api/users/4d0197cc-2b4a-3f83-a26b-d77bc8423d3c",
			Body:   strings.NewReader(``),
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"4d0197cc-2b4a-3f83-a26b-d77bc8423d3c"`,
				`"email":"test@example.com"`,
			},
			ExpectedEvents: map[string]int{
				"OnUserBeforeUpdateRequest": 1,
				"OnUserAfterUpdateRequest":  1,
				"OnModelBeforeUpdate":       1,
				"OnModelAfterUpdate":        1,
			},
		},
		{
			Name:   "authorized as admin - invalid data",
			Method: http.MethodPatch,
			Url:    "/api/users/4d0197cc-2b4a-3f83-a26b-d77bc8423d3c",
			Body:   strings.NewReader(`{"email":"test2@example.com"}`),
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"email":{"code":"validation_user_email_exists"`,
			},
		},
		{
			Name:   "authorized as admin - valid data",
			Method: http.MethodPatch,
			Url:    "/api/users/4d0197cc-2b4a-3f83-a26b-d77bc8423d3c",
			Body:   strings.NewReader(`{"email":"new@example.com"}`),
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"4d0197cc-2b4a-3f83-a26b-d77bc8423d3c"`,
				`"email":"new@example.com"`,
			},
			ExpectedEvents: map[string]int{
				"OnUserBeforeUpdateRequest": 1,
				"OnUserAfterUpdateRequest":  1,
				"OnModelBeforeUpdate":       1,
				"OnModelAfterUpdate":        1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
