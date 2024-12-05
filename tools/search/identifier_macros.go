package search

import (
	"fmt"
	"time"

	"github.com/pocketbase/pocketbase/tools/types"
)

// note: used primarily for the tests
var timeNow = func() time.Time {
	return time.Now()
}

var identifierMacros = map[string]func() (any, error){
	"@now": func() (any, error) {
		today := timeNow().UTC()

		d, err := types.ParseDateTime(today)
		if err != nil {
			return "", fmt.Errorf("@now: %w", err)
		}

		return d.String(), nil
	},
	"@yesterday": func() (any, error) {
		yesterday := timeNow().UTC().AddDate(0, 0, -1)

		d, err := types.ParseDateTime(yesterday)
		if err != nil {
			return "", fmt.Errorf("@yesterday: %w", err)
		}

		return d.String(), nil
	},
	"@tomorrow": func() (any, error) {
		tomorrow := timeNow().UTC().AddDate(0, 0, 1)

		d, err := types.ParseDateTime(tomorrow)
		if err != nil {
			return "", fmt.Errorf("@tomorrow: %w", err)
		}

		return d.String(), nil
	},
	"@second": func() (any, error) {
		return timeNow().UTC().Second(), nil
	},
	"@minute": func() (any, error) {
		return timeNow().UTC().Minute(), nil
	},
	"@hour": func() (any, error) {
		return timeNow().UTC().Hour(), nil
	},
	"@day": func() (any, error) {
		return timeNow().UTC().Day(), nil
	},
	"@month": func() (any, error) {
		return int(timeNow().UTC().Month()), nil
	},
	"@weekday": func() (any, error) {
		return int(timeNow().UTC().Weekday()), nil
	},
	"@year": func() (any, error) {
		return timeNow().UTC().Year(), nil
	},
	"@todayStart": func() (any, error) {
		today := timeNow().UTC()
		start := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)

		d, err := types.ParseDateTime(start)
		if err != nil {
			return "", fmt.Errorf("@todayStart: %w", err)
		}

		return d.String(), nil
	},
	"@todayEnd": func() (any, error) {
		today := timeNow().UTC()

		start := time.Date(today.Year(), today.Month(), today.Day(), 23, 59, 59, 999999999, time.UTC)

		d, err := types.ParseDateTime(start)
		if err != nil {
			return "", fmt.Errorf("@todayEnd: %w", err)
		}

		return d.String(), nil
	},
	"@monthStart": func() (any, error) {
		today := timeNow().UTC()
		start := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, time.UTC)

		d, err := types.ParseDateTime(start)
		if err != nil {
			return "", fmt.Errorf("@monthStart: %w", err)
		}

		return d.String(), nil
	},
	"@monthEnd": func() (any, error) {
		today := timeNow().UTC()
		start := time.Date(today.Year(), today.Month(), 1, 23, 59, 59, 999999999, time.UTC)
		end := start.AddDate(0, 1, -1)

		d, err := types.ParseDateTime(end)
		if err != nil {
			return "", fmt.Errorf("@monthEnd: %w", err)
		}

		return d.String(), nil
	},
	"@yearStart": func() (any, error) {
		today := timeNow().UTC()
		start := time.Date(today.Year(), 1, 1, 0, 0, 0, 0, time.UTC)

		d, err := types.ParseDateTime(start)
		if err != nil {
			return "", fmt.Errorf("@yearStart: %w", err)
		}

		return d.String(), nil
	},
	"@yearEnd": func() (any, error) {
		today := timeNow().UTC()
		end := time.Date(today.Year(), 12, 31, 23, 59, 59, 999999999, time.UTC)

		d, err := types.ParseDateTime(end)
		if err != nil {
			return "", fmt.Errorf("@yearEnd: %w", err)
		}

		return d.String(), nil
	},
}
