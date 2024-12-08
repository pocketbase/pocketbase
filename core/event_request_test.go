package core_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

func TestEventRequestRealIP(t *testing.T) {
	t.Parallel()

	headers := map[string][]string{
		"CF-Connecting-IP": {"1.2.3.4", "1.1.1.1"},
		"Fly-Client-IP":    {"1.2.3.4", "1.1.1.2"},
		"X-Real-IP":        {"1.2.3.4", "1.1.1.3,1.1.1.4"},
		"X-Forwarded-For":  {"1.2.3.4", "invalid,1.1.1.5,1.1.1.6,invalid"},
	}

	scenarios := []struct {
		name           string
		headers        map[string][]string
		trustedHeaders []string
		useLeftmostIP  bool
		expected       string
	}{
		{
			"no trusted headers",
			headers,
			nil,
			false,
			"127.0.0.1",
		},
		{
			"non-matching trusted header",
			headers,
			[]string{"header1", "header2"},
			false,
			"127.0.0.1",
		},
		{
			"trusted X-Real-IP (rightmost)",
			headers,
			[]string{"header1", "x-real-ip", "x-forwarded-for"},
			false,
			"1.1.1.4",
		},
		{
			"trusted X-Real-IP (leftmost)",
			headers,
			[]string{"header1", "x-real-ip", "x-forwarded-for"},
			true,
			"1.1.1.3",
		},
		{
			"trusted X-Forwarded-For (rightmost)",
			headers,
			[]string{"header1", "x-forwarded-for"},
			false,
			"1.1.1.6",
		},
		{
			"trusted X-Forwarded-For (leftmost)",
			headers,
			[]string{"header1", "x-forwarded-for"},
			true,
			"1.1.1.5",
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			app, err := tests.NewTestApp()
			if err != nil {
				t.Fatal(err)
			}
			defer app.Cleanup()

			app.Settings().TrustedProxy.Headers = s.trustedHeaders
			app.Settings().TrustedProxy.UseLeftmostIP = s.useLeftmostIP

			event := core.RequestEvent{}
			event.App = app

			event.Request, err = http.NewRequest(http.MethodGet, "/", nil)
			if err != nil {
				t.Fatal(err)
			}
			event.Request.RemoteAddr = "127.0.0.1:80" // fallback

			for k, values := range s.headers {
				for _, v := range values {
					event.Request.Header.Add(k, v)
				}
			}

			result := event.RealIP()

			if result != s.expected {
				t.Fatalf("Expected ip %q, got %q", s.expected, result)
			}
		})
	}
}

func TestEventRequestHasSuperUserAuth(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	user, err := app.FindAuthRecordByEmail("users", "test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	superuser, err := app.FindAuthRecordByEmail(core.CollectionNameSuperusers, "test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		name     string
		record   *core.Record
		expected bool
	}{
		{"nil record", nil, false},
		{"regular user record", user, false},
		{"superuser record", superuser, true},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			e := core.RequestEvent{}
			e.Auth = s.record

			result := e.HasSuperuserAuth()

			if result != s.expected {
				t.Fatalf("Expected %v, got %v", s.expected, result)
			}
		})
	}
}

func TestRequestEventRequestInfo(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	userCol, err := app.FindCollectionByNameOrId("users")
	if err != nil {
		t.Fatal(err)
	}

	user1 := core.NewRecord(userCol)
	user1.Id = "user1"
	user1.SetEmail("test1@example.com")

	user2 := core.NewRecord(userCol)
	user2.Id = "user2"
	user2.SetEmail("test2@example.com")

	testBody := `{"a":123,"b":"test"}`

	event := core.RequestEvent{}
	event.Request, err = http.NewRequest("POST", "/test?q1=123&q2=456", strings.NewReader(testBody))
	if err != nil {
		t.Fatal(err)
	}
	event.Request.Header.Add("content-type", "application/json")
	event.Request.Header.Add("x-test", "test")
	event.Set(core.RequestEventKeyInfoContext, "test")
	event.Auth = user1

	t.Run("init", func(t *testing.T) {
		info, err := event.RequestInfo()
		if err != nil {
			t.Fatalf("Failed to resolve request info: %v", err)
		}

		raw, err := json.Marshal(info)
		if err != nil {
			t.Fatalf("Failed to serialize request info: %v", err)
		}
		rawStr := string(raw)

		expected := `{"query":{"q1":"123","q2":"456"},"headers":{"content_type":"application/json","x_test":"test"},"body":{"a":123,"b":"test"},"auth":{"avatar":"","collectionId":"_pb_users_auth_","collectionName":"users","created":"","emailVisibility":false,"file":[],"id":"user1","name":"","rel":"","updated":"","username":"","verified":false},"method":"POST","context":"test"}`

		if expected != rawStr {
			t.Fatalf("Expected\n%v\ngot\n%v", expected, rawStr)
		}
	})

	t.Run("change user and context", func(t *testing.T) {
		event.Set(core.RequestEventKeyInfoContext, "test2")
		event.Auth = user2

		info, err := event.RequestInfo()
		if err != nil {
			t.Fatalf("Failed to resolve request info: %v", err)
		}

		raw, err := json.Marshal(info)
		if err != nil {
			t.Fatalf("Failed to serialize request info: %v", err)
		}
		rawStr := string(raw)

		expected := `{"query":{"q1":"123","q2":"456"},"headers":{"content_type":"application/json","x_test":"test"},"body":{"a":123,"b":"test"},"auth":{"avatar":"","collectionId":"_pb_users_auth_","collectionName":"users","created":"","emailVisibility":false,"file":[],"id":"user2","name":"","rel":"","updated":"","username":"","verified":false},"method":"POST","context":"test2"}`

		if expected != rawStr {
			t.Fatalf("Expected\n%v\ngot\n%v", expected, rawStr)
		}
	})
}

func TestRequestInfoHasSuperuserAuth(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	user, err := app.FindAuthRecordByEmail("users", "test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	superuser, err := app.FindAuthRecordByEmail(core.CollectionNameSuperusers, "test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	event := core.RequestEvent{}
	event.Request, err = http.NewRequest("POST", "/test?q1=123&q2=456", strings.NewReader(`{"a":123,"b":"test"}`))
	if err != nil {
		t.Fatal(err)
	}
	event.Request.Header.Add("content-type", "application/json")

	scenarios := []struct {
		name     string
		record   *core.Record
		expected bool
	}{
		{"nil record", nil, false},
		{"regular user record", user, false},
		{"superuser record", superuser, true},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			event.Auth = s.record

			info, err := event.RequestInfo()
			if err != nil {
				t.Fatalf("Failed to resolve request info: %v", err)
			}

			result := info.HasSuperuserAuth()

			if result != s.expected {
				t.Fatalf("Expected %v, got %v", s.expected, result)
			}
		})
	}
}

func TestRequestInfoClone(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	userCol, err := app.FindCollectionByNameOrId("users")
	if err != nil {
		t.Fatal(err)
	}

	user := core.NewRecord(userCol)
	user.Id = "user1"
	user.SetEmail("test1@example.com")

	event := core.RequestEvent{}
	event.Request, err = http.NewRequest("POST", "/test?q1=123&q2=456", strings.NewReader(`{"a":123,"b":"test"}`))
	if err != nil {
		t.Fatal(err)
	}
	event.Request.Header.Add("content-type", "application/json")
	event.Auth = user

	info, err := event.RequestInfo()
	if err != nil {
		t.Fatalf("Failed to resolve request info: %v", err)
	}

	clone := info.Clone()

	// modify the clone fields to ensure that it is a shallow copy
	clone.Headers["new_header"] = "test"
	clone.Query["new_query"] = "test"
	clone.Body["new_body"] = "test"
	clone.Auth.Id = "user2" // should be a Fresh copy of the record

	// check the original data
	// ---
	originalRaw, err := json.Marshal(info)
	if err != nil {
		t.Fatalf("Failed to serialize original request info: %v", err)
	}
	originalRawStr := string(originalRaw)

	expectedRawStr := `{"query":{"q1":"123","q2":"456"},"headers":{"content_type":"application/json"},"body":{"a":123,"b":"test"},"auth":{"avatar":"","collectionId":"_pb_users_auth_","collectionName":"users","created":"","emailVisibility":false,"file":[],"id":"user1","name":"","rel":"","updated":"","username":"","verified":false},"method":"POST","context":"default"}`
	if expectedRawStr != originalRawStr {
		t.Fatalf("Expected original info\n%v\ngot\n%v", expectedRawStr, originalRawStr)
	}

	// check the clone data
	// ---
	cloneRaw, err := json.Marshal(clone)
	if err != nil {
		t.Fatalf("Failed to serialize clone request info: %v", err)
	}
	cloneRawStr := string(cloneRaw)

	expectedCloneStr := `{"query":{"new_query":"test","q1":"123","q2":"456"},"headers":{"content_type":"application/json","new_header":"test"},"body":{"a":123,"b":"test","new_body":"test"},"auth":{"avatar":"","collectionId":"_pb_users_auth_","collectionName":"users","created":"","emailVisibility":false,"file":[],"id":"user2","name":"","rel":"","updated":"","username":"","verified":false},"method":"POST","context":"default"}`
	if expectedCloneStr != cloneRawStr {
		t.Fatalf("Expected clone info\n%v\ngot\n%v", expectedCloneStr, cloneRawStr)
	}
}
