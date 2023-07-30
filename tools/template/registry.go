// Package template is a thin wrapper around the standard html/template
// and text/template packages that implements a convenient registry to
// load and cache templates on the fly concurrently.
//
// It was created to assist the JSVM plugin HTML rendering, but could be used in other Go code.
//
// Example:
//
// 	registry := template.NewRegistry()
//
// 	html1, err := registry.LoadFiles(
// 		// the files set wil be parsed only once and then cached
// 		"layout.html",
// 		"content.html",
// 	).Render(map[string]any{"name": "John"})
//
// 	html2, err := registry.LoadFiles(
// 		// reuse the already parsed and cached files set
// 		"layout.html",
// 		"content.html",
// 	).Render(map[string]any{"name": "Jane"})
package template

import (
	"fmt"
	"html/template"
	"io/fs"
	"strings"

	"github.com/pocketbase/pocketbase/tools/store"
)

// NewRegistry creates and initializes a new blank templates registry.
//
// Use the Registry.Load* methods to load templates into the registry.
func NewRegistry() *Registry {
	return &Registry{
		cache: store.New[*Renderer](nil),
	}
}

// Registry defines a templates registry that is safe to be used by multiple goroutines.
//
// Use the Registry.Load* methods to load templates into the registry.
type Registry struct {
	cache *store.Store[*Renderer]
}

// LoadFiles caches (if not already) the specified filenames set as a
// single template and returns a ready to use Renderer instance.
//
// There must be at least 1 filename specified.
func (r *Registry) LoadFiles(filenames ...string) *Renderer {
	key := strings.Join(filenames, ",")

	found := r.cache.Get(key)

	if found == nil {
		// parse and cache
		tpl, err := template.ParseFiles(filenames...)
		found = &Renderer{template: tpl, parseError: err}
		r.cache.Set(key, found)
	}

	return found
}

// LoadString caches (if not already) the specified inline string as a
// single template and returns a ready to use Renderer instance.
func (r *Registry) LoadString(text string) *Renderer {
	found := r.cache.Get(text)

	if found == nil {
		// parse and cache (using the text as key)
		tpl, err := template.New("").Parse(text)
		found = &Renderer{template: tpl, parseError: err}
		r.cache.Set(text, found)
	}

	return found
}

// LoadString caches (if not already) the specified fs and globPatterns
// pair as single template and returns a ready to use Renderer instance.
//
// There must be at least 1 file matching the provided globPattern(s)
// (note that most file names serves as glob patterns matching themselves).
func (r *Registry) LoadFS(fs fs.FS, globPatterns ...string) *Renderer {
	key := fmt.Sprintf("%v%v", fs, globPatterns)

	found := r.cache.Get(key)

	if found == nil {
		// parse and cache
		tpl, err := template.ParseFS(fs, globPatterns...)
		found = &Renderer{template: tpl, parseError: err}
		r.cache.Set(key, found)
	}

	return found
}
