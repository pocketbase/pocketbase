package core_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

func TestCollectionViewOptionsValidate(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		name           string
		collection     func(app core.App) (*core.Collection, error)
		expectedErrors []string
	}{
		{
			name: "view with empty query",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewViewCollection("new_auth")
				return c, nil
			},
			expectedErrors: []string{"fields", "viewQuery"},
		},
		{
			name: "view with invalid query",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewViewCollection("new_auth")
				c.ViewQuery = "invalid"
				return c, nil
			},
			expectedErrors: []string{"fields", "viewQuery"},
		},
		{
			name: "view with valid query but missing id",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewViewCollection("new_auth")
				c.ViewQuery = "select 1"
				return c, nil
			},
			expectedErrors: []string{"fields", "viewQuery"},
		},
		{
			name: "view with valid query",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewViewCollection("new_auth")
				c.ViewQuery = "select demo1.id, text as example from demo1"
				return c, nil
			},
			expectedErrors: []string{},
		},
		{
			name: "update view query ",
			collection: func(app core.App) (*core.Collection, error) {
				c, _ := app.FindCollectionByNameOrId("view2")
				c.ViewQuery = "select demo1.id, text as example from demo1"
				return c, nil
			},
			expectedErrors: []string{},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			app, _ := tests.NewTestApp()
			defer app.Cleanup()

			collection, err := s.collection(app)
			if err != nil {
				t.Fatalf("Failed to retrieve test collection: %v", err)
			}

			result := app.Validate(collection)

			tests.TestValidationErrors(t, result, s.expectedErrors)
		})
	}
}
