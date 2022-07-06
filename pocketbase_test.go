package pocketbase

import (
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	testDir := "./pb_test_data_dir"
	defer os.RemoveAll(testDir)

	// reset os.Args
	os.Args = os.Args[0:1]
	os.Args = append(
		os.Args,
		"--dir="+testDir,
		"--encryptionEnv=test_encryption_env",
		"--debug=true",
	)

	app := New()

	if app == nil {
		t.Fatal("Expected initialized PocketBase instance, got nil")
	}

	if app.RootCmd == nil {
		t.Fatal("Expected RootCmd to be initialized, got nil")
	}

	if app.appWrapper == nil {
		t.Fatal("Expected appWrapper to be initialized, got nil")
	}

	if app.DataDir() != testDir {
		t.Fatalf("Expected app.DataDir() %q, got %q", testDir, app.DataDir())
	}

	if app.EncryptionEnv() != "test_encryption_env" {
		t.Fatalf("Expected app.DataDir() test_encryption_env, got %q", app.EncryptionEnv())
	}

	if app.IsDebug() != true {
		t.Fatal("Expected app.IsDebug() true, got false")
	}
}

func TestDefaultDebug(t *testing.T) {
	app := New()

	app.DefaultDebug(true)
	if app.defaultDebug != true {
		t.Fatalf("Expected defaultDebug %v, got %v", true, app.defaultDebug)
	}

	app.DefaultDebug(false)
	if app.defaultDebug != false {
		t.Fatalf("Expected defaultDebug %v, got %v", false, app.defaultDebug)
	}
}

func TestDefaultDataDir(t *testing.T) {
	app := New()

	expected := "test_default"

	app.DefaultDataDir(expected)
	if app.defaultDataDir != expected {
		t.Fatalf("Expected defaultDataDir %v, got %v", expected, app.defaultDataDir)
	}
}

func TestDefaultEncryptionEnv(t *testing.T) {
	app := New()

	expected := "test_env"

	app.DefaultEncryptionEnv(expected)
	if app.defaultEncryptionEnv != expected {
		t.Fatalf("Expected defaultEncryptionEnv %v, got %v", expected, app.defaultEncryptionEnv)
	}
}

func TestShowStartBanner(t *testing.T) {
	app := New()

	app.ShowStartBanner(true)
	if app.showStartBanner != true {
		t.Fatalf("Expected showStartBanner %v, got %v", true, app.showStartBanner)
	}

	app.ShowStartBanner(false)
	if app.showStartBanner != false {
		t.Fatalf("Expected showStartBanner %v, got %v", false, app.showStartBanner)
	}
}
