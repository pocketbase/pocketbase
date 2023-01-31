package main

import (
	"github.com/pocketbase/pocketbase/plugins/cloudcode"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/jsvm"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
)

func main() {
	app := pocketbase.New()

	// ---------------------------------------------------------------
	// Optional plugin flags:
	// ---------------------------------------------------------------

	var migrationsDir string
	app.RootCmd.PersistentFlags().StringVar(
		&migrationsDir,
		"migrationsDir",
		"",
		"the directory with the user defined migrations",
	)

	var automigrate bool
	app.RootCmd.PersistentFlags().BoolVar(
		&automigrate,
		"automigrate",
		true,
		"enable/disable auto migrations",
	)

	var cloudCode string
	app.RootCmd.PersistentFlags().StringVar(
		&cloudCode,
		"cloudCode",
		defaultCloudCode(),
		"the path to user-defined cloud code",
	)

	var cloudCodeInit string
	app.RootCmd.PersistentFlags().StringVar(
		&cloudCode,
		"cloudCodeInit",
		"",
		"the path to a lua file run before the main cloud code. useful for sandboxing or pre-loading libraries.",
	)

	var publicDir string
	app.RootCmd.PersistentFlags().StringVar(
		&publicDir,
		"publicDir",
		defaultPublicDir(),
		"the directory to serve static files",
	)

	var indexFallback bool
	app.RootCmd.PersistentFlags().BoolVar(
		&indexFallback,
		"indexFallback",
		true,
		"fallback the request to index.html on missing static path (eg. when pretty urls are used with SPA)",
	)

	app.RootCmd.ParseFlags(os.Args[1:])

	// ---------------------------------------------------------------
	// Plugins and hooks:
	// ---------------------------------------------------------------

	// load js pb_migrations
	jsvm.MustRegisterMigrations(app, &jsvm.MigrationsOptions{
		Dir: migrationsDir,
	})

	// migrate command (with js templates)
	migratecmd.MustRegister(app, app.RootCmd, &migratecmd.Options{
		TemplateLang: migratecmd.TemplateLangJS,
		Automigrate:  automigrate,
		Dir:          migrationsDir,
	})

	// register cloud code.
	cloudcode.MustRegister(app, cloudCode, cloudCodeInit)

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// serves static files from the provided public dir (if exists)
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS(publicDir), indexFallback))
		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

// the default pb_public dir location is relative to the executable
func defaultPublicDir() string {
	if strings.HasPrefix(os.Args[0], os.TempDir()) {
		// most likely ran with go run
		return "./pb_public"
	}
	return filepath.Join(os.Args[0], "../pb_public")
}

func defaultCloudCode() string {
	if strings.HasPrefix(os.Args[0], os.TempDir()) {
		// most likely ran with go run
		return "./pb_cloudcode/main.lua"
	}
	return filepath.Join(os.Args[0], "../pb_cloudcode/main.lua")
}
