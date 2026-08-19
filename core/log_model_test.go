package core_test

import (
	"bytes"
	"encoding/json/v2"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestLogTableName(t *testing.T) {
	t.Parallel()

	var log core.Log

	if name := log.TableName(); name != core.LogsTableName {
		t.Fatalf("Expected Log table name %q, got %q", core.LogsTableName, name)
	}
}

func TestLogDBExport(t *testing.T) {
	t.Parallel()

	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	date, err := types.ParseDateTime("2026-08-18 10:20:30.456Z")
	if err != nil {
		t.Fatal(err)
	}

	messageLimit := 8000
	dataLimit := 16 << 10

	scenarios := []struct {
		name       string
		log        core.Log
		limit      int64 // 0 -> use default
		expectJSON string
	}{
		{
			"empty log",
			core.Log{},
			0,
			`{"created":"","data":{},"id":"","level":0,"message":""}`,
		},
		{
			"with message and data below the default limits",
			core.Log{
				BaseModel: core.BaseModel{Id: "test_id"},
				Created:   date,
				Level:     123,
				Message:   "test_message",
				Data:      types.JSONMap[any]{"a": "test1", "b": "test2"},
			},
			0,
			`{"created":"2026-08-18 10:20:30.456Z","data":{"a":"test1","b":"test2"},"id":"test_id","level":123,"message":"test_message"}`,
		},
		{
			"with message and data exactly the default limits",
			core.Log{
				BaseModel: core.BaseModel{Id: "test_id"},
				Created:   date,
				Level:     123,
				Message:   strings.Repeat("a", messageLimit),
				Data:      types.JSONMap[any]{"a": "test1", "b": "test2", "c": strings.Repeat("a", dataLimit-32)},
			},
			0,
			`{"created":"2026-08-18 10:20:30.456Z","data":{"a":"test1","b":"test2","c":"` + strings.Repeat("a", dataLimit-32) + `"},"id":"test_id","level":123,"message":"` + strings.Repeat("a", messageLimit) + `"}`,
		},
		{
			"with message and data above the default limits",
			core.Log{
				BaseModel: core.BaseModel{Id: "test_id"},
				Created:   date,
				Level:     123,
				Message:   strings.Repeat("a", messageLimit) + "x",                                                      // "x" should be omitted
				Data:      types.JSONMap[any]{"a": "test1", "b": "test2", "c": strings.Repeat("a", dataLimit-32) + "x"}, // the end will be incomplete and something like `"c":"...aaaaaax`
			},
			0,
			`{"created":"2026-08-18 10:20:30.456Z","data":{"__pb_truncated__":true,"a":"test1","b":"test2","c":"` + strings.Repeat("a", dataLimit-32) + `x"},"id":"test_id","level":123,"message":"` + strings.Repeat("a", messageLimit) + `"}`,
		},
		{
			"with data above custom limit",
			core.Log{
				BaseModel: core.BaseModel{Id: "test_id"},
				Created:   date,
				Level:     123,
				Message:   "test_message",
				Data:      types.JSONMap[any]{"a": "test1", "b": "test2", "c": strings.Repeat("a", (2<<10)-32) + "x"}, // the end will be incomplete and something like `"c":"...aaaaaax`
			},
			2 << 10,
			`{"created":"2026-08-18 10:20:30.456Z","data":{"__pb_truncated__":true,"a":"test1","b":"test2","c":"` + strings.Repeat("a", (2<<10)-32) + `x"},"id":"test_id","level":123,"message":"test_message"}`,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			testApp.Settings().Logs.MaxDataSize = s.limit

			result, err := s.log.DBExport(testApp)
			if err != nil {
				t.Fatal(err)
			}

			raw, err := json.Marshal(result, json.Deterministic(true))
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(raw, []byte(s.expectJSON)) {
				t.Fatalf("Expected export data\n%s\ngot\n%s", s.expectJSON, raw)
			}
		})
	}
}
