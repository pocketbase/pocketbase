package apis_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/subscriptions"
)

func TestRealtimeConnect(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Method:         http.MethodGet,
			Url:            "/api/realtime",
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`id:`,
				`event:PB_CONNECT`,
				`data:{"clientId":`,
			},
			ExpectedEvents: map[string]int{
				"OnRealtimeConnectRequest": 1,
			},
			AfterFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				if len(app.SubscriptionsBroker().Clients()) != 0 {
					t.Errorf("Expected the subscribers to be removed after connection close, found %d", len(app.SubscriptionsBroker().Clients()))
				}
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRealtimeSubscribe(t *testing.T) {
	client := subscriptions.NewDefaultClient()

	resetClient := func() {
		client.Unsubscribe()
		client.Set(apis.ContextAdminKey, nil)
		client.Set(apis.ContextUserKey, nil)
	}

	scenarios := []tests.ApiScenario{
		{
			Name:            "missing client",
			Method:          http.MethodPost,
			Url:             "/api/realtime",
			Body:            strings.NewReader(`{"clientId":"missing","subscriptions":["test1", "test2"]}`),
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:           "existing client - empty subscriptions",
			Method:         http.MethodPost,
			Url:            "/api/realtime",
			Body:           strings.NewReader(`{"clientId":"` + client.Id() + `","subscriptions":[]}`),
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnRealtimeBeforeSubscribeRequest": 1,
				"OnRealtimeAfterSubscribeRequest":  1,
			},
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				client.Subscribe("test0")
				app.SubscriptionsBroker().Register(client)
			},
			AfterFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				if len(client.Subscriptions()) != 0 {
					t.Errorf("Expected no subscriptions, got %v", client.Subscriptions())
				}
				resetClient()
			},
		},
		{
			Name:           "existing client - 2 new subscriptions",
			Method:         http.MethodPost,
			Url:            "/api/realtime",
			Body:           strings.NewReader(`{"clientId":"` + client.Id() + `","subscriptions":["test1", "test2"]}`),
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnRealtimeBeforeSubscribeRequest": 1,
				"OnRealtimeAfterSubscribeRequest":  1,
			},
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				client.Subscribe("test0")
				app.SubscriptionsBroker().Register(client)
			},
			AfterFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				expectedSubs := []string{"test1", "test2"}
				if len(expectedSubs) != len(client.Subscriptions()) {
					t.Errorf("Expected subscriptions %v, got %v", expectedSubs, client.Subscriptions())
				}

				for _, s := range expectedSubs {
					if !client.HasSubscription(s) {
						t.Errorf("Cannot find %q subscription in %v", s, client.Subscriptions())
					}
				}
				resetClient()
			},
		},
		{
			Name:   "existing client - authorized admin",
			Method: http.MethodPost,
			Url:    "/api/realtime",
			Body:   strings.NewReader(`{"clientId":"` + client.Id() + `","subscriptions":["test1", "test2"]}`),
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnRealtimeBeforeSubscribeRequest": 1,
				"OnRealtimeAfterSubscribeRequest":  1,
			},
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				app.SubscriptionsBroker().Register(client)
			},
			AfterFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				admin, _ := client.Get(apis.ContextAdminKey).(*models.Admin)
				if admin == nil {
					t.Errorf("Expected admin auth model, got nil")
				}
				resetClient()
			},
		},
		{
			Name:   "existing client - authorized user",
			Method: http.MethodPost,
			Url:    "/api/realtime",
			Body:   strings.NewReader(`{"clientId":"` + client.Id() + `","subscriptions":["test1", "test2"]}`),
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnRealtimeBeforeSubscribeRequest": 1,
				"OnRealtimeAfterSubscribeRequest":  1,
			},
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				app.SubscriptionsBroker().Register(client)
			},
			AfterFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				user, _ := client.Get(apis.ContextUserKey).(*models.User)
				if user == nil {
					t.Errorf("Expected user auth model, got nil")
				}
				resetClient()
			},
		},
		{
			Name:   "existing client - mismatched auth",
			Method: http.MethodPost,
			Url:    "/api/realtime",
			Body:   strings.NewReader(`{"clientId":"` + client.Id() + `","subscriptions":["test1", "test2"]}`),
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				initialAuth := &models.User{}
				initialAuth.RefreshId()
				client.Set(apis.ContextUserKey, initialAuth)

				app.SubscriptionsBroker().Register(client)
			},
			AfterFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				user, _ := client.Get(apis.ContextUserKey).(*models.User)
				if user == nil {
					t.Errorf("Expected user auth model, got nil")
				}
				resetClient()
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRealtimeUserDeleteEvent(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	apis.InitApi(testApp)

	user, err := testApp.Dao().FindUserByEmail("test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	client := subscriptions.NewDefaultClient()
	client.Set(apis.ContextUserKey, user)
	testApp.SubscriptionsBroker().Register(client)

	testApp.OnModelAfterDelete().Trigger(&core.ModelEvent{Dao: testApp.Dao(), Model: user})

	if len(testApp.SubscriptionsBroker().Clients()) != 0 {
		t.Fatalf("Expected no subscription clients, found %d", len(testApp.SubscriptionsBroker().Clients()))
	}
}

func TestRealtimeUserUpdateEvent(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	apis.InitApi(testApp)

	user1, err := testApp.Dao().FindUserByEmail("test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	client := subscriptions.NewDefaultClient()
	client.Set(apis.ContextUserKey, user1)
	testApp.SubscriptionsBroker().Register(client)

	// refetch the user and change its email
	user2, err := testApp.Dao().FindUserByEmail("test@example.com")
	if err != nil {
		t.Fatal(err)
	}
	user2.Email = "new@example.com"

	testApp.OnModelAfterUpdate().Trigger(&core.ModelEvent{Dao: testApp.Dao(), Model: user2})

	clientUser, _ := client.Get(apis.ContextUserKey).(*models.User)
	if clientUser.Email != user2.Email {
		t.Fatalf("Expected user with email %q, got %q", user2.Email, clientUser.Email)
	}
}

func TestRealtimeAdminDeleteEvent(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	apis.InitApi(testApp)

	admin, err := testApp.Dao().FindAdminByEmail("test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	client := subscriptions.NewDefaultClient()
	client.Set(apis.ContextAdminKey, admin)
	testApp.SubscriptionsBroker().Register(client)

	testApp.OnModelAfterDelete().Trigger(&core.ModelEvent{Dao: testApp.Dao(), Model: admin})

	if len(testApp.SubscriptionsBroker().Clients()) != 0 {
		t.Fatalf("Expected no subscription clients, found %d", len(testApp.SubscriptionsBroker().Clients()))
	}
}

func TestRealtimeAdminUpdateEvent(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	apis.InitApi(testApp)

	admin1, err := testApp.Dao().FindAdminByEmail("test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	client := subscriptions.NewDefaultClient()
	client.Set(apis.ContextAdminKey, admin1)
	testApp.SubscriptionsBroker().Register(client)

	// refetch the user and change its email
	admin2, err := testApp.Dao().FindAdminByEmail("test@example.com")
	if err != nil {
		t.Fatal(err)
	}
	admin2.Email = "new@example.com"

	testApp.OnModelAfterUpdate().Trigger(&core.ModelEvent{Dao: testApp.Dao(), Model: admin2})

	clientAdmin, _ := client.Get(apis.ContextAdminKey).(*models.Admin)
	if clientAdmin.Email != admin2.Email {
		t.Fatalf("Expected user with email %q, got %q", admin2.Email, clientAdmin.Email)
	}
}
