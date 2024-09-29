package core_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"golang.org/x/crypto/bcrypt"
)

func TestPasswordFieldBaseMethods(t *testing.T) {
	testFieldBaseMethods(t, core.FieldTypePassword)
}

func TestPasswordFieldColumnType(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	f := &core.PasswordField{}

	expected := "TEXT DEFAULT '' NOT NULL"

	if v := f.ColumnType(app); v != expected {
		t.Fatalf("Expected\n%q\ngot\n%q", expected, v)
	}
}

func TestPasswordFieldPrepareValue(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	f := &core.PasswordField{}
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

			pv, ok := v.(*core.PasswordFieldValue)
			if !ok {
				t.Fatalf("Expected PasswordFieldValue instance, got %T", v)
			}

			if pv.Hash != s.expected {
				t.Fatalf("Expected %q, got %q", s.expected, v)
			}
		})
	}
}

func TestPasswordFieldDriverValue(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	f := &core.PasswordField{Name: "test"}

	err := errors.New("example_err")

	scenarios := []struct {
		raw      any
		expected *core.PasswordFieldValue
	}{
		{123, &core.PasswordFieldValue{}},
		{"abc", &core.PasswordFieldValue{}},
		{"$2abc", &core.PasswordFieldValue{Hash: "$2abc"}},
		{&core.PasswordFieldValue{Hash: "test", LastError: err}, &core.PasswordFieldValue{Hash: "test", LastError: err}},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%v", i, s.raw), func(t *testing.T) {
			record := core.NewRecord(core.NewBaseCollection("test"))
			record.SetRaw(f.GetName(), s.raw)

			v, err := f.DriverValue(record)

			vStr, ok := v.(string)
			if !ok {
				t.Fatalf("Expected string instance, got %T", v)
			}

			var errStr string
			if err != nil {
				errStr = err.Error()
			}

			var expectedErrStr string
			if s.expected.LastError != nil {
				expectedErrStr = s.expected.LastError.Error()
			}

			if errStr != expectedErrStr {
				t.Fatalf("Expected error %q, got %q", expectedErrStr, errStr)
			}

			if vStr != s.expected.Hash {
				t.Fatalf("Expected hash %q, got %q", s.expected.Hash, vStr)
			}
		})
	}
}

func TestPasswordFieldValidateValue(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewBaseCollection("test_collection")

	scenarios := []struct {
		name        string
		field       *core.PasswordField
		record      func() *core.Record
		expectError bool
	}{
		{
			"invalid raw value",
			&core.PasswordField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "123")
				return record
			},
			true,
		},
		{
			"zero field value (not required)",
			&core.PasswordField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", &core.PasswordFieldValue{})
				return record
			},
			false,
		},
		{
			"zero field value (required)",
			&core.PasswordField{Name: "test", Required: true},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", &core.PasswordFieldValue{})
				return record
			},
			true,
		},
		{
			"empty hash but non-empty plain password (required)",
			&core.PasswordField{Name: "test", Required: true},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", &core.PasswordFieldValue{Plain: "test"})
				return record
			},
			true,
		},
		{
			"non-empty hash (required)",
			&core.PasswordField{Name: "test", Required: true},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", &core.PasswordFieldValue{Hash: "test"})
				return record
			},
			false,
		},
		{
			"with LastError",
			&core.PasswordField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", &core.PasswordFieldValue{LastError: errors.New("test")})
				return record
			},
			true,
		},
		{
			"< Min",
			&core.PasswordField{Name: "test", Min: 3},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", &core.PasswordFieldValue{Plain: "аб"}) // multi-byte chars test
				return record
			},
			true,
		},
		{
			">= Min",
			&core.PasswordField{Name: "test", Min: 3},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", &core.PasswordFieldValue{Plain: "абв"}) // multi-byte chars test
				return record
			},
			false,
		},
		{
			"> default Max",
			&core.PasswordField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", &core.PasswordFieldValue{Plain: strings.Repeat("a", 72)})
				return record
			},
			true,
		},
		{
			"<= default Max",
			&core.PasswordField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", &core.PasswordFieldValue{Plain: strings.Repeat("a", 71)})
				return record
			},
			false,
		},
		{
			"> Max",
			&core.PasswordField{Name: "test", Max: 2},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", &core.PasswordFieldValue{Plain: "абв"}) // multi-byte chars test
				return record
			},
			true,
		},
		{
			"<= Max",
			&core.PasswordField{Name: "test", Max: 2},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", &core.PasswordFieldValue{Plain: "аб"}) // multi-byte chars test
				return record
			},
			false,
		},
		{
			"non-matching pattern",
			&core.PasswordField{Name: "test", Pattern: `\d+`},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", &core.PasswordFieldValue{Plain: "abc"})
				return record
			},
			true,
		},
		{
			"matching pattern",
			&core.PasswordField{Name: "test", Pattern: `\d+`},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", &core.PasswordFieldValue{Plain: "123"})
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

func TestPasswordFieldValidateSettings(t *testing.T) {
	testDefaultFieldIdValidation(t, core.FieldTypePassword)
	testDefaultFieldNameValidation(t, core.FieldTypePassword)

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		name         string
		field        func(col *core.Collection) *core.PasswordField
		expectErrors []string
	}{
		{
			"zero minimal",
			func(col *core.Collection) *core.PasswordField {
				return &core.PasswordField{
					Id:   "test",
					Name: "test",
				}
			},
			[]string{},
		},
		{
			"invalid pattern",
			func(col *core.Collection) *core.PasswordField {
				return &core.PasswordField{
					Id:      "test",
					Name:    "test",
					Pattern: "(invalid",
				}
			},
			[]string{"pattern"},
		},
		{
			"valid pattern",
			func(col *core.Collection) *core.PasswordField {
				return &core.PasswordField{
					Id:      "test",
					Name:    "test",
					Pattern: `\d+`,
				}
			},
			[]string{},
		},
		{
			"Min < 0",
			func(col *core.Collection) *core.PasswordField {
				return &core.PasswordField{
					Id:   "test",
					Name: "test",
					Min:  -1,
				}
			},
			[]string{"min"},
		},
		{
			"Min > 71",
			func(col *core.Collection) *core.PasswordField {
				return &core.PasswordField{
					Id:   "test",
					Name: "test",
					Min:  72,
				}
			},
			[]string{"min"},
		},
		{
			"valid Min",
			func(col *core.Collection) *core.PasswordField {
				return &core.PasswordField{
					Id:   "test",
					Name: "test",
					Min:  5,
				}
			},
			[]string{},
		},
		{
			"Max < Min",
			func(col *core.Collection) *core.PasswordField {
				return &core.PasswordField{
					Id:   "test",
					Name: "test",
					Min:  2,
					Max:  1,
				}
			},
			[]string{"max"},
		},
		{
			"Min > Min",
			func(col *core.Collection) *core.PasswordField {
				return &core.PasswordField{
					Id:   "test",
					Name: "test",
					Min:  2,
					Max:  3,
				}
			},
			[]string{},
		},
		{
			"Max > 71",
			func(col *core.Collection) *core.PasswordField {
				return &core.PasswordField{
					Id:   "test",
					Name: "test",
					Max:  72,
				}
			},
			[]string{"max"},
		},
		{
			"cost < bcrypt.MinCost",
			func(col *core.Collection) *core.PasswordField {
				return &core.PasswordField{
					Id:   "test",
					Name: "test",
					Cost: bcrypt.MinCost - 1,
				}
			},
			[]string{"cost"},
		},
		{
			"cost > bcrypt.MaxCost",
			func(col *core.Collection) *core.PasswordField {
				return &core.PasswordField{
					Id:   "test",
					Name: "test",
					Cost: bcrypt.MaxCost + 1,
				}
			},
			[]string{"cost"},
		},
		{
			"valid cost",
			func(col *core.Collection) *core.PasswordField {
				return &core.PasswordField{
					Id:   "test",
					Name: "test",
					Cost: 12,
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

func TestPasswordFieldFindSetter(t *testing.T) {
	scenarios := []struct {
		name      string
		key       string
		value     any
		field     *core.PasswordField
		hasSetter bool
		expected  string
	}{
		{
			"no match",
			"example",
			"abc",
			&core.PasswordField{Name: "test"},
			false,
			"",
		},
		{
			"exact match",
			"test",
			"abc",
			&core.PasswordField{Name: "test"},
			true,
			`"abc"`,
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

func TestPasswordFieldFindGetter(t *testing.T) {
	scenarios := []struct {
		name      string
		key       string
		field     *core.PasswordField
		hasGetter bool
		expected  string
	}{
		{
			"no match",
			"example",
			&core.PasswordField{Name: "test"},
			false,
			"",
		},
		{
			"field name match",
			"test",
			&core.PasswordField{Name: "test"},
			true,
			"test_plain",
		},
		{
			"field name hash modifier",
			"test:hash",
			&core.PasswordField{Name: "test"},
			true,
			"test_hash",
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			collection := core.NewBaseCollection("test_collection")
			collection.Fields.Add(s.field)

			getter := s.field.FindGetter(s.key)

			hasGetter := getter != nil
			if hasGetter != s.hasGetter {
				t.Fatalf("Expected hasGetter %v, got %v", s.hasGetter, hasGetter)
			}

			if !hasGetter {
				return
			}

			record := core.NewRecord(collection)
			record.SetRaw(s.field.GetName(), &core.PasswordFieldValue{Hash: "test_hash", Plain: "test_plain"})

			result := getter(record)

			if result != s.expected {
				t.Fatalf("Expected %q, got %#v", s.expected, result)
			}
		})
	}
}
