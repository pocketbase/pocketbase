package search

import (
	"testing"
	"time"
)

func TestIdentifierMacros(t *testing.T) {
	originalTimeNow := timeNow

	timeNow = func() time.Time {
		return time.Date(2023, 2, 3, 4, 5, 6, 7, time.UTC)
	}

	testMacros := map[string]any{
		"@now":        "2023-02-03 04:05:06.000Z",
		"@yesterday":  "2023-02-02 04:05:06.000Z",
		"@tomorrow":   "2023-02-04 04:05:06.000Z",
		"@second":     6,
		"@minute":     5,
		"@hour":       4,
		"@day":        3,
		"@month":      2,
		"@weekday":    5,
		"@year":       2023,
		"@todayStart": "2023-02-03 00:00:00.000Z",
		"@todayEnd":   "2023-02-03 23:59:59.999Z",
		"@monthStart": "2023-02-01 00:00:00.000Z",
		"@monthEnd":   "2023-02-28 23:59:59.999Z",
		"@yearStart":  "2023-01-01 00:00:00.000Z",
		"@yearEnd":    "2023-12-31 23:59:59.999Z",
	}

	if len(testMacros) != len(identifierMacros) {
		t.Fatalf("Expected %d macros, got %d", len(testMacros), len(identifierMacros))
	}

	for key, expected := range testMacros {
		t.Run(key, func(t *testing.T) {
			macro, ok := identifierMacros[key]
			if !ok {
				t.Fatalf("Missing macro %s", key)
			}

			result, err := macro()
			if err != nil {
				t.Fatal(err)
			}

			if result != expected {
				t.Fatalf("Expected %q, got %q", expected, result)
			}
		})
	}

	// restore
	timeNow = originalTimeNow
}
