package core_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestFileFieldBaseMethods(t *testing.T) {
	testFieldBaseMethods(t, core.FieldTypeFile)
}

func TestFileFieldColumnType(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		name     string
		field    *core.FileField
		expected string
	}{
		{
			"single (zero)",
			&core.FileField{},
			"TEXT DEFAULT '' NOT NULL",
		},
		{
			"single",
			&core.FileField{MaxSelect: 1},
			"TEXT DEFAULT '' NOT NULL",
		},
		{
			"multiple",
			&core.FileField{MaxSelect: 2},
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

func TestFileFieldIsMultiple(t *testing.T) {
	scenarios := []struct {
		name     string
		field    *core.FileField
		expected bool
	}{
		{
			"zero",
			&core.FileField{},
			false,
		},
		{
			"single",
			&core.FileField{MaxSelect: 1},
			false,
		},
		{
			"multiple",
			&core.FileField{MaxSelect: 2},
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

func TestFileFieldPrepareValue(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	record := core.NewRecord(core.NewBaseCollection("test"))

	f1, err := filesystem.NewFileFromBytes([]byte("test"), "test1.txt")
	if err != nil {
		t.Fatal(err)
	}
	f1Raw, err := json.Marshal(f1)
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		raw      any
		field    *core.FileField
		expected string
	}{
		// single
		{nil, &core.FileField{MaxSelect: 1}, `""`},
		{"", &core.FileField{MaxSelect: 1}, `""`},
		{123, &core.FileField{MaxSelect: 1}, `"123"`},
		{"a", &core.FileField{MaxSelect: 1}, `"a"`},
		{`["a"]`, &core.FileField{MaxSelect: 1}, `"a"`},
		{*f1, &core.FileField{MaxSelect: 1}, string(f1Raw)},
		{f1, &core.FileField{MaxSelect: 1}, string(f1Raw)},
		{[]string{}, &core.FileField{MaxSelect: 1}, `""`},
		{[]string{"a", "b"}, &core.FileField{MaxSelect: 1}, `"b"`},

		// multiple
		{nil, &core.FileField{MaxSelect: 2}, `[]`},
		{"", &core.FileField{MaxSelect: 2}, `[]`},
		{123, &core.FileField{MaxSelect: 2}, `["123"]`},
		{"a", &core.FileField{MaxSelect: 2}, `["a"]`},
		{`["a"]`, &core.FileField{MaxSelect: 2}, `["a"]`},
		{[]any{f1}, &core.FileField{MaxSelect: 2}, `[` + string(f1Raw) + `]`},
		{[]*filesystem.File{f1}, &core.FileField{MaxSelect: 2}, `[` + string(f1Raw) + `]`},
		{[]filesystem.File{*f1}, &core.FileField{MaxSelect: 2}, `[` + string(f1Raw) + `]`},
		{[]string{}, &core.FileField{MaxSelect: 2}, `[]`},
		{[]string{"a", "b", "c"}, &core.FileField{MaxSelect: 2}, `["a","b","c"]`},
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

func TestFileFieldDriverValue(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	f1, err := filesystem.NewFileFromBytes([]byte("test"), "test.txt")
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		raw      any
		field    *core.FileField
		expected string
	}{
		// single
		{nil, &core.FileField{MaxSelect: 1}, `""`},
		{"", &core.FileField{MaxSelect: 1}, `""`},
		{123, &core.FileField{MaxSelect: 1}, `"123"`},
		{"a", &core.FileField{MaxSelect: 1}, `"a"`},
		{`["a"]`, &core.FileField{MaxSelect: 1}, `"a"`},
		{f1, &core.FileField{MaxSelect: 1}, `"` + f1.Name + `"`},
		{[]string{}, &core.FileField{MaxSelect: 1}, `""`},
		{[]string{"a", "b"}, &core.FileField{MaxSelect: 1}, `"b"`},

		// multiple
		{nil, &core.FileField{MaxSelect: 2}, `[]`},
		{"", &core.FileField{MaxSelect: 2}, `[]`},
		{123, &core.FileField{MaxSelect: 2}, `["123"]`},
		{"a", &core.FileField{MaxSelect: 2}, `["a"]`},
		{`["a"]`, &core.FileField{MaxSelect: 2}, `["a"]`},
		{[]any{"a", f1}, &core.FileField{MaxSelect: 2}, `["a","` + f1.Name + `"]`},
		{[]string{}, &core.FileField{MaxSelect: 2}, `[]`},
		{[]string{"a", "b", "c"}, &core.FileField{MaxSelect: 2}, `["a","b","c"]`},
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

func TestFileFieldValidateValue(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewBaseCollection("test_collection")

	f1, err := filesystem.NewFileFromBytes([]byte("test"), "test1.txt")
	if err != nil {
		t.Fatal(err)
	}

	f2, err := filesystem.NewFileFromBytes([]byte("test"), "test2.txt")
	if err != nil {
		t.Fatal(err)
	}

	f3, err := filesystem.NewFileFromBytes([]byte("test_abc"), "test3.txt")
	if err != nil {
		t.Fatal(err)
	}

	f4, err := filesystem.NewFileFromBytes(make([]byte, core.DefaultFileFieldMaxSize+1), "test4.txt")
	if err != nil {
		t.Fatal(err)
	}

	f5, err := filesystem.NewFileFromBytes(make([]byte, core.DefaultFileFieldMaxSize), "test5.txt")
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		name        string
		field       *core.FileField
		record      func() *core.Record
		expectError bool
	}{
		// single
		{
			"zero field value (not required)",
			&core.FileField{Name: "test", MaxSize: 9999, MaxSelect: 1},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "")
				return record
			},
			false,
		},
		{
			"zero field value (required)",
			&core.FileField{Name: "test", MaxSize: 9999, MaxSelect: 1, Required: true},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "")
				return record
			},
			true,
		},
		{
			"new plain filename", // new files must be *filesystem.File
			&core.FileField{Name: "test", MaxSize: 9999, MaxSelect: 1},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "a")
				return record
			},
			true,
		},
		{
			"new file",
			&core.FileField{Name: "test", MaxSize: 9999, MaxSelect: 1},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", f1)
				return record
			},
			false,
		},
		{
			"new files > MaxSelect",
			&core.FileField{Name: "test", MaxSize: 9999, MaxSelect: 1},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", []any{f1, f2})
				return record
			},
			true,
		},
		{
			"new files <= MaxSelect",
			&core.FileField{Name: "test", MaxSize: 9999, MaxSelect: 2},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", []any{f1, f2})
				return record
			},
			false,
		},
		{
			"> default MaxSize",
			&core.FileField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", f4)
				return record
			},
			true,
		},
		{
			"<= default MaxSize",
			&core.FileField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", f5)
				return record
			},
			false,
		},
		{
			"> MaxSize",
			&core.FileField{Name: "test", MaxSize: 4, MaxSelect: 3},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", []any{f1, f2, f3}) // f3=8
				return record
			},
			true,
		},
		{
			"<= MaxSize",
			&core.FileField{Name: "test", MaxSize: 8, MaxSelect: 3},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", []any{f1, f2, f3})
				return record
			},
			false,
		},
		{
			"non-matching MimeType",
			&core.FileField{Name: "test", MaxSize: 999, MaxSelect: 3, MimeTypes: []string{"a", "b"}},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", []any{f1, f2})
				return record
			},
			true,
		},
		{
			"matching MimeType",
			&core.FileField{Name: "test", MaxSize: 999, MaxSelect: 3, MimeTypes: []string{"text/plain", "b"}},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", []any{f1, f2})
				return record
			},
			false,
		},
		{
			"existing files > MaxSelect",
			&core.FileField{Name: "file_many", MaxSize: 999, MaxSelect: 2},
			func() *core.Record {
				record, _ := app.FindRecordById("demo1", "84nmscqy84lsi1t") // 5 files
				return record
			},
			true,
		},
		{
			"existing files should ignore the MaxSize and Mimetypes checks",
			&core.FileField{Name: "file_many", MaxSize: 1, MaxSelect: 5, MimeTypes: []string{"a", "b"}},
			func() *core.Record {
				record, _ := app.FindRecordById("demo1", "84nmscqy84lsi1t")
				return record
			},
			false,
		},
		{
			"existing + new file > MaxSelect (5+2)",
			&core.FileField{Name: "file_many", MaxSize: 999, MaxSelect: 6},
			func() *core.Record {
				record, _ := app.FindRecordById("demo1", "84nmscqy84lsi1t")
				record.Set("file_many+", []any{f1, f2})
				return record
			},
			true,
		},
		{
			"existing + new file <= MaxSelect (5+2)",
			&core.FileField{Name: "file_many", MaxSize: 999, MaxSelect: 7},
			func() *core.Record {
				record, _ := app.FindRecordById("demo1", "84nmscqy84lsi1t")
				record.Set("file_many+", []any{f1, f2})
				return record
			},
			false,
		},
		{
			"existing + new filename",
			&core.FileField{Name: "file_many", MaxSize: 999, MaxSelect: 99},
			func() *core.Record {
				record, _ := app.FindRecordById("demo1", "84nmscqy84lsi1t")
				record.Set("file_many+", "test123.png")
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

func TestFileFieldValidateSettings(t *testing.T) {
	testDefaultFieldIdValidation(t, core.FieldTypeFile)
	testDefaultFieldNameValidation(t, core.FieldTypeFile)

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		name         string
		field        func() *core.FileField
		expectErrors []string
	}{
		{
			"zero minimal",
			func() *core.FileField {
				return &core.FileField{
					Id:   "test",
					Name: "test",
				}
			},
			[]string{},
		},
		{
			"0x0 thumb",
			func() *core.FileField {
				return &core.FileField{
					Id:        "test",
					Name:      "test",
					MaxSelect: 1,
					Thumbs:    []string{"100x200", "0x0"},
				}
			},
			[]string{"thumbs"},
		},
		{
			"0x0t thumb",
			func() *core.FileField {
				return &core.FileField{
					Id:        "test",
					Name:      "test",
					MaxSize:   1,
					MaxSelect: 1,
					Thumbs:    []string{"100x200", "0x0t"},
				}
			},
			[]string{"thumbs"},
		},
		{
			"0x0b thumb",
			func() *core.FileField {
				return &core.FileField{
					Id:        "test",
					Name:      "test",
					MaxSize:   1,
					MaxSelect: 1,
					Thumbs:    []string{"100x200", "0x0b"},
				}
			},
			[]string{"thumbs"},
		},
		{
			"0x0f thumb",
			func() *core.FileField {
				return &core.FileField{
					Id:        "test",
					Name:      "test",
					MaxSize:   1,
					MaxSelect: 1,
					Thumbs:    []string{"100x200", "0x0f"},
				}
			},
			[]string{"thumbs"},
		},
		{
			"invalid format",
			func() *core.FileField {
				return &core.FileField{
					Id:        "test",
					Name:      "test",
					MaxSize:   1,
					MaxSelect: 1,
					Thumbs:    []string{"100x200", "100x"},
				}
			},
			[]string{"thumbs"},
		},
		{
			"valid thumbs",
			func() *core.FileField {
				return &core.FileField{
					Id:        "test",
					Name:      "test",
					MaxSize:   1,
					MaxSelect: 1,
					Thumbs:    []string{"100x200", "100x40", "100x200"},
				}
			},
			[]string{},
		},
		{
			"MaxSize > safe json int",
			func() *core.FileField {
				return &core.FileField{
					Id:      "test",
					Name:    "test",
					MaxSize: 1 << 53,
				}
			},
			[]string{"maxSize"},
		},
		{
			"MaxSize < 0",
			func() *core.FileField {
				return &core.FileField{
					Id:      "test",
					Name:    "test",
					MaxSize: -1,
				}
			},
			[]string{"maxSize"},
		},
		{
			"MaxSelect > safe json int",
			func() *core.FileField {
				return &core.FileField{
					Id:        "test",
					Name:      "test",
					MaxSelect: 1 << 53,
				}
			},
			[]string{"maxSelect"},
		},
		{
			"MaxSelect < 0",
			func() *core.FileField {
				return &core.FileField{
					Id:        "test",
					Name:      "test",
					MaxSelect: -1,
				}
			},
			[]string{"maxSelect"},
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

func TestFileFieldCalculateMaxBodySize(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	scenarios := []struct {
		field    *core.FileField
		expected int64
	}{
		{&core.FileField{}, core.DefaultFileFieldMaxSize},
		{&core.FileField{MaxSelect: 2}, 2 * core.DefaultFileFieldMaxSize},
		{&core.FileField{MaxSize: 10}, 10},
		{&core.FileField{MaxSize: 10, MaxSelect: 1}, 10},
		{&core.FileField{MaxSize: 10, MaxSelect: 2}, 20},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%d_%d", i, s.field.MaxSelect, s.field.MaxSize), func(t *testing.T) {
			result := s.field.CalculateMaxBodySize()

			if result != s.expected {
				t.Fatalf("Expected %d, got %d", s.expected, result)
			}
		})
	}
}

func TestFileFieldFindGetter(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	f1, err := filesystem.NewFileFromBytes([]byte("test"), "f1")
	if err != nil {
		t.Fatal(err)
	}
	f1.Name = "f1"

	f2, err := filesystem.NewFileFromBytes([]byte("test"), "f2")
	if err != nil {
		t.Fatal(err)
	}
	f2.Name = "f2"

	record, err := app.FindRecordById("demo3", "lcl9d87w22ml6jy")
	if err != nil {
		t.Fatal(err)
	}
	record.Set("files+", []any{f1, f2})
	record.Set("files-", "test_FLurQTgrY8.txt")

	field, ok := record.Collection().Fields.GetByName("files").(*core.FileField)
	if !ok {
		t.Fatalf("Expected *core.FileField, got %T", record.Collection().Fields.GetByName("files"))
	}

	scenarios := []struct {
		name      string
		key       string
		hasGetter bool
		expected  string
	}{
		{
			"no match",
			"example",
			false,
			"",
		},
		{
			"exact match",
			field.GetName(),
			true,
			`["300_UhLKX91HVb.png",{"name":"f1","originalName":"f1","size":4},{"name":"f2","originalName":"f2","size":4}]`,
		},
		{
			"unsaved",
			field.GetName() + ":unsaved",
			true,
			`[{"name":"f1","originalName":"f1","size":4},{"name":"f2","originalName":"f2","size":4}]`,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			getter := field.FindGetter(s.key)

			hasGetter := getter != nil
			if hasGetter != s.hasGetter {
				t.Fatalf("Expected hasGetter %v, got %v", s.hasGetter, hasGetter)
			}

			if !hasGetter {
				return
			}

			v := getter(record)

			raw, err := json.Marshal(v)
			if err != nil {
				t.Fatal(err)
			}
			rawStr := string(raw)

			if rawStr != s.expected {
				t.Fatalf("Expected\n%v\ngot\n%v", s.expected, rawStr)
			}
		})
	}
}

func TestFileFieldFindSetter(t *testing.T) {
	scenarios := []struct {
		name      string
		key       string
		value     any
		field     *core.FileField
		hasSetter bool
		expected  string
	}{
		{
			"no match",
			"example",
			"b",
			&core.FileField{Name: "test", MaxSelect: 1},
			false,
			"",
		},
		{
			"exact match (single)",
			"test",
			"b",
			&core.FileField{Name: "test", MaxSelect: 1},
			true,
			`"b"`,
		},
		{
			"exact match (multiple)",
			"test",
			[]string{"a", "b", "b"},
			&core.FileField{Name: "test", MaxSelect: 2},
			true,
			`["a","b"]`,
		},
		{
			"append (single)",
			"test+",
			"b",
			&core.FileField{Name: "test", MaxSelect: 1},
			true,
			`"b"`,
		},
		{
			"append (multiple)",
			"test+",
			[]string{"a"},
			&core.FileField{Name: "test", MaxSelect: 2},
			true,
			`["c","d","a"]`,
		},
		{
			"prepend (single)",
			"+test",
			"b",
			&core.FileField{Name: "test", MaxSelect: 1},
			true,
			`"d"`, // the last of the existing values
		},
		{
			"prepend (multiple)",
			"+test",
			[]string{"a"},
			&core.FileField{Name: "test", MaxSelect: 2},
			true,
			`["a","c","d"]`,
		},
		{
			"subtract (single)",
			"test-",
			"d",
			&core.FileField{Name: "test", MaxSelect: 1},
			true,
			`"c"`,
		},
		{
			"subtract (multiple)",
			"test-",
			[]string{"unknown", "c"},
			&core.FileField{Name: "test", MaxSelect: 2},
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

func TestFileFieldIntercept(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	demo1, err := testApp.FindCollectionByNameOrId("demo1")
	if err != nil {
		t.Fatal(err)
	}
	demo1.Fields.GetByName("text").(*core.TextField).Required = true // trigger validation error

	f1, err := filesystem.NewFileFromBytes([]byte("test"), "new1.txt")
	if err != nil {
		t.Fatal(err)
	}

	f2, err := filesystem.NewFileFromBytes([]byte("test"), "new2.txt")
	if err != nil {
		t.Fatal(err)
	}

	f3, err := filesystem.NewFileFromBytes([]byte("test"), "new3.txt")
	if err != nil {
		t.Fatal(err)
	}

	f4, err := filesystem.NewFileFromBytes([]byte("test"), "new4.txt")
	if err != nil {
		t.Fatal(err)
	}

	record := core.NewRecord(demo1)

	ok := t.Run("1. create - with validation error", func(t *testing.T) {
		record.Set("file_many", []any{f1, f2})

		err := testApp.Save(record)

		tests.TestValidationErrors(t, err, []string{"text"})

		value, _ := record.GetRaw("file_many").([]any)
		if len(value) != 2 {
			t.Fatalf("Expected the file field value to be unchanged, got %v", value)
		}
	})
	if !ok {
		return
	}

	ok = t.Run("2. create - fixing the validation error", func(t *testing.T) {
		record.Set("text", "abc")

		err := testApp.Save(record)
		if err != nil {
			t.Fatalf("Expected save to succeed, got %v", err)
		}

		expectedKeys := []string{f1.Name, f2.Name}

		raw := record.GetRaw("file_many")

		// ensure that the value was replaced with the file names
		value := list.ToUniqueStringSlice(raw)
		if len(value) != len(expectedKeys) {
			t.Fatalf("Expected the file field to be updated with the %d file names, got\n%v", len(expectedKeys), raw)
		}
		for _, name := range expectedKeys {
			if !slices.Contains(value, name) {
				t.Fatalf("Missing file %q in %v", name, value)
			}
		}

		checkRecordFiles(t, testApp, record, expectedKeys)
	})
	if !ok {
		return
	}

	ok = t.Run("3. update - validation error", func(t *testing.T) {
		record.Set("text", "")
		record.Set("file_many+", f3)
		record.Set("file_many-", f2.Name)

		err := testApp.Save(record)

		tests.TestValidationErrors(t, err, []string{"text"})

		raw, _ := json.Marshal(record.GetRaw("file_many"))
		expectedRaw, _ := json.Marshal([]any{f1.Name, f3})
		if !bytes.Equal(expectedRaw, raw) {
			t.Fatalf("Expected file field value\n%s\ngot\n%s", expectedRaw, raw)
		}

		checkRecordFiles(t, testApp, record, []string{f1.Name, f2.Name})
	})
	if !ok {
		return
	}

	ok = t.Run("4. update - fixing the validation error", func(t *testing.T) {
		record.Set("text", "abc2")

		err := testApp.Save(record)
		if err != nil {
			t.Fatalf("Expected save to succeed, got %v", err)
		}

		raw, _ := json.Marshal(record.GetRaw("file_many"))
		expectedRaw, _ := json.Marshal([]any{f1.Name, f3.Name})
		if !bytes.Equal(expectedRaw, raw) {
			t.Fatalf("Expected file field value\n%s\ngot\n%s", expectedRaw, raw)
		}

		checkRecordFiles(t, testApp, record, []string{f1.Name, f3.Name})
	})
	if !ok {
		return
	}

	t.Run("5. update - second time update", func(t *testing.T) {
		record.Set("file_many-", f1.Name)
		record.Set("file_many+", f4)

		err := testApp.Save(record)
		if err != nil {
			t.Fatalf("Expected save to succeed, got %v", err)
		}

		raw, _ := json.Marshal(record.GetRaw("file_many"))
		expectedRaw, _ := json.Marshal([]any{f3.Name, f4.Name})
		if !bytes.Equal(expectedRaw, raw) {
			t.Fatalf("Expected file field value\n%s\ngot\n%s", expectedRaw, raw)
		}

		checkRecordFiles(t, testApp, record, []string{f3.Name, f4.Name})
	})
}

func TestFileFieldInterceptTx(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	demo1, err := testApp.FindCollectionByNameOrId("demo1")
	if err != nil {
		t.Fatal(err)
	}
	demo1.Fields.GetByName("text").(*core.TextField).Required = true // trigger validation error

	f1, err := filesystem.NewFileFromBytes([]byte("test"), "new1.txt")
	if err != nil {
		t.Fatal(err)
	}

	f2, err := filesystem.NewFileFromBytes([]byte("test"), "new2.txt")
	if err != nil {
		t.Fatal(err)
	}

	f3, err := filesystem.NewFileFromBytes([]byte("test"), "new3.txt")
	if err != nil {
		t.Fatal(err)
	}

	f4, err := filesystem.NewFileFromBytes([]byte("test"), "new4.txt")
	if err != nil {
		t.Fatal(err)
	}

	var record *core.Record

	tx := func(succeed bool) func(txApp core.App) error {
		var txErr error
		if !succeed {
			txErr = errors.New("tx error")
		}

		return func(txApp core.App) error {
			record = core.NewRecord(demo1)
			ok := t.Run(fmt.Sprintf("[tx_%v] create with validation error", succeed), func(t *testing.T) {
				record.Set("text", "")
				record.Set("file_many", []any{f1, f2})

				err := txApp.Save(record)
				tests.TestValidationErrors(t, err, []string{"text"})

				checkRecordFiles(t, txApp, record, []string{}) // no uploaded files
			})
			if !ok {
				return txErr
			}

			// ---

			ok = t.Run(fmt.Sprintf("[tx_%v] create with fixed validation error", succeed), func(t *testing.T) {
				record.Set("text", "abc")

				err = txApp.Save(record)
				if err != nil {
					t.Fatalf("Expected save to succeed, got %v", err)
				}

				checkRecordFiles(t, txApp, record, []string{f1.Name, f2.Name})
			})
			if !ok {
				return txErr
			}

			// ---

			ok = t.Run(fmt.Sprintf("[tx_%v] update with validation error", succeed), func(t *testing.T) {
				record.Set("text", "")
				record.Set("file_many+", f3)
				record.Set("file_many-", f2.Name)

				err = txApp.Save(record)
				tests.TestValidationErrors(t, err, []string{"text"})

				raw, _ := json.Marshal(record.GetRaw("file_many"))
				expectedRaw, _ := json.Marshal([]any{f1.Name, f3})
				if !bytes.Equal(expectedRaw, raw) {
					t.Fatalf("Expected file field value\n%s\ngot\n%s", expectedRaw, raw)
				}

				checkRecordFiles(t, txApp, record, []string{f1.Name, f2.Name}) // no file changes
			})
			if !ok {
				return txErr
			}

			// ---

			ok = t.Run(fmt.Sprintf("[tx_%v] update with fixed validation error", succeed), func(t *testing.T) {
				record.Set("text", "abc2")

				err = txApp.Save(record)
				if err != nil {
					t.Fatalf("Expected save to succeed, got %v", err)
				}

				raw, _ := json.Marshal(record.GetRaw("file_many"))
				expectedRaw, _ := json.Marshal([]any{f1.Name, f3.Name})
				if !bytes.Equal(expectedRaw, raw) {
					t.Fatalf("Expected file field value\n%s\ngot\n%s", expectedRaw, raw)
				}

				checkRecordFiles(t, txApp, record, []string{f1.Name, f3.Name, f2.Name}) // f2 shouldn't be deleted yet
			})
			if !ok {
				return txErr
			}

			// ---

			ok = t.Run(fmt.Sprintf("[tx_%v] second time update", succeed), func(t *testing.T) {
				record.Set("file_many-", f1.Name)
				record.Set("file_many+", f4)

				err := txApp.Save(record)
				if err != nil {
					t.Fatalf("Expected save to succeed, got %v", err)
				}

				raw, _ := json.Marshal(record.GetRaw("file_many"))
				expectedRaw, _ := json.Marshal([]any{f3.Name, f4.Name})
				if !bytes.Equal(expectedRaw, raw) {
					t.Fatalf("Expected file field value\n%s\ngot\n%s", expectedRaw, raw)
				}

				checkRecordFiles(t, txApp, record, []string{f3.Name, f4.Name, f1.Name, f2.Name}) // f1 and f2 shouldn't be deleted yet
			})
			if !ok {
				return txErr
			}

			// ---

			return txErr
		}
	}

	// failed transaction
	txErr := testApp.RunInTransaction(tx(false))
	if txErr == nil {
		t.Fatal("Expected transaction error")
	}
	// there shouldn't be any fails associated with the record id
	checkRecordFiles(t, testApp, record, []string{})

	txErr = testApp.RunInTransaction(tx(true))
	if txErr != nil {
		t.Fatalf("Expected transaction to succeed, got %v", txErr)
	}
	// only the last updated files should remain
	checkRecordFiles(t, testApp, record, []string{f3.Name, f4.Name})
}

// -------------------------------------------------------------------

func checkRecordFiles(t *testing.T, testApp core.App, record *core.Record, expectedKeys []string) {
	fsys, err := testApp.NewFilesystem()
	if err != nil {
		t.Fatal(err)
	}
	defer fsys.Close()

	objects, err := fsys.List(record.BaseFilesPath() + "/")
	if err != nil {
		t.Fatal(err)
	}
	objectKeys := make([]string, 0, len(objects))
	for _, obj := range objects {
		// exclude thumbs
		if !strings.Contains(obj.Key, "/thumbs_") {
			objectKeys = append(objectKeys, obj.Key)
		}
	}

	if len(objectKeys) != len(expectedKeys) {
		t.Fatalf("Expected files:\n%v\ngot\n%v", expectedKeys, objectKeys)
	}
	for _, key := range expectedKeys {
		fullKey := record.BaseFilesPath() + "/" + key
		if !slices.Contains(objectKeys, fullKey) {
			t.Fatalf("Missing expected file key\n%q\nin\n%v", fullKey, objectKeys)
		}
	}
}
