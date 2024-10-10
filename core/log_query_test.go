package core_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestFindLogById(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	tests.StubLogsData(app)

	scenarios := []struct {
		id          string
		expectError bool
	}{
		{"", true},
		{"invalid", true},
		{"00000000-9f38-44fb-bf82-c8f53b310d91", true},
		{"873f2133-9f38-44fb-bf82-c8f53b310d91", false},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%s", i, s.id), func(t *testing.T) {
			log, err := app.FindLogById(s.id)

			hasErr := err != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr to be %v, got %v (%v)", s.expectError, hasErr, err)
			}

			if log != nil && log.Id != s.id {
				t.Fatalf("Expected log with id %q, got %q", s.id, log.Id)
			}
		})
	}
}

func TestLogsStats(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	tests.StubLogsData(app)

	expected := `[{"date":"2022-05-01 10:00:00.000Z","total":1},{"date":"2022-05-02 10:00:00.000Z","total":1}]`

	now := time.Now().UTC().Format(types.DefaultDateLayout)
	exp := dbx.NewExp("[[created]] <= {:date}", dbx.Params{"date": now})
	result, err := app.LogsStats(exp)
	if err != nil {
		t.Fatal(err)
	}

	encoded, _ := json.Marshal(result)
	if string(encoded) != expected {
		t.Fatalf("Expected\n%q\ngot\n%q", expected, string(encoded))
	}
}

func TestDeleteOldLogs(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	tests.StubLogsData(app)

	scenarios := []struct {
		date          string
		expectedTotal int
	}{
		{"2022-01-01 10:00:00.000Z", 2}, // no logs to delete before that time
		{"2022-05-01 11:00:00.000Z", 1}, // only 1 log should have left
		{"2022-05-03 11:00:00.000Z", 0}, // no more logs should have left
		{"2022-05-04 11:00:00.000Z", 0}, // no more logs should have left
	}

	for _, s := range scenarios {
		t.Run(s.date, func(t *testing.T) {
			date, dateErr := time.Parse(types.DefaultDateLayout, s.date)
			if dateErr != nil {
				t.Fatalf("Date error %v", dateErr)
			}

			deleteErr := app.DeleteOldLogs(date)
			if deleteErr != nil {
				t.Fatalf("Delete error %v", deleteErr)
			}

			// check total remaining logs
			var total int
			countErr := app.AuxModelQuery(&core.Log{}).Select("count(*)").Row(&total)
			if countErr != nil {
				t.Errorf("Count error %v", countErr)
			}

			if total != s.expectedTotal {
				t.Errorf("Expected %d remaining logs, got %d", s.expectedTotal, total)
			}
		})
	}
}
