package models_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/models"
)

func TestRequestInfoHasModifierDataKeys(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		name        string
		requestInfo *models.RequestInfo
		expected    bool
	}{
		{
			"empty",
			&models.RequestInfo{},
			false,
		},
		{
			"Data with regular fields",
			&models.RequestInfo{
				Query: map[string]any{"data+": "demo"}, // should be ignored
				Data:  map[string]any{"a": 123, "b": "test", "c.d": false},
			},
			false,
		},
		{
			"Data with +modifier fields",
			&models.RequestInfo{
				Data: map[string]any{"a+": 123, "b": "test", "c.d": false},
			},
			true,
		},
		{
			"Data with -modifier fields",
			&models.RequestInfo{
				Data: map[string]any{"a": 123, "b-": "test", "c.d": false},
			},
			true,
		},
		{
			"Data with mixed modifier fields",
			&models.RequestInfo{
				Data: map[string]any{"a": 123, "b-": "test", "c.d+": false},
			},
			true,
		},
	}

	for _, s := range scenarios {
		result := s.requestInfo.HasModifierDataKeys()

		if result != s.expected {
			t.Fatalf("[%s] Expected %v, got %v", s.name, s.expected, result)
		}
	}
}
