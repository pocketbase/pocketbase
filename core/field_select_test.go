package core_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestSelectFieldBaseMethods(t *testing.T) {
	testFieldBaseMethods(t, core.FieldTypeSelect)
}

func TestSelectFieldColumnType(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		name     string
		field    *core.SelectField
		expected string
	}{
		{
			"single (zero)",
			&core.SelectField{},
			"TEXT DEFAULT '' NOT NULL",
		},
		{
			"single",
			&core.SelectField{MaxSelect: 1},
			"TEXT DEFAULT '' NOT NULL",
		},
		{
			"multiple",
			&core.SelectField{MaxSelect: 2},
			"JSON DEFAULT '[]' NOT NULL",
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			if v := s.field.ColumnType(app); v != s.expected {
				t.Fatalf("Expected\n%q\ngot\n%q", s.expected, v)
			}
		})
	}
}

func TestSelectFieldIsMultiple(t *testing.T) {
	scenarios := []struct {
		name     string
		field    *core.SelectField
		expected bool
	}{
		{
			"single (zero)",
			&core.SelectField{},
			false,
		},
		{
			"single",
			&core.SelectField{MaxSelect: 1},
			false,
		},
		{
			"multiple (>1)",
			&core.SelectField{MaxSelect: 2},
			true,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			if v := s.field.IsMultiple(); v != s.expected {
				t.Fatalf("Expected %v, got %v", s.expected, v)
			}
		})
	}
}

func TestSelectFieldPrepareValue(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	record := core.NewRecord(core.NewBaseCollection("test"))

	scenarios := []struct {
		raw      any
		field    *core.SelectField
		expected string
	}{
		// single
		{nil, &core.SelectField{}, `""`},
		{"", &core.SelectField{}, `""`},
		{123, &core.SelectField{}, `"123"`},
		{"a", &core.SelectField{}, `"a"`},
		{`["a"]`, &core.SelectField{}, `"a"`},
		{[]string{}, &core.SelectField{}, `""`},
		{[]string{"a", "b"}, &core.SelectField{}, `"b"`},

		// multiple
		{nil, &core.SelectField{MaxSelect: 2}, `[]`},
		{"", &core.SelectField{MaxSelect: 2}, `[]`},
		{123, &core.SelectField{MaxSelect: 2}, `["123"]`},
		{"a", &core.SelectField{MaxSelect: 2}, `["a"]`},
		{`["a"]`, &core.SelectField{MaxSelect: 2}, `["a"]`},
		{[]string{}, &core.SelectField{MaxSelect: 2}, `[]`},
		{[]string{"a", "b", "c"}, &core.SelectField{MaxSelect: 2}, `["a","b","c"]`},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#v_%v", i, s.raw, s.field.IsMultiple()), func(t *testing.T) {
			v, err := s.field.PrepareValue(record, s.raw)
			if err != nil {
				t.Fatal(err)
			}

			vRaw, err := json.Marshal(v)
			if err != nil {
				t.Fatal(err)
			}

			if string(vRaw) != s.expected {
				t.Fatalf("Expected %q, got %q", s.expected, vRaw)
			}
		})
	}
}

func TestSelectFieldDriverValue(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		raw      any
		field    *core.SelectField
		expected string
	}{
		// single
		{nil, &core.SelectField{}, `""`},
		{"", &core.SelectField{}, `""`},
		{123, &core.SelectField{}, `"123"`},
		{"a", &core.SelectField{}, `"a"`},
		{`["a"]`, &core.SelectField{}, `"a"`},
		{[]string{}, &core.SelectField{}, `""`},
		{[]string{"a", "b"}, &core.SelectField{}, `"b"`},

		// multiple
		{nil, &core.SelectField{MaxSelect: 2}, `[]`},
		{"", &core.SelectField{MaxSelect: 2}, `[]`},
		{123, &core.SelectField{MaxSelect: 2}, `["123"]`},
		{"a", &core.SelectField{MaxSelect: 2}, `["a"]`},
		{`["a"]`, &core.SelectField{MaxSelect: 2}, `["a"]`},
		{[]string{}, &core.SelectField{MaxSelect: 2}, `[]`},
		{[]string{"a", "b", "c"}, &core.SelectField{MaxSelect: 2}, `["a","b","c"]`},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#v_%v", i, s.raw, s.field.IsMultiple()), func(t *testing.T) {
			record := core.NewRecord(core.NewBaseCollection("test"))
			record.SetRaw(s.field.GetName(), s.raw)

			v, err := s.field.DriverValue(record)
			if err != nil {
				t.Fatal(err)
			}

			if s.field.IsMultiple() {
				_, ok := v.(types.JSONArray[string])
				if !ok {
					t.Fatalf("Expected types.JSONArray value, got %T", v)
				}
			} else {
				_, ok := v.(string)
				if !ok {
					t.Fatalf("Expected string value, got %T", v)
				}
			}

			vRaw, err := json.Marshal(v)
			if err != nil {
				t.Fatal(err)
			}

			if string(vRaw) != s.expected {
				t.Fatalf("Expected %q, got %q", s.expected, vRaw)
			}
		})
	}
}

func TestSelectFieldValidateValue(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewBaseCollection("test_collection")

	values := []string{"a", "b", "c"}

	scenarios := []struct {
		name        string
		field       *core.SelectField
		record      func() *core.Record
		expectError bool
	}{
		// single
		{
			"[single] zero field value (not required)",
			&core.SelectField{Name: "test", Values: values, MaxSelect: 1},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "")
				return record
			},
			false,
		},
		{
			"[single] zero field value (required)",
			&core.SelectField{Name: "test", Values: values, MaxSelect: 1, Required: true},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "")
				return record
			},
			true,
		},
		{
			"[single] unknown value",
			&core.SelectField{Name: "test", Values: values, MaxSelect: 1},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "unknown")
				return record
			},
			true,
		},
		{
			"[single] known value",
			&core.SelectField{Name: "test", Values: values, MaxSelect: 1},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "a")
				return record
			},
			false,
		},
		{
			"[single] > MaxSelect",
			&core.SelectField{Name: "test", Values: values, MaxSelect: 1},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", []string{"a", "b"})
				return record
			},
			true,
		},

		// multiple
		{
			"[multiple] zero field value (not required)",
			&core.SelectField{Name: "test", Values: values, MaxSelect: 2},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", []string{})
				return record
			},
			false,
		},
		{
			"[multiple] zero field value (required)",
			&core.SelectField{Name: "test", Values: values, MaxSelect: 2, Required: true},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", []string{})
				return record
			},
			true,
		},
		{
			"[multiple] unknown value",
			&core.SelectField{Name: "test", Values: values, MaxSelect: 2},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", []string{"a", "unknown"})
				return record
			},
			true,
		},
		{
			"[multiple] known value",
			&core.SelectField{Name: "test", Values: values, MaxSelect: 2},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", []string{"a", "b"})
				return record
			},
			false,
		},
		{
			"[multiple] > MaxSelect",
			&core.SelectField{Name: "test", Values: values, MaxSelect: 2},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", []string{"a", "b", "c"})
				return record
			},
			true,
		},
		{
			"[multiple] > MaxSelect (duplicated values)",
			&core.SelectField{Name: "test", Values: values, MaxSelect: 2},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", []string{"a", "b", "b", "a"})
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

func TestSelectFieldValidateSettings(t *testing.T) {
	testDefaultFieldIdValidation(t, core.FieldTypeSelect)
	testDefaultFieldNameValidation(t, core.FieldTypeSelect)

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		name         string
		field        func() *core.SelectField
		expectErrors []string
	}{
		{
			"zero minimal",
			func() *core.SelectField {
				return &core.SelectField{
					Id:   "test",
					Name: "test",
				}
			},
			[]string{"values"},
		},
		{
			"MaxSelect > Values length",
			func() *core.SelectField {
				return &core.SelectField{
					Id:        "test",
					Name:      "test",
					Values:    []string{"a", "b"},
					MaxSelect: 3,
				}
			},
			[]string{"maxSelect"},
		},
		{
			"MaxSelect <= Values length",
			func() *core.SelectField {
				return &core.SelectField{
					Id:        "test",
					Name:      "test",
					Values:    []string{"a", "b"},
					MaxSelect: 2,
				}
			},
			[]string{},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			field := s.field()

			collection := core.NewBaseCollection("test_collection")
			collection.Fields.Add(field)

			errs := field.ValidateSettings(context.Background(), app, collection)

			tests.TestValidationErrors(t, errs, s.expectErrors)
		})
	}
}

func TestSelectFieldFindSetter(t *testing.T) {
	values := []string{"a", "b", "c", "d"}

	scenarios := []struct {
		name      string
		key       string
		value     any
		field     *core.SelectField
		hasSetter bool
		expected  string
	}{
		{
			"no match",
			"example",
			"b",
			&core.SelectField{Name: "test", MaxSelect: 1, Values: values},
			false,
			"",
		},
		{
			"exact match (single)",
			"test",
			"b",
			&core.SelectField{Name: "test", MaxSelect: 1, Values: values},
			true,
			`"b"`,
		},
		{
			"exact match (multiple)",
			"test",
			[]string{"a", "b"},
			&core.SelectField{Name: "test", MaxSelect: 2, Values: values},
			true,
			`["a","b"]`,
		},
		{
			"append (single)",
			"test+",
			"b",
			&core.SelectField{Name: "test", MaxSelect: 1, Values: values},
			true,
			`"b"`,
		},
		{
			"append (multiple)",
			"test+",
			[]string{"a"},
			&core.SelectField{Name: "test", MaxSelect: 2, Values: values},
			true,
			`["c","d","a"]`,
		},
		{
			"prepend (single)",
			"+test",
			"b",
			&core.SelectField{Name: "test", MaxSelect: 1, Values: values},
			true,
			`"d"`, // the last of the existing values
		},
		{
			"prepend (multiple)",
			"+test",
			[]string{"a"},
			&core.SelectField{Name: "test", MaxSelect: 2, Values: values},
			true,
			`["a","c","d"]`,
		},
		{
			"subtract (single)",
			"test-",
			"d",
			&core.SelectField{Name: "test", MaxSelect: 1, Values: values},
			true,
			`"c"`,
		},
		{
			"subtract (multiple)",
			"test-",
			[]string{"unknown", "c"},
			&core.SelectField{Name: "test", MaxSelect: 2, Values: values},
			true,
			`["d"]`,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			collection := core.NewBaseCollection("test_collection")
			collection.Fields.Add(s.field)

			setter := s.field.FindSetter(s.key)

			hasSetter := setter != nil
			if hasSetter != s.hasSetter {
				t.Fatalf("Expected hasSetter %v, got %v", s.hasSetter, hasSetter)
			}

			if !hasSetter {
				return
			}

			record := core.NewRecord(collection)
			record.SetRaw(s.field.GetName(), []string{"c", "d"})

			setter(record, s.value)

			raw, err := json.Marshal(record.Get(s.field.GetName()))
			if err != nil {
				t.Fatal(err)
			}
			rawStr := string(raw)

			if rawStr != s.expected {
				t.Fatalf("Expected %q, got %q", s.expected, rawStr)
			}
		})
	}
}
