package core_test

import (
	"os"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestReloadSettings(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// cleanup all stored settings
	// ---
	if _, err := app.DB().NewQuery("DELETE from _params;").Execute(); err != nil {
		t.Fatalf("Failed to delete all test settings: %v", err)
	}

	// check if the new settings are saved in the db
	// ---
	app.Settings().Meta.AppName = "test_name_after_delete"

	app.ResetEventCalls()
	if err := app.ReloadSettings(); err != nil {
		t.Fatalf("Failed to reload the settings after delete: %v", err)
	}
	testEventCalls(t, app, map[string]int{
		"OnModelCreate":             1,
		"OnModelCreateExecute":      1,
		"OnModelAfterCreateSuccess": 1,
		"OnModelValidate":           1,
		"OnSettingsReload":          1,
	})

	param := &core.Param{}
	err := app.ModelQuery(param).Model("settings", param)
	if err != nil {
		t.Fatalf("Expected new settings to be persisted, got %v", err)
	}

	if !strings.Contains(param.Value.String(), "test_name_after_delete") {
		t.Fatalf("Expected to find AppName test_name_after_delete in\n%s", param.Value.String())
	}

	// change the db entry and reload the app settings (ensure that there was no db update)
	// ---
	param.Value = types.JSONRaw([]byte(`{"meta": {"appName":"test_name_after_update"}}`))
	if err := app.Save(param); err != nil {
		t.Fatalf("Failed to update the test settings: %v", err)
	}

	app.ResetEventCalls()
	if err := app.ReloadSettings(); err != nil {
		t.Fatalf("Failed to reload app settings: %v", err)
	}
	testEventCalls(t, app, map[string]int{
		"OnSettingsReload": 1,
	})

	// try to reload again without doing any changes
	// ---
	app.ResetEventCalls()
	if err := app.ReloadSettings(); err != nil {
		t.Fatalf("Failed to reload app settings without change: %v", err)
	}
	testEventCalls(t, app, map[string]int{
		"OnSettingsReload": 1,
	})

	if app.Settings().Meta.AppName != "test_name_after_update" {
		t.Fatalf("Expected AppName %q, got %q", "test_name_after_update", app.Settings().Meta.AppName)
	}
}

func TestReloadSettingsWithEncryption(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	os.Setenv("pb_test_env", strings.Repeat("a", 32))

	// cleanup all stored settings
	// ---
	if _, err := app.DB().NewQuery("DELETE from _params;").Execute(); err != nil {
		t.Fatalf("Failed to delete all test settings: %v", err)
	}

	// check if the new settings are saved in the db
	// ---
	app.Settings().Meta.AppName = "test_name_after_delete"

	app.ResetEventCalls()
	if err := app.ReloadSettings(); err != nil {
		t.Fatalf("Failed to reload the settings after delete: %v", err)
	}
	testEventCalls(t, app, map[string]int{
		"OnModelCreate":             1,
		"OnModelCreateExecute":      1,
		"OnModelAfterCreateSuccess": 1,
		"OnModelValidate":           1,
		"OnSettingsReload":          1,
	})

	param := &core.Param{}
	err := app.ModelQuery(param).Model("settings", param)
	if err != nil {
		t.Fatalf("Expected new settings to be persisted, got %v", err)
	}
	rawValue := param.Value.String()
	if rawValue == "" || strings.Contains(rawValue, "test_name") {
		t.Fatalf("Expected inserted settings to be encrypted, found\n%s", rawValue)
	}

	// change and reload the app settings (ensure that there was no db update)
	// ---
	app.Settings().Meta.AppName = "test_name_after_update"
	if err := app.Save(app.Settings()); err != nil {
		t.Fatalf("Failed to update app settings: %v", err)
	}

	// try to reload again without doing any changes
	// ---
	app.ResetEventCalls()
	if err := app.ReloadSettings(); err != nil {
		t.Fatalf("Failed to reload app settings: %v", err)
	}
	testEventCalls(t, app, map[string]int{
		"OnSettingsReload": 1,
	})

	// refetch the settings param to ensure that the new value was stored encrypted
	err = app.ModelQuery(param).Model("settings", param)
	if err != nil {
		t.Fatalf("Expected new settings to be persisted, got %v", err)
	}
	rawValue = param.Value.String()
	if rawValue == "" || strings.Contains(rawValue, "test_name") {
		t.Fatalf("Expected updated settings to be encrypted, found\n%s", rawValue)
	}

	if app.Settings().Meta.AppName != "test_name_after_update" {
		t.Fatalf("Expected AppName %q, got %q", "test_name_after_update", app.Settings().Meta.AppName)
	}
}

func testEventCalls(t *testing.T, app *tests.TestApp, events map[string]int) {
	if len(events) != len(app.EventCalls) {
		t.Fatalf("Expected events doesn't match:\n%v\ngot\n%v", events, app.EventCalls)
	}

	for name, total := range events {
		if v, ok := app.EventCalls[name]; !ok || v != total {
			t.Fatalf("Expected events doesn't exist or match:\n%v\ngot\n%v", events, app.EventCalls)
		}
	}
}
