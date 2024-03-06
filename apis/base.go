// Package apis implements the default PocketBase api services and middlewares.
package apis

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/rest"
	"github.com/pocketbase/pocketbase/ui"
	"github.com/spf13/cast"
)

const trailedAdminPath = "/_/"

// InitApi creates a configured echo instance with registered
// system and app specific routes and middlewares.
func InitApi(app core.App) (*echo.Echo, error) {
	e := echo.New()
	e.Debug = false
	e.Binder = &rest.MultiBinder{}
	e.JSONSerializer = &rest.Serializer{
		FieldsParam: fieldsQueryParam,
	}

	// configure a custom router
	e.ResetRouterCreator(func(ec *echo.Echo) echo.Router {
		return echo.NewRouter(echo.RouterConfig{
			UnescapePathParamValues: true,
			AllowOverwritingRoute:   true,
		})
	})

	// default middlewares
	e.Pre(middleware.RemoveTrailingSlashWithConfig(middleware.RemoveTrailingSlashConfig{
		Skipper: func(c echo.Context) bool {
			// enable by default only for the API routes
			return !strings.HasPrefix(c.Request().URL.Path, "/api/")
		},
	}))
	e.Pre(LoadAuthContext(app))
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(ContextExecStartKey, time.Now())

			return next(c)
		}
	})

	// custom error handler
	e.HTTPErrorHandler = func(c echo.Context, err error) {
		if err == nil {
			return // no error
		}

		var apiErr *ApiError

		if errors.As(err, &apiErr) {
			// already an api error...
		} else if v := new(echo.HTTPError); errors.As(err, &v) {
			msg := fmt.Sprintf("%v", v.Message)
			apiErr = NewApiError(v.Code, msg, v)
		} else {
			if errors.Is(err, sql.ErrNoRows) {
				apiErr = NewNotFoundError("", err)
			} else {
				apiErr = NewBadRequestError("", err)
			}
		}

		logRequest(app, c, apiErr)

		if c.Response().Committed {
			return // already committed
		}

		event := new(core.ApiErrorEvent)
		event.HttpContext = c
		event.Error = apiErr

		// send error response
		hookErr := app.OnBeforeApiError().Trigger(event, func(e *core.ApiErrorEvent) error {
			if e.HttpContext.Response().Committed {
				return nil
			}

			// @see https://github.com/labstack/echo/issues/608
			if e.HttpContext.Request().Method == http.MethodHead {
				return e.HttpContext.NoContent(apiErr.Code)
			}

			return e.HttpContext.JSON(apiErr.Code, apiErr)
		})

		if hookErr == nil {
			if err := app.OnAfterApiError().Trigger(event); err != nil {
				app.Logger().Debug("OnAfterApiError failure", slog.String("error", err.Error()))
			}
		} else {
			app.Logger().Debug("OnBeforeApiError error (truly rare case, eg. client already disconnected)", slog.String("error", hookErr.Error()))
		}
	}

	// admin ui routes
	bindStaticAdminUI(app, e)

	// default routes
	api := e.Group("/api", eagerRequestInfoCache(app))
	bindSettingsApi(app, api)
	bindAdminApi(app, api)
	bindCollectionApi(app, api)
	bindRecordCrudApi(app, api)
	bindRecordAuthApi(app, api)
	bindFileApi(app, api)
	bindRealtimeApi(app, api)
	bindLogsApi(app, api)
	bindHealthApi(app, api)
	bindBackupApi(app, api)

	// catch all any route
	api.Any("/*", func(c echo.Context) error {
		return echo.ErrNotFound
	}, ActivityLogger(app))

	return e, nil
}

// StaticDirectoryHandler is similar to `echo.StaticDirectoryHandler`
// but without the directory redirect which conflicts with RemoveTrailingSlash middleware.
//
// If a file resource is missing and indexFallback is set, the request
// will be forwarded to the base index.html (useful also for SPA).
//
// @see https://github.com/labstack/echo/issues/2211
func StaticDirectoryHandler(fileSystem fs.FS, indexFallback bool) echo.HandlerFunc {
	return func(c echo.Context) error {
		p := c.PathParam("*")

		// escape url path
		tmpPath, err := url.PathUnescape(p)
		if err != nil {
			return fmt.Errorf("failed to unescape path variable: %w", err)
		}
		p = tmpPath

		// fs.FS.Open() already assumes that file names are relative to FS root path and considers name with prefix `/` as invalid
		name := filepath.ToSlash(filepath.Clean(strings.TrimPrefix(p, "/")))

		fileErr := c.FileFS(name, fileSystem)

		if fileErr != nil && indexFallback && errors.Is(fileErr, echo.ErrNotFound) {
			return c.FileFS("index.html", fileSystem)
		}

		return fileErr
	}
}

// bindStaticAdminUI registers the endpoints that serves the static admin UI.
func bindStaticAdminUI(app core.App, e *echo.Echo) error {
	// redirect to trailing slash to ensure that relative urls will still work properly
	e.GET(
		strings.TrimRight(trailedAdminPath, "/"),
		func(c echo.Context) error {
			return c.Redirect(http.StatusTemporaryRedirect, strings.TrimLeft(trailedAdminPath, "/"))
		},
	)

	// serves static files from the /ui/dist directory
	// (similar to echo.StaticFS but with gzip middleware enabled)
	e.GET(
		trailedAdminPath+"*",
		echo.StaticDirectoryHandler(ui.DistDirFS, false),
		installerRedirect(app),
		uiCacheControl(),
		middleware.Gzip(),
	)

	return nil
}

func uiCacheControl() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// add default Cache-Control header for all Admin UI resources
			// (ignoring the root admin path)
			if c.Request().URL.Path != trailedAdminPath {
				c.Response().Header().Set("Cache-Control", "max-age=1209600, stale-while-revalidate=86400")
			}

			return next(c)
		}
	}
}

const hasAdminsCacheKey = "@hasAdmins"

func updateHasAdminsCache(app core.App) error {
	total, err := app.Dao().TotalAdmins()
	if err != nil {
		return err
	}

	app.Store().Set(hasAdminsCacheKey, total > 0)

	return nil
}

// installerRedirect redirects the user to the installer admin UI page
// when the application needs some preliminary configurations to be done.
func installerRedirect(app core.App) echo.MiddlewareFunc {
	// keep hasAdminsCacheKey value up-to-date
	app.OnAdminAfterCreateRequest().Add(func(data *core.AdminCreateEvent) error {
		return updateHasAdminsCache(app)
	})

	app.OnAdminAfterDeleteRequest().Add(func(data *core.AdminDeleteEvent) error {
		return updateHasAdminsCache(app)
	})

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// skip redirect checks for non-root level index.html requests
			path := c.Request().URL.Path
			if path != trailedAdminPath && path != trailedAdminPath+"index.html" {
				return next(c)
			}

			hasAdmins := cast.ToBool(app.Store().Get(hasAdminsCacheKey))

			if !hasAdmins {
				// update the cache to make sure that the admin wasn't created by another process
				if err := updateHasAdminsCache(app); err != nil {
					return err
				}
				hasAdmins = cast.ToBool(app.Store().Get(hasAdminsCacheKey))
			}

			_, hasInstallerParam := c.Request().URL.Query()["installer"]

			if !hasAdmins && !hasInstallerParam {
				// redirect to the installer page
				return c.Redirect(http.StatusTemporaryRedirect, "?installer#")
			}

			if hasAdmins && hasInstallerParam {
				// clear the installer param
				return c.Redirect(http.StatusTemporaryRedirect, "?")
			}

			return next(c)
		}
	}
}
