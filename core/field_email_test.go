package core_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

func TestEmailFieldBaseMethods(t *testing.T) {
	testFieldBaseMethods(t, core.FieldTypeEmail)
}

func TestEmailFieldColumnType(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	f := &core.EmailField{}

	expected := "TEXT DEFAULT '' NOT NULL"

	if v := f.ColumnType(app); v != expected {
		t.Fatalf("Expected\n%q\ngot\n%q", expected, v)
	}
}

func TestEmailFieldPrepareValue(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	f := &core.EmailField{}
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

func TestEmailFieldValidateValue(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewBaseCollection("test_collection")

	scenarios := []struct {
		name        string
		field       *core.EmailField
		record      func() *core.Record
		expectError bool
	}{
		{
			"invalid raw value",
			&core.EmailField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", 123)
				return record
			},
			true,
		},
		{
			"zero field value (not required)",
			&core.EmailField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "")
				return record
			},
			false,
		},
		{
			"zero field value (required)",
			&core.EmailField{Name: "test", Required: true},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "")
				return record
			},
			true,
		},
		{
			"non-zero field value (required)",
			&core.EmailField{Name: "test", Required: true},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "test@example.com")
				return record
			},
			false,
		},
		{
			"invalid email",
			&core.EmailField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "invalid")
				return record
			},
			true,
		},
		{
			"failed onlyDomains",
			&core.EmailField{Name: "test", OnlyDomains: []string{"example.org", "example.net"}},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "test@example.com")
				return record
			},
			true,
		},
		{
			"success onlyDomains",
			&core.EmailField{Name: "test", OnlyDomains: []string{"example.org", "example.com"}},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "test@example.com")
				return record
			},
			false,
		},
		{
			"failed exceptDomains",
			&core.EmailField{Name: "test", ExceptDomains: []string{"example.org", "example.com"}},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "test@example.com")
				return record
			},
			true,
		},
		{
			"success exceptDomains",
			&core.EmailField{Name: "test", ExceptDomains: []string{"example.org", "example.net"}},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "test@example.com")
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

func TestEmailFieldValidateSettings(t *testing.T) {
	testDefaultFieldIdValidation(t, core.FieldTypeEmail)
	testDefaultFieldNameValidation(t, core.FieldTypeEmail)

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewBaseCollection("test_collection")

	scenarios := []struct {
		name         string
		field        func() *core.EmailField
		expectErrors []string
	}{
		{
			"zero minimal",
			func() *core.EmailField {
				return &core.EmailField{
					Id:   "test",
					Name: "test",
				}
			},
			[]string{},
		},
		{
			"both onlyDomains and exceptDomains",
			func() *core.EmailField {
				return &core.EmailField{
					Id:            "test",
					Name:          "test",
					OnlyDomains:   []string{"example.com"},
					ExceptDomains: []string{"example.org"},
				}
			},
			[]string{"onlyDomains", "exceptDomains"},
		},
		{
			"invalid onlyDomains",
			func() *core.EmailField {
				return &core.EmailField{
					Id:          "test",
					Name:        "test",
					OnlyDomains: []string{"example.com", "invalid"},
				}
			},
			[]string{"onlyDomains"},
		},
		{
			"valid onlyDomains",
			func() *core.EmailField {
				return &core.EmailField{
					Id:          "test",
					Name:        "test",
					OnlyDomains: []string{"example.com", "example.org"},
				}
			},
			[]string{},
		},
		{
			"invalid exceptDomains",
			func() *core.EmailField {
				return &core.EmailField{
					Id:            "test",
					Name:          "test",
					ExceptDomains: []string{"example.com", "invalid"},
				}
			},
			[]string{"exceptDomains"},
		},
		{
			"valid exceptDomains",
			func() *core.EmailField {
				return &core.EmailField{
					Id:            "test",
					Name:          "test",
					ExceptDomains: []string{"example.com", "example.org"},
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
