package update

import "testing"

func TestReleaseFindFileBySuffix(t *testing.T) {
	r := release{
		Files: []string{
			"test1.zip",
			"test2.zip",
			"test22.zip",
			"test3.zip",
		},
	}

	filename, err := r.findFileBySuffix("2.zip")
	if err != nil {
		t.Fatalf("Expected nil, got err: %v", err)
	}

	expected := "test2.zip"
	if filename != expected {
		t.Fatalf("Expected file %q, got %q", expected, filename)
	}
}

func TestReleaseBuildFilePath(t *testing.T) {
	r := release{
		Id:             "test id",
		CollectionName: "test collection",
	}

	result := r.buildFilePath("test file")

	expected := "/api/files/test%20collection/test%20id/test%20file"
	if result != expected {
		t.Fatalf("Expected path %q, got %q", expected, result)
	}
}
