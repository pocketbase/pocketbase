package jsvm

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/pocketbase/pocketbase/tools/list"
)

func TestFormDataAppendAndSet(t *testing.T) {
	t.Parallel()

	data := FormData{}

	data.Append("a", 1)
	data.Append("a", 2)

	data.Append("b", 3)
	data.Append("b", 4)
	data.Set("b", 5) // should overwrite the previous 2 calls

	data.Set("c", 6)
	data.Set("c", 7)

	if len(data["a"]) != 2 {
		t.Fatalf("Expected 2 'a' values, got %v", data["a"])
	}
	if data["a"][0] != 1 || data["a"][1] != 2 {
		t.Fatalf("Expected 1 and 2 'a' key values, got %v", data["a"])
	}

	if len(data["b"]) != 1 {
		t.Fatalf("Expected 1 'b' values, got %v", data["b"])
	}
	if data["b"][0] != 5 {
		t.Fatalf("Expected 5 as 'b' key value, got %v", data["b"])
	}

	if len(data["c"]) != 1 {
		t.Fatalf("Expected 1 'c' values, got %v", data["c"])
	}
	if data["c"][0] != 7 {
		t.Fatalf("Expected 7 as 'c' key value, got %v", data["c"])
	}
}

func TestFormDataDelete(t *testing.T) {
	t.Parallel()

	data := FormData{}
	data.Append("a", 1)
	data.Append("a", 2)
	data.Append("b", 3)

	data.Delete("missing") // should do nothing
	data.Delete("a")

	if len(data) != 1 {
		t.Fatalf("Expected exactly 1 data remaining key, got %v", data)
	}

	if data["b"][0] != 3 {
		t.Fatalf("Expected 3 as 'b' key value, got %v", data["b"])
	}
}

func TestFormDataGet(t *testing.T) {
	t.Parallel()

	data := FormData{}
	data.Append("a", 1)
	data.Append("a", 2)

	if v := data.Get("missing"); v != nil {
		t.Fatalf("Expected %v for key 'missing', got %v", nil, v)
	}

	if v := data.Get("a"); v != 1 {
		t.Fatalf("Expected %v for key 'a', got %v", 1, v)
	}
}

func TestFormDataGetAll(t *testing.T) {
	t.Parallel()

	data := FormData{}
	data.Append("a", 1)
	data.Append("a", 2)

	if v := data.GetAll("missing"); v != nil {
		t.Fatalf("Expected %v for key 'a', got %v", nil, v)
	}

	values := data.GetAll("a")
	if len(values) != 2 || values[0] != 1 || values[1] != 2 {
		t.Fatalf("Expected 1 and 2 values, got %v", values)
	}
}

func TestFormDataHas(t *testing.T) {
	t.Parallel()

	data := FormData{}
	data.Append("a", 1)

	if v := data.Has("missing"); v {
		t.Fatalf("Expected key 'missing' to not exist: %v", v)
	}

	if v := data.Has("a"); !v {
		t.Fatalf("Expected key 'a' to exist: %v", v)
	}
}

func TestFormDataKeys(t *testing.T) {
	t.Parallel()

	data := FormData{}
	data.Append("a", 1)
	data.Append("b", 1)
	data.Append("c", 1)
	data.Append("a", 1)

	keys := data.Keys()

	expectedKeys := []string{"a", "b", "c"}

	for _, expected := range expectedKeys {
		if !list.ExistInSlice(expected, keys) {
			t.Fatalf("Expected key %s to exists in %v", expected, keys)
		}
	}
}

func TestFormDataValues(t *testing.T) {
	t.Parallel()

	data := FormData{}
	data.Append("a", 1)
	data.Append("b", 2)
	data.Append("c", 3)
	data.Append("a", 4)

	values := data.Values()

	expectedKeys := []any{1, 2, 3, 4}

	for _, expected := range expectedKeys {
		if !list.ExistInSlice(expected, values) {
			t.Fatalf("Expected key %s to exists in %v", expected, values)
		}
	}
}

func TestFormDataEntries(t *testing.T) {
	t.Parallel()

	data := FormData{}
	data.Append("a", 1)
	data.Append("b", 2)
	data.Append("c", 3)
	data.Append("a", 4)

	entries := data.Entries()

	rawEntries, err := json.Marshal(entries)
	if err != nil {
		t.Fatal(err)
	}

	if len(entries) != 4 {
		t.Fatalf("Expected 4 entries")
	}

	expectedEntries := []string{`["a",1]`, `["a",4]`, `["b",2]`, `["c",3]`}
	for _, expected := range expectedEntries {
		if !bytes.Contains(rawEntries, []byte(expected)) {
			t.Fatalf("Expected entry %s to exists in %s", expected, rawEntries)
		}
	}
}

func TestFormDataToMultipart(t *testing.T) {
	t.Parallel()

	f, err := filesystem.NewFileFromBytes([]byte("abc"), "test")
	if err != nil {
		t.Fatal(err)
	}

	data := FormData{}
	data.Append("a", 1) // should be casted
	data.Append("b", "test1")
	data.Append("b", "test2")
	data.Append("c", f)

	body, mp, err := data.toMultipart()
	if err != nil {
		t.Fatal(err)
	}
	bodyStr := body.String()

	// content type checks
	contentType := mp.FormDataContentType()
	expectedContentType := "multipart/form-data; boundary="
	if !strings.Contains(contentType, expectedContentType) {
		t.Fatalf("Expected to find content-type %s in %s", expectedContentType, contentType)
	}

	// body checks
	expectedBodyParts := []string{
		"Content-Disposition: form-data; name=\"a\"\r\n\r\n1",
		"Content-Disposition: form-data; name=\"b\"\r\n\r\ntest1",
		"Content-Disposition: form-data; name=\"b\"\r\n\r\ntest2",
		"Content-Disposition: form-data; name=\"c\"; filename=\"test\"\r\nContent-Type: application/octet-stream\r\n\r\nabc",
	}
	for _, part := range expectedBodyParts {
		if !strings.Contains(bodyStr, part) {
			t.Fatalf("Expected to find %s in body\n%s", part, bodyStr)
		}
	}
}
