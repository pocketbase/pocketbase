package types_test

import (
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
		dt, err := types.ParseDateTime(s.value)
		if err != nil {
			t.Errorf("(%d) Failed to parse %v: %v", i, s.value, err)
			continue
		}

		if dt.String() != s.expected {
			t.Errorf("(%d) Expected %q, got %q", i, s.expected, dt.String())
		}
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
		t.Errorf("Expected time %v, got %v", expected, result)
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
		dt, err := types.ParseDateTime(s.date)
		if err != nil {
			t.Errorf("(%d) %v", i, err)
		}

		result, err := dt.MarshalJSON()
		if err != nil {
			t.Errorf("(%d) %v", i, err)
		}

		if string(result) != s.expected {
			t.Errorf("(%d) Expected %q, got %q", i, s.expected, string(result))
		}
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
		dt := types.DateTime{}
		dt.UnmarshalJSON([]byte(s.date))

		if dt.String() != s.expected {
			t.Errorf("(%d) Expected %q, got %q", i, s.expected, dt.String())
		}
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
		dt, _ := types.ParseDateTime(s.value)
		result, err := dt.Value()
		if err != nil {
			t.Errorf("(%d) %v", i, err)
			continue
		}

		if result != s.expected {
			t.Errorf("(%d) Expected %q, got %q", i, s.expected, result)
		}
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
		dt := types.DateTime{}

		err := dt.Scan(s.value)
		if err != nil {
			t.Errorf("(%d) Failed to parse %v: %v", i, s.value, err)
			continue
		}

		if !strings.Contains(dt.String(), s.expected) {
			t.Errorf("(%d) Expected %q, got %q", i, s.expected, dt.String())
		}
	}
}
