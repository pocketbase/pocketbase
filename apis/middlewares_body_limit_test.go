package apis_test

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

func TestBodyLimitMiddleware(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	pbRouter, err := apis.NewRouter(app)
	if err != nil {
		t.Fatal(err)
	}
	pbRouter.POST("/a", func(e *core.RequestEvent) error {
		return e.String(200, "a")
	}) // default global BodyLimit check

	pbRouter.POST("/b", func(e *core.RequestEvent) error {
		return e.String(200, "b")
	}).Bind(apis.BodyLimit(20))

	mux, err := pbRouter.BuildMux()
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		url            string
		size           int64
		expectedStatus int
	}{
		{"/a", 21, 200},
		{"/a", apis.DefaultMaxBodySize + 1, 413},
		{"/b", 20, 200},
		{"/b", 21, 413},
	}

	for _, s := range scenarios {
		t.Run(fmt.Sprintf("%s_%d", s.url, s.size), func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest("POST", s.url, bytes.NewReader(make([]byte, s.size)))
			mux.ServeHTTP(rec, req)

			result := rec.Result()
			defer result.Body.Close()

			if result.StatusCode != s.expectedStatus {
				t.Fatalf("Expected response status %d, got %d", s.expectedStatus, result.StatusCode)
			}
		})
	}
}
