package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// GeoPoint defines a struct for storing geo coordinates as serialized json object
// (e.g. {lon:0,lat:0}).
//
// Note: using object notation and not a plain array to avoid the confusion
// as there doesn't seem to be a fixed standard for the coordinates order.
type GeoPoint struct {
	Lon float64 `form:"lon" json:"lon"`
	Lat float64 `form:"lat" json:"lat"`
}

// String returns the string representation of the current GeoPoint instance.
func (p GeoPoint) String() string {
	raw, _ := json.Marshal(p)
	return string(raw)
}

// AsMap implements [core.mapExtractor] and returns a value suitable
// to be used in an API rule expression.
func (p GeoPoint) AsMap() map[string]any {
	return map[string]any{
		"lon": p.Lon,
		"lat": p.Lat,
	}
}

// Value implements the [driver.Valuer] interface.
func (p GeoPoint) Value() (driver.Value, error) {
	data, err := json.Marshal(p)
	return string(data), err
}

// Scan implements [sql.Scanner] interface to scan the provided value
// into the current GeoPoint instance.
//
// The value argument could be nil (no-op), another GeoPoint instance,
// map or serialized json object with lat-lon props.
func (p *GeoPoint) Scan(value any) error {
	var err error

	switch v := value.(type) {
	case nil:
		// no cast needed
	case *GeoPoint:
		p.Lon = v.Lon
		p.Lat = v.Lat
	case GeoPoint:
		p.Lon = v.Lon
		p.Lat = v.Lat
	case JSONRaw:
		if len(v) != 0 {
			err = json.Unmarshal(v, p)
		}
	case []byte:
		if len(v) != 0 {
			err = json.Unmarshal(v, p)
		}
	case string:
		if len(v) != 0 {
			err = json.Unmarshal([]byte(v), p)
		}
	default:
		var raw []byte
		raw, err = json.Marshal(v)
		if err != nil {
			err = fmt.Errorf("unable to marshalize value for scanning: %w", err)
		} else {
			err = json.Unmarshal(raw, p)
		}
	}

	if err != nil {
		return fmt.Errorf("[GeoPoint] unable to scan value %v: %w", value, err)
	}

	return nil
}
