package core_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestBaseAppRefreshSettings(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// cleanup all stored settings
	if _, err := app.DB().NewQuery("DELETE from _params;").Execute(); err != nil {
		t.Fatalf("Failed to delete all test settings: %v", err)
	}

	// check if the new settings are saved in the db
	app.ResetEventCalls()
	if err := app.RefreshSettings(); err != nil {
		t.Fatalf("Failed to refresh the settings after delete: %v", err)
	}
	testEventCalls(t, app, map[string]int{
		"OnModelBeforeCreate": 1,
		"OnModelAfterCreate":  1,
	})
	param, err := app.Dao().FindParamByKey(models.ParamAppSettings)
	if err != nil {
		t.Fatalf("Expected new settings to be persisted, got %v", err)
	}

	// change the db entry and refresh the app settings (ensure that there was no db update)
	param.Value = types.JsonRaw([]byte(`{"example": 123}`))
	if err := app.Dao().SaveParam(param.Key, param.Value); err != nil {
		t.Fatalf("Failed to update the test settings: %v", err)
	}
	app.ResetEventCalls()
	if err := app.RefreshSettings(); err != nil {
		t.Fatalf("Failed to refresh the app settings: %v", err)
	}
	testEventCalls(t, app, nil)

	// try to refresh again without doing any changes
	app.ResetEventCalls()
	if err := app.RefreshSettings(); err != nil {
		t.Fatalf("Failed to refresh the app settings without change: %v", err)
	}
	testEventCalls(t, app, nil)
}

func testEventCalls(t *testing.T, app *tests.TestApp, events map[string]int) {
	if len(events) != len(app.EventCalls) {
		t.Fatalf("Expected events doesn't match: \n%v, \ngot \n%v", events, app.EventCalls)
	}

	for name, total := range events {
		if v, ok := app.EventCalls[name]; !ok || v != total {
			t.Fatalf("Expected events doesn't exist or match: \n%v, \ngot \n%v", events, app.EventCalls)
		}
	}
}
