package apis_test

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/subscriptions"
)

func TestRecordAuthWithOAuth2Redirect(t *testing.T) {
	t.Parallel()

	clientStubs := make([]map[string]subscriptions.Client, 0, 10)

	for i := 0; i < 10; i++ {
		c1 := subscriptions.NewDefaultClient()

		c2 := subscriptions.NewDefaultClient()
		c2.Subscribe("@oauth2")

		c3 := subscriptions.NewDefaultClient()
		c3.Subscribe("test1", "@oauth2")

		c4 := subscriptions.NewDefaultClient()
		c4.Subscribe("test1", "test2")

		c5 := subscriptions.NewDefaultClient()
		c5.Subscribe("@oauth2")
		c5.Discard()

		clientStubs = append(clientStubs, map[string]subscriptions.Client{
			"c1": c1,
			"c2": c2,
			"c3": c3,
			"c4": c4,
			"c5": c5,
		})
	}

	checkFailureRedirect := func(t testing.TB, app *tests.TestApp, res *http.Response) {
		loc := res.Header.Get("Location")
		if !strings.Contains(loc, "/oauth2-redirect-failure") {
			t.Fatalf("Expected failure redirect, got %q", loc)
		}
	}

	checkSuccessRedirect := func(t testing.TB, app *tests.TestApp, res *http.Response) {
		loc := res.Header.Get("Location")
		if !strings.Contains(loc, "/oauth2-redirect-success") {
			t.Fatalf("Expected success redirect, got %q", loc)
		}
	}

	// note: don't exit because it is usually called as part of a separate goroutine
	checkClientMessages := func(t testing.TB, clientId string, msg subscriptions.Message, expectedMessages map[string][]string) {
		if len(expectedMessages[clientId]) == 0 {
			t.Errorf("Unexpected client %q message, got %q:\n%q", clientId, msg.Name, msg.Data)
			return
		}

		if msg.Name != "@oauth2" {
			t.Errorf("Expected @oauth2 msg.Name, got %q", msg.Name)
			return
		}

		for _, txt := range expectedMessages[clientId] {
			if !strings.Contains(string(msg.Data), txt) {
				t.Errorf("Failed to find %q in \n%s", txt, msg.Data)
				return
			}
		}
	}

	beforeTestFunc := func(
		clients map[string]subscriptions.Client,
		expectedMessages map[string][]string,
	) func(testing.TB, *tests.TestApp, *core.ServeEvent) {
		return func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
			for _, client := range clients {
				app.SubscriptionsBroker().Register(client)
			}

			ctx, cancelFunc := context.WithTimeout(context.Background(), 100*time.Millisecond)

			// add to the app store so that it can be cancelled manually after test completion
			app.Store().Set("cancelFunc", cancelFunc)

			go func() {
				defer cancelFunc()

				for {
					select {
					case msg, ok := <-clients["c1"].Channel():
						if ok {
							checkClientMessages(t, "c1", msg, expectedMessages)
						} else {
							t.Errorf("Unexpected c1 closed channel")
						}
					case msg, ok := <-clients["c2"].Channel():
						if ok {
							checkClientMessages(t, "c2", msg, expectedMessages)
						} else {
							t.Errorf("Unexpected c2 closed channel")
						}
					case msg, ok := <-clients["c3"].Channel():
						if ok {
							checkClientMessages(t, "c3", msg, expectedMessages)
						} else {
							t.Errorf("Unexpected c3 closed channel")
						}
					case msg, ok := <-clients["c4"].Channel():
						if ok {
							checkClientMessages(t, "c4", msg, expectedMessages)
						} else {
							t.Errorf("Unexpected c4 closed channel")
						}
					case _, ok := <-clients["c5"].Channel():
						if ok {
							t.Errorf("Expected c5 channel to be closed")
						}
					case <-ctx.Done():
						for _, c := range clients {
							c.Discard()
						}
						return
					}
				}
			}()
		}
	}

	scenarios := []tests.ApiScenario{
		{
			Name:           "no state query param",
			Method:         http.MethodGet,
			URL:            "/api/oauth2-redirect?code=123",
			BeforeTestFunc: beforeTestFunc(clientStubs[0], nil),
			ExpectedStatus: http.StatusTemporaryRedirect,
			ExpectedEvents: map[string]int{"*": 0},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				app.Store().Get("cancelFunc").(context.CancelFunc)()

				checkFailureRedirect(t, app, res)
			},
		},
		{
			Name:           "invalid or missing client",
			Method:         http.MethodGet,
			URL:            "/api/oauth2-redirect?code=123&state=missing",
			BeforeTestFunc: beforeTestFunc(clientStubs[1], nil),
			ExpectedStatus: http.StatusTemporaryRedirect,
			ExpectedEvents: map[string]int{"*": 0},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				app.Store().Get("cancelFunc").(context.CancelFunc)()

				checkFailureRedirect(t, app, res)
			},
		},
		{
			Name:   "no code query param",
			Method: http.MethodGet,
			URL:    "/api/oauth2-redirect?state=" + clientStubs[2]["c3"].Id(),
			BeforeTestFunc: beforeTestFunc(clientStubs[2], map[string][]string{
				"c3": {`"state":"` + clientStubs[2]["c3"].Id(), `"code":""`},
			}),
			ExpectedStatus: http.StatusTemporaryRedirect,
			ExpectedEvents: map[string]int{"*": 0},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				app.Store().Get("cancelFunc").(context.CancelFunc)()

				checkFailureRedirect(t, app, res)

				if clientStubs[2]["c3"].HasSubscription("@oauth2") {
					t.Fatalf("Expected oauth2 subscription to be removed")
				}
			},
		},
		{
			Name:   "error query param",
			Method: http.MethodGet,
			URL:    "/api/oauth2-redirect?error=example&code=123&state=" + clientStubs[3]["c3"].Id(),
			BeforeTestFunc: beforeTestFunc(clientStubs[3], map[string][]string{
				"c3": {`"state":"` + clientStubs[3]["c3"].Id(), `"code":"123"`, `"error":"example"`},
			}),
			ExpectedStatus: http.StatusTemporaryRedirect,
			ExpectedEvents: map[string]int{"*": 0},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				app.Store().Get("cancelFunc").(context.CancelFunc)()

				checkFailureRedirect(t, app, res)

				if clientStubs[3]["c3"].HasSubscription("@oauth2") {
					t.Fatalf("Expected oauth2 subscription to be removed")
				}
			},
		},
		{
			Name:           "discarded client with @oauth2 subscription",
			Method:         http.MethodGet,
			URL:            "/api/oauth2-redirect?code=123&state=" + clientStubs[4]["c5"].Id(),
			BeforeTestFunc: beforeTestFunc(clientStubs[4], nil),
			ExpectedStatus: http.StatusTemporaryRedirect,
			ExpectedEvents: map[string]int{"*": 0},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				app.Store().Get("cancelFunc").(context.CancelFunc)()

				checkFailureRedirect(t, app, res)
			},
		},
		{
			Name:           "client without @oauth2 subscription",
			Method:         http.MethodGet,
			URL:            "/api/oauth2-redirect?code=123&state=" + clientStubs[4]["c4"].Id(),
			BeforeTestFunc: beforeTestFunc(clientStubs[5], nil),
			ExpectedStatus: http.StatusTemporaryRedirect,
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				app.Store().Get("cancelFunc").(context.CancelFunc)()

				checkFailureRedirect(t, app, res)
			},
		},
		{
			Name:   "client with @oauth2 subscription",
			Method: http.MethodGet,
			URL:    "/api/oauth2-redirect?code=123&state=" + clientStubs[6]["c3"].Id(),
			BeforeTestFunc: beforeTestFunc(clientStubs[6], map[string][]string{
				"c3": {`"state":"` + clientStubs[6]["c3"].Id(), `"code":"123"`},
			}),
			ExpectedStatus: http.StatusTemporaryRedirect,
			ExpectedEvents: map[string]int{"*": 0},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				app.Store().Get("cancelFunc").(context.CancelFunc)()

				checkSuccessRedirect(t, app, res)

				if clientStubs[6]["c3"].HasSubscription("@oauth2") {
					t.Fatalf("Expected oauth2 subscription to be removed")
				}
			},
		},
		{
			Name:   "(POST) client with @oauth2 subscription",
			Method: http.MethodPost,
			URL:    "/api/oauth2-redirect",
			Body:   strings.NewReader("code=123&state=" + clientStubs[7]["c3"].Id()),
			Headers: map[string]string{
				"content-type": "application/x-www-form-urlencoded",
			},
			BeforeTestFunc: beforeTestFunc(clientStubs[7], map[string][]string{
				"c3": {`"state":"` + clientStubs[7]["c3"].Id(), `"code":"123"`},
			}),
			ExpectedStatus: http.StatusSeeOther,
			ExpectedEvents: map[string]int{"*": 0},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				app.Store().Get("cancelFunc").(context.CancelFunc)()

				checkSuccessRedirect(t, app, res)

				if clientStubs[7]["c3"].HasSubscription("@oauth2") {
					t.Fatalf("Expected oauth2 subscription to be removed")
				}
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
