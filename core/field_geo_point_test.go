package core_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestGeoPointFieldBaseMethods(t *testing.T) {
	testFieldBaseMethods(t, core.FieldTypeGeoPoint)
}

func TestGeoPointFieldColumnType(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	f := &core.GeoPointField{}

	expected := `JSON DEFAULT '{"lon":0,"lat":0}' NOT NULL`

	if v := f.ColumnType(app); v != expected {
		t.Fatalf("Expected\n%q\ngot\n%q", expected, v)
	}
}

func TestGeoPointFieldPrepareValue(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	f := &core.GeoPointField{}
	record := core.NewRecord(core.NewBaseCollection("test"))

	scenarios := []struct {
		raw      any
		expected string
	}{
		{nil, `{"lon":0,"lat":0}`},
		{"", `{"lon":0,"lat":0}`},
		{[]byte{}, `{"lon":0,"lat":0}`},
		{map[string]any{}, `{"lon":0,"lat":0}`},
		{types.GeoPoint{Lon: 10, Lat: 20}, `{"lon":10,"lat":20}`},
		{&types.GeoPoint{Lon: 10, Lat: 20}, `{"lon":10,"lat":20}`},
		{[]byte(`{"lon": 10, "lat": 20}`), `{"lon":10,"lat":20}`},
		{map[string]any{"lon": 10, "lat": 20}, `{"lon":10,"lat":20}`},
		{map[string]float64{"lon": 10, "lat": 20}, `{"lon":10,"lat":20}`},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#v", i, s.raw), func(t *testing.T) {
			v, err := f.PrepareValue(record, s.raw)
			if err != nil {
				t.Fatal(err)
			}

			raw, err := json.Marshal(v)
			if err != nil {
				t.Fatal(err)
			}
			rawStr := string(raw)

			if rawStr != s.expected {
				t.Fatalf("Expected\n%s\ngot\n%s", s.expected, rawStr)
			}
		})
	}
}

func TestGeoPointFieldValidateValue(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewBaseCollection("test_collection")

	scenarios := []struct {
		name        string
		field       *core.GeoPointField
		record      func() *core.Record
		expectError bool
	}{
		{
			"invalid raw value",
			&core.GeoPointField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", 123)
				return record
			},
			true,
		},
		{
			"zero field value (non-required)",
			&core.GeoPointField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", types.GeoPoint{})
				return record
			},
			false,
		},
		{
			"zero field value (required)",
			&core.GeoPointField{Name: "test", Required: true},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", types.GeoPoint{})
				return record
			},
			true,
		},
		{
			"non-zero Lat field value (required)",
			&core.GeoPointField{Name: "test", Required: true},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", types.GeoPoint{Lat: 1})
				return record
			},
			false,
		},
		{
			"non-zero Lon field value (required)",
			&core.GeoPointField{Name: "test", Required: true},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", types.GeoPoint{Lon: 1})
				return record
			},
			false,
		},
		{
			"non-zero Lat-Lon field value (required)",
			&core.GeoPointField{Name: "test", Required: true},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", types.GeoPoint{Lon: -1, Lat: -2})
				return record
			},
			false,
		},
		{
			"lat < -90",
			&core.GeoPointField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", types.GeoPoint{Lat: -90.1})
				return record
			},
			true,
		},
		{
			"lat > 90",
			&core.GeoPointField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", types.GeoPoint{Lat: 90.1})
				return record
			},
			true,
		},
		{
			"lon < -180",
			&core.GeoPointField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", types.GeoPoint{Lon: -180.1})
				return record
			},
			true,
		},
		{
			"lon > 180",
			&core.GeoPointField{Name: "test"},
			func() *core.Record {
				record := core.NewRecord(collection)
				record.SetRaw("test", types.GeoPoint{Lon: 180.1})
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

func TestGeoPointFieldValidateSettings(t *testing.T) {
	testDefaultFieldIdValidation(t, core.FieldTypeGeoPoint)
	testDefaultFieldNameValidation(t, core.FieldTypeGeoPoint)
}
