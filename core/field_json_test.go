package core_test

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/search"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestJSONFieldBaseMethods(t *testing.T) {
	testFieldBaseMethods(t, core.FieldTypeJSON)
}

func TestJSONFieldColumnType(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	f := &core.JSONField{}

	/* SQLite:
	expected := "JSON DEFAULT NULL"
	*/
	// PostgreSQL:
	expected := "JSONB DEFAULT NULL"

	if v := f.ColumnType(app); v != expected {
		t.Fatalf("Expected\n%q\ngot\n%q", expected, v)
	}
}

func TestJSONFieldPrepareValue(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	f := &core.JSONField{}
	record := core.NewRecord(core.NewBaseCollection("test"))

	scenarios := []struct {
		raw      any
		expected string
	}{
		{"null", `null`},
		{"", `""`},
		{"true", `true`},
		{"false", `false`},
		{"test", `"test"`},
		{"123", `123`},
		{"-456", `-456`},
		{"[1,2,3]", `[1,2,3]`},
		{"[1,2,3", `"[1,2,3"`},
		{`{"a":1,"b":2}`, `{"a":1,"b":2}`},
		{`{"a":1,"b":2`, `"{\"a\":1,\"b\":2"`},
		{[]int{1, 2, 3}, `[1,2,3]`},
		{map[string]int{"a": 1, "b": 2}, `{"a":1,"b":2}`},
		{nil, `null`},
		{false, `false`},
		{true, `true`},
		{-78, `-78`},
		{123.456, `123.456`},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#v", i, s.raw), func(t *testing.T) {
			v, err := f.PrepareValue(record, s.raw)
			if err != nil {
				t.Fatal(err)
			}

			raw, ok := v.(types.JSONRaw)
			if !ok {
				t.Fatalf("Expected string instance, got %T", v)
			}
			rawStr := raw.String()

			if rawStr != s.expected {
				t.Fatalf("Expected\n%#v\ngot\n%#v", s.expected, rawStr)
			}
		})
	}
}

func TestJSONFieldValidateValue(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewBaseCollection("test_collection")

	scenarios := []struct {
		name        string
		field       *core.JSONField
		record      func() *core.Record
		expectError bool
	}{
		{
			"invalid raw value",
			&core.JSONField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", 123)
				return record
			},
			true,
		},
		{
			"zero field value (not required)",
			&core.JSONField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", types.JSONRaw{})
				return record
			},
			false,
		},
		{
			"zero field value (required)",
			&core.JSONField{Name: "test", Required: true},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", types.JSONRaw{})
				return record
			},
			true,
		},
		{
			"non-zero field value (required)",
			&core.JSONField{Name: "test", Required: true},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", types.JSONRaw("[1,2,3]"))
				return record
			},
			false,
		},
		{
			"non-zero field value (required)",
			&core.JSONField{Name: "test", Required: true},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", types.JSONRaw(`"aaa"`))
				return record
			},
			false,
		},
		{
			"> default MaxSize",
			&core.JSONField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", types.JSONRaw(`"`+strings.Repeat("a", (1<<20))+`"`))
				return record
			},
			true,
		},
		{
			"> MaxSize",
			&core.JSONField{Name: "test", MaxSize: 5},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", types.JSONRaw(`"aaaa"`))
				return record
			},
			true,
		},
		{
			"<= MaxSize",
			&core.JSONField{Name: "test", MaxSize: 5},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", types.JSONRaw(`"aaa"`))
				return record
			},
			false,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			err := s.field.ValidateValue(context.Background(), app, s.record())

			hasErr := err != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr %v, got %v (%v)", s.expectError, hasErr, err)
			}
		})
	}
}

func TestJSONFieldValidateSettings(t *testing.T) {
	testDefaultFieldIdValidation(t, core.FieldTypeJSON)
	testDefaultFieldNameValidation(t, core.FieldTypeJSON)

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewBaseCollection("test_collection")

	scenarios := []struct {
		name         string
		field        func() *core.JSONField
		expectErrors []string
	}{
		{
			"MaxSize < 0",
			func() *core.JSONField {
				return &core.JSONField{
					Id:      "test",
					Name:    "test",
					MaxSize: -1,
				}
			},
			[]string{"maxSize"},
		},
		{
			"MaxSize = 0",
			func() *core.JSONField {
				return &core.JSONField{
					Id:   "test",
					Name: "test",
				}
			},
			[]string{},
		},
		{
			"MaxSize > 0",
			func() *core.JSONField {
				return &core.JSONField{
					Id:      "test",
					Name:    "test",
					MaxSize: 1,
				}
			},
			[]string{},
		},
		{
			"MaxSize > safe json int",
			func() *core.JSONField {
				return &core.JSONField{
					Id:      "test",
					Name:    "test",
					MaxSize: 1 << 53,
				}
			},
			[]string{"maxSize"},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			errs := s.field().ValidateSettings(context.Background(), app, collection)

			tests.TestValidationErrors(t, errs, s.expectErrors)
		})
	}
}

func TestJSONFieldCalculateMaxBodySize(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	scenarios := []struct {
		field    *core.JSONField
		expected int64
	}{
		{&core.JSONField{}, core.DefaultJSONFieldMaxSize},
		{&core.JSONField{MaxSize: 10}, 10},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%d", i, s.field.MaxSize), func(t *testing.T) {
			result := s.field.CalculateMaxBodySize()

			if result != s.expected {
				t.Fatalf("Expected %d, got %d", s.expected, result)
			}
		})
	}
}

func TestJsonFilter_TypeCasting(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	jsonField := &core.JSONField{Name: "json_col"}
	boolField := &core.BoolField{Name: "bool_col"}
	collection := core.NewBaseCollection("test")
	collection.Fields.Add(jsonField, boolField)
	resolver := core.NewRecordFieldResolver(testApp, collection, nil, false)

	// Append 5 rows data to the db
	sql := `
		CREATE TABLE test (
			id		 INTEGER,
			text_col TEXT,
			bool_col BOOLEAN,
			json_col JSON
		);
		INSERT INTO test (id, text_col, bool_col, json_col) VALUES
			(0, 'text a', true, '{"string_key": "value a", "int_key": 0, "bool_key": true}'),
			(1, 'test b', false, '{"string_key": "value b", "int_key": 1, "bool_key": false}'),
			(2, 'test c', true, '{"string_key": "value c", "int_key": 2, "bool_key": true}'),
			(3, 'test d', false, '{"string_key": "value d", "int_key": 3, "bool_key": false}'),
			(4, '123', true, '{"string_key": "value e", "int_key": 4, "bool_key": true}')
	`
	_, err := testApp.DB().NewQuery(sql).Execute()
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		name           string
		filterData     search.FilterData
		expectedRowIds []int
	}{
		{
			name:           "bool filters",
			filterData:     search.FilterData(`bool_col = true`),
			expectedRowIds: []int{0, 2, 4},
		},
		{
			name:           "json filter type cast to string",
			filterData:     search.FilterData(`json_col.string_key = 'value a'`),
			expectedRowIds: []int{0},
		},
		{
			name:           "json filter type cast to number",
			filterData:     search.FilterData(`json_col.int_key = 1`),
			expectedRowIds: []int{1},
		},
		{
			name:           "json filter type cast to boolean",
			filterData:     search.FilterData(`json_col.bool_key = true`),
			expectedRowIds: []int{0, 2, 4},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			expr, err := scenario.filterData.BuildExpr(resolver)
			if err != nil {
				t.Fatal(err)
			}

			var rows []struct {
				ID int `db:"id"`
			}
			err = testApp.DB().Select("id").From("test").Where(expr).OrderBy("id").All(&rows)
			if err != nil {
				t.Fatal(err)
			}

			rowIds := make([]int, len(rows))
			for i, row := range rows {
				rowIds[i] = row.ID
			}

			if !reflect.DeepEqual(rowIds, scenario.expectedRowIds) {
				t.Fatalf("Expected rows %v, got %v", scenario.expectedRowIds, rowIds)
			}
		})
	}
}
