package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cast"
)

// ParsePoint creates a new Point from the provided value
func ParsePoint(value any) (Point, error) {
	p := Point{}
	err := p.Scan(value)
	return p, err
}

// Point represents a geographic point on Earth,
// serialized as a comma-separated pair of float64s.
type Point struct {
	lat  float64
	long float64
	// if set is true, the point is not the zero value
	set bool
}

// Lat returns the internal latitude value.
func (p Point) Lat() float64 {
	return p.lat
}

// Long returns the internal longitude value.
func (p Point) Long() float64 {
	return p.long
}

// Equal reports whether the two points are equal.
func (a Point) Equal(b Point) bool {
	return a.lat == b.lat && a.long == b.long && a.set == b.set
}

// String serializes the current point instance into a formatted coordinate pair.
//
// The zero value is serialized to an empty string.
func (p Point) String() string {
	if !p.set {
		return ""
	}
	return fmt.Sprintf("%s, %s", strconv.FormatFloat(p.lat, "f"[0], -1, 64), strconv.FormatFloat(p.long, "f"[0], -1, 64))
}

// MarshalJSON implements the [json.Marshaler] interface.
func (p Point) MarshalJSON() ([]byte, error) {
	return []byte(`"` + p.String() + `"`), nil
}

// UnmarshalJSON implements the [json.Unmarshaler] interface.
func (p *Point) UnmarshalJSON(b []byte) error {
	var raw string
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	return p.Scan(raw)
}

// Value implements the [driver.Valuer] interface.
func (p Point) Value() (driver.Value, error) {
	return p.String(), nil
}

// IsEmpty checks whether the current Point instance has been set.
func (p Point) IsEmpty() bool {
	return !p.set
}

// Scan implements [sql.Scanner] interface to scan the provided value
// into the current Point instance.
func (p *Point) Scan(value any) error {
	switch v := value.(type) {
	case Point:
		*p = v
		return nil
	case string:
		return p.parsePointString(v)
	case []byte:
		return p.parsePointString(string(v))
	default:
		return p.parsePointString(cast.ToString(v))
	}
}

func (p *Point) parsePointString(pair string) error {
	if pair == "" {
		*p = Point{}
		return nil
	}
	latStr, longStr, found := strings.Cut(pair, ",")
	if !found {
		*p = Point{}
		return fmt.Errorf("point must have a comma-separated latitude and longitude, got: %s", pair)
	}

	latStr = strings.TrimSpace(latStr)
	longStr = strings.TrimSpace(longStr)

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return err
	}
	long, err := strconv.ParseFloat(longStr, 64)
	if err != nil {
		return err
	}
	*p = Point{lat: lat, long: long, set: true}
	return nil
}
