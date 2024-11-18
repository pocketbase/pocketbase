package apis_test

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
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
		wait           float64
		authenticated  bool
		expectedStatus int
	}{
		{"/norate", 0, false, 200},
		{"/norate", 0, false, 200},
		{"/norate", 0, false, 200},
		{"/norate", 0, false, 200},
		{"/norate", 0, false, 200},

		{"/rate/a", 0, false, 200},
		{"/rate/a", 0, false, 200},
		{"/rate/a", 0, false, 429},
		{"/rate/a", 0, false, 429},
		{"/rate/a", 1.1, false, 200},
		{"/rate/a", 0, false, 200},
		{"/rate/a", 0, false, 429},

		{"/rate/b", 0, false, 200},
		{"/rate/b", 0, false, 200},
		{"/rate/b", 0, false, 200},
		{"/rate/b", 0, false, 429},
		{"/rate/b", 1.1, false, 200},
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
		{"/rate/guest", 1, true, 200},
		{"/rate/guest", 0, true, 200},
		{"/rate/guest", 0, true, 429},
		{"/rate/guest", 0, true, 429},
	}

	for _, s := range scenarios {
		t.Run(s.url, func(t *testing.T) {
			if s.wait > 0 {
				time.Sleep(time.Duration(s.wait) * time.Second)
			}

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", s.url, nil)

			if s.authenticated {
				auth, err := app.FindAuthRecordByEmail("users", "test@example.com")
				if err != nil {
					t.Fatal(err)
				}

				token, err := auth.NewAuthToken()
				if err != nil {
					t.Fatal(err)
				}

				req.Header.Add("Authorization", token)
			}

			mux.ServeHTTP(rec, req)

			result := rec.Result()

			if result.StatusCode != s.expectedStatus {
				t.Fatalf("Expected response status %d, got %d", s.expectedStatus, result.StatusCode)
			}
		})
	}
}
