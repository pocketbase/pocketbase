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

	mux, err := pbRouter.BuildMux()
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		url            string
		wait           float64
		expectedStatus int
	}{
		{"/norate", 0, 200},
		{"/norate", 0, 200},
		{"/norate", 0, 200},
		{"/norate", 0, 200},
		{"/norate", 0, 200},

		{"/rate/a", 0, 200},
		{"/rate/a", 0, 200},
		{"/rate/a", 0, 429},
		{"/rate/a", 0, 429},
		{"/rate/a", 1.1, 200},
		{"/rate/a", 0, 200},
		{"/rate/a", 0, 429},

		{"/rate/b", 0, 200},
		{"/rate/b", 0, 200},
		{"/rate/b", 0, 200},
		{"/rate/b", 0, 429},
		{"/rate/b", 1.1, 200},
		{"/rate/b", 0, 200},
		{"/rate/b", 0, 200},
		{"/rate/b", 0, 429},
	}

	for _, s := range scenarios {
		t.Run(s.url, func(t *testing.T) {
			if s.wait > 0 {
				time.Sleep(time.Duration(s.wait) * time.Second)
			}

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", s.url, nil)
			mux.ServeHTTP(rec, req)

			result := rec.Result()

			if result.StatusCode != s.expectedStatus {
				t.Fatalf("Expected response status %d, got %d", s.expectedStatus, result.StatusCode)
			}
		})
	}
}
