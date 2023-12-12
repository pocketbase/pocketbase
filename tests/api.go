package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

// ApiScenario defines a single api request test case/scenario.
type ApiScenario struct {
	Name           string
	Method         string
	Url            string
	Body           io.Reader
	RequestHeaders map[string]string

	// Delay adds a delay before checking the expectations usually
	// to ensure that all fired non-awaited go routines have finished
	Delay time.Duration

	// Timeout specifies how long to wait before cancelling the request context.
	//
	// A zero or negative value means that there will be no timeout.
	Timeout time.Duration

	// expectations
	// ---
	ExpectedStatus     int
	ExpectedContent    []string
	NotExpectedContent []string
	ExpectedEvents     map[string]int

	// test hooks
	// ---
	TestAppFactory func(t *testing.T) *TestApp
	BeforeTestFunc func(t *testing.T, app *TestApp, e *echo.Echo)
	AfterTestFunc  func(t *testing.T, app *TestApp, res *http.Response)
}

// Test executes the test scenario.
func (scenario *ApiScenario) Test(t *testing.T) {
	var name = scenario.Name
	if name == "" {
		name = fmt.Sprintf("%s:%s", scenario.Method, scenario.Url)
	}

	t.Run(name, scenario.test)
}

func (scenario *ApiScenario) test(t *testing.T) {
	var testApp *TestApp
	if scenario.TestAppFactory != nil {
		testApp = scenario.TestAppFactory(t)
		if testApp == nil {
			t.Fatal("TestAppFactory must return a non-nill app instance")
		}
	} else {
		var testAppErr error
		testApp, testAppErr = NewTestApp()
		if testAppErr != nil {
			t.Fatalf("Failed to initialize the test app instance: %v", testAppErr)
		}
	}
	defer testApp.Cleanup()

	e, err := apis.InitApi(testApp)
	if err != nil {
		t.Fatal(err)
	}

	// manually trigger the serve event to ensure that custom app routes and middlewares are registered
	testApp.OnBeforeServe().Trigger(&core.ServeEvent{
		App:    testApp,
		Router: e,
	})

	if scenario.BeforeTestFunc != nil {
		scenario.BeforeTestFunc(t, testApp, e)
	}

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(scenario.Method, scenario.Url, scenario.Body)

	// add middleware to timeout long-running requests (eg. keep-alive routes)
	e.Pre(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			slowTimer := time.AfterFunc(3*time.Second, func() {
				t.Logf("[WARN] Long running test %q", scenario.Name)
			})
			defer slowTimer.Stop()

			if scenario.Timeout > 0 {
				ctx, cancelFunc := context.WithTimeout(c.Request().Context(), scenario.Timeout)
				defer cancelFunc()
				c.SetRequest(c.Request().Clone(ctx))
			}

			return next(c)
		}
	})

	// set default header
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// set scenario headers
	for k, v := range scenario.RequestHeaders {
		req.Header.Set(k, v)
	}

	// execute request
	e.ServeHTTP(recorder, req)

	res := recorder.Result()

	if res.StatusCode != scenario.ExpectedStatus {
		t.Errorf("Expected status code %d, got %d", scenario.ExpectedStatus, res.StatusCode)
	}

	if scenario.Delay > 0 {
		time.Sleep(scenario.Delay)
	}

	if len(scenario.ExpectedContent) == 0 && len(scenario.NotExpectedContent) == 0 {
		if len(recorder.Body.Bytes()) != 0 {
			t.Errorf("Expected empty body, got \n%v", recorder.Body.String())
		}
	} else {
		// normalize json response format
		buffer := new(bytes.Buffer)
		err := json.Compact(buffer, recorder.Body.Bytes())
		var normalizedBody string
		if err != nil {
			// not a json...
			normalizedBody = recorder.Body.String()
		} else {
			normalizedBody = buffer.String()
		}

		for _, item := range scenario.ExpectedContent {
			if !strings.Contains(normalizedBody, item) {
				t.Errorf("Cannot find %v in response body \n%v", item, normalizedBody)
				break
			}
		}

		for _, item := range scenario.NotExpectedContent {
			if strings.Contains(normalizedBody, item) {
				t.Errorf("Didn't expect %v in response body \n%v", item, normalizedBody)
				break
			}
		}
	}

	// to minimize the breaking changes we always expect the error
	// events to be called on API error
	if res.StatusCode >= 400 {
		if scenario.ExpectedEvents == nil {
			scenario.ExpectedEvents = map[string]int{}
		}
		if _, ok := scenario.ExpectedEvents["OnBeforeApiError"]; !ok {
			scenario.ExpectedEvents["OnBeforeApiError"] = 1
		}
		if _, ok := scenario.ExpectedEvents["OnAfterApiError"]; !ok {
			scenario.ExpectedEvents["OnAfterApiError"] = 1
		}
	}

	if len(testApp.EventCalls) > len(scenario.ExpectedEvents) {
		t.Errorf("Expected events %v, got %v", scenario.ExpectedEvents, testApp.EventCalls)
	}

	for event, expectedCalls := range scenario.ExpectedEvents {
		actualCalls := testApp.EventCalls[event]
		if actualCalls != expectedCalls {
			t.Errorf("Expected event %s to be called %d, got %d", event, expectedCalls, actualCalls)
		}
	}

	if scenario.AfterTestFunc != nil {
		scenario.AfterTestFunc(t, testApp, res)
	}
}
