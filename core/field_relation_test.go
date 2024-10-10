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

func TestRelationFieldBaseMethods(t *testing.T) {
	testFieldBaseMethods(t, core.FieldTypeRelation)
}

func TestRelationFieldColumnType(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		name     string
		field    *core.RelationField
		expected string
	}{
		{
			"single (zero)",
			&core.RelationField{},
			"TEXT DEFAULT '' NOT NULL",
		},
		{
			"single",
			&core.RelationField{MaxSelect: 1},
			"TEXT DEFAULT '' NOT NULL",
		},
		{
			"multiple",
			&core.RelationField{MaxSelect: 2},
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

func TestRelationFieldIsMultiple(t *testing.T) {
	scenarios := []struct {
		name     string
		field    *core.RelationField
		expected bool
	}{
		{
			"zero",
			&core.RelationField{},
			false,
		},
		{
			"single",
			&core.RelationField{MaxSelect: 1},
			false,
		},
		{
			"multiple",
			&core.RelationField{MaxSelect: 2},
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

func TestRelationFieldPrepareValue(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	record := core.NewRecord(core.NewBaseCollection("test"))

	scenarios := []struct {
		raw      any
		field    *core.RelationField
		expected string
	}{
		// single
		{nil, &core.RelationField{MaxSelect: 1}, `""`},
		{"", &core.RelationField{MaxSelect: 1}, `""`},
		{123, &core.RelationField{MaxSelect: 1}, `"123"`},
		{"a", &core.RelationField{MaxSelect: 1}, `"a"`},
		{`["a"]`, &core.RelationField{MaxSelect: 1}, `"a"`},
		{[]string{}, &core.RelationField{MaxSelect: 1}, `""`},
		{[]string{"a", "b"}, &core.RelationField{MaxSelect: 1}, `"b"`},

		// multiple
		{nil, &core.RelationField{MaxSelect: 2}, `[]`},
		{"", &core.RelationField{MaxSelect: 2}, `[]`},
		{123, &core.RelationField{MaxSelect: 2}, `["123"]`},
		{"a", &core.RelationField{MaxSelect: 2}, `["a"]`},
		{`["a"]`, &core.RelationField{MaxSelect: 2}, `["a"]`},
		{[]string{}, &core.RelationField{MaxSelect: 2}, `[]`},
		{[]string{"a", "b", "c"}, &core.RelationField{MaxSelect: 2}, `["a","b","c"]`},
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

func TestRelationFieldDriverValue(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		raw      any
		field    *core.RelationField
		expected string
	}{
		// single
		{nil, &core.RelationField{MaxSelect: 1}, `""`},
		{"", &core.RelationField{MaxSelect: 1}, `""`},
		{123, &core.RelationField{MaxSelect: 1}, `"123"`},
		{"a", &core.RelationField{MaxSelect: 1}, `"a"`},
		{`["a"]`, &core.RelationField{MaxSelect: 1}, `"a"`},
		{[]string{}, &core.RelationField{MaxSelect: 1}, `""`},
		{[]string{"a", "b"}, &core.RelationField{MaxSelect: 1}, `"b"`},

		// multiple
		{nil, &core.RelationField{MaxSelect: 2}, `[]`},
		{"", &core.RelationField{MaxSelect: 2}, `[]`},
		{123, &core.RelationField{MaxSelect: 2}, `["123"]`},
		{"a", &core.RelationField{MaxSelect: 2}, `["a"]`},
		{`["a"]`, &core.RelationField{MaxSelect: 2}, `["a"]`},
		{[]string{}, &core.RelationField{MaxSelect: 2}, `[]`},
		{[]string{"a", "b", "c"}, &core.RelationField{MaxSelect: 2}, `["a","b","c"]`},
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

func TestRelationFieldValidateValue(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	demo1, err := app.FindCollectionByNameOrId("demo1")
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		name        string
		field       *core.RelationField
		record      func() *core.Record
		expectError bool
	}{
		// single
		{
			"[single] zero field value (not required)",
			&core.RelationField{Name: "test", MaxSelect: 1, CollectionId: demo1.Id},
			func() *core.Record {
				record := core.NewRecord(core.NewBaseCollection("test_collection"))
				record.SetRaw("test", "")
				return record
			},
			false,
		},
		{
			"[single] zero field value (required)",
			&core.RelationField{Name: "test", MaxSelect: 1, CollectionId: demo1.Id, Required: true},
			func() *core.Record {
				record := core.NewRecord(core.NewBaseCollection("test_collection"))
				record.SetRaw("test", "")
				return record
			},
			true,
		},
		{
			"[single] id from other collection",
			&core.RelationField{Name: "test", MaxSelect: 1, CollectionId: demo1.Id},
			func() *core.Record {
				record := core.NewRecord(core.NewBaseCollection("test_collection"))
				record.SetRaw("test", "achvryl401bhse3")
				return record
			},
			true,
		},
		{
			"[single] valid id",
			&core.RelationField{Name: "test", MaxSelect: 1, CollectionId: demo1.Id},
			func() *core.Record {
				record := core.NewRecord(core.NewBaseCollection("test_collection"))
				record.SetRaw("test", "84nmscqy84lsi1t")
				return record
			},
			false,
		},
		{
			"[single] > MaxSelect",
			&core.RelationField{Name: "test", MaxSelect: 1, CollectionId: demo1.Id},
			func() *core.Record {
				record := core.NewRecord(core.NewBaseCollection("test_collection"))
				record.SetRaw("test", []string{"84nmscqy84lsi1t", "al1h9ijdeojtsjy"})
				return record
			},
			true,
		},

		// multiple
		{
			"[multiple] zero field value (not required)",
			&core.RelationField{Name: "test", MaxSelect: 2, CollectionId: demo1.Id},
			func() *core.Record {
				record := core.NewRecord(core.NewBaseCollection("test_collection"))
				record.SetRaw("test", []string{})
				return record
			},
			false,
		},
		{
			"[multiple] zero field value (required)",
			&core.RelationField{Name: "test", MaxSelect: 2, CollectionId: demo1.Id, Required: true},
			func() *core.Record {
				record := core.NewRecord(core.NewBaseCollection("test_collection"))
				record.SetRaw("test", []string{})
				return record
			},
			true,
		},
		{
			"[multiple] id from other collection",
			&core.RelationField{Name: "test", MaxSelect: 2, CollectionId: demo1.Id},
			func() *core.Record {
				record := core.NewRecord(core.NewBaseCollection("test_collection"))
				record.SetRaw("test", []string{"84nmscqy84lsi1t", "achvryl401bhse3"})
				return record
			},
			true,
		},
		{
			"[multiple] valid id",
			&core.RelationField{Name: "test", MaxSelect: 2, CollectionId: demo1.Id},
			func() *core.Record {
				record := core.NewRecord(core.NewBaseCollection("test_collection"))
				record.SetRaw("test", []string{"84nmscqy84lsi1t", "al1h9ijdeojtsjy"})
				return record
			},
			false,
		},
		{
			"[multiple] > MaxSelect",
			&core.RelationField{Name: "test", MaxSelect: 2, CollectionId: demo1.Id},
			func() *core.Record {
				record := core.NewRecord(core.NewBaseCollection("test_collection"))
				record.SetRaw("test", []string{"84nmscqy84lsi1t", "al1h9ijdeojtsjy", "imy661ixudk5izi"})
				return record
			},
			true,
		},
		{
			"[multiple] < MinSelect",
			&core.RelationField{Name: "test", MinSelect: 2, MaxSelect: 99, CollectionId: demo1.Id},
			func() *core.Record {
				record := core.NewRecord(core.NewBaseCollection("test_collection"))
				record.SetRaw("test", []string{"84nmscqy84lsi1t"})
				return record
			},
			true,
		},
		{
			"[multiple] >= MinSelect",
			&core.RelationField{Name: "test", MinSelect: 2, MaxSelect: 99, CollectionId: demo1.Id},
			func() *core.Record {
				record := core.NewRecord(core.NewBaseCollection("test_collection"))
				record.SetRaw("test", []string{"84nmscqy84lsi1t", "al1h9ijdeojtsjy", "imy661ixudk5izi"})
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

func TestRelationFieldValidateSettings(t *testing.T) {
	testDefaultFieldIdValidation(t, core.FieldTypeRelation)
	testDefaultFieldNameValidation(t, core.FieldTypeRelation)

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	demo1, err := app.FindCollectionByNameOrId("demo1")
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		name         string
		field        func(col *core.Collection) *core.RelationField
		expectErrors []string
	}{
		{
			"zero minimal",
			func(col *core.Collection) *core.RelationField {
				return &core.RelationField{
					Id:   "test",
					Name: "test",
				}
			},
			[]string{"collectionId"},
		},
		{
			"invalid collectionId",
			func(col *core.Collection) *core.RelationField {
				return &core.RelationField{
					Id:           "test",
					Name:         "test",
					CollectionId: demo1.Name,
				}
			},
			[]string{"collectionId"},
		},
		{
			"valid collectionId",
			func(col *core.Collection) *core.RelationField {
				return &core.RelationField{
					Id:           "test",
					Name:         "test",
					CollectionId: demo1.Id,
				}
			},
			[]string{},
		},
		{
			"base->view",
			func(col *core.Collection) *core.RelationField {
				return &core.RelationField{
					Id:           "test",
					Name:         "test",
					CollectionId: "v9gwnfh02gjq1q0",
				}
			},
			[]string{"collectionId"},
		},
		{
			"view->view",
			func(col *core.Collection) *core.RelationField {
				col.Type = core.CollectionTypeView
				return &core.RelationField{
					Id:           "test",
					Name:         "test",
					CollectionId: "v9gwnfh02gjq1q0",
				}
			},
			[]string{},
		},
		{
			"MinSelect < 0",
			func(col *core.Collection) *core.RelationField {
				return &core.RelationField{
					Id:           "test",
					Name:         "test",
					CollectionId: demo1.Id,
					MinSelect:    -1,
				}
			},
			[]string{"minSelect"},
		},
		{
			"MinSelect > 0",
			func(col *core.Collection) *core.RelationField {
				return &core.RelationField{
					Id:           "test",
					Name:         "test",
					CollectionId: demo1.Id,
					MinSelect:    1,
				}
			},
			[]string{"maxSelect"},
		},
		{
			"MaxSelect < MinSelect",
			func(col *core.Collection) *core.RelationField {
				return &core.RelationField{
					Id:           "test",
					Name:         "test",
					CollectionId: demo1.Id,
					MinSelect:    2,
					MaxSelect:    1,
				}
			},
			[]string{"maxSelect"},
		},
		{
			"MaxSelect >= MinSelect",
			func(col *core.Collection) *core.RelationField {
				return &core.RelationField{
					Id:           "test",
					Name:         "test",
					CollectionId: demo1.Id,
					MinSelect:    2,
					MaxSelect:    2,
				}
			},
			[]string{},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			collection := core.NewBaseCollection("test_collection")
			collection.Fields.GetByName("id").SetId("test") // set a dummy known id so that it can be replaced

			field := s.field(collection)

			collection.Fields.Add(field)

			errs := field.ValidateSettings(context.Background(), app, collection)

			tests.TestValidationErrors(t, errs, s.expectErrors)
		})
	}
}

func TestRelationFieldFindSetter(t *testing.T) {
	scenarios := []struct {
		name      string
		key       string
		value     any
		field     *core.RelationField
		hasSetter bool
		expected  string
	}{
		{
			"no match",
			"example",
			"b",
			&core.RelationField{Name: "test", MaxSelect: 1},
			false,
			"",
		},
		{
			"exact match (single)",
			"test",
			"b",
			&core.RelationField{Name: "test", MaxSelect: 1},
			true,
			`"b"`,
		},
		{
			"exact match (multiple)",
			"test",
			[]string{"a", "b"},
			&core.RelationField{Name: "test", MaxSelect: 2},
			true,
			`["a","b"]`,
		},
		{
			"append (single)",
			"test+",
			"b",
			&core.RelationField{Name: "test", MaxSelect: 1},
			true,
			`"b"`,
		},
		{
			"append (multiple)",
			"test+",
			[]string{"a"},
			&core.RelationField{Name: "test", MaxSelect: 2},
			true,
			`["c","d","a"]`,
		},
		{
			"prepend (single)",
			"+test",
			"b",
			&core.RelationField{Name: "test", MaxSelect: 1},
			true,
			`"d"`, // the last of the existing values
		},
		{
			"prepend (multiple)",
			"+test",
			[]string{"a"},
			&core.RelationField{Name: "test", MaxSelect: 2},
			true,
			`["a","c","d"]`,
		},
		{
			"subtract (single)",
			"test-",
			"d",
			&core.RelationField{Name: "test", MaxSelect: 1},
			true,
			`"c"`,
		},
		{
			"subtract (multiple)",
			"test-",
			[]string{"unknown", "c"},
			&core.RelationField{Name: "test", MaxSelect: 2},
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
