// Package template is a thin wrapper around the standard html/template
// and text/template packages that implements a convenient registry to
// load and cache templates on the fly concurrently.
//
// It was created to assist the JSVM plugin HTML rendering, but could be used in other Go code.
//
// Example:
//
//	registry := template.NewRegistry()
//
//	html1, err := registry.LoadFiles([]string{
//		// the files set will be parsed only once and then cached
//		"layout.html",
//		"content.html",
//	}).Render(map[string]any{"name": "John"})
//
//	html2, err := registry.LoadFiles([]string{
//		// reuse the already parsed and cached files set
//		"layout.html",
//		"content.html",
//	}).Render(map[string]any{"name": "Jane"})
//
//	html3, err := registry.LoadFiles([]string{
//		// newly parsed files with cache update
//		"layout.html",
//		"content.html",
//	}, false).Render(map[string]any{"name": "Juan"})
package template

import (
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/pocketbase/pocketbase/tools/store"
)

// NewRegistry creates and initializes a new templates registry with
// some defaults (eg. global "raw" template function for unescaped HTML).
//
// Use the Registry.Load* methods to load templates into the registry.
func NewRegistry() *Registry {
	return &Registry{
		cache: store.New[*Renderer](nil),
		forceCache: true,
		funcs: template.FuncMap{
			"raw": func(str string) template.HTML {
				return template.HTML(str)
			},
		},
	}
}


// Registry defines a templates registry that is safe to be used by multiple goroutines.
//
// Use the Registry.Load* methods to load templates into the registry.
type Registry struct {
	cache *store.Store[*Renderer]
	forceCache bool
	funcs template.FuncMap
}

// AddFuncs registers new global template functions.
//
// The key of each map entry is the function name that will be used in the templates.
// If a function with the map entry name already exists it will be replaced with the new one.
//
// The value of each map entry is a function that must have either a
// single return value, or two return values of which the second has type error.
//
// Example:
//
//	r.AddFuncs(map[string]any{
//	  "toUpper": func(str string) string {
//	      return strings.ToUppser(str)
//	  },
//	  ...
//	})
func (r *Registry) AddFuncs(funcs map[string]any) *Registry {
	for name, f := range funcs {
		r.funcs[name] = f
	}

	return r
}

// ForceCache allows to turn on/off systematic reuse from cache.
// Default Registry forceCache state is true.
// 
// In any case cache is updated on each Registry.Load* methods invocation.
func (r *Registry) ForceCache(state bool) *Registry {
	r.forceCache = state
	return r
}

// ClearCache clears Registry cache
func (r *Registry) ClearCache() *Registry {
	r.cache = store.New[*Renderer](nil)
	return r
}

// LoadFiles caches (if not already) the specified filenames set as a
// single template and returns a ready to use Renderer instance.
//
// There must be at least 1 filename specified.
func (r *Registry) LoadFiles(filenames []string, requestCache ...bool) *Renderer {
	key := strings.Join(filenames, ",")

	requestFromCache := r.forceCache

	if len(requestCache) > 0 {
			requestFromCache = requestCache[0]
	}

	var found *Renderer

	if requestFromCache {
		found = r.cache.Get(key)
	}

	if found == nil {
		// parse and cache
		tpl, err := template.New(filepath.Base(filenames[0])).Funcs(r.funcs).ParseFiles(filenames...)
		found = &Renderer{template: tpl, parseError: err}
		r.cache.Set(key, found)
	}

	return found
}

// LoadString caches (if not already) the specified inline string as a
// single template and returns a ready to use Renderer instance.
func (r *Registry) LoadString(text string, requestCache ...bool) *Renderer {
	requestFromCache := r.forceCache

	if len(requestCache) > 0 {
		requestFromCache = requestCache[0]
	}

	var found *Renderer

	if requestFromCache {
		found = r.cache.Get(text)
	}

	if found == nil {
		// parse and cache (using the text as key)
		tpl, err := template.New("").Funcs(r.funcs).Parse(text)
		found = &Renderer{template: tpl, parseError: err}
		r.cache.Set(text, found)
	}

	return found
}

// LoadFS caches (if not already) the specified fs and globPatterns
// pair as single template and returns a ready to use Renderer instance.
//
// There must be at least 1 file matching the provided globPattern(s)
// (note that most file names serves as glob patterns matching themselves).
func (r *Registry) LoadFS(fsys fs.FS, globPatterns []string, requestCache ...bool) *Renderer {
	key := fmt.Sprintf("%v%v", fsys, globPatterns)

	requestFromCache := r.forceCache

	if len(requestCache) > 0 {
		requestFromCache = requestCache[0]
	}

	var found *Renderer

	if requestFromCache {
		found = r.cache.Get(key)
	}

	if found == nil {
		// find the first file to use as template name (it is required when specifying Funcs)
		var firstFilename string
		if len(globPatterns) > 0 {
			list, _ := fs.Glob(fsys, globPatterns[0])
			if len(list) > 0 {
				firstFilename = filepath.Base(list[0])
			}
		}

		tpl, err := template.New(firstFilename).Funcs(r.funcs).ParseFS(fsys, globPatterns...)
		found = &Renderer{template: tpl, parseError: err}
		r.cache.Set(key, found)
	}

	return found
}
