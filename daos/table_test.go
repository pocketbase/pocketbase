package daos_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/list"
)

func TestHasTable(t *testing.T) {
	t.Parallel()

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
		{"view1", true}, // view
	}

	for i, scenario := range scenarios {
		result := app.Dao().HasTable(scenario.tableName)
		if result != scenario.expected {
			t.Errorf("[%d] Expected %v, got %v", i, scenario.expected, result)
		}
	}
}

func TestTableColumns(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		tableName string
		expected  []string
	}{
		{"", nil},
		{"_params", []string{"id", "key", "value", "created", "updated"}},
	}

	for i, s := range scenarios {
		columns, _ := app.Dao().TableColumns(s.tableName)

		if len(columns) != len(s.expected) {
			t.Errorf("[%d] Expected columns %v, got %v", i, s.expected, columns)
			continue
		}

		for _, c := range columns {
			if !list.ExistInSlice(c, s.expected) {
				t.Errorf("[%d] Didn't expect column %s", i, c)
			}
		}
	}
}

func TestTableInfo(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		tableName string
		expected  string
	}{
		{"", "null"},
		{"missing", "null"},
		{
			"_admins",
			`[{"PK":1,"Index":0,"Name":"id","Type":"TEXT","NotNull":false,"DefaultValue":null},{"PK":0,"Index":1,"Name":"avatar","Type":"INTEGER","NotNull":true,"DefaultValue":0},{"PK":0,"Index":2,"Name":"email","Type":"TEXT","NotNull":true,"DefaultValue":null},{"PK":0,"Index":3,"Name":"tokenKey","Type":"TEXT","NotNull":true,"DefaultValue":null},{"PK":0,"Index":4,"Name":"passwordHash","Type":"TEXT","NotNull":true,"DefaultValue":null},{"PK":0,"Index":5,"Name":"lastResetSentAt","Type":"TEXT","NotNull":true,"DefaultValue":""},{"PK":0,"Index":6,"Name":"created","Type":"TEXT","NotNull":true,"DefaultValue":""},{"PK":0,"Index":7,"Name":"updated","Type":"TEXT","NotNull":true,"DefaultValue":""}]`,
		},
	}

	for i, s := range scenarios {
		rows, _ := app.Dao().TableInfo(s.tableName)

		raw, _ := json.Marshal(rows)
		rawStr := string(raw)

		if rawStr != s.expected {
			t.Errorf("[%d] Expected \n%v, \ngot \n%v", i, s.expected, rawStr)
		}
	}
}

func TestDeleteTable(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		tableName   string
		expectError bool
	}{
		{"", true},
		{"test", false}, // missing tables are ignored
		{"_admins", false},
		{"demo3", false},
	}

	for i, s := range scenarios {
		err := app.Dao().DeleteTable(s.tableName)

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("[%d] Expected hasErr %v, got %v", i, s.expectError, hasErr)
		}
	}
}

func TestVacuum(t *testing.T) {
	t.Parallel()

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

func TestTableIndexes(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		table         string
		expectError   bool
		expectIndexes []string
	}{
		{
			"missing",
			false,
			nil,
		},
		{
			"demo2",
			false,
			[]string{"idx_demo2_created", "idx_unique_demo2_title", "idx_demo2_active"},
		},
	}

	for _, s := range scenarios {
		result, err := app.Dao().TableIndexes(s.table)

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("[%s] Expected hasErr %v, got %v", s.table, s.expectError, hasErr)
		}

		if len(s.expectIndexes) != len(result) {
			t.Errorf("[%s] Expected %d indexes, got %d:\n%v", s.table, len(s.expectIndexes), len(result), result)
			continue
		}

		for _, name := range s.expectIndexes {
			if result[name] == "" {
				t.Errorf("[%s] Missing index %q in \n%v", s.table, name, result)
			}
		}
	}
}
