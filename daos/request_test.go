package daos_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestRequestQuery(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	expected := "SELECT {{_requests}}.* FROM `_requests`"

	sql := app.Dao().RequestQuery().Build().SQL()
	if sql != expected {
		t.Errorf("Expected sql %s, got %s", expected, sql)
	}
}

func TestFindRequestById(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	tests.MockRequestLogsData(app)

	scenarios := []struct {
		id          string
		expectError bool
	}{
		{"", true},
		{"invalid", true},
		{"00000000-9f38-44fb-bf82-c8f53b310d91", true},
		{"873f2133-9f38-44fb-bf82-c8f53b310d91", false},
	}

	for i, scenario := range scenarios {
		admin, err := app.LogsDao().FindRequestById(scenario.id)

		hasErr := err != nil
		if hasErr != scenario.expectError {
			t.Errorf("(%d) Expected hasErr to be %v, got %v (%v)", i, scenario.expectError, hasErr, err)
		}

		if admin != nil && admin.Id != scenario.id {
			t.Errorf("(%d) Expected admin with id %s, got %s", i, scenario.id, admin.Id)
		}
	}
}

func TestRequestsStats(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	tests.MockRequestLogsData(app)

	expected := `[{"total":1,"date":"2022-05-01 10:00:00.000Z"},{"total":1,"date":"2022-05-02 10:00:00.000Z"}]`

	now := time.Now().UTC().Format(types.DefaultDateLayout)
	exp := dbx.NewExp("[[created]] <= {:date}", dbx.Params{"date": now})
	result, err := app.LogsDao().RequestsStats(exp)
	if err != nil {
		t.Fatal(err)
	}

	encoded, _ := json.Marshal(result)
	if string(encoded) != expected {
		t.Fatalf("Expected %s, got %s", expected, string(encoded))
	}
}

func TestDeleteOldRequests(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	tests.MockRequestLogsData(app)

	scenarios := []struct {
		date          string
		expectedTotal int
	}{
		{"2022-01-01 10:00:00.000Z", 2}, // no requests to delete before that time
		{"2022-05-01 11:00:00.000Z", 1}, // only 1 request should have left
		{"2022-05-03 11:00:00.000Z", 0}, // no more requests should have left
		{"2022-05-04 11:00:00.000Z", 0}, // no more requests should have left
	}

	for i, scenario := range scenarios {
		date, dateErr := time.Parse(types.DefaultDateLayout, scenario.date)
		if dateErr != nil {
			t.Errorf("(%d) Date error %v", i, dateErr)
		}

		deleteErr := app.LogsDao().DeleteOldRequests(date)
		if deleteErr != nil {
			t.Errorf("(%d) Delete error %v", i, deleteErr)
		}

		// check total remaining requests
		var total int
		countErr := app.LogsDao().RequestQuery().Select("count(*)").Row(&total)
		if countErr != nil {
			t.Errorf("(%d) Count error %v", i, countErr)
		}

		if total != scenario.expectedTotal {
			t.Errorf("(%d) Expected %d remaining requests, got %d", i, scenario.expectedTotal, total)
		}
	}
}

func TestSaveRequest(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	tests.MockRequestLogsData(app)

	// create new request
	newRequest := &models.Request{}
	newRequest.Method = "get"
	newRequest.Meta = types.JsonMap{}
	createErr := app.LogsDao().SaveRequest(newRequest)
	if createErr != nil {
		t.Fatal(createErr)
	}

	// check if it was really created
	existingRequest, fetchErr := app.LogsDao().FindRequestById(newRequest.Id)
	if fetchErr != nil {
		t.Fatal(fetchErr)
	}

	existingRequest.Method = "post"
	updateErr := app.LogsDao().SaveRequest(existingRequest)
	if updateErr != nil {
		t.Fatal(updateErr)
	}
	// refresh instance to check if it was really updated
	existingRequest, _ = app.LogsDao().FindRequestById(existingRequest.Id)
	if existingRequest.Method != "post" {
		t.Fatalf("Expected request method to be %s, got %s", "post", existingRequest.Method)
	}
}
