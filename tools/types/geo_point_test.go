package types_test

import (
	"fmt"
	"testing"

	"github.com/pocketbase/pocketbase/tools/types"
)

func TestGeoPointAsMap(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		name     string
		point    types.GeoPoint
		expected map[string]any
	}{
		{"zero", types.GeoPoint{}, map[string]any{"lon": 0.0, "lat": 0.0}},
		{"non-zero", types.GeoPoint{Lon: -10, Lat: 20.123}, map[string]any{"lon": -10.0, "lat": 20.123}},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			result := s.point.AsMap()

			if len(result) != len(s.expected) {
				t.Fatalf("Expected %d keys, got %d: %v", len(s.expected), len(result), result)
			}

			for k, v := range s.expected {
				found, ok := result[k]
				if !ok {
					t.Fatalf("Missing expected %q key: %v", k, result)
				}

				if found != v {
					t.Fatalf("Expected %q key value %v, got %v", k, v, found)
				}
			}
		})
	}
}

func TestGeoPointStringAndValue(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		name     string
		point    types.GeoPoint
		expected string
	}{
		{"zero", types.GeoPoint{}, `{"lon":0,"lat":0}`},
		{"non-zero", types.GeoPoint{Lon: -10, Lat: 20.123}, `{"lon":-10,"lat":20.123}`},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			str := s.point.String()

			val, err := s.point.Value()
			if err != nil {
				t.Fatal(err)
			}

			if str != val {
				t.Fatalf("Expected String and Value to return the same value")
			}

			if str != s.expected {
				t.Fatalf("Expected\n%s\ngot\n%s", s.expected, str)
			}
		})
	}
}

func TestGeoPointScan(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		value     any
		expectErr bool
		expectStr string
	}{
		{nil, false, `{"lon":1,"lat":2}`},
		{"", false, `{"lon":1,"lat":2}`},
		{types.JSONRaw{}, false, `{"lon":1,"lat":2}`},
		{[]byte{}, false, `{"lon":1,"lat":2}`},
		{`{}`, false, `{"lon":1,"lat":2}`},
		{`[]`, true, `{"lon":1,"lat":2}`},
		{0, true, `{"lon":1,"lat":2}`},
		{`{"lon":"1.23","lat":"4.56"}`, true, `{"lon":1,"lat":2}`},
		{`{"lon":1.23,"lat":4.56}`, false, `{"lon":1.23,"lat":4.56}`},
		{[]byte(`{"lon":1.23,"lat":4.56}`), false, `{"lon":1.23,"lat":4.56}`},
		{types.JSONRaw(`{"lon":1.23,"lat":4.56}`), false, `{"lon":1.23,"lat":4.56}`},
		{types.GeoPoint{}, false, `{"lon":0,"lat":0}`},
		{types.GeoPoint{Lon: 1.23, Lat: 4.56}, false, `{"lon":1.23,"lat":4.56}`},
		{&types.GeoPoint{Lon: 1.23, Lat: 4.56}, false, `{"lon":1.23,"lat":4.56}`},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#v", i, s.value), func(t *testing.T) {
			point := types.GeoPoint{Lon: 1, Lat: 2}

			err := point.Scan(s.value)

			hasErr := err != nil
			if hasErr != s.expectErr {
				t.Errorf("Expected hasErr %v, got %v (%v)", s.expectErr, hasErr, err)
			}

			if str := point.String(); str != s.expectStr {
				t.Errorf("Expected\n%s\ngot\n%s", s.expectStr, str)
			}
		})
	}
}
