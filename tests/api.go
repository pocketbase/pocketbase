package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
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

	// expectations
	// ---
	ExpectedStatus     int
	ExpectedContent    []string
	NotExpectedContent []string
	ExpectedEvents     map[string]int

	// test hooks
	// ---
	TestAppFactory func() (*TestApp, error)
	BeforeTestFunc func(t *testing.T, app *TestApp, e *echo.Echo)
	AfterTestFunc  func(t *testing.T, app *TestApp, e *echo.Echo)
}

// Test executes the test case/scenario.
func (scenario *ApiScenario) Test(t *testing.T) {
	var testApp *TestApp
	var testAppErr error
	if scenario.TestAppFactory != nil {
		testApp, testAppErr = scenario.TestAppFactory()
	} else {
		testApp, testAppErr = NewTestApp()
	}
	if testAppErr != nil {
		t.Fatalf("Failed to initialize the test app instance: %v", testAppErr)
	}
	defer testApp.Cleanup()

	e, err := apis.InitApi(testApp)
	if err != nil {
		t.Fatal(err)
	}

	if scenario.BeforeTestFunc != nil {
		scenario.BeforeTestFunc(t, testApp, e)
	}

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(scenario.Method, scenario.Url, scenario.Body)

	// add middleware to timeout long-running requests (eg. keep-alive routes)
	e.Pre(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx, cancelFunc := context.WithTimeout(c.Request().Context(), 100*time.Millisecond)
			defer cancelFunc()
			c.SetRequest(c.Request().Clone(ctx))
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

	var prefix = scenario.Name
	if prefix == "" {
		prefix = fmt.Sprintf("%s:%s", scenario.Method, scenario.Url)
	}

	if res.StatusCode != scenario.ExpectedStatus {
		t.Errorf("[%s] Expected status code %d, got %d", prefix, scenario.ExpectedStatus, res.StatusCode)
	}

	if scenario.Delay > 0 {
		time.Sleep(scenario.Delay)
	}

	if len(scenario.ExpectedContent) == 0 && len(scenario.NotExpectedContent) == 0 {
		if len(recorder.Body.Bytes()) != 0 {
			t.Errorf("[%s] Expected empty body, got \n%v", prefix, recorder.Body.String())
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
				t.Errorf("[%s] Cannot find %v in response body \n%v", prefix, item, normalizedBody)
				break
			}
		}

		for _, item := range scenario.NotExpectedContent {
			if strings.Contains(normalizedBody, item) {
				t.Errorf("[%s] Didn't expect %v in response body \n%v", prefix, item, normalizedBody)
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
		t.Errorf("[%s] Expected events %v, got %v", prefix, scenario.ExpectedEvents, testApp.EventCalls)
	}

	for event, expectedCalls := range scenario.ExpectedEvents {
		actualCalls := testApp.EventCalls[event]
		if actualCalls != expectedCalls {
			t.Errorf("[%s] Expected event %s to be called %d, got %d", prefix, event, expectedCalls, actualCalls)
		}
	}

	// @todo consider adding the response body to the AfterTestFunc args
	if scenario.AfterTestFunc != nil {
		scenario.AfterTestFunc(t, testApp, e)
	}
}
