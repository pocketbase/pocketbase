package apis_test

import (
	"encoding/json"
	"net/http"
	"path"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/pocketbase/pocketbase/tools/subscriptions"
)

func TestAuthRecordFromJson(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	id := "4q1xlclmfloku33"

	sql := `
	    SELECT row_to_json(users) FROM users WHERE id = '4q1xlclmfloku33' LIMIT 1
	`
	var jsonStr string
	err := app.DB().NewQuery(sql).Row(&jsonStr)
	if err != nil {
		t.Fatalf("Failed to get JSON string: %v", err)
	}

	authRecord, err := apis.AuthRecordFromJson(app, "users", jsonStr)
	if err != nil {
		t.Fatalf("Failed to create AuthRecord from JSON: %v", err)
	}
	realAuthRecord, err := app.FindRecordById("users", id)
	if err != nil {
		t.Fatalf("Failed to find record by ID: %v", err)
	}

	json1, _ := json.Marshal(authRecord.Clone().Unhide(authRecord.Collection().Fields.FieldNames()...).IgnoreEmailVisibility(true).PublicExport())
	json2, _ := json.Marshal(realAuthRecord.Clone().Unhide(realAuthRecord.Collection().Fields.FieldNames()...).IgnoreEmailVisibility(true).PublicExport())
	if string(json1) != string(json2) {
		t.Fatalf("AuthRecord JSON does not match: expected \n%s\n, got \n%s", string(json2), string(json1))
	}
}

func TestRealtimeBridge(t *testing.T) {
	_, currentFile, _, _ := runtime.Caller(0)
	config := core.BaseAppConfig{
		DataDir:          filepath.Join(path.Dir(currentFile), "..", "tests", "data"),
		EncryptionEnv:    "pb_test_env",
		PostgresURL:      "postgres://user:pass@127.0.0.1:5432/postgres?sslmode=disable",
		PostgresDataDB:   "pb_test_" + security.RandomString(5),
		PostgresAuxDB:    "pb_test_" + security.RandomString(5) + "_aux",
		IsRealtimeBridge: true,
		IsDev:            false, // Turn it on in unit tests to see sql erros.
	}
	app1, _ := tests.NewTestAppWithConfig(config)
	// run second app with the same config (So that they share the same postgres db)
	config.DataDir = app1.DataDir() // use the temp dir of the first app
	app2 := &tests.TestApp{
		BaseApp:    core.NewBaseApp(config),
		EventCalls: make(map[string]int),
		TestMailer: &tests.TestMailer{},
	}
	err := app2.Bootstrap()
	if err != nil {
		t.Fatalf("Failed to bootstrap app2: %v", err)
	}
	defer app1.Cleanup()
	defer app2.Cleanup()

	// simulate we have two running pocketbase instances
	bridgeA, ok := newBridge(app1, ":8090")
	if !ok {
		t.Fatalf("Failed to create bridgeA")
	}
	bridgeB, ok := newBridge(app2, ":8091")
	if !ok {
		t.Fatalf("Failed to create bridgeB")
	}

	// wait 0.5 seconds for the bridge to initialize
	time.Sleep(time.Second / 2)

	local_client1_in_channelA := apis.NewBridgedClient(bridgeA)
	// test client is used to capture messages sent to the ws clients
	testClient := &TestClient{Client: local_client1_in_channelA.Client}
	local_client1_in_channelA.Client = testClient

	// client1 online
	app1.SubscriptionsBroker().Register(local_client1_in_channelA)
	local_client1_in_channelA.Subscribe("sub1")
	local_client1_in_channelA.BroadcastChanges()

	// wait 0.5 second for online event to be broadcasted
	time.Sleep(time.Second / 2)

	remote_client1_in_channelB_, _ := bridgeB.App().SubscriptionsBroker().ClientById(local_client1_in_channelA.Id())
	remote_client1_in_channelB, _ := remote_client1_in_channelB_.(*apis.BridgedClient)

	t.Run("Test client online events broadcasted to other servers", func(t *testing.T) {
		// bridgeB in another pocketbase instance should receive the client online event
		if remote_client1_in_channelB == nil {
			t.Fatalf("Failed to find remote client by ID: %v", local_client1_in_channelA.Id())
		}
		if !remote_client1_in_channelB.IsRemoteClient() {
			t.Fatalf("Expected remote client to be true, but got false")
		}
		if local_client1_in_channelA.IsRemoteClient() {
			t.Fatalf("Expected local client to be false, but got true")
		}
		if remote_client1_in_channelB.Id() != local_client1_in_channelA.Id() {
			t.Fatalf("Expected the client ID to be the same, but got %s and %s", remote_client1_in_channelB.Id(), local_client1_in_channelA.Id())
		}
	})

	t.Run("Test local client message sent to local channel", func(t *testing.T) {
		// Send message from to client1 from the same pocketbase instance
		message := subscriptions.Message{Name: "test1", Data: []byte("Hello World!")}
		local_client1_in_channelA.Send(message)
		if len(testClient.SentMessages) != 1 {
			t.Fatalf("Expected local client channel to have first message, but got empty")
		} else {
			msg := testClient.SentMessages[0]
			if msg.Name != "test1" || string(msg.Data) != "Hello World!" {
				t.Fatalf("Expected message to be 'Hello World!', but got %s", string(msg.Data))
			}
		}
		if len(remote_client1_in_channelB.Channel()) != 0 {
			t.Fatalf("Expected remote client channel to be empty, but got %d messages", len(remote_client1_in_channelB.Channel()))
		}
	})

	t.Run("Test remote client message sent to remote channel", func(t *testing.T) {
		// send message to client1 from different pocketbase instance
		message := subscriptions.Message{Name: "test2", Data: []byte("Hello World!")}
		remote_client1_in_channelB.Send(message)
		time.Sleep(time.Second * 5) // wait for the message to be sent
		if len(testClient.SentMessages) != 2 {
			t.Fatalf("Expected local client channel to have second message, but got empty")
		} else {
			msg := testClient.SentMessages[1]
			if msg.Name != "test2" || string(msg.Data) != "Hello World!" {
				t.Fatalf("Expected message to be 'Hello World!', but got %s", string(msg.Data))
			}
		}
		if len(remote_client1_in_channelB.Channel()) != 0 {
			t.Fatalf("Expected remote client channel to be empty, but got %d messages", len(remote_client1_in_channelB.Channel()))
		}
	})

	t.Run("Test client auth state broadcasted to other servers", func(t *testing.T) {
		// auth state updated in app2
		authRecord1, err := app2.FindRecordById("users", "4q1xlclmfloku33")
		if err != nil {
			t.Fatalf("Failed to find record by ID: %v", err)
		}
		remote_client1_in_channelB.Set(apis.RealtimeClientAuthKey, authRecord1)
		remote_client1_in_channelB.BroadcastChanges()

		time.Sleep(time.Second / 2)

		// then auth state should be updated in app1
		authRecord2, ok := local_client1_in_channelA.Get(apis.RealtimeClientAuthKey).(*core.Record)
		if !ok {
			t.Fatalf("Failed to get auth record from local client")
		}
		if authRecord2.Id != authRecord1.Id {
			t.Fatalf("Expected auth record ID to be %s, but got %s", authRecord1.Id, authRecord2.Id)
		}

		json1, _ := json.Marshal(authRecord1.Clone().Unhide(authRecord1.Collection().Fields.FieldNames()...).IgnoreEmailVisibility(true).PublicExport())
		json2, _ := json.Marshal(authRecord2.Clone().Unhide(authRecord2.Collection().Fields.FieldNames()...).IgnoreEmailVisibility(true).PublicExport())
		if string(json1) != string(json2) {
			t.Fatalf("AuthRecord JSON does not match: expected \n%s\n, got \n%s", string(json2), string(json1))
		}
	})

	t.Run("Test client offline event broadcasted to other servers", func(t *testing.T) {
		// client1 offline
		local_client1_in_channelA.BroadcastGoOffline()

		time.Sleep(1 * time.Second)

		// assume bridgeB receives the client offline notification
		client1_not_found_in_channelB, _ := app2.SubscriptionsBroker().ClientById(local_client1_in_channelA.Id())
		if client1_not_found_in_channelB != nil {
			t.Fatalf("Remote client should be nil, but got: %v", client1_not_found_in_channelB)
		}
	})

	t.Logf("Test passed")
}

func TestRealtimeBridge_ChannelOffline(t *testing.T) {
	_, currentFile, _, _ := runtime.Caller(0)
	config := core.BaseAppConfig{
		DataDir:          filepath.Join(path.Dir(currentFile), "..", "tests", "data"),
		EncryptionEnv:    "pb_test_env",
		PostgresURL:      "postgres://user:pass@127.0.0.1:5432/postgres?sslmode=disable",
		PostgresDataDB:   "pb_test_" + security.RandomString(5),
		PostgresAuxDB:    "pb_test_" + security.RandomString(5) + "_aux",
		IsRealtimeBridge: true,
		IsDev:            false, // Turn it on in unit tests to see sql erros.
	}
	app1, _ := tests.NewTestAppWithConfig(config)
	// run second app with the same config (So that they share the same postgres db)
	config.DataDir = app1.DataDir() // use the temp dir of the first app
	app2 := &tests.TestApp{
		BaseApp:    core.NewBaseApp(config),
		EventCalls: make(map[string]int),
		TestMailer: &tests.TestMailer{},
	}
	err := app2.Bootstrap()
	if err != nil {
		t.Fatalf("Failed to bootstrap app2: %v", err)
	}
	defer app1.Cleanup()
	defer app2.Cleanup()

	// simulate we have two running pocketbase instances
	bridgeA, ok := newBridge(app1, ":8090")
	if !ok {
		t.Fatalf("Failed to create bridgeA")
	}
	bridgeB, ok := newBridge(app2, ":8091")
	if !ok {
		t.Fatalf("Failed to create bridgeB")
	}
	_ = bridgeB

	// wait 0.5 seconds for the bridge to initialize
	time.Sleep(time.Second / 2)

	local_client1_in_channelA := apis.NewBridgedClient(bridgeA)
	// test client is used to capture messages sent to the ws clients
	testClient := &TestClient{Client: local_client1_in_channelA.Client}
	local_client1_in_channelA.Client = testClient

	// client1 online
	app1.SubscriptionsBroker().Register(local_client1_in_channelA)
	local_client1_in_channelA.Subscribe("sub1")
	local_client1_in_channelA.BroadcastChanges()

	// wait 0.5 second for online event to be broadcasted
	time.Sleep(time.Second / 2)

	t.Run("Test channel offline event broadcasted to other servers", func(t *testing.T) {
		if app1.SubscriptionsBroker().TotalClients() != 1 {
			t.Fatalf("Expected 1 client in channelA, but got %d", app1.SubscriptionsBroker().TotalClients())
		}
		if app2.SubscriptionsBroker().TotalClients() != 1 {
			t.Fatalf("Expected 1 client in channelB, but got %d", app2.SubscriptionsBroker().TotalClients())
		}

		// Pretend that heartbeat event in app1 stopped 1 hour ago
		_, err := app1.DB().NewQuery(`
			UPDATE "_realtimeChannels"
			SET "validUntil" = now() - interval '1 hour'
		`).Execute()
		if err != nil {
			t.Fatalf("Failed to update client heartbeat: %v", err)
		}

		// Now a new pocketbase instance (app3) is started.
		// It should send heartbeat events to postgres as well as
		// cleaning the dead realtime channels (eg: app1).
		app3 := &tests.TestApp{
			BaseApp:    core.NewBaseApp(config),
			EventCalls: make(map[string]int),
			TestMailer: &tests.TestMailer{},
		}
		err = app3.Bootstrap()
		if err != nil {
			t.Fatalf("Failed to bootstrap app2: %v", err)
		}
		defer app3.Cleanup()
		_, _ = newBridge(app3, ":8092")

		// wait 1 second. app3 should clean the dead channels
		// in its heartbeat loop operations.
		time.Sleep(time.Second)

		// check if the dead channels are cleaned
		if app1.SubscriptionsBroker().TotalClients() != 0 {
			t.Fatalf("Expected 0 clients in app1, but got %d", app1.SubscriptionsBroker().TotalClients())
		}
		if app2.SubscriptionsBroker().TotalClients() != 0 {
			t.Fatalf("Expected 0 clients in app2, but got %d", app2.SubscriptionsBroker().TotalClients())
		}
		if app3.SubscriptionsBroker().TotalClients() != 0 {
			t.Fatalf("Expected 0 clients in app3, but got %d", app3.SubscriptionsBroker().TotalClients())
		}
	})

	t.Logf("Test passed")
}

func newBridge(app core.App, addr string) (*apis.RealtimeBridge, bool) {

	// Create a new RealtimeBridge instance
	serveEvent := new(core.ServeEvent)
	serveEvent.App = app
	serveEvent.Server = &http.Server{Addr: addr}
	serveEvent.Router, _ = apis.NewRouter(app) // Set to nil for this test

	app.OnServe().Trigger(serveEvent, func(se *core.ServeEvent) error {
		return se.Next()
	})

	bridge, ok := app.Store().Get(apis.RealtimeBridgeInstanceKey).(*apis.RealtimeBridge)
	return bridge, ok
}

type TestClient struct {
	subscriptions.Client
	SentMessages []subscriptions.Message
}

func (c *TestClient) Send(message subscriptions.Message) {
	// c.Client.Send(message)
	c.SentMessages = append(c.SentMessages, message)
}
