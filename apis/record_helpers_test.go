package apis_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/router"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestEnrichRecords(t *testing.T) {
	t.Parallel()

	// mock test data
	// ---
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	freshRecords := func(records []*core.Record) []*core.Record {
		result := make([]*core.Record, len(records))
		for i, r := range records {
			result[i] = r.Fresh()
		}
		return result
	}

	user, err := app.FindAuthRecordByEmail("users", "test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	superuser, err := app.FindAuthRecordByEmail(core.CollectionNameSuperusers, "test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	usersRecords, err := app.FindRecordsByIds("users", []string{"4q1xlclmfloku33", "bgs820n361vj1qd"})
	if err != nil {
		t.Fatal(err)
	}

	nologinRecords, err := app.FindRecordsByIds("nologin", []string{"dc49k6jgejn40h3", "oos036e9xvqeexy"})
	if err != nil {
		t.Fatal(err)
	}

	demo1Records, err := app.FindRecordsByIds("demo1", []string{"al1h9ijdeojtsjy", "84nmscqy84lsi1t"})
	if err != nil {
		t.Fatal(err)
	}

	demo5Records, err := app.FindRecordsByIds("demo5", []string{"la4y2w4o98acwuj", "qjeql998mtp1azp"})
	if err != nil {
		t.Fatal(err)
	}

	// temp update the view rule to ensure that request context is set to "expand"
	demo4, err := app.FindCollectionByNameOrId("demo4")
	if err != nil {
		t.Fatal(err)
	}
	demo4.ViewRule = types.Pointer("@request.context = 'expand'")
	if err := app.Save(demo4); err != nil {
		t.Fatal(err)
	}
	// ---

	scenarios := []struct {
		name           string
		auth           *core.Record
		records        []*core.Record
		queryExpand    string
		defaultExpands []string
		expected       []string
		notExpected    []string
	}{
		// email visibility checks
		{
			name:           "[emailVisibility] guest",
			auth:           nil,
			records:        freshRecords(usersRecords),
			queryExpand:    "",
			defaultExpands: nil,
			expected: []string{
				`"customField":"123"`,
				`"test3@example.com"`, // emailVisibility=true
			},
			notExpected: []string{
				`"test@example.com"`,
			},
		},
		{
			name:           "[emailVisibility] owner",
			auth:           user,
			records:        freshRecords(usersRecords),
			queryExpand:    "",
			defaultExpands: nil,
			expected: []string{
				`"customField":"123"`,
				`"test3@example.com"`, // emailVisibility=true
				`"test@example.com"`,  // owner
			},
		},
		{
			name:           "[emailVisibility] manager",
			auth:           user,
			records:        freshRecords(nologinRecords),
			queryExpand:    "",
			defaultExpands: nil,
			expected: []string{
				`"customField":"123"`,
				`"test3@example.com"`,
				`"test@example.com"`,
			},
		},
		{
			name:           "[emailVisibility] superuser",
			auth:           superuser,
			records:        freshRecords(nologinRecords),
			queryExpand:    "",
			defaultExpands: nil,
			expected: []string{
				`"customField":"123"`,
				`"test3@example.com"`,
				`"test@example.com"`,
			},
		},
		{
			name:           "[emailVisibility + expand] recursive auth rule checks (regular user)",
			auth:           user,
			records:        freshRecords(demo1Records),
			queryExpand:    "",
			defaultExpands: []string{"rel_many"},
			expected: []string{
				`"customField":"123"`,
				`"expand":{"rel_many"`,
				`"expand":{}`,
				`"test@example.com"`,
			},
			notExpected: []string{
				`"id":"bgs820n361vj1qd"`,
				`"id":"oap640cot4yru2s"`,
			},
		},
		{
			name:           "[emailVisibility + expand] recursive auth rule checks (superuser)",
			auth:           superuser,
			records:        freshRecords(demo1Records),
			queryExpand:    "",
			defaultExpands: []string{"rel_many"},
			expected: []string{
				`"customField":"123"`,
				`"test@example.com"`,
				`"expand":{"rel_many"`,
				`"id":"bgs820n361vj1qd"`,
				`"id":"4q1xlclmfloku33"`,
				`"id":"oap640cot4yru2s"`,
			},
			notExpected: []string{
				`"expand":{}`,
			},
		},

		// expand checks
		{
			name:           "[expand] guest (query)",
			auth:           nil,
			records:        freshRecords(usersRecords),
			queryExpand:    "rel",
			defaultExpands: nil,
			expected: []string{
				`"customField":"123"`,
				`"expand":{"rel"`,
				`"id":"llvuca81nly1qls"`,
				`"id":"0yxhwia2amd8gec"`,
			},
			notExpected: []string{
				`"expand":{}`,
			},
		},
		{
			name:           "[expand] guest (default expands)",
			auth:           nil,
			records:        freshRecords(usersRecords),
			queryExpand:    "",
			defaultExpands: []string{"rel"},
			expected: []string{
				`"customField":"123"`,
				`"expand":{"rel"`,
				`"id":"llvuca81nly1qls"`,
				`"id":"0yxhwia2amd8gec"`,
			},
		},
		{
			name:           "[expand] @request.context=expand check",
			auth:           nil,
			records:        freshRecords(demo5Records),
			queryExpand:    "rel_one",
			defaultExpands: []string{"rel_many"},
			expected: []string{
				`"customField":"123"`,
				`"expand":{}`,
				`"expand":{"`,
				`"rel_many":[{`,
				`"rel_one":{`,
				`"id":"i9naidtvr6qsgb4"`,
				`"id":"qzaqccwrmva4o1n"`,
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			app, _ := tests.NewTestApp()
			defer app.Cleanup()

			app.OnRecordEnrich().BindFunc(func(e *core.RecordEnrichEvent) error {
				e.Record.WithCustomData(true)
				e.Record.Set("customField", "123")
				return e.Next()
			})

			req := httptest.NewRequest(http.MethodGet, "/?expand="+s.queryExpand, nil)
			rec := httptest.NewRecorder()

			requestEvent := new(core.RequestEvent)
			requestEvent.App = app
			requestEvent.Request = req
			requestEvent.Response = rec
			requestEvent.Auth = s.auth

			err := apis.EnrichRecords(requestEvent, s.records, s.defaultExpands...)
			if err != nil {
				t.Fatal(err)
			}

			raw, err := json.Marshal(s.records)
			if err != nil {
				t.Fatal(err)
			}
			rawStr := string(raw)

			for _, str := range s.expected {
				if !strings.Contains(rawStr, str) {
					t.Fatalf("Expected\n%q\nin\n%v", str, rawStr)
				}
			}

			for _, str := range s.notExpected {
				if strings.Contains(rawStr, str) {
					t.Fatalf("Didn't expected\n%q\nin\n%v", str, rawStr)
				}
			}
		})
	}
}

func TestRecordAuthResponseAuthRuleCheck(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	event := new(core.RequestEvent)
	event.App = app
	event.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	event.Response = httptest.NewRecorder()

	user, err := app.FindAuthRecordByEmail("users", "test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		name        string
		rule        *string
		expectError bool
	}{
		{
			"admin only rule",
			nil,
			true,
		},
		{
			"empty rule",
			types.Pointer(""),
			false,
		},
		{
			"false rule",
			types.Pointer("1=2"),
			true,
		},
		{
			"true rule",
			types.Pointer("1=1"),
			false,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			user.Collection().AuthRule = s.rule

			err := apis.RecordAuthResponse(event, user, "", nil)

			hasErr := err != nil
			if s.expectError != hasErr {
				t.Fatalf("Expected hasErr %v, got %v (%v)", s.expectError, hasErr, err)
			}

			// in all cases login alert shouldn't be send because of the empty auth method
			if app.TestMailer.TotalSend() != 0 {
				t.Fatalf("Expected no emails send, got %d:\n%v", app.TestMailer.TotalSend(), app.TestMailer.LastMessage().HTML)
			}

			if !hasErr {
				return
			}

			apiErr, ok := err.(*router.ApiError)

			if !ok || apiErr == nil {
				t.Fatalf("Expected ApiError, got %v", apiErr)
			}

			if apiErr.Status != http.StatusForbidden {
				t.Fatalf("Expected ApiError.Status %d, got %d", http.StatusForbidden, apiErr.Status)
			}
		})
	}
}

func TestRecordAuthResponseAuthAlertCheck(t *testing.T) {
	const testFingerprint = "d0f88d6c87767262ba8e93d6acccd784"

	scenarios := []struct {
		name          string
		devices       []string // mock existing device fingerprints
		expectDevices []string
		enabled       bool
		expectEmail   bool
	}{
		{
			name:          "first login",
			devices:       nil,
			expectDevices: []string{testFingerprint},
			enabled:       true,
			expectEmail:   false,
		},
		{
			name:          "existing device",
			devices:       []string{"1", testFingerprint},
			expectDevices: []string{"1", testFingerprint},
			enabled:       true,
			expectEmail:   false,
		},
		{
			name:          "new device (< 5)",
			devices:       []string{"1", "2"},
			expectDevices: []string{"1", "2", testFingerprint},
			enabled:       true,
			expectEmail:   true,
		},
		{
			name:          "new device (>= 5)",
			devices:       []string{"1", "2", "3", "4", "5"},
			expectDevices: []string{"2", "3", "4", "5", testFingerprint},
			enabled:       true,
			expectEmail:   true,
		},
		{
			name:          "with disabled auth alert collection flag",
			devices:       []string{"1", "2"},
			expectDevices: []string{"1", "2"},
			enabled:       false,
			expectEmail:   false,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			app, _ := tests.NewTestApp()
			defer app.Cleanup()

			event := new(core.RequestEvent)
			event.App = app
			event.Request = httptest.NewRequest(http.MethodGet, "/", nil)
			event.Response = httptest.NewRecorder()

			user, err := app.FindAuthRecordByEmail("users", "test@example.com")
			if err != nil {
				t.Fatal(err)
			}

			user.Collection().MFA.Enabled = false
			user.Collection().AuthRule = types.Pointer("")
			user.Collection().AuthAlert.Enabled = s.enabled

			// ensure that there are no other auth origins
			err = app.DeleteAllAuthOriginsByRecord(user)
			if err != nil {
				t.Fatal(err)
			}

			mockCreated := types.NowDateTime().Add(-time.Duration(len(s.devices)+1) * time.Second)
			// insert the mock devices
			for _, fingerprint := range s.devices {
				mockCreated = mockCreated.Add(1 * time.Second)
				d := core.NewAuthOrigin(app)
				d.SetCollectionRef(user.Collection().Id)
				d.SetRecordRef(user.Id)
				d.SetFingerprint(fingerprint)
				d.SetRaw("created", mockCreated)
				d.SetRaw("updated", mockCreated)
				if err = app.Save(d); err != nil {
					t.Fatal(err)
				}
			}

			err = apis.RecordAuthResponse(event, user, "example", nil)
			if err != nil {
				t.Fatalf("Failed to resolve auth response: %v", err)
			}

			var expectTotalSend int
			if s.expectEmail {
				expectTotalSend = 1
			}
			if total := app.TestMailer.TotalSend(); total != expectTotalSend {
				t.Fatalf("Expected %d sent emails, got %d", expectTotalSend, total)
			}

			devices, err := app.FindAllAuthOriginsByRecord(user)
			if err != nil {
				t.Fatalf("Failed to retrieve auth origins: %v", err)
			}

			if len(devices) != len(s.expectDevices) {
				t.Fatalf("Expected %d devices, got %d", len(s.expectDevices), len(devices))
			}

			for _, fingerprint := range s.expectDevices {
				var exists bool
				fingerprints := make([]string, 0, len(devices))
				for _, d := range devices {
					if d.Fingerprint() == fingerprint {
						exists = true
						break
					}
					fingerprints = append(fingerprints, d.Fingerprint())
				}
				if !exists {
					t.Fatalf("Missing device with fingerprint %q:\n%v", fingerprint, fingerprints)
				}
			}
		})
	}
}

func TestRecordAuthResponseMFACheck(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	user, err := app.FindAuthRecordByEmail("users", "test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	user2, err := app.FindAuthRecordByEmail("users", "test2@example.com")
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	event := new(core.RequestEvent)
	event.App = app
	event.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	event.Response = rec

	resetMFAs := func(authRecord *core.Record) {
		// ensure that mfa is enabled
		user.Collection().MFA.Enabled = true
		user.Collection().MFA.Duration = 5
		user.Collection().MFA.Rule = ""

		mfas, err := app.FindAllMFAsByRecord(authRecord)
		if err != nil {
			t.Fatalf("Failed to retrieve mfas: %v", err)
		}
		for _, mfa := range mfas {
			if err := app.Delete(mfa); err != nil {
				t.Fatalf("Failed to delete mfa %q: %v", mfa.Id, err)
			}
		}

		// reset response
		rec = httptest.NewRecorder()
		event.Response = rec
	}

	totalMFAs := func(authRecord *core.Record) int {
		mfas, err := app.FindAllMFAsByRecord(authRecord)
		if err != nil {
			t.Fatalf("Failed to retrieve mfas: %v", err)
		}
		return len(mfas)
	}

	t.Run("no collection MFA enabled", func(t *testing.T) {
		resetMFAs(user)

		user.Collection().MFA.Enabled = false

		err = apis.RecordAuthResponse(event, user, "example", nil)
		if err != nil {
			t.Fatalf("Expected nil, got error: %v", err)
		}

		body := rec.Body.String()
		if strings.Contains(body, "mfaId") {
			t.Fatalf("Expected no mfaId in the response body, got\n%v", body)
		}
		if !strings.Contains(body, "token") {
			t.Fatalf("Expected auth token in the response body, got\n%v", body)
		}

		if total := totalMFAs(user); total != 0 {
			t.Fatalf("Expected no mfa records to be created, got %d", total)
		}
	})

	t.Run("no explicit auth method", func(t *testing.T) {
		resetMFAs(user)

		err = apis.RecordAuthResponse(event, user, "", nil)
		if err != nil {
			t.Fatalf("Expected nil, got error: %v", err)
		}

		body := rec.Body.String()
		if strings.Contains(body, "mfaId") {
			t.Fatalf("Expected no mfaId in the response body, got\n%v", body)
		}
		if !strings.Contains(body, "token") {
			t.Fatalf("Expected auth token in the response body, got\n%v", body)
		}

		if total := totalMFAs(user); total != 0 {
			t.Fatalf("Expected no mfa records to be created, got %d", total)
		}
	})

	t.Run("no mfa wanted (mfa rule check failure)", func(t *testing.T) {
		resetMFAs(user)
		user.Collection().MFA.Rule = "1=2"

		err = apis.RecordAuthResponse(event, user, "example", nil)
		if err != nil {
			t.Fatalf("Expected nil, got error: %v", err)
		}

		body := rec.Body.String()
		if strings.Contains(body, "mfaId") {
			t.Fatalf("Expected no mfaId in the response body, got\n%v", body)
		}
		if !strings.Contains(body, "token") {
			t.Fatalf("Expected auth token in the response body, got\n%v", body)
		}

		if total := totalMFAs(user); total != 0 {
			t.Fatalf("Expected no mfa records to be created, got %d", total)
		}
	})

	t.Run("mfa wanted (mfa rule check success)", func(t *testing.T) {
		resetMFAs(user)
		user.Collection().MFA.Rule = "1=1"

		err = apis.RecordAuthResponse(event, user, "example", nil)
		if !errors.Is(err, apis.ErrMFA) {
			t.Fatalf("Expected ErrMFA, got: %v", err)
		}

		body := rec.Body.String()
		if !strings.Contains(body, "mfaId") {
			t.Fatalf("Expected the created mfaId to be returned in the response body, got\n%v", body)
		}

		if total := totalMFAs(user); total != 1 {
			t.Fatalf("Expected a single mfa record to be created, got %d", total)
		}
	})

	t.Run("mfa first-time", func(t *testing.T) {
		resetMFAs(user)

		err = apis.RecordAuthResponse(event, user, "example", nil)
		if !errors.Is(err, apis.ErrMFA) {
			t.Fatalf("Expected ErrMFA, got: %v", err)
		}

		body := rec.Body.String()
		if !strings.Contains(body, "mfaId") {
			t.Fatalf("Expected the created mfaId to be returned in the response body, got\n%v", body)
		}

		if total := totalMFAs(user); total != 1 {
			t.Fatalf("Expected a single mfa record to be created, got %d", total)
		}
	})

	t.Run("mfa second-time with the same auth method", func(t *testing.T) {
		resetMFAs(user)

		// create a dummy mfa record
		mfa := core.NewMFA(app)
		mfa.SetCollectionRef(user.Collection().Id)
		mfa.SetRecordRef(user.Id)
		mfa.SetMethod("example")
		if err = app.Save(mfa); err != nil {
			t.Fatal(err)
		}

		event.Request = httptest.NewRequest(http.MethodGet, "/?mfaId="+mfa.Id, nil)

		err = apis.RecordAuthResponse(event, user, "example", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if total := totalMFAs(user); total != 1 {
			t.Fatalf("Expected only 1 mfa record (the existing one), got %d", total)
		}
	})

	t.Run("mfa second-time with the different auth method (query param)", func(t *testing.T) {
		resetMFAs(user)

		// create a dummy mfa record
		mfa := core.NewMFA(app)
		mfa.SetCollectionRef(user.Collection().Id)
		mfa.SetRecordRef(user.Id)
		mfa.SetMethod("example1")
		if err = app.Save(mfa); err != nil {
			t.Fatal(err)
		}

		event.Request = httptest.NewRequest(http.MethodGet, "/?mfaId="+mfa.Id, nil)

		err = apis.RecordAuthResponse(event, user, "example2", nil)
		if err != nil {
			t.Fatalf("Expected nil, got error: %v", err)
		}

		if total := totalMFAs(user); total != 0 {
			t.Fatalf("Expected the dummy mfa record to be deleted, found %d", total)
		}
	})

	t.Run("mfa second-time with the different auth method (body param)", func(t *testing.T) {
		resetMFAs(user)

		// create a dummy mfa record
		mfa := core.NewMFA(app)
		mfa.SetCollectionRef(user.Collection().Id)
		mfa.SetRecordRef(user.Id)
		mfa.SetMethod("example1")
		if err = app.Save(mfa); err != nil {
			t.Fatal(err)
		}

		event.Request = httptest.NewRequest(http.MethodGet, "/", strings.NewReader(`{"mfaId":"`+mfa.Id+`"}`))
		event.Request.Header.Add("content-type", "application/json")

		err = apis.RecordAuthResponse(event, user, "example2", nil)
		if err != nil {
			t.Fatalf("Expected nil, got error: %v", err)
		}

		if total := totalMFAs(user); total != 0 {
			t.Fatalf("Expected the dummy mfa record to be deleted, found %d", total)
		}
	})

	t.Run("missing mfa", func(t *testing.T) {
		resetMFAs(user)

		event.Request = httptest.NewRequest(http.MethodGet, "/?mfaId=missing", nil)

		err = apis.RecordAuthResponse(event, user, "example2", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if total := totalMFAs(user); total != 0 {
			t.Fatalf("Expected 0 mfa records, got %d", total)
		}
	})

	t.Run("expired mfa", func(t *testing.T) {
		resetMFAs(user)

		// create a dummy expired mfa record
		mfa := core.NewMFA(app)
		mfa.SetCollectionRef(user.Collection().Id)
		mfa.SetRecordRef(user.Id)
		mfa.SetMethod("example1")
		mfa.SetRaw("created", types.NowDateTime().Add(-1*time.Hour))
		mfa.SetRaw("updated", types.NowDateTime().Add(-1*time.Hour))
		if err = app.Save(mfa); err != nil {
			t.Fatal(err)
		}

		event.Request = httptest.NewRequest(http.MethodGet, "/?mfaId="+mfa.Id, nil)

		err = apis.RecordAuthResponse(event, user, "example2", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if totalMFAs(user) != 0 {
			t.Fatal("Expected the expired mfa record to be deleted")
		}
	})

	t.Run("mfa for different auth record", func(t *testing.T) {
		resetMFAs(user)

		// create a dummy expired mfa record
		mfa := core.NewMFA(app)
		mfa.SetCollectionRef(user2.Collection().Id)
		mfa.SetRecordRef(user2.Id)
		mfa.SetMethod("example1")
		if err = app.Save(mfa); err != nil {
			t.Fatal(err)
		}

		event.Request = httptest.NewRequest(http.MethodGet, "/?mfaId="+mfa.Id, nil)

		err = apis.RecordAuthResponse(event, user, "example2", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if total := totalMFAs(user); total != 0 {
			t.Fatalf("Expected no user mfas, got %d", total)
		}

		if total := totalMFAs(user2); total != 1 {
			t.Fatalf("Expected only 1 user2 mfa, got %d", total)
		}
	})
}
