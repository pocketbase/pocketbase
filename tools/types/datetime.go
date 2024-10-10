package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/spf13/cast"
)

// DefaultDateLayout specifies the default app date strings layout.
const DefaultDateLayout = "2006-01-02 15:04:05.000Z"

// NowDateTime returns new DateTime instance with the current local time.
func NowDateTime() DateTime {
	return DateTime{t: time.Now()}
}

// ParseDateTime creates a new DateTime from the provided value
// (could be [cast.ToTime] supported string, [time.Time], etc.).
func ParseDateTime(value any) (DateTime, error) {
	d := DateTime{}
	err := d.Scan(value)
	return d, err
}

// DateTime represents a [time.Time] instance in UTC that is wrapped
// and serialized using the app default date layout.
type DateTime struct {
	t time.Time
}

// Time returns the internal [time.Time] instance.
func (d DateTime) Time() time.Time {
	return d.t
}

// Add returns a new DateTime based on the current DateTime + the specified duration.
func (d DateTime) Add(duration time.Duration) DateTime {
	d.t = d.t.Add(duration)
	return d
}

// Sub returns a [time.Duration] by subtracting the specified DateTime from the current one.
//
// If the result exceeds the maximum (or minimum) value that can be stored in a [time.Duration],
// the maximum (or minimum) duration will be returned.
func (d DateTime) Sub(u DateTime) time.Duration {
	return d.Time().Sub(u.Time())
}

// AddDate returns a new DateTime based on the current one + duration.
//
// It follows the same rules as [time.AddDate].
func (d DateTime) AddDate(years, months, days int) DateTime {
	d.t = d.t.AddDate(years, months, days)
	return d
}

// After reports whether the current DateTime instance is after u.
func (d DateTime) After(u DateTime) bool {
	return d.Time().After(u.Time())
}

// Before reports whether the current DateTime instance is before u.
func (d DateTime) Before(u DateTime) bool {
	return d.Time().Before(u.Time())
}

// Compare compares the current DateTime instance with u.
// If the current instance is before u, it returns -1.
// If the current instance is after u, it returns +1.
// If they're the same, it returns 0.
func (d DateTime) Compare(u DateTime) int {
	return d.Time().Compare(u.Time())
}

// Equal reports whether the current DateTime and u represent the same time instant.
// Two DateTime can be equal even if they are in different locations.
// For example, 6:00 +0200 and 4:00 UTC are Equal.
func (d DateTime) Equal(u DateTime) bool {
	return d.Time().Equal(u.Time())
}

// Unix returns the current DateTime as a Unix time, aka.
// the number of seconds elapsed since January 1, 1970 UTC.
func (d DateTime) Unix() int64 {
	return d.Time().Unix()
}

// IsZero checks whether the current DateTime instance has zero time value.
func (d DateTime) IsZero() bool {
	return d.Time().IsZero()
}

// String serializes the current DateTime instance into a formatted
// UTC date string.
//
// The zero value is serialized to an empty string.
func (d DateTime) String() string {
	t := d.Time()
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(DefaultDateLayout)
}

// MarshalJSON implements the [json.Marshaler] interface.
func (d DateTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.String() + `"`), nil
}

// UnmarshalJSON implements the [json.Unmarshaler] interface.
func (d *DateTime) UnmarshalJSON(b []byte) error {
	var raw string
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	return d.Scan(raw)
}

// Value implements the [driver.Valuer] interface.
func (d DateTime) Value() (driver.Value, error) {
	return d.String(), nil
}

// Scan implements [sql.Scanner] interface to scan the provided value
// into the current DateTime instance.
func (d *DateTime) Scan(value any) error {
	switch v := value.(type) {
	case time.Time:
		d.t = v
	case DateTime:
		d.t = v.Time()
	case string:
		if v == "" {
			d.t = time.Time{}
		} else {
			t, err := time.Parse(DefaultDateLayout, v)
			if err != nil {
				// check for other common date layouts
				t = cast.ToTime(v)
			}
			d.t = t
		}
	case int, int64, int32, uint, uint64, uint32:
		d.t = cast.ToTime(v)
	default:
		str := cast.ToString(v)
		if str == "" {
			d.t = time.Time{}
		} else {
			d.t = cast.ToTime(str)
		}
	}

	return nil
}
