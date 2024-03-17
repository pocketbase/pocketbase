package archive_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pocketbase/pocketbase/tools/archive"
)

func TestCreateFailure(t *testing.T) {
	testDir := createTestDir(t)
	defer os.RemoveAll(testDir)

	zipPath := filepath.Join(os.TempDir(), "pb_test.zip")
	defer os.RemoveAll(zipPath)

	missingDir := filepath.Join(os.TempDir(), "missing")

	if err := archive.Create(missingDir, zipPath); err == nil {
		t.Fatal("Expected to fail due to missing directory or file")
	}

	if _, err := os.Stat(zipPath); err == nil {
		t.Fatalf("Expected the zip file not to be created")
	}
}

func TestCreateSuccess(t *testing.T) {
	testDir := createTestDir(t)
	defer os.RemoveAll(testDir)

	zipName := "pb_test.zip"
	zipPath := filepath.Join(os.TempDir(), zipName)
	defer os.RemoveAll(zipPath)

	// zip testDir content (excluding test and a/b/c dir)
	if err := archive.Create(testDir, zipPath, "a/b/c", "test"); err != nil {
		t.Fatalf("Failed to create archive: %v", err)
	}

	info, err := os.Stat(zipPath)
	if err != nil {
		t.Fatalf("Failed to retrieve the generated zip file: %v", err)
	}

	if name := info.Name(); name != zipName {
		t.Fatalf("Expected zip with name %q, got %q", zipName, name)
	}

	expectedSize := int64(544)
	if size := info.Size(); size != expectedSize {
		t.Fatalf("Expected zip with size %d, got %d", expectedSize, size)
	}
}

// -------------------------------------------------------------------

// note: make sure to call os.RemoveAll(dir) after you are done
// working with the created test dir.
func createTestDir(t *testing.T) string {
	dir, err := os.MkdirTemp(os.TempDir(), "pb_zip_test")
	if err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(dir, "a/b/c"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	{
		f, err := os.OpenFile(filepath.Join(dir, "test"), os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			t.Fatal(err)
		}
		f.Close()
	}

	{
		f, err := os.OpenFile(filepath.Join(dir, "test2"), os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			t.Fatal(err)
		}
		f.Close()
	}

	{
		f, err := os.OpenFile(filepath.Join(dir, "a/test"), os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			t.Fatal(err)
		}
		f.Close()
	}

	{
		f, err := os.OpenFile(filepath.Join(dir, "a/b/sub1"), os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			t.Fatal(err)
		}
		f.Close()
	}

	{
		f, err := os.OpenFile(filepath.Join(dir, "a/b/c/sub2"), os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			t.Fatal(err)
		}
		f.Close()
	}

	{
		f, err := os.OpenFile(filepath.Join(dir, "a/b/c/sub3"), os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			t.Fatal(err)
		}
		f.Close()
	}

	// symbolic link
	if err := os.Symlink(filepath.Join(dir, "test"), filepath.Join(dir, "test_symlink")); err != nil {
		t.Fatal(err)
	}

	return dir
}
