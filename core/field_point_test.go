package core_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestPointFieldBaseMethods(t *testing.T) {
	testFieldBaseMethods(t, core.FieldTypePoint)
}

func TestPointFieldColumnType(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	f := &core.PointField{}
	expected := "TEXT DEFAULT '' NOT NULL"

	if v := f.ColumnType(app); v != expected {
		t.Fatalf("Expected\n%q\ngot\n%q", expected, v)
	}
}

func TestPointFieldPrepareValue(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	f := &core.PointField{}
	record := core.NewRecord(core.NewBaseCollection("test"))

	scenarios := []struct {
		raw      any
		expected string
	}{
		{"", ""},
		{"invalid", ""},
		{"42.3631, -71.0574", "42.3631, -71.0574"},
		{types.Point{}, ""},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#v", i, s.raw), func(t *testing.T) {
			v, err := f.PrepareValue(record, s.raw)
			if err != nil && s.expected != "" {
				t.Fatal(err)
			}

			vPoint, ok := v.(types.Point)
			if !ok {
				t.Fatalf("Expected types.Point instance, got %T", v)
			}

			if vPoint.String() != s.expected {
				t.Fatalf("Expected %v, got %v", s.expected, v)
			}
		})
	}
}

func TestPointFieldValidateValue(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewBaseCollection("test_collection")

	scenarios := []struct {
		name        string
		field       *core.PointField
		record      func() *core.Record
		expectError bool
	}{
		{
			"invalid raw value",
			&core.PointField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", 123)
				return record
			},
			true,
		},
		{
			"empty field value (not required)",
			&core.PointField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", types.Point{})
				return record
			},
			false,
		},
		{
			"empty field value (required)",
			&core.PointField{Name: "test", Required: true},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", types.Point{})
				return record
			},
			true,
		},
		{
			"valid coordinate pair",
			&core.PointField{Name: "test", Required: true},
			func() *core.Record {
				record := core.NewRecord(collection)
				p, _ := types.ParsePoint("42.3631, -71.0574")
				record.SetRaw("test", p)
				return record
			},
			false,
		},
		{
			"invalid latitude (>90)",
			&core.PointField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				p, _ := types.ParsePoint("91, 0")
				record.SetRaw("test", p)
				return record
			},
			true,
		},
		{
			"invalid latitude (<-90)",
			&core.PointField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				p, _ := types.ParsePoint("-91, 0")
				record.SetRaw("test", p)
				return record
			},
			true,
		},
		{
			"invalid longitude (>180)",
			&core.PointField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				p, _ := types.ParsePoint("0, 181")
				record.SetRaw("test", p)
				return record
			},
			true,
		},
		{
			"invalid longitude (<-180)",
			&core.PointField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				p, _ := types.ParsePoint("0, -181")
				record.SetRaw("test", p)
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

func TestPointFieldValidateSettings(t *testing.T) {
	testDefaultFieldIdValidation(t, core.FieldTypePoint)
	testDefaultFieldNameValidation(t, core.FieldTypePoint)

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewBaseCollection("test_collection")

	scenarios := []struct {
		name         string
		field        func() *core.PointField
		expectErrors []string
	}{
		{
			"valid settings",
			func() *core.PointField {
				return &core.PointField{
					Id:   "test",
					Name: "test",
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
