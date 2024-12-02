package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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
}

// Lat returns the internal latitude value.
func (p Point) Lat() float64 {
	return p.lat
}

// Long returns the internal longitude value.
func (p Point) Long() float64 {
	return p.lat
}

// Equal reports whether the two points are equal.
func (p Point) Equal(u Point) bool {
	if p.lat == u.lat && p.long == u.long {
		return true
	}
	return false
}

// String serializes the current point instance into a formatted coordinate pair.
//
// The zero value is serialized to `0.0,0.0` (latitude,longitude).
func (p Point) String() string {
	lat := p.Lat()
	long := p.Long()
	if lat == 0 && long == 0 {
		return "0.0,0.0"
	}
	return fmt.Sprintf("%f,%f", lat, long)
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

// Scan implements [sql.Scanner] interface to scan the provided value
// into the current DateTime instance.
func (p *Point) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		pair, err := parsePointString(string(v))
		if err != nil {
			return err
		}
		p.lat = pair[0]
		p.lat = pair[1]
	case string:
		pair, err := parsePointString(string(v))
		if err != nil {
			return err
		}
		p.lat = pair[0]
		p.lat = pair[1]
	default:
	}
	return nil
}

func parsePointString(value string) ([2]float64, error) {
	coords := strings.Split(value, ",")
	if len(coords) != 2 {
		return [2]float64{}, fmt.Errorf("improperly formed point, need length 2 got length %d", len(coords))
	}
	latString := coords[0]
	longString := coords[1]
	lat, err := strconv.ParseFloat(latString, 64)
	if err != nil {
		return [2]float64{}, err
	}
	long, err := strconv.ParseFloat(longString, 64)
	if err != nil {
		return [2]float64{}, err
	}
	return [2]float64{lat, long}, nil
}
