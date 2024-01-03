package daos_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/security"
)

func TestSaveAndFindSettings(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	encryptionKey := security.PseudorandomString(32)

	// change unencrypted app settings
	app.Settings().Meta.AppName = "save_unencrypted"
	if err := app.Dao().SaveSettings(app.Settings()); err != nil {
		t.Fatal(err)
	}

	// check if the change was persisted
	s1, err := app.Dao().FindSettings()
	if err != nil {
		t.Fatalf("Failed to fetch settings: %v", err)
	}
	if s1.Meta.AppName != "save_unencrypted" {
		t.Fatalf("Expected settings to be changed with app name %q, got \n%v", "save_unencrypted", s1)
	}

	// make another change but this time provide an encryption key
	app.Settings().Meta.AppName = "save_encrypted"
	if err := app.Dao().SaveSettings(app.Settings(), encryptionKey); err != nil {
		t.Fatal(err)
	}

	// try to fetch the settings without encryption key (should fail)
	if s2, err := app.Dao().FindSettings(); err == nil {
		t.Fatalf("Expected FindSettings to fail without an encryption key, got \n%v", s2)
	}

	// try again but this time with an encryption key
	s3, err := app.Dao().FindSettings(encryptionKey)
	if err != nil {
		t.Fatalf("Failed to fetch settings with an encryption key %s: %v", encryptionKey, err)
	}
	if s3.Meta.AppName != "save_encrypted" {
		t.Fatalf("Expected settings to be changed with app name %q, got \n%v", "save_encrypted", s3)
	}
}
