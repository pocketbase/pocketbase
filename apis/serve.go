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

// ServeConfig defines a configuration struct for apis.Serve().
type ServeConfig struct {
	// ShowStartBanner indicates whether to show or hide the server start console message.
	ShowStartBanner bool

	// HttpAddr is the HTTP server address to bind (eg. `127.0.0.1:80`).
	HttpAddr string

	// HttpsAddr is the HTTPS server address to bind (eg. `127.0.0.1:443`).
	HttpsAddr string

	// AllowedOrigins is an optional list of CORS origins (default to "*").
	AllowedOrigins []string
}

// Serve starts a new app web server.
func Serve(app core.App, config ServeConfig) (*http.Server, error) {
	if len(config.AllowedOrigins) == 0 {
		config.AllowedOrigins = []string{"*"}
	}

	// ensure that the latest migrations are applied before starting the server
	if err := runMigrations(app); err != nil {
		return nil, err
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
		return nil, err
	}

	// configure cors
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: config.AllowedOrigins,
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	// start http server
	// ---
	mainAddr := config.HttpAddr
	if config.HttpsAddr != "" {
		mainAddr = config.HttpsAddr
	}

	mainHost, _, _ := net.SplitHostPort(mainAddr)

	certManager := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache(filepath.Join(app.DataDir(), ".autocert_cache")),
		HostPolicy: autocert.HostWhitelist(mainHost, "www."+mainHost),
	}

	server := &http.Server{
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

	serveEvent := &core.ServeEvent{
		App:         app,
		Router:      router,
		Server:      server,
		CertManager: certManager,
	}
	if err := app.OnBeforeServe().Trigger(serveEvent); err != nil {
		return nil, err
	}

	if config.ShowStartBanner {
		schema := "http"
		if config.HttpsAddr != "" {
			schema = "https"
		}

		date := new(strings.Builder)
		log.New(date, "", log.LstdFlags).Print()

		bold := color.New(color.Bold).Add(color.FgGreen)
		bold.Printf(
			"%s Server started at %s\n",
			strings.TrimSpace(date.String()),
			color.CyanString("%s://%s", schema, server.Addr),
		)

		regular := color.New()
		regular.Printf("├─ REST API: %s\n", color.CyanString("%s://%s/api/", schema, server.Addr))
		regular.Printf("└─ Admin UI: %s\n", color.CyanString("%s://%s/_/", schema, server.Addr))
	}

	// try to gracefully shutdown the server on app termination
	app.OnTerminate().Add(func(e *core.TerminateEvent) error {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		server.Shutdown(ctx)
		return nil
	})

	// start HTTPS server
	if config.HttpsAddr != "" {
		// if httpAddr is set, start an HTTP server to redirect the traffic to the HTTPS version
		if config.HttpAddr != "" {
			go http.ListenAndServe(config.HttpAddr, certManager.HTTPHandler(nil))
		}

		return server, server.ListenAndServeTLS("", "")
	}

	// OR start HTTP server
	return server, server.ListenAndServe()
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
