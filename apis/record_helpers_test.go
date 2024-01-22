package apis_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tests"
)

func TestRequestInfo(t *testing.T) {
	t.Parallel()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/?test=123", strings.NewReader(`{"test":456}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("X-Token-Test", "123")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	dummyRecord := &models.Record{}
	dummyRecord.Id = "id1"
	c.Set(apis.ContextAuthRecordKey, dummyRecord)

	dummyAdmin := &models.Admin{}
	dummyAdmin.Id = "id2"
	c.Set(apis.ContextAdminKey, dummyAdmin)

	result := apis.RequestInfo(c)

	if result == nil {
		t.Fatal("Expected *models.RequestInfo instance, got nil")
	}

	if result.Method != http.MethodPost {
		t.Fatalf("Expected Method %v, got %v", http.MethodPost, result.Method)
	}

	rawHeaders, _ := json.Marshal(result.Headers)
	expectedHeaders := `{"content_type":"application/json","x_token_test":"123"}`
	if v := string(rawHeaders); v != expectedHeaders {
		t.Fatalf("Expected Query %v, got %v", expectedHeaders, v)
	}

	rawQuery, _ := json.Marshal(result.Query)
	expectedQuery := `{"test":"123"}`
	if v := string(rawQuery); v != expectedQuery {
		t.Fatalf("Expected Query %v, got %v", expectedQuery, v)
	}

	rawData, _ := json.Marshal(result.Data)
	expectedData := `{"test":456}`
	if v := string(rawData); v != expectedData {
		t.Fatalf("Expected Data %v, got %v", expectedData, v)
	}

	if result.AuthRecord == nil || result.AuthRecord.Id != dummyRecord.Id {
		t.Fatalf("Expected AuthRecord %v, got %v", dummyRecord, result.AuthRecord)
	}

	if result.Admin == nil || result.Admin.Id != dummyAdmin.Id {
		t.Fatalf("Expected Admin %v, got %v", dummyAdmin, result.Admin)
	}
}

func TestRecordAuthResponse(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	dummyAdmin := &models.Admin{}
	dummyAdmin.Id = "id1"

	nonAuthRecord, err := app.Dao().FindRecordById("demo1", "al1h9ijdeojtsjy")
	if err != nil {
		t.Fatal(err)
	}

	authRecord, err := app.Dao().FindRecordById("users", "4q1xlclmfloku33")
	if err != nil {
		t.Fatal(err)
	}

	unverifiedAuthRecord, err := app.Dao().FindRecordById("clients", "o1y0dd0spd786md")
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		name               string
		record             *models.Record
		meta               any
		expectError        bool
		expectedContent    []string
		notExpectedContent []string
		expectedEvents     map[string]int
	}{
		{
			name:        "non auth record",
			record:      nonAuthRecord,
			expectError: true,
		},
		{
			name:        "valid auth record but with unverified email in onlyVerified collection",
			record:      unverifiedAuthRecord,
			expectError: true,
		},
		{
			name:        "valid auth record - without meta",
			record:      authRecord,
			expectError: false,
			expectedContent: []string{
				`"token":"`,
				`"record":{`,
				`"id":"`,
				`"expand":{"rel":{`,
			},
			notExpectedContent: []string{
				`"meta":`,
			},
			expectedEvents: map[string]int{
				"OnRecordAuthRequest": 1,
			},
		},
		{
			name:        "valid auth record - with meta",
			record:      authRecord,
			meta:        map[string]any{"meta_test": 123},
			expectError: false,
			expectedContent: []string{
				`"token":"`,
				`"record":{`,
				`"id":"`,
				`"expand":{"rel":{`,
				`"meta":{"meta_test":123`,
			},
			expectedEvents: map[string]int{
				"OnRecordAuthRequest": 1,
			},
		},
	}

	for _, s := range scenarios {
		app.ResetEventCalls()

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?expand=rel", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set(apis.ContextAdminKey, dummyAdmin)

		responseErr := apis.RecordAuthResponse(app, c, s.record, s.meta)

		hasErr := responseErr != nil
		if hasErr != s.expectError {
			t.Fatalf("[%s] Expected hasErr to be %v, got %v (%v)", s.name, s.expectError, hasErr, responseErr)
		}

		if len(app.EventCalls) != len(s.expectedEvents) {
			t.Fatalf("[%s] Expected events \n%v, \ngot \n%v", s.name, s.expectedEvents, app.EventCalls)
		}
		for k, v := range s.expectedEvents {
			if app.EventCalls[k] != v {
				t.Fatalf("[%s] Expected event %s to be called %d times, got %d", s.name, k, v, app.EventCalls[k])
			}
		}

		if hasErr {
			continue
		}

		response := rec.Body.String()

		for _, v := range s.expectedContent {
			if !strings.Contains(response, v) {
				t.Fatalf("[%s] Missing %v in response \n%v", s.name, v, response)
			}
		}

		for _, v := range s.notExpectedContent {
			if strings.Contains(response, v) {
				t.Fatalf("[%s] Unexpected %v in response \n%v", s.name, v, response)
			}
		}
	}
}

func TestEnrichRecords(t *testing.T) {
	t.Parallel()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/?expand=rel_many", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	dummyAdmin := &models.Admin{}
	dummyAdmin.Id = "test_id"
	c.Set(apis.ContextAdminKey, dummyAdmin)

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	records, err := app.Dao().FindRecordsByIds("demo1", []string{"al1h9ijdeojtsjy", "84nmscqy84lsi1t"})
	if err != nil {
		t.Fatal(err)
	}

	apis.EnrichRecords(c, app.Dao(), records, "rel_one")

	for _, record := range records {
		expand := record.Expand()
		if len(expand) == 0 {
			t.Fatalf("Expected non-empty expand, got nil for record %v", record)
		}

		if len(record.GetStringSlice("rel_one")) != 0 {
			if _, ok := expand["rel_one"]; !ok {
				t.Fatalf("Expected rel_one to be expanded for record %v, got \n%v", record, expand)
			}
		}

		if len(record.GetStringSlice("rel_many")) != 0 {
			if _, ok := expand["rel_many"]; !ok {
				t.Fatalf("Expected rel_many to be expanded for record %v, got \n%v", record, expand)
			}
		}
	}
}
