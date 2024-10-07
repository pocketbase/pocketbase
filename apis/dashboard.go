package apis

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
)

const installerParam = "pbinstal"

var wildcardPlaceholderRegex = regexp.MustCompile(`/{.+\.\.\.}$`)

func stripWildcard(pattern string) string {
	return wildcardPlaceholderRegex.ReplaceAllString(pattern, "")
}

// installerRedirect redirects the user to the installer dashboard UI page
// when the application needs some preliminary configurations to be done.
func installerRedirect(app core.App, cpPath string) func(*core.RequestEvent) error {
	// note: to avoid locks contention it is not concurrent safe but it
	// is expected to be updated only once during initialization
	var hasSuperuser bool

	// strip named wildcard
	cpPath = stripWildcard(cpPath)

	updateHasSuperuser := func(app core.App) error {
		total, err := app.CountRecords(core.CollectionNameSuperusers)
		if err != nil {
			return err
		}

		hasSuperuser = total > 0

		return nil
	}

	// load initial state on app init
	app.OnBootstrap().BindFunc(func(e *core.BootstrapEvent) error {
		err := e.Next()
		if err != nil {
			return err
		}

		err = updateHasSuperuser(e.App)
		if err != nil {
			return fmt.Errorf("failed to check for existing superuser: %w", err)
		}

		return nil
	})

	// update on superuser create
	app.OnRecordCreateRequest(core.CollectionNameSuperusers).BindFunc(func(e *core.RecordRequestEvent) error {
		err := e.Next()
		if err != nil {
			return err
		}

		if !hasSuperuser {
			hasSuperuser = true
		}

		return nil
	})

	return func(e *core.RequestEvent) error {
		if hasSuperuser {
			return e.Next()
		}

		isAPI := strings.HasPrefix(e.Request.URL.Path, "/api/")
		isControlPanel := strings.HasPrefix(e.Request.URL.Path, cpPath)
		wildcard := e.Request.PathValue(StaticWildcardParam)

		// skip redirect checks for API and non-root level dashboard index.html requests (css, images, etc.)
		if isAPI || (isControlPanel && wildcard != "" && wildcard != router.IndexPage) {
			return e.Next()
		}

		// check again in case the superuser was created by some other process
		if err := updateHasSuperuser(e.App); err != nil {
			return err
		}

		if hasSuperuser {
			return e.Next()
		}

		_, hasInstallerParam := e.Request.URL.Query()[installerParam]

		// redirect to the installer page
		if !hasInstallerParam {
			return e.Redirect(http.StatusTemporaryRedirect, cpPath+"?"+installerParam+"#")
		}

		return e.Next()
	}
}

// dashboardRemoveInstallerParam redirects to a non-installer
// query param in case there is already a superuser created.
//
// Note: intended to be registered only for the dashboard route
// to prevent excessive checks for every other route in installerRedirect.
func dashboardRemoveInstallerParam() func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		_, hasInstallerParam := e.Request.URL.Query()[installerParam]
		if !hasInstallerParam {
			return e.Next() // nothing to remove
		}

		// clear installer param
		total, _ := e.App.CountRecords(core.CollectionNameSuperusers)
		if total > 0 {
			return e.Redirect(http.StatusTemporaryRedirect, "?")
		}

		return e.Next()
	}
}

// dashboardCacheControl adds default Cache-Control header for all
// dashboard UI resources (ignoring the root index.html path)
func dashboardCacheControl() func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		if e.Request.PathValue(StaticWildcardParam) != "" {
			e.Response.Header().Set("Cache-Control", "max-age=1209600, stale-while-revalidate=86400")
		}

		return e.Next()
	}
}
