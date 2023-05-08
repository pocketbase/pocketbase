package apis

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/labstack/echo/v5/middleware"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/migrations/logs"
	"github.com/pocketbase/pocketbase/tools/migrate"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

// ServeOptions defines an optional struct for apis.Serve().
type ServeOptions struct {
	ShowStartBanner bool
	HttpAddr        string
	HttpsAddr       string
	AllowedOrigins  []string // optional list of CORS origins (default to "*")
	BeforeServeFunc func(server *http.Server) error
}

// Serve starts a new app web server.
func Serve(app core.App, options *ServeOptions) error {
	if options == nil {
		options = &ServeOptions{}
	}

	if len(options.AllowedOrigins) == 0 {
		options.AllowedOrigins = []string{"*"}
	}

	// ensure that the latest migrations are applied before starting the server
	if err := runMigrations(app); err != nil {
		return err
	}

	// reload app settings in case a new default value was set with a migration
	// (or if this is the first time the init migration was executed)
	if err := app.RefreshSettings(); err != nil {
		color.Yellow("=====================================")
		color.Yellow("WARNING: Settings load error! \n%v", err)
		color.Yellow("Fallback to the application defaults.")
		color.Yellow("=====================================")
	}

	router, err := InitApi(app)
	if err != nil {
		return err
	}

	// configure cors
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: options.AllowedOrigins,
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	// start http server
	// ---
	mainAddr := options.HttpAddr
	if options.HttpsAddr != "" {
		mainAddr = options.HttpsAddr
	}

	mainHost, _, _ := net.SplitHostPort(mainAddr)

	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache(filepath.Join(app.DataDir(), ".autocert_cache")),
		HostPolicy: autocert.HostWhitelist(mainHost, "www."+mainHost),
	}

	serverConfig := &http.Server{
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
			NextProtos:     []string{acme.ALPNProto},
		},
		ReadTimeout:       10 * time.Minute,
		ReadHeaderTimeout: 30 * time.Second,
		// WriteTimeout: 60 * time.Second, // breaks sse!
		Handler: router,
		Addr:    mainAddr,
	}

	if options.BeforeServeFunc != nil {
		if err := options.BeforeServeFunc(serverConfig); err != nil {
			return err
		}
	}

	if options.ShowStartBanner {
		schema := "http"
		if options.HttpsAddr != "" {
			schema = "https"
		}

		date := new(strings.Builder)
		log.New(date, "", log.LstdFlags).Print()

		bold := color.New(color.Bold).Add(color.FgGreen)
		bold.Printf(
			"%s Server started at %s\n",
			strings.TrimSpace(date.String()),
			color.CyanString("%s://%s", schema, serverConfig.Addr),
		)

		regular := color.New()
		regular.Printf(" ➜ REST API: %s\n", color.CyanString("%s://%s/api/", schema, serverConfig.Addr))
		regular.Printf(" ➜ Admin UI: %s\n", color.CyanString("%s://%s/_/", schema, serverConfig.Addr))
	}

	// try to gracefully shutdown the server on app termination
	app.OnTerminate().Add(func(e *core.TerminateEvent) error {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		serverConfig.Shutdown(ctx)
		return nil
	})

	// start HTTPS server
	if options.HttpsAddr != "" {
		// if httpAddr is set, start an HTTP server to redirect the traffic to the HTTPS version
		if options.HttpAddr != "" {
			go http.ListenAndServe(options.HttpAddr, certManager.HTTPHandler(nil))
		}

		return serverConfig.ListenAndServeTLS("", "")
	}

	// OR start HTTP server
	return serverConfig.ListenAndServe()
}

type migrationsConnection struct {
	DB             *dbx.DB
	MigrationsList migrate.MigrationsList
}

func runMigrations(app core.App) error {
	connections := []migrationsConnection{
		{
			DB:             app.DB(),
			MigrationsList: migrations.AppMigrations,
		},
		{
			DB:             app.LogsDB(),
			MigrationsList: logs.LogsMigrations,
		},
	}

	for _, c := range connections {
		runner, err := migrate.NewRunner(c.DB, c.MigrationsList)
		if err != nil {
			return err
		}

		if _, err := runner.Up(); err != nil {
			return err
		}
	}

	return nil
}
