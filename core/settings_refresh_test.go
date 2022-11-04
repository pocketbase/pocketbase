package core_test

import (
	"bytes"
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
		t.Fatal("Failed to refresh the settings after delete")
	}
	testEventCalls(t, app, map[string]int{
		"OnModelBeforeCreate": 1,
		"OnModelAfterCreate":  1,
	})
	param, err := app.Dao().FindParamByKey(models.ParamAppSettings)
	if err != nil {
		t.Fatalf("Expected new settings to be persisted, got %v", err)
	}

	// change the db entry and refresh the app settings
	param.Value = types.JsonRaw([]byte(`{"example": 123}`))
	if err := app.Dao().SaveParam(param.Key, param.Value); err != nil {
		t.Fatalf("Failed to update the test settings: %v", err)
	}
	app.ResetEventCalls()
	if err := app.RefreshSettings(); err != nil {
		t.Fatalf("Failed to refresh the app settings: %v", err)
	}
	testEventCalls(t, app, map[string]int{
		"OnModelBeforeUpdate": 1,
		"OnModelAfterUpdate":  1,
	})

	// make sure that the newly merged settings were actually saved
	newParam, err := app.Dao().FindParamByKey(models.ParamAppSettings)
	if err != nil {
		t.Fatalf("Failed to fetch new settings param: %v", err)
	}
	if bytes.Equal(param.Value, newParam.Value) {
		t.Fatalf("Expected the new refreshed settings to be different, got: \n%v", string(newParam.Value))
	}

	// try to refresh again and ensure that there was no db update
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
