package archive_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pocketbase/pocketbase/tools/archive"
)

func TestExtractFailure(t *testing.T) {
	testDir := createTestDir(t)
	defer os.RemoveAll(testDir)

	missingZipPath := filepath.Join(os.TempDir(), "pb_missing_test.zip")
	extractPath := filepath.Join(os.TempDir(), "pb_zip_extract")
	defer os.RemoveAll(extractPath)

	if err := archive.Extract(missingZipPath, extractPath); err == nil {
		t.Fatal("Expected Extract to fail due to missing zipPath")
	}

	if _, err := os.Stat(extractPath); err == nil {
		t.Fatalf("Expected %q to not be created", extractPath)
	}
}

func TestExtractSuccess(t *testing.T) {
	testDir := createTestDir(t)
	defer os.RemoveAll(testDir)

	zipPath := filepath.Join(os.TempDir(), "pb_test.zip")
	defer os.RemoveAll(zipPath)

	extractPath := filepath.Join(os.TempDir(), "pb_zip_extract")
	defer os.RemoveAll(extractPath)

	// zip testDir content
	if err := archive.Create(testDir, zipPath); err != nil {
		t.Fatalf("Failed to create archive: %v", err)
	}

	if err := archive.Extract(zipPath, extractPath); err != nil {
		t.Fatalf("Failed to extract %q in %q", zipPath, extractPath)
	}

	pathsToCheck := []string{
		filepath.Join(extractPath, "a/sub1.txt"),
		filepath.Join(extractPath, "a/b/c/sub2.txt"),
	}

	for _, p := range pathsToCheck {
		if _, err := os.Stat(p); err != nil {
			t.Fatalf("Failed to retrieve extracted file %q: %v", p, err)
		}
	}
}
