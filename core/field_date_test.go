package core_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestDateFieldBaseMethods(t *testing.T) {
	testFieldBaseMethods(t, core.FieldTypeDate)
}

func TestDateFieldColumnType(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	f := &core.DateField{}

	expected := "TEXT DEFAULT '' NOT NULL"

	if v := f.ColumnType(app); v != expected {
		t.Fatalf("Expected\n%q\ngot\n%q", expected, v)
	}
}

func TestDateFieldPrepareValue(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	f := &core.DateField{}
	record := core.NewRecord(core.NewBaseCollection("test"))

	scenarios := []struct {
		raw      any
		expected string
	}{
		{"", ""},
		{"invalid", ""},
		{"2024-01-01 00:11:22.345Z", "2024-01-01 00:11:22.345Z"},
		{time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC), "2024-01-02 03:04:05.000Z"},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#v", i, s.raw), func(t *testing.T) {
			v, err := f.PrepareValue(record, s.raw)
			if err != nil {
				t.Fatal(err)
			}

			vDate, ok := v.(types.DateTime)
			if !ok {
				t.Fatalf("Expected types.DateTime instance, got %T", v)
			}

			if vDate.String() != s.expected {
				t.Fatalf("Expected %v, got %v", s.expected, v)
			}
		})
	}
}

func TestDateFieldValidateValue(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewBaseCollection("test_collection")

	scenarios := []struct {
		name        string
		field       *core.DateField
		record      func() *core.Record
		expectError bool
	}{
		{
			"invalid raw value",
			&core.DateField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", 123)
				return record
			},
			true,
		},
		{
			"zero field value (not required)",
			&core.DateField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", types.DateTime{})
				return record
			},
			false,
		},
		{
			"zero field value (required)",
			&core.DateField{Name: "test", Required: true},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", types.DateTime{})
				return record
			},
			true,
		},
		{
			"non-zero field value (required)",
			&core.DateField{Name: "test", Required: true},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", types.NowDateTime())
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

func TestDateFieldValidateSettings(t *testing.T) {
	testDefaultFieldIdValidation(t, core.FieldTypeDate)
	testDefaultFieldNameValidation(t, core.FieldTypeDate)

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewBaseCollection("test_collection")

	scenarios := []struct {
		name         string
		field        func() *core.DateField
		expectErrors []string
	}{
		{
			"zero Min/Max",
			func() *core.DateField {
				return &core.DateField{
					Id:   "test",
					Name: "test",
				}
			},
			[]string{},
		},
		{
			"non-empty Min with empty Max",
			func() *core.DateField {
				return &core.DateField{
					Id:   "test",
					Name: "test",
					Min:  types.NowDateTime(),
				}
			},
			[]string{},
		},
		{
			"empty Min non-empty Max",
			func() *core.DateField {
				return &core.DateField{
					Id:   "test",
					Name: "test",
					Max:  types.NowDateTime(),
				}
			},
			[]string{},
		},
		{
			"Min = Max",
			func() *core.DateField {
				date := types.NowDateTime()
				return &core.DateField{
					Id:   "test",
					Name: "test",
					Min:  date,
					Max:  date,
				}
			},
			[]string{},
		},
		{
			"Min > Max",
			func() *core.DateField {
				min := types.NowDateTime()
				max := min.Add(-5 * time.Second)
				return &core.DateField{
					Id:   "test",
					Name: "test",
					Min:  min,
					Max:  max,
				}
			},
			[]string{},
		},
		{
			"Min < Max",
			func() *core.DateField {
				max := types.NowDateTime()
				min := max.Add(-5 * time.Second)
				return &core.DateField{
					Id:   "test",
					Name: "test",
					Min:  min,
					Max:  max,
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
