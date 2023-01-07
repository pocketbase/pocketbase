package models_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/models"
)

func TestRequestDataHasModifierDataKeys(t *testing.T) {
	scenarios := []struct {
		name        string
		requestData *models.RequestData
		expected    bool
	}{
		{
			"empty",
			&models.RequestData{},
			false,
		},
		{
			"Data with regular fields",
			&models.RequestData{
				Query: map[string]any{"data+": "demo"}, // should be ignored
				Data:  map[string]any{"a": 123, "b": "test", "c.d": false},
			},
			false,
		},
		{
			"Data with +modifier fields",
			&models.RequestData{
				Data: map[string]any{"a+": 123, "b": "test", "c.d": false},
			},
			true,
		},
		{
			"Data with -modifier fields",
			&models.RequestData{
				Data: map[string]any{"a": 123, "b-": "test", "c.d": false},
			},
			true,
		},
		{
			"Data with mixed modifier fields",
			&models.RequestData{
				Data: map[string]any{"a": 123, "b-": "test", "c.d+": false},
			},
			true,
		},
	}

	for _, s := range scenarios {
		result := s.requestData.HasModifierDataKeys()

		if result != s.expected {
			t.Fatalf("[%s] Expected %v, got %v", s.name, s.expected, result)
		}
	}
}
