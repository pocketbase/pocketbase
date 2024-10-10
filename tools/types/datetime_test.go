package types_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/pocketbase/pocketbase/tools/types"
)

func TestNowDateTime(t *testing.T) {
	now := time.Now().UTC().Format("2006-01-02 15:04:05") // without ms part for test consistency
	dt := types.NowDateTime()

	if !strings.Contains(dt.String(), now) {
		t.Fatalf("Expected %q, got %q", now, dt.String())
	}
}

func TestParseDateTime(t *testing.T) {
	nowTime := time.Now().UTC()
	nowDateTime, _ := types.ParseDateTime(nowTime)
	nowStr := nowTime.Format(types.DefaultDateLayout)

	scenarios := []struct {
		value    any
		expected string
	}{
		{nil, ""},
		{"", ""},
		{"invalid", ""},
		{nowDateTime, nowStr},
		{nowTime, nowStr},
		{1641024040, "2022-01-01 08:00:40.000Z"},
		{int32(1641024040), "2022-01-01 08:00:40.000Z"},
		{int64(1641024040), "2022-01-01 08:00:40.000Z"},
		{uint(1641024040), "2022-01-01 08:00:40.000Z"},
		{uint64(1641024040), "2022-01-01 08:00:40.000Z"},
		{uint32(1641024040), "2022-01-01 08:00:40.000Z"},
		{"2022-01-01 11:23:45.678", "2022-01-01 11:23:45.678Z"},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#v", i, s.value), func(t *testing.T) {
			dt, err := types.ParseDateTime(s.value)
			if err != nil {
				t.Fatalf("Failed to parse %v: %v", s.value, err)
			}

			if dt.String() != s.expected {
				t.Fatalf("Expected %q, got %q", s.expected, dt.String())
			}
		})
	}
}

func TestDateTimeTime(t *testing.T) {
	str := "2022-01-01 11:23:45.678Z"

	expected, err := time.Parse(types.DefaultDateLayout, str)
	if err != nil {
		t.Fatal(err)
	}

	dt, err := types.ParseDateTime(str)
	if err != nil {
		t.Fatal(err)
	}

	result := dt.Time()

	if !expected.Equal(result) {
		t.Fatalf("Expected time %v, got %v", expected, result)
	}
}

func TestDateTimeAdd(t *testing.T) {
	t.Parallel()

	d1, _ := types.ParseDateTime("2024-01-01 10:00:00.123Z")

	d2 := d1.Add(1 * time.Hour)

	if d1.String() != "2024-01-01 10:00:00.123Z" {
		t.Fatalf("Expected d1 to remain unchanged, got %s", d1.String())
	}

	expected := "2024-01-01 11:00:00.123Z"
	if d2.String() != expected {
		t.Fatalf("Expected d2 %s, got %s", expected, d2.String())
	}
}

func TestDateTimeSub(t *testing.T) {
	t.Parallel()

	d1, _ := types.ParseDateTime("2024-01-01 10:00:00.123Z")
	d2, _ := types.ParseDateTime("2024-01-01 10:30:00.123Z")

	result := d2.Sub(d1)

	if result.Minutes() != 30 {
		t.Fatalf("Expected %v minutes diff, got %v", 30, result.Minutes())
	}
}

func TestDateTimeAddDate(t *testing.T) {
	t.Parallel()

	d1, _ := types.ParseDateTime("2024-01-01 10:00:00.123Z")

	d2 := d1.AddDate(1, 2, 3)

	if d1.String() != "2024-01-01 10:00:00.123Z" {
		t.Fatalf("Expected d1 to remain unchanged, got %s", d1.String())
	}

	expected := "2025-03-04 10:00:00.123Z"
	if d2.String() != expected {
		t.Fatalf("Expected d2 %s, got %s", expected, d2.String())
	}
}

func TestDateTimeAfter(t *testing.T) {
	t.Parallel()

	d1, _ := types.ParseDateTime("2024-01-01 10:00:00.123Z")
	d2, _ := types.ParseDateTime("2024-01-02 10:00:00.123Z")
	d3, _ := types.ParseDateTime("2024-01-03 10:00:00.123Z")

	scenarios := []struct {
		a      types.DateTime
		b      types.DateTime
		expect bool
	}{
		// d1
		{d1, d1, false},
		{d1, d2, false},
		{d1, d3, false},
		// d2
		{d2, d1, true},
		{d2, d2, false},
		{d2, d3, false},
		// d3
		{d3, d1, true},
		{d3, d2, true},
		{d3, d3, false},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("after_%d", i), func(t *testing.T) {
			if v := s.a.After(s.b); v != s.expect {
				t.Fatalf("Expected %v, got %v", s.expect, v)
			}
		})
	}
}

func TestDateTimeBefore(t *testing.T) {
	t.Parallel()

	d1, _ := types.ParseDateTime("2024-01-01 10:00:00.123Z")
	d2, _ := types.ParseDateTime("2024-01-02 10:00:00.123Z")
	d3, _ := types.ParseDateTime("2024-01-03 10:00:00.123Z")

	scenarios := []struct {
		a      types.DateTime
		b      types.DateTime
		expect bool
	}{
		// d1
		{d1, d1, false},
		{d1, d2, true},
		{d1, d3, true},
		// d2
		{d2, d1, false},
		{d2, d2, false},
		{d2, d3, true},
		// d3
		{d3, d1, false},
		{d3, d2, false},
		{d3, d3, false},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("before_%d", i), func(t *testing.T) {
			if v := s.a.Before(s.b); v != s.expect {
				t.Fatalf("Expected %v, got %v", s.expect, v)
			}
		})
	}
}

func TestDateTimeCompare(t *testing.T) {
	t.Parallel()

	d1, _ := types.ParseDateTime("2024-01-01 10:00:00.123Z")
	d2, _ := types.ParseDateTime("2024-01-02 10:00:00.123Z")
	d3, _ := types.ParseDateTime("2024-01-03 10:00:00.123Z")

	scenarios := []struct {
		a      types.DateTime
		b      types.DateTime
		expect int
	}{
		// d1
		{d1, d1, 0},
		{d1, d2, -1},
		{d1, d3, -1},
		// d2
		{d2, d1, 1},
		{d2, d2, 0},
		{d2, d3, -1},
		// d3
		{d3, d1, 1},
		{d3, d2, 1},
		{d3, d3, 0},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("compare_%d", i), func(t *testing.T) {
			if v := s.a.Compare(s.b); v != s.expect {
				t.Fatalf("Expected %v, got %v", s.expect, v)
			}
		})
	}
}

func TestDateTimeEqual(t *testing.T) {
	t.Parallel()

	d1, _ := types.ParseDateTime("2024-01-01 10:00:00.123Z")
	d2, _ := types.ParseDateTime("2024-01-01 10:00:00.123Z")
	d3, _ := types.ParseDateTime("2024-01-01 10:00:00.124Z")

	scenarios := []struct {
		a      types.DateTime
		b      types.DateTime
		expect bool
	}{
		{d1, d1, true},
		{d1, d2, true},
		{d1, d3, false},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("equal_%d", i), func(t *testing.T) {
			if v := s.a.Equal(s.b); v != s.expect {
				t.Fatalf("Expected %v, got %v", s.expect, v)
			}
		})
	}
}

func TestDateTimeUnix(t *testing.T) {
	scenarios := []struct {
		date     string
		expected int64
	}{
		{"", -62135596800},
		{"2022-01-01 11:23:45.678Z", 1641036225},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%s", i, s.date), func(t *testing.T) {
			dt, err := types.ParseDateTime(s.date)
			if err != nil {
				t.Fatal(err)
			}

			v := dt.Unix()

			if v != s.expected {
				t.Fatalf("Expected %d, got %d", s.expected, v)
			}
		})
	}
}

func TestDateTimeIsZero(t *testing.T) {
	dt0 := types.DateTime{}
	if !dt0.IsZero() {
		t.Fatalf("Expected zero datatime, got %v", dt0)
	}

	dt1 := types.NowDateTime()
	if dt1.IsZero() {
		t.Fatalf("Expected non-zero datatime, got %v", dt1)
	}
}

func TestDateTimeString(t *testing.T) {
	dt0 := types.DateTime{}
	if dt0.String() != "" {
		t.Fatalf("Expected empty string for zer datetime, got %q", dt0.String())
	}

	expected := "2022-01-01 11:23:45.678Z"
	dt1, _ := types.ParseDateTime(expected)
	if dt1.String() != expected {
		t.Fatalf("Expected %q, got %v", expected, dt1)
	}
}

func TestDateTimeMarshalJSON(t *testing.T) {
	scenarios := []struct {
		date     string
		expected string
	}{
		{"", `""`},
		{"2022-01-01 11:23:45.678", `"2022-01-01 11:23:45.678Z"`},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%s", i, s.date), func(t *testing.T) {
			dt, err := types.ParseDateTime(s.date)
			if err != nil {
				t.Fatal(err)
			}

			result, err := dt.MarshalJSON()
			if err != nil {
				t.Fatal(err)
			}

			if string(result) != s.expected {
				t.Fatalf("Expected %q, got %q", s.expected, string(result))
			}
		})
	}
}

func TestDateTimeUnmarshalJSON(t *testing.T) {
	scenarios := []struct {
		date     string
		expected string
	}{
		{"", ""},
		{"invalid_json", ""},
		{"'123'", ""},
		{"2022-01-01 11:23:45.678", ""},
		{`"2022-01-01 11:23:45.678"`, "2022-01-01 11:23:45.678Z"},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%s", i, s.date), func(t *testing.T) {
			dt := types.DateTime{}
			dt.UnmarshalJSON([]byte(s.date))

			if dt.String() != s.expected {
				t.Fatalf("Expected %q, got %q", s.expected, dt.String())
			}
		})
	}
}

func TestDateTimeValue(t *testing.T) {
	scenarios := []struct {
		value    any
		expected string
	}{
		{"", ""},
		{"invalid", ""},
		{1641024040, "2022-01-01 08:00:40.000Z"},
		{"2022-01-01 11:23:45.678", "2022-01-01 11:23:45.678Z"},
		{types.NowDateTime(), types.NowDateTime().String()},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%s", i, s.value), func(t *testing.T) {
			dt, _ := types.ParseDateTime(s.value)

			result, err := dt.Value()
			if err != nil {
				t.Fatal(err)
			}

			if result != s.expected {
				t.Fatalf("Expected %q, got %q", s.expected, result)
			}
		})
	}
}

func TestDateTimeScan(t *testing.T) {
	now := time.Now().UTC().Format("2006-01-02 15:04:05") // without ms part for test consistency

	scenarios := []struct {
		value    any
		expected string
	}{
		{nil, ""},
		{"", ""},
		{"invalid", ""},
		{types.NowDateTime(), now},
		{time.Now(), now},
		{1.0, ""},
		{1641024040, "2022-01-01 08:00:40.000Z"},
		{"2022-01-01 11:23:45.678", "2022-01-01 11:23:45.678Z"},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#v", i, s.value), func(t *testing.T) {
			dt := types.DateTime{}

			err := dt.Scan(s.value)
			if err != nil {
				t.Fatalf("Failed to parse %v: %v", s.value, err)
			}

			if !strings.Contains(dt.String(), s.expected) {
				t.Fatalf("Expected %q, got %q", s.expected, dt.String())
			}
		})
	}
}
