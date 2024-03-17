package archive_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/pocketbase/pocketbase/tools/archive"
)

func TestExtractFailure(t *testing.T) {
	testDir := createTestDir(t)
	defer os.RemoveAll(testDir)

	missingZipPath := filepath.Join(os.TempDir(), "pb_missing_test.zip")
	extractedPath := filepath.Join(os.TempDir(), "pb_zip_extract")
	defer os.RemoveAll(extractedPath)

	if err := archive.Extract(missingZipPath, extractedPath); err == nil {
		t.Fatal("Expected Extract to fail due to missing zipPath")
	}

	if _, err := os.Stat(extractedPath); err == nil {
		t.Fatalf("Expected %q to not be created", extractedPath)
	}
}

func TestExtractSuccess(t *testing.T) {
	testDir := createTestDir(t)
	defer os.RemoveAll(testDir)

	zipPath := filepath.Join(os.TempDir(), "pb_test.zip")
	defer os.RemoveAll(zipPath)

	extractedPath := filepath.Join(os.TempDir(), "pb_zip_extract")
	defer os.RemoveAll(extractedPath)

	// zip testDir content (with exclude)
	if err := archive.Create(testDir, zipPath, "a/b/c", "test2", "sub2"); err != nil {
		t.Fatalf("Failed to create archive: %v", err)
	}

	if err := archive.Extract(zipPath, extractedPath); err != nil {
		t.Fatalf("Failed to extract %q in %q", zipPath, extractedPath)
	}

	availableFiles := []string{}

	walkErr := filepath.WalkDir(extractedPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		availableFiles = append(availableFiles, path)

		return nil
	})
	if walkErr != nil {
		t.Fatalf("Failed to read the extracted dir: %v", walkErr)
	}

	// (note: symbolic links and other regular files should be missing)
	expectedFiles := []string{
		filepath.Join(extractedPath, "test"),
		filepath.Join(extractedPath, "a/test"),
		filepath.Join(extractedPath, "a/b/sub1"),
	}

	if len(availableFiles) != len(expectedFiles) {
		t.Fatalf("Expected \n%v, \ngot \n%v", expectedFiles, availableFiles)
	}

ExpectedLoop:
	for _, expected := range expectedFiles {
		for _, available := range availableFiles {
			if available == expected {
				continue ExpectedLoop
			}
		}

		t.Fatalf("Missing file %q in \n%v", expected, availableFiles)
	}
}
