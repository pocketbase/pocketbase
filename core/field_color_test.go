package core_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

func TestColorFieldBaseMethods(t *testing.T) {
	testFieldBaseMethods(t, core.FieldTypeColor)
}

func TestColorFieldColumnType(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	f := &core.ColorField{}

	expected := "TEXT DEFAULT '' NOT NULL"

	if v := f.ColumnType(app); v != expected {
		t.Fatalf("Expected\n%q\ngot\n%q", expected, v)
	}
}

func TestColorFieldPrepareValue(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	f := &core.ColorField{}
	record := core.NewRecord(core.NewBaseCollection("test"))

	scenarios := []struct {
		raw      any
		expected string
	}{
		{"", ""},
		{"#ff0000", "#ff0000"},
		{"  #abc  ", "#abc"},
		{"rgb(255, 0, 0)", "rgb(255, 0, 0)"},
		{"  hsl(120deg, 100%, 50%)  ", "hsl(120deg, 100%, 50%)"},
		{123, "123"},
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

func TestColorFieldValidateValue(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewBaseCollection("test_collection")

	scenarios := []struct {
		name        string
		field       *core.ColorField
		record      func() *core.Record
		expectError bool
	}{
		{
			"invalid raw value",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", 123)
				return record
			},
			true,
		},
		{
			"zero field value (not required)",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "")
				return record
			},
			false,
		},
		{
			"zero field value (required)",
			&core.ColorField{Name: "test", Required: true},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "")
				return record
			},
			true,
		},
		// Hex color tests
		{
			"valid hex 3-digit",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "#f00")
				return record
			},
			false,
		},
		{
			"valid hex 6-digit",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "#ff0000")
				return record
			},
			false,
		},
		{
			"valid hex 8-digit with alpha",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "#ff0000ff")
				return record
			},
			false,
		},
		{
			"valid hex 4-digit with alpha",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "#f00f")
				return record
			},
			false,
		},
		{
			"invalid hex (5 digits)",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "#ff000")
				return record
			},
			true,
		},
		{
			"invalid hex (missing #)",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "ff0000")
				return record
			},
			true,
		},
		{
			"invalid hex (invalid characters)",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "#gggggg")
				return record
			},
			true,
		},
		// RGB color tests
		{
			"valid rgb",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "rgb(255, 0, 0)")
				return record
			},
			false,
		},
		{
			"valid rgb without commas (modern syntax)",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "rgb(255 0 0)")
				return record
			},
			false,
		},
		{
			"valid rgba",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "rgba(255, 0, 0, 0.5)")
				return record
			},
			false,
		},
		{
			"valid rgba with slash (modern syntax)",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "rgb(255 0 0 / 0.5)")
				return record
			},
			false,
		},
		{
			"invalid rgb (out of range)",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "rgb(256, 0, 0)")
				return record
			},
			true,
		},
		// HSL color tests
		{
			"valid hsl",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "hsl(120, 100%, 50%)")
				return record
			},
			false,
		},
		{
			"valid hsl with deg",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "hsl(120deg, 100%, 50%)")
				return record
			},
			false,
		},
		{
			"valid hsl without commas (modern syntax)",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "hsl(120deg 100% 50%)")
				return record
			},
			false,
		},
		{
			"valid hsla",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "hsla(120, 100%, 50%, 0.5)")
				return record
			},
			false,
		},
		{
			"valid hsla with slash (modern syntax)",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "hsl(120deg 100% 50% / 0.5)")
				return record
			},
			false,
		},
		{
			"invalid hsl (missing %)",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "hsl(120, 100, 50)")
				return record
			},
			true,
		},
		// HWB color tests
		{
			"valid hwb",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "hwb(120deg 30% 50%)")
				return record
			},
			false,
		},
		{
			"valid hwb with alpha",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "hwb(120deg 30% 50% / 0.5)")
				return record
			},
			false,
		},
		{
			"invalid hwb (missing %)",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "hwb(120deg 30 50)")
				return record
			},
			true,
		},
		// LAB color tests
		{
			"valid lab",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "lab(50% 40 30)")
				return record
			},
			false,
		},
		{
			"valid lab with alpha",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "lab(50% 40 30 / 0.5)")
				return record
			},
			false,
		},
		{
			"valid lab with negative values",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "lab(50% -40 -30)")
				return record
			},
			false,
		},
		// LCH color tests
		{
			"valid lch",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "lch(50% 40 120deg)")
				return record
			},
			false,
		},
		{
			"valid lch with alpha",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "lch(50% 40 120deg / 0.5)")
				return record
			},
			false,
		},
		// OKLAB color tests
		{
			"valid oklab",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "oklab(0.5 0.1 0.1)")
				return record
			},
			false,
		},
		{
			"valid oklab with alpha",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "oklab(0.5 0.1 0.1 / 0.5)")
				return record
			},
			false,
		},
		{
			"valid oklab with negative values",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "oklab(0.5 -0.1 -0.1)")
				return record
			},
			false,
		},
		// OKLCH color tests
		{
			"valid oklch",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "oklch(0.5 0.1 120deg)")
				return record
			},
			false,
		},
		{
			"valid oklch with alpha",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "oklch(0.5 0.1 120deg / 0.5)")
				return record
			},
			false,
		},
		{
			"valid oklch without deg",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "oklch(0.5 0.1 120)")
				return record
			},
			false,
		},
		// Named color tests
		{
			"valid named color - red",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "red")
				return record
			},
			false,
		},
		{
			"valid named color - rebeccapurple",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "rebeccapurple")
				return record
			},
			false,
		},
		{
			"valid named color - transparent",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "transparent")
				return record
			},
			false,
		},
		{
			"valid named color - case insensitive",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "DarkSlateGray")
				return record
			},
			false,
		},
		{
			"invalid named color",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "notacolor")
				return record
			},
			true,
		},
		// Edge cases
		{
			"color with extra whitespace",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "  rgb(255, 0, 0)  ")
				return record
			},
			false,
		},
		{
			"completely invalid format",
			&core.ColorField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", "this is not a color")
				return record
			},
			true,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			record := s.record()
			err := s.field.ValidateValue(context.Background(), app, record)

			hasErr := err != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr %v, got %v (%v)", s.expectError, hasErr, err)
			}
		})
	}
}

func TestColorFieldValidateSettings(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		name         string
		field        func() *core.ColorField
		expectErrors []string
	}{
		{
			"minimal valid",
			func() *core.ColorField {
				return &core.ColorField{
					Id:   "test",
					Name: "test",
				}
			},
			[]string{},
		},
		{
			"with required",
			func() *core.ColorField {
				return &core.ColorField{
					Id:       "test",
					Name:     "test",
					Required: true,
				}
			},
			[]string{},
		},
		{
			"invalid name",
			func() *core.ColorField {
				return &core.ColorField{
					Id:   "test",
					Name: "",
				}
			},
			[]string{"name"},
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

func TestColorFieldFindSetter(t *testing.T) {
	scenarios := []struct {
		name      string
		key       string
		value     any
		field     *core.ColorField
		hasSetter bool
		expected  string
	}{
		{
			"no match",
			"example",
			"#ff0000",
			&core.ColorField{Name: "test"},
			false,
			"",
		},
		{
			"exact match",
			"test",
			"#ff0000",
			&core.ColorField{Name: "test"},
			true,
			"#ff0000",
		},
		{
			"exact match with whitespace",
			"test",
			"  rgb(255, 0, 0)  ",
			&core.ColorField{Name: "test"},
			true,
			"rgb(255, 0, 0)",
		},
		{
			"exact match with named color",
			"test",
			"red",
			&core.ColorField{Name: "test"},
			true,
			"red",
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