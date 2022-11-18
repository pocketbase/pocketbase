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

func TestRequestData(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/?test=123", strings.NewReader(`{"test":456}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	dummyRecord := &models.Record{}
	dummyRecord.Id = "id1"
	c.Set(apis.ContextAuthRecordKey, dummyRecord)

	dummyAdmin := &models.Admin{}
	dummyAdmin.Id = "id2"
	c.Set(apis.ContextAdminKey, dummyAdmin)

	result := apis.RequestData(c)

	if result == nil {
		t.Fatal("Expected *models.RequestData instance, got nil")
	}

	if result.Method != http.MethodPost {
		t.Fatalf("Expected Method %v, got %v", http.MethodPost, result.Method)
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

func TestEnrichRecords(t *testing.T) {
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
