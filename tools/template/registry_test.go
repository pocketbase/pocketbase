package template

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func checkRegistryFuncs(t *testing.T, r *Registry, expectedFuncs ...string) {
	if v := len(r.funcs); v != len(expectedFuncs) {
		t.Fatalf("Expected total %d funcs, got %d", len(expectedFuncs), v)
	}

	for _, name := range expectedFuncs {
		if _, ok := r.funcs[name]; !ok {
			t.Fatalf("Missing %q func", name)
		}
	}
}

func TestNewRegistry(t *testing.T) {
	r := NewRegistry()

	if r.cache == nil {
		t.Fatalf("Expected cache store to be initialized, got nil")
	}

	if v := r.cache.Length(); v != 0 {
		t.Fatalf("Expected cache store length to be 0, got %d", v)
	}

	checkRegistryFuncs(t, r, "raw")
}

func TestRegistryAddFuncs(t *testing.T) {
	r := NewRegistry()

	r.AddFuncs(map[string]any{
		"test": func(a string) string { return a + "-TEST" },
	})

	checkRegistryFuncs(t, r, "raw", "test")

	result, err := r.LoadString(`{{.|test}}`).Render("example")
	if err != nil {
		t.Fatalf("Unexpected Render() error, got %v", err)
	}

	expected := "example-TEST"
	if result != expected {
		t.Fatalf("Expected Render() result %q, got %q", expected, result)
	}
}

func TestRegistryLoadFiles(t *testing.T) {
	r := NewRegistry()

	t.Run("invalid or missing files", func(t *testing.T) {
		r.LoadFiles("file1.missing", "file2.missing")

		key := "file1.missing,file2.missing"
		renderer := r.cache.Get(key)

		if renderer == nil {
			t.Fatal("Expected renderer to be initialized even if invalid, got nil")
		}

		if renderer.template != nil {
			t.Fatalf("Expected renderer template to be nil, got %v", renderer.template)
		}

		if renderer.parseError == nil {
			t.Fatalf("Expected renderer parseError to be set, got nil")
		}
	})

	t.Run("valid files", func(t *testing.T) {
		// create test templates
		dir, err := os.MkdirTemp(os.TempDir(), "template_test")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "base.html"), []byte(`Base:{{template "content" .}}`), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "content.html"), []byte(`{{define "content"}}Content:{{.|raw}}{{end}}`), 0644); err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(dir)

		files := []string{filepath.Join(dir, "base.html"), filepath.Join(dir, "content.html")}

		r.LoadFiles(files...)

		renderer := r.cache.Get(strings.Join(files, ","))

		if renderer == nil {
			t.Fatal("Expected renderer to be initialized even if invalid, got nil")
		}

		if renderer.template == nil {
			t.Fatal("Expected renderer template to be set, got nil")
		}

		if renderer.parseError != nil {
			t.Fatalf("Expected renderer parseError to be nil, got %v", renderer.parseError)
		}

		result, err := renderer.Render("<h1>123</h1>")
		if err != nil {
			t.Fatalf("Unexpected Render() error, got %v", err)
		}

		expected := "Base:Content:<h1>123</h1>"
		if result != expected {
			t.Fatalf("Expected Render() result %q, got %q", expected, result)
		}
	})
}

func TestRegistryLoadString(t *testing.T) {
	r := NewRegistry()

	t.Run("invalid template string", func(t *testing.T) {
		txt := `test {{define "content"}}`

		r.LoadString(txt)

		renderer := r.cache.Get(txt)

		if renderer == nil {
			t.Fatal("Expected renderer to be initialized even if invalid, got nil")
		}

		if renderer.template != nil {
			t.Fatalf("Expected renderer template to be nil, got %v", renderer.template)
		}

		if renderer.parseError == nil {
			t.Fatalf("Expected renderer parseError to be set, got nil")
		}
	})

	t.Run("valid template string", func(t *testing.T) {
		txt := `test {{.|raw}}`

		r.LoadString(txt)

		renderer := r.cache.Get(txt)

		if renderer == nil {
			t.Fatal("Expected renderer to be initialized even if invalid, got nil")
		}

		if renderer.template == nil {
			t.Fatal("Expected renderer template to be set, got nil")
		}

		if renderer.parseError != nil {
			t.Fatalf("Expected renderer parseError to be nil, got %v", renderer.parseError)
		}

		result, err := renderer.Render("<h1>123</h1>")
		if err != nil {
			t.Fatalf("Unexpected Render() error, got %v", err)
		}

		expected := "test <h1>123</h1>"
		if result != expected {
			t.Fatalf("Expected Render() result %q, got %q", expected, result)
		}
	})
}

func TestRegistryLoadFS(t *testing.T) {
	r := NewRegistry()

	t.Run("invalid fs", func(t *testing.T) {
		fs := os.DirFS("__missing__")

		files := []string{"missing1", "missing2"}

		key := fmt.Sprintf("%v%v", fs, files)

		r.LoadFS(fs, files...)

		renderer := r.cache.Get(key)

		if renderer == nil {
			t.Fatal("Expected renderer to be initialized even if invalid, got nil")
		}

		if renderer.template != nil {
			t.Fatalf("Expected renderer template to be nil, got %v", renderer.template)
		}

		if renderer.parseError == nil {
			t.Fatalf("Expected renderer parseError to be set, got nil")
		}
	})

	t.Run("valid fs", func(t *testing.T) {
		// create test templates
		dir, err := os.MkdirTemp(os.TempDir(), "template_test2")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "base.html"), []byte(`Base:{{template "content" .}}`), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "content.html"), []byte(`{{define "content"}}Content:{{.|raw}}{{end}}`), 0644); err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(dir)

		fs := os.DirFS(dir)

		files := []string{"base.html", "content.html"}

		key := fmt.Sprintf("%v%v", fs, files)

		r.LoadFS(fs, files...)

		renderer := r.cache.Get(key)

		if renderer == nil {
			t.Fatal("Expected renderer to be initialized even if invalid, got nil")
		}

		if renderer.template == nil {
			t.Fatal("Expected renderer template to be set, got nil")
		}

		if renderer.parseError != nil {
			t.Fatalf("Expected renderer parseError to be nil, got %v", renderer.parseError)
		}

		result, err := renderer.Render("<h1>123</h1>")
		if err != nil {
			t.Fatalf("Unexpected Render() error, got %v", err)
		}

		expected := "Base:Content:<h1>123</h1>"
		if result != expected {
			t.Fatalf("Expected Render() result %q, got %q", expected, result)
		}
	})
}
