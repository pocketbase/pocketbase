package core_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

func TestBoolFieldBaseMethods(t *testing.T) {
	testFieldBaseMethods(t, core.FieldTypeBool)
}

func TestBoolFieldColumnType(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	f := &core.BoolField{}

	expected := "BOOLEAN DEFAULT FALSE NOT NULL"

	if v := f.ColumnType(app); v != expected {
		t.Fatalf("Expected\n%q\ngot\n%q", expected, v)
	}
}

func TestBoolFieldPrepareValue(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	f := &core.BoolField{}
	record := core.NewRecord(core.NewBaseCollection("test"))

	scenarios := []struct {
		raw      any
		expected bool
	}{
		{"", false},
		{"f", false},
		{"t", true},
		{1, true},
		{0, false},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#v", i, s.raw), func(t *testing.T) {
			v, err := f.PrepareValue(record, s.raw)
			if err != nil {
				t.Fatal(err)
			}

			if v != s.expected {
				t.Fatalf("Expected %v, got %v", s.expected, v)
			}
		})
	}
}

func TestBoolFieldValidateValue(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewBaseCollection("test_collection")

	scenarios := []struct {
		name        string
		field       *core.BoolField
		record      func() *core.Record
		expectError bool
	}{
		{
			"invalid raw value",
			&core.BoolField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", 123)
				return record
			},
			true,
		},
		{
			"missing field value (non-required)",
			&core.BoolField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("abc", true)
				return record
			},
			true, // because of failed nil.(bool) cast
		},
		{
			"missing field value (required)",
			&core.BoolField{Name: "test", Required: true},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("abc", true)
				return record
			},
			true,
		},
		{
			"false field value (non-required)",
			&core.BoolField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", false)
				return record
			},
			false,
		},
		{
			"false field value (required)",
			&core.BoolField{Name: "test", Required: true},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", false)
				return record
			},
			true,
		},
		{
			"true field value (required)",
			&core.BoolField{Name: "test", Required: true},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", true)
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

func TestBoolFieldValidateSettings(t *testing.T) {
	testDefaultFieldIdValidation(t, core.FieldTypeBool)
	testDefaultFieldNameValidation(t, core.FieldTypeBool)
}
