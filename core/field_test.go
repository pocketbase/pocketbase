package core_test

import (
	"context"
	"strings"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

func testFieldBaseMethods(t *testing.T, fieldType string) {
	factory, ok := core.Fields[fieldType]
	if !ok {
		t.Fatalf("Missing %q field factory", fieldType)
	}

	f := factory()
	if f == nil {
		t.Fatal("Expected non-nil Field instance")
	}

	t.Run("type", func(t *testing.T) {
		if v := f.Type(); v != fieldType {
			t.Fatalf("Expected type %q, got %q", fieldType, v)
		}
	})

	t.Run("id", func(t *testing.T) {
		testValues := []string{"new_id", ""}
		for _, expected := range testValues {
			f.SetId(expected)
			if v := f.GetId(); v != expected {
				t.Fatalf("Expected id %q, got %q", expected, v)
			}
		}
	})

	t.Run("name", func(t *testing.T) {
		testValues := []string{"new_name", ""}
		for _, expected := range testValues {
			f.SetName(expected)
			if v := f.GetName(); v != expected {
				t.Fatalf("Expected name %q, got %q", expected, v)
			}
		}
	})

	t.Run("system", func(t *testing.T) {
		testValues := []bool{false, true}
		for _, expected := range testValues {
			f.SetSystem(expected)
			if v := f.GetSystem(); v != expected {
				t.Fatalf("Expected system %v, got %v", expected, v)
			}
		}
	})

	t.Run("hidden", func(t *testing.T) {
		testValues := []bool{false, true}
		for _, expected := range testValues {
			f.SetHidden(expected)
			if v := f.GetHidden(); v != expected {
				t.Fatalf("Expected hidden %v, got %v", expected, v)
			}
		}
	})
}

func testDefaultFieldIdValidation(t *testing.T, fieldType string) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewBaseCollection("test_collection")

	scenarios := []struct {
		name        string
		field       func() core.Field
		expectError bool
	}{
		{
			"empty value",
			func() core.Field {
				f := core.Fields[fieldType]()
				return f
			},
			true,
		},
		{
			"invalid length",
			func() core.Field {
				f := core.Fields[fieldType]()
				f.SetId(strings.Repeat("a", 101))
				return f
			},
			true,
		},
		{
			"valid length",
			func() core.Field {
				f := core.Fields[fieldType]()
				f.SetId(strings.Repeat("a", 100))
				return f
			},
			false,
		},
	}

	for _, s := range scenarios {
		t.Run("[id] "+s.name, func(t *testing.T) {
			errs, _ := s.field().ValidateSettings(context.Background(), app, collection).(validation.Errors)

			hasErr := errs["id"] != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr %v, got %v", s.expectError, hasErr)
			}
		})
	}
}

func testDefaultFieldNameValidation(t *testing.T, fieldType string) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewBaseCollection("test_collection")

	scenarios := []struct {
		name        string
		field       func() core.Field
		expectError bool
	}{
		{
			"empty value",
			func() core.Field {
				f := core.Fields[fieldType]()
				return f
			},
			true,
		},
		{
			"invalid length",
			func() core.Field {
				f := core.Fields[fieldType]()
				f.SetName(strings.Repeat("a", 101))
				return f
			},
			true,
		},
		{
			"valid length",
			func() core.Field {
				f := core.Fields[fieldType]()
				f.SetName(strings.Repeat("a", 100))
				return f
			},
			false,
		},
		{
			"invalid regex",
			func() core.Field {
				f := core.Fields[fieldType]()
				f.SetName("test(")
				return f
			},
			true,
		},
		{
			"valid regex",
			func() core.Field {
				f := core.Fields[fieldType]()
				f.SetName("test_123")
				return f
			},
			false,
		},
		{
			"_via_",
			func() core.Field {
				f := core.Fields[fieldType]()
				f.SetName("a_via_b")
				return f
			},
			true,
		},
		{
			"system reserved - null",
			func() core.Field {
				f := core.Fields[fieldType]()
				f.SetName("null")
				return f
			},
			true,
		},
		{
			"system reserved - false",
			func() core.Field {
				f := core.Fields[fieldType]()
				f.SetName("false")
				return f
			},
			true,
		},
		{
			"system reserved - true",
			func() core.Field {
				f := core.Fields[fieldType]()
				f.SetName("true")
				return f
			},
			true,
		},
		{
			"system reserved - _rowid_",
			func() core.Field {
				f := core.Fields[fieldType]()
				f.SetName("_rowid_")
				return f
			},
			true,
		},
		{
			"system reserved - expand",
			func() core.Field {
				f := core.Fields[fieldType]()
				f.SetName("expand")
				return f
			},
			true,
		},
		{
			"system reserved - collectionId",
			func() core.Field {
				f := core.Fields[fieldType]()
				f.SetName("collectionId")
				return f
			},
			true,
		},
		{
			"system reserved - collectionName",
			func() core.Field {
				f := core.Fields[fieldType]()
				f.SetName("collectionName")
				return f
			},
			true,
		},
	}

	for _, s := range scenarios {
		t.Run("[name] "+s.name, func(t *testing.T) {
			errs, _ := s.field().ValidateSettings(context.Background(), app, collection).(validation.Errors)

			hasErr := errs["name"] != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr %v, got %v", s.expectError, hasErr)
			}
		})
	}
}
