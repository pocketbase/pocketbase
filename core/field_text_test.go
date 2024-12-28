package core_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

func TestTextFieldBaseMethods(t *testing.T) {
	testFieldBaseMethods(t, core.FieldTypeText)
}

func TestTextFieldColumnType(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	f := &core.TextField{}

	expected := "TEXT DEFAULT '' NOT NULL"

	if v := f.ColumnType(app); v != expected {
		t.Fatalf("Expected\n%q\ngot\n%q", expected, v)
	}
}

func TestTextFieldPrepareValue(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	f := &core.TextField{}
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

func TestTextFieldValidateValue(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, err := app.FindCollectionByNameOrId("demo1")
	if err != nil {
		t.Fatal(err)
	}

	existingRecord, err := app.FindFirstRecordByFilter(collection, "id != ''")
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		name        string
		field       *core.TextField
		record      func() *core.Record
		expectError bool
	}{
		{
			"invalid raw value",
			&core.TextField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", 123)
				return record
			},
			true,
		},
		{
			"zero field value (not required)",
			&core.TextField{Name: "test", Pattern: `\d+`, Min: 10, Max: 100}, // other fields validators should be ignored
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "")
				return record
			},
			false,
		},
		{
			"zero field value (required)",
			&core.TextField{Name: "test", Required: true},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "")
				return record
			},
			true,
		},
		{
			"non-zero field value (required)",
			&core.TextField{Name: "test", Required: true},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "abc")
				return record
			},
			false,
		},
		{
			"special forbidden character / (non-primaryKey)",
			&core.TextField{Name: "test", PrimaryKey: false},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "/")
				return record
			},
			false,
		},
		{
			"special forbidden character \\ (non-primaryKey)",
			&core.TextField{Name: "test", PrimaryKey: false},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "\\")
				return record
			},
			false,
		},
		{
			"special forbidden character / (primaryKey)",
			&core.TextField{Name: "test", PrimaryKey: true},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "/")
				return record
			},
			true,
		},
		{
			"special forbidden character \\ (primaryKey)",
			&core.TextField{Name: "test", PrimaryKey: true},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "\\")
				return record
			},
			true,
		},
		{
			"zero field value (primaryKey)",
			&core.TextField{Name: "test", PrimaryKey: true},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "")
				return record
			},
			true,
		},
		{
			"non-zero field value (primaryKey)",
			&core.TextField{Name: "test", PrimaryKey: true},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "abcd")
				return record
			},
			false,
		},
		{
			"case-insensitive duplicated primary key check",
			&core.TextField{Name: "test", PrimaryKey: true},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", strings.ToUpper(existingRecord.Id))
				return record
			},
			true,
		},
		{
			"< min",
			&core.TextField{Name: "test", Min: 4},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "абв") // multi-byte
				return record
			},
			true,
		},
		{
			">= min",
			&core.TextField{Name: "test", Min: 3},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "абв") // multi-byte
				return record
			},
			false,
		},
		{
			"> default max",
			&core.TextField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", strings.Repeat("a", 5001))
				return record
			},
			true,
		},
		{
			"<= default max",
			&core.TextField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", strings.Repeat("a", 500))
				return record
			},
			false,
		},
		{
			"> max",
			&core.TextField{Name: "test", Max: 2},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "абв") // multi-byte
				return record
			},
			true,
		},
		{
			"<= max",
			&core.TextField{Name: "test", Min: 3},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "абв") // multi-byte
				return record
			},
			false,
		},
		{
			"mismatched pattern",
			&core.TextField{Name: "test", Pattern: `\d+`},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "abc")
				return record
			},
			true,
		},
		{
			"matched pattern",
			&core.TextField{Name: "test", Pattern: `\d+`},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "123")
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

func TestTextFieldValidateSettings(t *testing.T) {
	testDefaultFieldIdValidation(t, core.FieldTypeText)
	testDefaultFieldNameValidation(t, core.FieldTypeText)

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		name         string
		field        func() *core.TextField
		expectErrors []string
	}{
		{
			"zero minimal",
			func() *core.TextField {
				return &core.TextField{
					Id:   "test",
					Name: "test",
				}
			},
			[]string{},
		},
		{
			"primaryKey without required",
			func() *core.TextField {
				return &core.TextField{
					Id:         "test",
					Name:       "id",
					PrimaryKey: true,
					Pattern:    `\d+`,
				}
			},
			[]string{"required"},
		},
		{
			"primaryKey without pattern",
			func() *core.TextField {
				return &core.TextField{
					Id:         "test",
					Name:       "id",
					PrimaryKey: true,
					Required:   true,
				}
			},
			[]string{"pattern"},
		},
		{
			"primaryKey with hidden",
			func() *core.TextField {
				return &core.TextField{
					Id:         "test",
					Name:       "id",
					Required:   true,
					PrimaryKey: true,
					Hidden:     true,
					Pattern:    `\d+`,
				}
			},
			[]string{"hidden"},
		},
		{
			"primaryKey with name != id",
			func() *core.TextField {
				return &core.TextField{
					Id:         "test",
					Name:       "test",
					PrimaryKey: true,
					Required:   true,
					Pattern:    `\d+`,
				}
			},
			[]string{"name"},
		},
		{
			"multiple primaryKey fields",
			func() *core.TextField {
				return &core.TextField{
					Id:         "test2",
					Name:       "id",
					PrimaryKey: true,
					Pattern:    `\d+`,
					Required:   true,
				}
			},
			[]string{"primaryKey"},
		},
		{
			"invalid pattern",
			func() *core.TextField {
				return &core.TextField{
					Id:      "test2",
					Name:    "id",
					Pattern: `(invalid`,
				}
			},
			[]string{"pattern"},
		},
		{
			"valid pattern",
			func() *core.TextField {
				return &core.TextField{
					Id:      "test2",
					Name:    "id",
					Pattern: `\d+`,
				}
			},
			[]string{},
		},
		{
			"invalid autogeneratePattern",
			func() *core.TextField {
				return &core.TextField{
					Id:                  "test2",
					Name:                "id",
					AutogeneratePattern: `(invalid`,
				}
			},
			[]string{"autogeneratePattern"},
		},
		{
			"valid autogeneratePattern",
			func() *core.TextField {
				return &core.TextField{
					Id:                  "test2",
					Name:                "id",
					AutogeneratePattern: `[a-z]+`,
				}
			},
			[]string{},
		},
		{
			"conflicting pattern and autogeneratePattern",
			func() *core.TextField {
				return &core.TextField{
					Id:                  "test2",
					Name:                "id",
					Pattern:             `\d+`,
					AutogeneratePattern: `[a-z]+`,
				}
			},
			[]string{"autogeneratePattern"},
		},
		{
			"Max > safe json int",
			func() *core.TextField {
				return &core.TextField{
					Id:   "test",
					Name: "test",
					Max:  1 << 53,
				}
			},
			[]string{"max"},
		},
		{
			"Max < 0",
			func() *core.TextField {
				return &core.TextField{
					Id:   "test",
					Name: "test",
					Max:  -1,
				}
			},
			[]string{"max"},
		},
		{
			"Min > safe json int",
			func() *core.TextField {
				return &core.TextField{
					Id:   "test",
					Name: "test",
					Min:  1 << 53,
				}
			},
			[]string{"min"},
		},
		{
			"Min < 0",
			func() *core.TextField {
				return &core.TextField{
					Id:   "test",
					Name: "test",
					Min:  -1,
				}
			},
			[]string{"min"},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			field := s.field()

			collection := core.NewBaseCollection("test_collection")
			collection.Fields.GetByName("id").SetId("test") // set a dummy known id so that it can be replaced
			collection.Fields.Add(field)

			errs := field.ValidateSettings(context.Background(), app, collection)

			tests.TestValidationErrors(t, errs, s.expectErrors)
		})
	}
}

func TestTextFieldAutogenerate(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewBaseCollection("test_collection")

	scenarios := []struct {
		name       string
		actionName string
		field      *core.TextField
		record     func() *core.Record
		expected   string
	}{
		{
			"non-matching action",
			core.InterceptorActionUpdate,
			&core.TextField{Name: "test", AutogeneratePattern: "abc"},
			func() *core.Record {
				return core.NewRecord(collection)
			},
			"",
		},
		{
			"matching action (create)",
			core.InterceptorActionCreate,
			&core.TextField{Name: "test", AutogeneratePattern: "abc"},
			func() *core.Record {
				return core.NewRecord(collection)
			},
			"abc",
		},
		{
			"matching action (validate)",
			core.InterceptorActionValidate,
			&core.TextField{Name: "test", AutogeneratePattern: "abc"},
			func() *core.Record {
				return core.NewRecord(collection)
			},
			"abc",
		},
		{
			"existing non-zero value",
			core.InterceptorActionCreate,
			&core.TextField{Name: "test", AutogeneratePattern: "abc"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "123")
				return record
			},
			"123",
		},
		{
			"non-new record",
			core.InterceptorActionValidate,
			&core.TextField{Name: "test", AutogeneratePattern: "abc"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.Id = "test"
				record.PostScan()
				return record
			},
			"",
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			actionCalls := 0
			record := s.record()

			err := s.field.Intercept(context.Background(), app, record, s.actionName, func() error {
				actionCalls++
				return nil
			})
			if err != nil {
				t.Fatal(err)
			}

			if actionCalls != 1 {
				t.Fatalf("Expected actionCalls %d, got %d", 1, actionCalls)
			}

			v := record.GetString(s.field.GetName())
			if v != s.expected {
				t.Fatalf("Expected value %q, got %q", s.expected, v)
			}
		})
	}
}

func TestTextFieldFindSetter(t *testing.T) {
	scenarios := []struct {
		name      string
		key       string
		value     any
		field     *core.TextField
		hasSetter bool
		expected  string
	}{
		{
			"no match",
			"example",
			"abc",
			&core.TextField{Name: "test", AutogeneratePattern: "test"},
			false,
			"",
		},
		{
			"exact match",
			"test",
			"abc",
			&core.TextField{Name: "test", AutogeneratePattern: "test"},
			true,
			"abc",
		},
		{
			"autogenerate modifier",
			"test:autogenerate",
			"abc",
			&core.TextField{Name: "test", AutogeneratePattern: "test"},
			true,
			"abctest",
		},
		{
			"autogenerate modifier without AutogeneratePattern option",
			"test:autogenerate",
			"abc",
			&core.TextField{Name: "test"},
			true,
			"abc",
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

			setter(record, s.value)

			result := record.GetString(s.field.Name)

			if result != s.expected {
				t.Fatalf("Expected %q, got %q", s.expected, result)
			}
		})
	}
}
