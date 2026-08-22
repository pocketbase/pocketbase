package apis_test

import (
	"fmt"
	"net/http/httptest"
	"testing"
	"testing/synctest"
	"time"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/hook"
)

func TestDefaultRateLimitMiddleware(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	app.Settings().RateLimits.Enabled = true
	app.Settings().RateLimits.Rules = []core.RateLimitRule{
		{
			Label:       "/rate/",
			MaxRequests: 2,
			Duration:    1,
		},
		{
			Label:       "/rate/b",
			MaxRequests: 3,
			Duration:    1,
		},
		{
			Label:       "POST /rate/b",
			MaxRequests: 1,
			Duration:    1,
		},
		{
			Label:       "/rate/guest",
			MaxRequests: 1,
			Duration:    1,
			Audience:    core.RateLimitRuleAudienceGuest,
		},
		{
			Label:       "/rate/auth",
			MaxRequests: 1,
			Duration:    1,
			Audience:    core.RateLimitRuleAudienceAuth,
		},
	}

	pbRouter, err := apis.NewRouter(app)
	if err != nil {
		t.Fatal(err)
	}
	pbRouter.GET("/norate", func(e *core.RequestEvent) error {
		return e.String(200, "norate")
	}).BindFunc(func(e *core.RequestEvent) error {
		return e.Next()
	})
	pbRouter.GET("/rate/a", func(e *core.RequestEvent) error {
		return e.String(200, "a")
	})
	pbRouter.GET("/rate/b", func(e *core.RequestEvent) error {
		return e.String(200, "b")
	})
	pbRouter.GET("/rate/guest", func(e *core.RequestEvent) error {
		return e.String(200, "guest")
	})
	pbRouter.GET("/rate/auth", func(e *core.RequestEvent) error {
		return e.String(200, "auth")
	})

	mux, err := pbRouter.BuildMux()
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		url            string
		wait           float64 // ms
		authenticated  bool
		expectedStatus int
	}{
		{"/norate", 0, false, 200},
		{"/norate", 0, false, 200},
		{"/norate", 0, false, 200},
		{"/norate", 0, false, 200},
		{"/norate", 0, false, 200},

		{"/rate/a", 0, false, 200},
		{"/rate/a", 900, false, 200}, // (fixed window check) wait enough to ensure that it can't fit more than 2 requests in 1s
		{"/rate/a", 900, false, 200},
		{"/rate/a", 0, false, 200},
		{"/rate/a", 0, false, 429},
		{"/rate/a", 0, false, 429},
		{"/rate/a", 1000, false, 200},
		{"/rate/a", 0, false, 200},
		{"/rate/a", 0, false, 429},

		{"/rate/b", 0, false, 200},
		{"/rate/b", 0, false, 200},
		{"/rate/b", 0, false, 200},
		{"/rate/b", 0, false, 429},
		{"/rate/b", 1000, false, 200},
		{"/rate/b", 0, false, 200},
		{"/rate/b", 0, false, 200},
		{"/rate/b", 0, false, 429},

		// "auth" with guest (should fallback to the /rate/ rule)
		{"/rate/auth", 0, false, 200},
		{"/rate/auth", 0, false, 200},
		{"/rate/auth", 0, false, 429},
		{"/rate/auth", 0, false, 429},

		// "auth" rule with regular user (should match the /rate/auth rule)
		{"/rate/auth", 0, true, 200},
		{"/rate/auth", 0, true, 429},
		{"/rate/auth", 0, true, 429},

		// "guest" with guest (should match the /rate/guest rule)
		{"/rate/guest", 0, false, 200},
		{"/rate/guest", 0, false, 429},
		{"/rate/guest", 0, false, 429},

		// "guest" rule with regular user (should fallback to the /rate/ rule)
		{"/rate/guest", 1000, true, 200},
		{"/rate/guest", 0, true, 200},
		{"/rate/guest", 0, true, 429},
		{"/rate/guest", 0, true, 429},
	}

	synctest.Test(t, func(t *testing.T) {
		for i, s := range scenarios {
			prefix := fmt.Sprintf("[%s:%d] ", s.url, i+1)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", s.url, nil)

			if s.authenticated {
				auth, err := app.FindAuthRecordByEmail("users", "test@example.com")
				if err != nil {
					t.Fatalf(prefix+"%v", err)
				}

				token, err := auth.NewAuthToken()
				if err != nil {
					t.Fatalf(prefix+"%v", err)
				}

				req.Header.Add("Authorization", token)
			}

			if s.wait > 0 {
				synctest.Sleep(time.Duration(s.wait) * time.Millisecond)
			}

			mux.ServeHTTP(rec, req)

			result := rec.Result()

			if result.StatusCode != s.expectedStatus {
				t.Fatalf(prefix+"Expected response status %d, got %d", s.expectedStatus, result.StatusCode)
			}
		}
	})
}

func TestDefaultRateLimitMiddlewareSkipChecks(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	app.Settings().RateLimits.Enabled = true
	app.Settings().RateLimits.Rules = []core.RateLimitRule{
		{
			Label:       "/rate",
			MaxRequests: 1,
			Duration:    5,
		},
	}

	pbRouter, err := apis.NewRouter(app)
	if err != nil {
		t.Fatal(err)
	}

	// just for the exclude tests - load the user IP from a query param
	pbRouter.Bind(&hook.Handler[*core.RequestEvent]{
		Priority: apis.DefaultRateLimitMiddlewarePriority - 1,
		Func: func(e *core.RequestEvent) error {
			testIp := e.Request.URL.Query().Get("testIP")
			if testIp != "" {
				e.Request.Header.Set("x-test-ip", testIp)
			}

			return e.Next()
		},
	})

	pbRouter.GET("/rate", func(e *core.RequestEvent) error {
		return e.String(200, "test")
	})

	mux, err := pbRouter.BuildMux()
	if err != nil {
		t.Fatal(err)
	}

	checkStatusCodes := func(t *testing.T, got []int, expected []int) {
		if len(expected) != len(got) {
			t.Fatalf("Expected status codes %v, got %v", expected, got)
		}

		for i, item := range expected {
			if got[i] != item {
				t.Fatalf("Expected %d status code to be %d, got %d:\n%v", i, item, got[i], got)
			}
		}
	}

	t.Run("base check", func(t *testing.T) {
		app.Settings().RateLimits.Enabled = true

		statusCodes := []int{}
		for range 3 {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/rate", nil)

			mux.ServeHTTP(rec, req)

			result := rec.Result()

			statusCodes = append(statusCodes, result.StatusCode)
		}

		checkStatusCodes(t, statusCodes, []int{200, 429, 429})
	})

	t.Run("disabled rate limiter", func(t *testing.T) {
		app.Settings().RateLimits.Enabled = false

		statusCodes := []int{}
		for range 3 {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/rate", nil)

			mux.ServeHTTP(rec, req)

			result := rec.Result()

			statusCodes = append(statusCodes, result.StatusCode)
		}

		checkStatusCodes(t, statusCodes, []int{200, 200, 200})
	})

	t.Run("authenticated as superuser", func(t *testing.T) {
		app.Settings().RateLimits.Enabled = true

		superuser, err := app.FindAuthRecordByEmail(core.CollectionNameSuperusers, "test@example.com")
		if err != nil {
			t.Fatal(err)
		}

		token, err := superuser.NewAuthToken()
		if err != nil {
			t.Fatal(err)
		}

		statusCodes := []int{}
		for range 3 {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/rate", nil)
			req.Header.Add("Authorization", token)

			mux.ServeHTTP(rec, req)

			result := rec.Result()

			statusCodes = append(statusCodes, result.StatusCode)
		}

		checkStatusCodes(t, statusCodes, []int{200, 200, 200})
	})

	t.Run("excludedIPs (different)", func(t *testing.T) {
		app.Settings().RateLimits.Enabled = true
		app.Settings().RateLimits.ExcludedIPs = []string{"10.0.0.0"}
		app.Settings().TrustedProxy.Headers = []string{"x-test-ip"}

		statusCodes := []int{}
		for range 3 {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/rate", nil)
			req.Header.Set("x-test-ip", "127.0.0.1")

			mux.ServeHTTP(rec, req)

			result := rec.Result()

			statusCodes = append(statusCodes, result.StatusCode)
		}

		checkStatusCodes(t, statusCodes, []int{200, 429, 429})
	})

	t.Run("excludedIPs (match)", func(t *testing.T) {
		app.Settings().RateLimits.Enabled = true
		app.Settings().RateLimits.ExcludedIPs = []string{"127.0.0.1"}
		app.Settings().TrustedProxy.Headers = []string{"x-test-ip"}

		statusCodes := []int{}
		for range 3 {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/rate", nil)
			req.Header.Set("x-test-ip", "127.0.0.1")

			mux.ServeHTTP(rec, req)

			result := rec.Result()

			statusCodes = append(statusCodes, result.StatusCode)
		}

		checkStatusCodes(t, statusCodes, []int{200, 200, 200})
	})
}
