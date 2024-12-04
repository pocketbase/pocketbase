package types_test

import (
	"fmt"
	"testing"

	"github.com/pocketbase/pocketbase/tools/types"
)

func TestParsePoint(t *testing.T) {
	scenarios := []struct {
		value    any
		expected string
	}{
		{nil, ""},
		{"", ""},
		{"invalid", ""},
		{"42.3631, -71.0574", "42.3631, -71.0574"},
		{types.Point{}, ""},
		{"38.8977,-77.0365", "38.8977, -77.0365"},
		{[]byte("121.141414,-8.12"), "121.141414, -8.12"},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#v", i, s.value), func(t *testing.T) {
			p, err := types.ParsePoint(s.value)
			if err != nil && s.expected != "" {
				t.Fatalf("Failed to parse %v: %v", s.value, err)
			}
			if p.String() != s.expected {
				t.Fatalf("Expected %q, got %q", s.expected, p.String())
			}
		})
	}
}

func TestPointLatLong(t *testing.T) {
	t.Parallel()

	lat := 42.3631
	long := -71.0574
	p, _ := types.ParsePoint(fmt.Sprintf("%f, %f", lat, long))

	if p.Lat() != lat {
		t.Fatalf("Expected lat %f, got %f", lat, p.Lat())
	}

	if p.Long() != long {
		t.Fatalf("Expected long %f, got %f", long, p.Long())
	}
}

func TestPointEqual(t *testing.T) {
	t.Parallel()

	p1, _ := types.ParsePoint("42.3631, -71.0574")
	p2, _ := types.ParsePoint("42.3631, -71.0574")
	p3, _ := types.ParsePoint("42.3632, -71.0574")
	p4 := types.Point{}

	scenarios := []struct {
		a      types.Point
		b      types.Point
		expect bool
	}{
		{p1, p1, true},
		{p1, p2, true},
		{p1, p3, false},
		{p4, p4, true},
		{p1, p4, false},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("equal_%d", i), func(t *testing.T) {
			if v := s.a.Equal(s.b); v != s.expect {
				t.Fatalf("Expected %v, got %v", s.expect, v)
			}
		})
	}
}

func TestPointString(t *testing.T) {
	p0 := types.Point{}
	if p0.String() != "" {
		t.Fatalf("Expected empty string for zero point, got %q", p0.String())
	}

	expected := "42.3631, -71.0574"
	p1, _ := types.ParsePoint(expected)
	if p1.String() != expected {
		t.Fatalf("Expected %q, got %v", expected, p1)
	}
}

func TestPointMarshalJSON(t *testing.T) {
	scenarios := []struct {
		point    string
		expected string
	}{
		{"", `""`},
		{"42.3631, -71.0574", `"42.3631, -71.0574"`},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%s", i, s.point), func(t *testing.T) {
			p, err := types.ParsePoint(s.point)
			if err != nil {
				t.Fatal(err)
			}

			result, err := p.MarshalJSON()
			if err != nil {
				t.Fatal(err)
			}

			if string(result) != s.expected {
				t.Fatalf("Expected %q, got %q", s.expected, string(result))
			}
		})
	}
}

func TestPointUnmarshalJSON(t *testing.T) {
	scenarios := []struct {
		json     string
		expected string
	}{
		{"", ""},
		{"invalid_json", ""},
		{"'123'", ""},
		{`"42.3631, -71.0574"`, "42.3631, -71.0574"},
		{`""`, ""},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%s", i, s.json), func(t *testing.T) {
			p := types.Point{}
			p.UnmarshalJSON([]byte(s.json))

			if p.String() != s.expected {
				t.Fatalf("Expected %q, got %q", s.expected, p.String())
			}
		})
	}
}

func TestPointValue(t *testing.T) {
	scenarios := []struct {
		value    any
		expected string
	}{
		{"", ""},
		{"invalid", ""},
		{"42.3631, -71.0574", "42.3631, -71.0574"},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%v", i, s.value), func(t *testing.T) {
			p, _ := types.ParsePoint(s.value)

			result, err := p.Value()
			if err != nil {
				t.Fatal(err)
			}

			if result != s.expected {
				t.Fatalf("Expected %q, got %q", s.expected, result)
			}
		})
	}
}
