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

func TestLogQuery(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	expected := "SELECT {{_logs}}.* FROM `_logs`"

	sql := app.Dao().LogQuery().Build().SQL()
	if sql != expected {
		t.Errorf("Expected sql %s, got %s", expected, sql)
	}
}

func TestFindLogById(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	tests.MockLogsData(app)

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
		admin, err := app.LogsDao().FindLogById(scenario.id)

		hasErr := err != nil
		if hasErr != scenario.expectError {
			t.Errorf("(%d) Expected hasErr to be %v, got %v (%v)", i, scenario.expectError, hasErr, err)
		}

		if admin != nil && admin.Id != scenario.id {
			t.Errorf("(%d) Expected admin with id %s, got %s", i, scenario.id, admin.Id)
		}
	}
}

func TestLogsStats(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	tests.MockLogsData(app)

	expected := `[{"total":1,"date":"2022-05-01 10:00:00.000Z"},{"total":1,"date":"2022-05-02 10:00:00.000Z"}]`

	now := time.Now().UTC().Format(types.DefaultDateLayout)
	exp := dbx.NewExp("[[created]] <= {:date}", dbx.Params{"date": now})
	result, err := app.LogsDao().LogsStats(exp)
	if err != nil {
		t.Fatal(err)
	}

	encoded, _ := json.Marshal(result)
	if string(encoded) != expected {
		t.Fatalf("Expected %s, got %s", expected, string(encoded))
	}
}

func TestDeleteOldLogs(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	tests.MockLogsData(app)

	scenarios := []struct {
		date          string
		expectedTotal int
	}{
		{"2022-01-01 10:00:00.000Z", 2}, // no logs to delete before that time
		{"2022-05-01 11:00:00.000Z", 1}, // only 1 log should have left
		{"2022-05-03 11:00:00.000Z", 0}, // no more logs should have left
		{"2022-05-04 11:00:00.000Z", 0}, // no more logs should have left
	}

	for i, scenario := range scenarios {
		date, dateErr := time.Parse(types.DefaultDateLayout, scenario.date)
		if dateErr != nil {
			t.Errorf("(%d) Date error %v", i, dateErr)
		}

		deleteErr := app.LogsDao().DeleteOldLogs(date)
		if deleteErr != nil {
			t.Errorf("(%d) Delete error %v", i, deleteErr)
		}

		// check total remaining logs
		var total int
		countErr := app.LogsDao().LogQuery().Select("count(*)").Row(&total)
		if countErr != nil {
			t.Errorf("(%d) Count error %v", i, countErr)
		}

		if total != scenario.expectedTotal {
			t.Errorf("(%d) Expected %d remaining logs, got %d", i, scenario.expectedTotal, total)
		}
	}
}

func TestSaveLog(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	tests.MockLogsData(app)

	// create new log
	newLog := &models.Log{}
	newLog.Level = -4
	newLog.Data = types.JsonMap{}
	createErr := app.LogsDao().SaveLog(newLog)
	if createErr != nil {
		t.Fatal(createErr)
	}

	// check if it was really created
	existingLog, fetchErr := app.LogsDao().FindLogById(newLog.Id)
	if fetchErr != nil {
		t.Fatal(fetchErr)
	}

	existingLog.Level = 4
	updateErr := app.LogsDao().SaveLog(existingLog)
	if updateErr != nil {
		t.Fatal(updateErr)
	}
	// refresh instance to check if it was really updated
	existingLog, _ = app.LogsDao().FindLogById(existingLog.Id)
	if existingLog.Level != 4 {
		t.Fatalf("Expected log level to be %d, got %d", 4, existingLog.Level)
	}
}
