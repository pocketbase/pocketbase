package ghupdate

import "testing"

func TestCompareVersions(t *testing.T) {
	scenarios := []struct {
		a        string
		b        string
		expected int
	}{
		{"", "", 0},
		{"0", "", 0},
		{"1", "1.0.0", 0},
		{"1.1", "1.1.0", 0},
		{"1.1", "1.1.1", 1},
		{"1.1", "1.0.1", -1},
		{"1.0", "1.0.1", 1},
		{"1.10", "1.9", -1},
		{"1.2", "1.12", 1},
		{"3.2", "1.6", -1},
		{"3.2", "4.0", 1},
		{"3.2.4", "3.2.3", -1},
	}

	for i, s := range scenarios {
		result := compareVersions(s.a, s.b)

		if result != s.expected {
			t.Fatalf("[%d] Expected %q vs %q to result in %d, got %d", i, s.a, s.b, s.expected, result)
		}
	}
}
