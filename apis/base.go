package apis

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
)

// StaticWildcardParam is the name of Static handler wildcard parameter.
const StaticWildcardParam = "path"

// NewRouter returns a new router instance loaded with the default app middlewares and api routes.
func NewRouter(app core.App) (*router.Router[*core.RequestEvent], error) {
	pbRouter := router.NewRouter(func(w http.ResponseWriter, r *http.Request) (*core.RequestEvent, router.EventCleanupFunc) {
		event := new(core.RequestEvent)
		event.Response = w
		event.Request = r
		event.App = app

		return event, nil
	})

	// register default middlewares
	pbRouter.Bind(activityLogger())
	pbRouter.Bind(panicRecover())
	pbRouter.Bind(rateLimit())
	pbRouter.Bind(loadAuthToken())
	pbRouter.Bind(securityHeaders())
	pbRouter.Bind(BodyLimit(DefaultMaxBodySize))

	apiGroup := pbRouter.Group("/api")
	bindSettingsApi(app, apiGroup)
	bindCollectionApi(app, apiGroup)
	bindRecordCrudApi(app, apiGroup)
	bindRecordAuthApi(app, apiGroup)
	bindLogsApi(app, apiGroup)
	bindBackupApi(app, apiGroup)
	bindCronApi(app, apiGroup)
	bindFileApi(app, apiGroup)
	bindBatchApi(app, apiGroup)
	bindRealtimeApi(app, apiGroup)
	bindHealthApi(app, apiGroup)

	return pbRouter, nil
}

// WrapStdHandler wraps Go [http.Handler] into a PocketBase handler func.
func WrapStdHandler(h http.Handler) func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		h.ServeHTTP(e.Response, e.Request)
		return nil
	}
}

// WrapStdMiddleware wraps Go [func(http.Handler) http.Handle] into a PocketBase middleware func.
func WrapStdMiddleware(m func(http.Handler) http.Handler) func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) (err error) {
		m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			e.Response = w
			e.Request = r
			err = e.Next()
		})).ServeHTTP(e.Response, e.Request)
		return err
	}
}

// MustSubFS returns an [fs.FS] corresponding to the subtree rooted at fsys's dir.
//
// This is similar to [fs.Sub] but panics on failure.
func MustSubFS(fsys fs.FS, dir string) fs.FS {
	dir = filepath.ToSlash(filepath.Clean(dir)) // ToSlash in case of Windows path

	sub, err := fs.Sub(fsys, dir)
	if err != nil {
		panic(fmt.Errorf("failed to create sub FS: %w", err))
	}

	return sub
}

// Static is a handler function to serve static directory content from fsys.
//
// If a file resource is missing and indexFallback is set, the request
// will be forwarded to the base index.html (useful for SPA with pretty urls).
//
// NB! Expects the route to have a "{path...}" wildcard parameter.
//
// Special redirects:
//   - if "path" is a file that ends in index.html, it is redirected to its non-index.html version (eg. /test/index.html -> /test/)
//   - if "path" is a directory that has index.html, the index.html file is rendered,
//     otherwise if missing - returns 404 or fallback to the root index.html if indexFallback is set
//
// Example:
//
//	fsys := os.DirFS("./pb_public")
//	router.GET("/files/{path...}", apis.Static(fsys, false))
func Static(fsys fs.FS, indexFallback bool) func(*core.RequestEvent) error {
	if fsys == nil {
		panic("Static: the provided fs.FS argument is nil")
	}

	return func(e *core.RequestEvent) error {
		// disable the activity logger to avoid flooding with messages
		//
		// note: errors are still logged
		if e.Get(requestEventKeySkipSuccessActivityLog) == nil {
			e.Set(requestEventKeySkipSuccessActivityLog, true)
		}

		filename := e.Request.PathValue(StaticWildcardParam)
		filename = filepath.ToSlash(filepath.Clean(strings.TrimPrefix(filename, "/")))

		// eagerly check for directory traversal
		//
		// note: this is just out of an abundance of caution because the fs.FS implementation could be non-std,
		// but usually shouldn't be necessary since os.DirFS.Open is expected to fail if the filename starts with dots
		if len(filename) > 2 && filename[0] == '.' && filename[1] == '.' && (filename[2] == '/' || filename[2] == '\\') {
			if indexFallback && filename != router.IndexPage {
				return e.FileFS(fsys, router.IndexPage)
			}
			return router.ErrFileNotFound
		}

		fi, err := fs.Stat(fsys, filename)
		if err != nil {
			if indexFallback && filename != router.IndexPage {
				return e.FileFS(fsys, router.IndexPage)
			}
			return router.ErrFileNotFound
		}

		if fi.IsDir() {
			// redirect to a canonical dir url, aka. with trailing slash
			if !strings.HasSuffix(e.Request.URL.Path, "/") {
				return e.Redirect(http.StatusMovedPermanently, safeRedirectPath(e.Request.URL.Path+"/"))
			}
		} else {
			urlPath := e.Request.URL.Path
			if strings.HasSuffix(urlPath, "/") {
				// redirect to a non-trailing slash file route
				urlPath = strings.TrimRight(urlPath, "/")
				if len(urlPath) > 0 {
					return e.Redirect(http.StatusMovedPermanently, safeRedirectPath(urlPath))
				}
			} else if stripped, ok := strings.CutSuffix(urlPath, router.IndexPage); ok {
				// redirect without the index.html
				return e.Redirect(http.StatusMovedPermanently, safeRedirectPath(stripped))
			}
		}

		fileErr := e.FileFS(fsys, filename)

		if fileErr != nil && indexFallback && filename != router.IndexPage && errors.Is(fileErr, router.ErrFileNotFound) {
			return e.FileFS(fsys, router.IndexPage)
		}

		return fileErr
	}
}

// safeRedirectPath normalizes the path string by replacing all beginning slashes
// (`\\`, `//`, `\/`) with a single forward slash to prevent open redirect attacks
func safeRedirectPath(path string) string {
	if len(path) > 1 && (path[0] == '\\' || path[0] == '/') && (path[1] == '\\' || path[1] == '/') {
		path = "/" + strings.TrimLeft(path, `/\`)
	}
	return path
}
