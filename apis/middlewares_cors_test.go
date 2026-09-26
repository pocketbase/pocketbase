package apis

import "testing"

func TestCorsExactWildcardMatch(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		domain   string
		pattern  string
		expected bool
	}{
		{"", "", false},
		{"http://example.com", "", false},
		{"", "http://example.com", false},
		{"http://example.com", "https://example.com", false},
		{"abc://example.com", "abc://example.com", true},
		{"https://example.com", "https://example.com", true},
		{"https://a.example.com", "https://example.com", false},
		{"https://example.com", "https://a.example.com", false},
		{"https://example.com", "https://*.example.com", false},
		{"https://a.example.com", "https://*.example.com", true},
		{"https://a.example.com", "https://a.*.example.com", false},
		{"https://a.b.example.com", "https://a.*.example.com", true},
		{"https://a.b.example.com", "https://a2.*.example.com", false},
		{"https://a.b.example.com", "https://a.*", false},
		{"https://a.b.example", "https://a.*", false},
		{"https://a.b", "https://a.*", true},
	}

	for _, s := range scenarios {
		t.Run(s.domain+":"+s.pattern, func(t *testing.T) {
			result := exactWildcardMatch(s.domain, s.pattern)
			if result != s.expected {
				t.Fatalf("Expected %v, got %v", s.expected, result)
			}
		})
	}
}
