package osutils_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/osutils"
	"github.com/pocketbase/pocketbase/tools/security"
)

func TestMoveDirContent(t *testing.T) {
	testDir := createTestDir(t)
	defer os.RemoveAll(testDir)

	exclude := []string{
		"missing",
		"test2",
		"b",
	}

	// missing dest path
	// ---
	dir1 := filepath.Join(filepath.Dir(testDir), "a", "b", "c", "d", "_pb_move_dir_content_test_"+security.PseudorandomString(4))
	defer os.RemoveAll(dir1)

	if err := osutils.MoveDirContent(testDir, dir1, exclude...); err == nil {
		t.Fatal("Expected path error, got nil")
	}

	// existing parent dir
	// ---
	dir2 := filepath.Join(filepath.Dir(testDir), "_pb_move_dir_content_test_"+security.PseudorandomString(4))
	defer os.RemoveAll(dir2)

	if err := osutils.MoveDirContent(testDir, dir2, exclude...); err != nil {
		t.Fatalf("Expected dir2 to be created, got error: %v", err)
	}

	// find all files
	files := []string{}
	filepath.WalkDir(dir2, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		files = append(files, path)

		return nil
	})

	expectedFiles := []string{
		filepath.Join(dir2, "test1"),
		filepath.Join(dir2, "a", "a1"),
		filepath.Join(dir2, "a", "a2"),
	}

	if len(files) != len(expectedFiles) {
		t.Fatalf("Expected %d files, got %d: \n%v", len(expectedFiles), len(files), files)
	}

	for _, expected := range expectedFiles {
		if !list.ExistInSlice(expected, files) {
			t.Fatalf("Missing expected file %q in \n%v", expected, files)
		}
	}
}

// -------------------------------------------------------------------

// note: make sure to call os.RemoveAll(dir) after you are done
// working with the created test dir.
func createTestDir(t *testing.T) string {
	dir, err := os.MkdirTemp(os.TempDir(), "test_dir")
	if err != nil {
		t.Fatal(err)
	}

	// create sub directories
	if err := os.MkdirAll(filepath.Join(dir, "a"), os.ModePerm); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "b"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	{
		f, err := os.OpenFile(filepath.Join(dir, "test1"), os.O_WRONLY|os.O_CREATE, 0644)
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
		f, err := os.OpenFile(filepath.Join(dir, "a/a1"), os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			t.Fatal(err)
		}
		f.Close()
	}

	{
		f, err := os.OpenFile(filepath.Join(dir, "a/a2"), os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			t.Fatal(err)
		}
		f.Close()
	}

	{
		f, err := os.OpenFile(filepath.Join(dir, "b/b2"), os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			t.Fatal(err)
		}
		f.Close()
	}

	{
		f, err := os.OpenFile(filepath.Join(dir, "b/b2"), os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			t.Fatal(err)
		}
		f.Close()
	}

	return dir
}
