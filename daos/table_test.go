package daos_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/list"
)

func TestHasTable(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		tableName string
		expected  bool
	}{
		{"", false},
		{"test", false},
		{"_admins", true},
		{"demo3", true},
		{"DEMO3", true}, // table names are case insensitives by default
	}

	for i, scenario := range scenarios {
		result := app.Dao().HasTable(scenario.tableName)
		if result != scenario.expected {
			t.Errorf("(%d) Expected %v, got %v", i, scenario.expected, result)
		}
	}
}

func TestGetTableColumns(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		tableName string
		expected  []string
	}{
		{"", nil},
		{"_params", []string{"id", "key", "value", "created", "updated"}},
	}

	for i, scenario := range scenarios {
		columns, _ := app.Dao().GetTableColumns(scenario.tableName)

		if len(columns) != len(scenario.expected) {
			t.Errorf("(%d) Expected columns %v, got %v", i, scenario.expected, columns)
		}

		for _, c := range columns {
			if !list.ExistInSlice(c, scenario.expected) {
				t.Errorf("(%d) Didn't expect column %s", i, c)
			}
		}
	}
}

func TestDeleteTable(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		tableName   string
		expectError bool
	}{
		{"", true},
		{"test", true},
		{"_admins", false},
		{"demo3", false},
	}

	for i, scenario := range scenarios {
		err := app.Dao().DeleteTable(scenario.tableName)
		hasErr := err != nil
		if hasErr != scenario.expectError {
			t.Errorf("(%d) Expected hasErr %v, got %v", i, scenario.expectError, hasErr)
		}
	}
}

func TestVacuum(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	calledQueries := []string{}
	app.DB().QueryLogFunc = func(ctx context.Context, t time.Duration, sql string, rows *sql.Rows, err error) {
		calledQueries = append(calledQueries, sql)
	}
	app.DB().ExecLogFunc = func(ctx context.Context, t time.Duration, sql string, result sql.Result, err error) {
		calledQueries = append(calledQueries, sql)
	}

	if err := app.Dao().Vacuum(); err != nil {
		t.Fatal(err)
	}

	if total := len(calledQueries); total != 1 {
		t.Fatalf("Expected 1 query, got %d", total)
	}

	if calledQueries[0] != "VACUUM" {
		t.Fatalf("Expected VACUUM query, got %s", calledQueries[0])
	}
}
