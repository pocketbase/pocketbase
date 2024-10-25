package core_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestNumberFieldBaseMethods(t *testing.T) {
	testFieldBaseMethods(t, core.FieldTypeNumber)
}

func TestNumberFieldColumnType(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	f := &core.NumberField{}

	expected := "NUMERIC DEFAULT 0 NOT NULL"

	if v := f.ColumnType(app); v != expected {
		t.Fatalf("Expected\n%q\ngot\n%q", expected, v)
	}
}

func TestNumberFieldPrepareValue(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	f := &core.NumberField{}
	record := core.NewRecord(core.NewBaseCollection("test"))

	scenarios := []struct {
		raw      any
		expected float64
	}{
		{"", 0},
		{"test", 0},
		{false, 0},
		{true, 1},
		{-2, -2},
		{123.456, 123.456},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#v", i, s.raw), func(t *testing.T) {
			vRaw, err := f.PrepareValue(record, s.raw)
			if err != nil {
				t.Fatal(err)
			}

			v, ok := vRaw.(float64)
			if !ok {
				t.Fatalf("Expected float64 instance, got %T", v)
			}

			if v != s.expected {
				t.Fatalf("Expected %f, got %f", s.expected, v)
			}
		})
	}
}

func TestNumberFieldValidateValue(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewBaseCollection("test_collection")

	scenarios := []struct {
		name        string
		field       *core.NumberField
		record      func() *core.Record
		expectError bool
	}{
		{
			"invalid raw value",
			&core.NumberField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "123")
				return record
			},
			true,
		},
		{
			"zero field value (not required)",
			&core.NumberField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", 0.0)
				return record
			},
			false,
		},
		{
			"zero field value (required)",
			&core.NumberField{Name: "test", Required: true},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", 0.0)
				return record
			},
			true,
		},
		{
			"non-zero field value (required)",
			&core.NumberField{Name: "test", Required: true},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", 123.0)
				return record
			},
			false,
		},
		{
			"decimal with onlyInt",
			&core.NumberField{Name: "test", OnlyInt: true},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", 123.456)
				return record
			},
			true,
		},
		{
			"int with onlyInt",
			&core.NumberField{Name: "test", OnlyInt: true},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", 123.0)
				return record
			},
			false,
		},
		{
			"< min",
			&core.NumberField{Name: "test", Min: types.Pointer(2.0)},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", 1.0)
				return record
			},
			true,
		},
		{
			">= min",
			&core.NumberField{Name: "test", Min: types.Pointer(2.0)},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", 2.0)
				return record
			},
			false,
		},
		{
			"> max",
			&core.NumberField{Name: "test", Max: types.Pointer(2.0)},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", 3.0)
				return record
			},
			true,
		},
		{
			"<= max",
			&core.NumberField{Name: "test", Max: types.Pointer(2.0)},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", 2.0)
				return record
			},
			false,
		},
		{
			"infinity",
			&core.NumberField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.Set("test", "Inf")
				return record
			},
			true,
		},
		{
			"NaN",
			&core.NumberField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.Set("test", "NaN")
				return record
			},
			true,
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

func TestNumberFieldValidateSettings(t *testing.T) {
	testDefaultFieldIdValidation(t, core.FieldTypeNumber)
	testDefaultFieldNameValidation(t, core.FieldTypeNumber)

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewBaseCollection("test_collection")

	scenarios := []struct {
		name         string
		field        func() *core.NumberField
		expectErrors []string
	}{
		{
			"zero",
			func() *core.NumberField {
				return &core.NumberField{
					Id:   "test",
					Name: "test",
				}
			},
			[]string{},
		},
		{
			"decumal min",
			func() *core.NumberField {
				return &core.NumberField{
					Id:   "test",
					Name: "test",
					Min:  types.Pointer(1.2),
				}
			},
			[]string{},
		},
		{
			"decumal min (onlyInt)",
			func() *core.NumberField {
				return &core.NumberField{
					Id:      "test",
					Name:    "test",
					OnlyInt: true,
					Min:     types.Pointer(1.2),
				}
			},
			[]string{"min"},
		},
		{
			"int min (onlyInt)",
			func() *core.NumberField {
				return &core.NumberField{
					Id:      "test",
					Name:    "test",
					OnlyInt: true,
					Min:     types.Pointer(1.0),
				}
			},
			[]string{},
		},
		{
			"decumal max",
			func() *core.NumberField {
				return &core.NumberField{
					Id:   "test",
					Name: "test",
					Max:  types.Pointer(1.2),
				}
			},
			[]string{},
		},
		{
			"decumal max (onlyInt)",
			func() *core.NumberField {
				return &core.NumberField{
					Id:      "test",
					Name:    "test",
					OnlyInt: true,
					Max:     types.Pointer(1.2),
				}
			},
			[]string{"max"},
		},
		{
			"int max (onlyInt)",
			func() *core.NumberField {
				return &core.NumberField{
					Id:      "test",
					Name:    "test",
					OnlyInt: true,
					Max:     types.Pointer(1.0),
				}
			},
			[]string{},
		},
		{
			"min > max",
			func() *core.NumberField {
				return &core.NumberField{
					Id:   "test",
					Name: "test",
					Min:  types.Pointer(2.0),
					Max:  types.Pointer(1.0),
				}
			},
			[]string{"max"},
		},
		{
			"min <= max",
			func() *core.NumberField {
				return &core.NumberField{
					Id:   "test",
					Name: "test",
					Min:  types.Pointer(2.0),
					Max:  types.Pointer(2.0),
				}
			},
			[]string{},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			errs := s.field().ValidateSettings(context.Background(), app, collection)

			tests.TestValidationErrors(t, errs, s.expectErrors)
		})
	}
}

func TestNumberFieldFindSetter(t *testing.T) {
	field := &core.NumberField{Name: "test"}

	collection := core.NewBaseCollection("test_collection")
	collection.Fields.Add(field)

	t.Run("no match", func(t *testing.T) {
		f := field.FindSetter("abc")
		if f != nil {
			t.Fatal("Expected nil setter")
		}
	})

	t.Run("direct name match", func(t *testing.T) {
		f := field.FindSetter("test")
		if f == nil {
			t.Fatal("Expected non-nil setter")
		}

		record := core.NewRecord(collection)
		record.SetRaw("test", 2.0)

		f(record, "123.456") // should be casted

		if v := record.Get("test"); v != 123.456 {
			t.Fatalf("Expected %f, got %f", 123.456, v)
		}
	})

	t.Run("name+ match", func(t *testing.T) {
		f := field.FindSetter("test+")
		if f == nil {
			t.Fatal("Expected non-nil setter")
		}

		record := core.NewRecord(collection)
		record.SetRaw("test", 2.0)

		f(record, "1.5") // should be casted and appended to the existing value

		if v := record.Get("test"); v != 3.5 {
			t.Fatalf("Expected %f, got %f", 3.5, v)
		}
	})

	t.Run("name- match", func(t *testing.T) {
		f := field.FindSetter("test-")
		if f == nil {
			t.Fatal("Expected non-nil setter")
		}

		record := core.NewRecord(collection)
		record.SetRaw("test", 2.0)

		f(record, "1.5") // should be casted and subtracted from the existing value

		if v := record.Get("test"); v != 0.5 {
			t.Fatalf("Expected %f, got %f", 0.5, v)
		}
	})
}
