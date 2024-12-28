package core_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

func TestEditorFieldBaseMethods(t *testing.T) {
	testFieldBaseMethods(t, core.FieldTypeEditor)
}

func TestEditorFieldColumnType(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	f := &core.EditorField{}

	expected := "TEXT DEFAULT '' NOT NULL"

	if v := f.ColumnType(app); v != expected {
		t.Fatalf("Expected\n%q\ngot\n%q", expected, v)
	}
}

func TestEditorFieldPrepareValue(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	f := &core.EditorField{}
	record := core.NewRecord(core.NewBaseCollection("test"))

	scenarios := []struct {
		raw      any
		expected string
	}{
		{"", ""},
		{"test", "test"},
		{false, "false"},
		{true, "true"},
		{123.456, "123.456"},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#v", i, s.raw), func(t *testing.T) {
			v, err := f.PrepareValue(record, s.raw)
			if err != nil {
				t.Fatal(err)
			}

			vStr, ok := v.(string)
			if !ok {
				t.Fatalf("Expected string instance, got %T", v)
			}

			if vStr != s.expected {
				t.Fatalf("Expected %q, got %q", s.expected, v)
			}
		})
	}
}

func TestEditorFieldValidateValue(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewBaseCollection("test_collection")

	scenarios := []struct {
		name        string
		field       *core.EditorField
		record      func() *core.Record
		expectError bool
	}{
		{
			"invalid raw value",
			&core.EditorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", 123)
				return record
			},
			true,
		},
		{
			"zero field value (not required)",
			&core.EditorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "")
				return record
			},
			false,
		},
		{
			"zero field value (required)",
			&core.EditorField{Name: "test", Required: true},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "")
				return record
			},
			true,
		},
		{
			"non-zero field value (required)",
			&core.EditorField{Name: "test", Required: true},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "abc")
				return record
			},
			false,
		},
		{
			"> default MaxSize",
			&core.EditorField{Name: "test", Required: true},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", strings.Repeat("a", 1+(5<<20)))
				return record
			},
			true,
		},
		{
			"> MaxSize",
			&core.EditorField{Name: "test", Required: true, MaxSize: 5},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "abcdef")
				return record
			},
			true,
		},
		{
			"<= MaxSize",
			&core.EditorField{Name: "test", Required: true, MaxSize: 5},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "abcde")
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

func TestEditorFieldValidateSettings(t *testing.T) {
	testDefaultFieldIdValidation(t, core.FieldTypeEditor)
	testDefaultFieldNameValidation(t, core.FieldTypeEditor)

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewBaseCollection("test_collection")

	scenarios := []struct {
		name         string
		field        func() *core.EditorField
		expectErrors []string
	}{
		{
			"< 0 MaxSize",
			func() *core.EditorField {
				return &core.EditorField{
					Id:      "test",
					Name:    "test",
					MaxSize: -1,
				}
			},
			[]string{"maxSize"},
		},
		{
			"= 0 MaxSize",
			func() *core.EditorField {
				return &core.EditorField{
					Id:   "test",
					Name: "test",
				}
			},
			[]string{},
		},
		{
			"> 0 MaxSize",
			func() *core.EditorField {
				return &core.EditorField{
					Id:      "test",
					Name:    "test",
					MaxSize: 1,
				}
			},
			[]string{},
		},
		{
			"MaxSize > safe json int",
			func() *core.EditorField {
				return &core.EditorField{
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

func TestEditorFieldCalculateMaxBodySize(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	scenarios := []struct {
		field    *core.EditorField
		expected int64
	}{
		{&core.EditorField{}, core.DefaultEditorFieldMaxSize},
		{&core.EditorField{MaxSize: 10}, 10},
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
