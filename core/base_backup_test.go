package core_test

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/thinkonmay/pocketbase/core"
	"github.com/thinkonmay/pocketbase/tests"
	"github.com/thinkonmay/pocketbase/tools/archive"
	"github.com/thinkonmay/pocketbase/tools/list"
)

func TestCreateBackup(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// set some long app name with spaces and special characters
	app.Settings().Meta.AppName = "test @! " + strings.Repeat("a", 100)

	expectedAppNamePrefix := "test_" + strings.Repeat("a", 45)

	// test pending error
	app.Store().Set(core.StoreKeyActiveBackup, "")
	if err := app.CreateBackup(context.Background(), "test.zip"); err == nil {
		t.Fatal("Expected pending error, got nil")
	}
	app.Store().Remove(core.StoreKeyActiveBackup)

	// create with auto generated name
	if err := app.CreateBackup(context.Background(), ""); err != nil {
		t.Fatal("Failed to create a backup with autogenerated name")
	}

	// create with custom name
	if err := app.CreateBackup(context.Background(), "custom"); err != nil {
		t.Fatal("Failed to create a backup with custom name")
	}

	// create new with the same name (aka. replace)
	if err := app.CreateBackup(context.Background(), "custom"); err != nil {
		t.Fatal("Failed to create and replace a backup with the same name")
	}

	backupsDir := filepath.Join(app.DataDir(), core.LocalBackupsDirName)

	entries, err := os.ReadDir(backupsDir)
	if err != nil {
		t.Fatal(err)
	}

	expectedFiles := []string{
		`^pb_backup_` + expectedAppNamePrefix + `_\w+\.zip$`,
		`^pb_backup_` + expectedAppNamePrefix + `_\w+\.zip.attrs$`,
		"custom",
		"custom.attrs",
	}

	if len(entries) != len(expectedFiles) {
		names := getEntryNames(entries)
		t.Fatalf("Expected %d backup files, got %d: \n%v", len(expectedFiles), len(entries), names)
	}

	for i, entry := range entries {
		if !list.ExistInSliceWithRegex(entry.Name(), expectedFiles) {
			t.Fatalf("[%d] Missing backup file %q", i, entry.Name())
		}

		if strings.HasSuffix(entry.Name(), ".attrs") {
			continue
		}

		path := filepath.Join(backupsDir, entry.Name())

		if err := verifyBackupContent(app, path); err != nil {
			t.Fatalf("[%d] Failed to verify backup content: %v", i, err)
		}
	}
}

func TestRestoreBackup(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// create a initial test backup to ensure that there are at least 1
	// backup file and that the generated zip doesn't contain the backups dir
	if err := app.CreateBackup(context.Background(), "initial"); err != nil {
		t.Fatal("Failed to create test initial backup")
	}

	// create test backup
	if err := app.CreateBackup(context.Background(), "test"); err != nil {
		t.Fatal("Failed to create test backup")
	}

	// test pending error
	app.Store().Set(core.StoreKeyActiveBackup, "")
	if err := app.RestoreBackup(context.Background(), "test"); err == nil {
		t.Fatal("Expected pending error, got nil")
	}
	app.Store().Remove(core.StoreKeyActiveBackup)

	// missing backup
	if err := app.RestoreBackup(context.Background(), "missing"); err == nil {
		t.Fatal("Expected missing error, got nil")
	}
}

// -------------------------------------------------------------------

func verifyBackupContent(app core.App, path string) error {
	dir, err := os.MkdirTemp("", "backup_test")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	if err := archive.Extract(path, dir); err != nil {
		return err
	}

	expectedRootEntries := []string{
		"storage",
		"data.db",
		"data.db-shm",
		"data.db-wal",
		"logs.db",
		"logs.db-shm",
		"logs.db-wal",
		".gitignore",
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	if len(entries) != len(expectedRootEntries) {
		names := getEntryNames(entries)
		return fmt.Errorf("Expected %d backup files, got %d: \n%v", len(expectedRootEntries), len(entries), names)
	}

	for _, entry := range entries {
		if !list.ExistInSliceWithRegex(entry.Name(), expectedRootEntries) {
			return fmt.Errorf("Didn't expect %q entry", entry.Name())
		}
	}

	return nil
}

func getEntryNames(entries []fs.DirEntry) []string {
	names := make([]string, len(entries))

	for i, entry := range entries {
		names[i] = entry.Name()
	}

	return names
}
