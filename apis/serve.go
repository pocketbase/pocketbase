package apis

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/migrations/logs"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/migrate"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

// ServeConfig defines a configuration struct for apis.Serve().
type ServeConfig struct {
	// ShowStartBanner indicates whether to show or hide the server start console message.
	ShowStartBanner bool

	// HttpAddr is the TCP address to listen for the HTTP server (eg. `127.0.0.1:80`).
	HttpAddr string

	// HttpsAddr is the TCP address to listen for the HTTPS server (eg. `127.0.0.1:443`).
	HttpsAddr string

	// Optional domains list to use when issuing the TLS certificate.
	//
	// If not set, the host from the bound server address will be used.
	//
	// For convenience, for each "non-www" domain a "www" entry and
	// redirect will be automatically added.
	CertificateDomains []string

	// AllowedOrigins is an optional list of CORS origins (default to "*").
	AllowedOrigins []string
}

// Serve starts a new app web server.
//
// NB! The app should be bootstrapped before starting the web server.
//
// Example:
//
//	app.Bootstrap()
//	apis.Serve(app, apis.ServeConfig{
//		HttpAddr:        "127.0.0.1:8080",
//		ShowStartBanner: false,
//	})
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

	var wwwRedirects []string

	// extract the host names for the certificate host policy
	hostNames := config.CertificateDomains
	if len(hostNames) == 0 {
		host, _, _ := net.SplitHostPort(mainAddr)
		hostNames = append(hostNames, host)
	}
	for _, host := range hostNames {
		if strings.HasPrefix(host, "www.") {
			continue // explicitly set www host
		}

		wwwHost := "www." + host
		if !list.ExistInSlice(wwwHost, hostNames) {
			hostNames = append(hostNames, wwwHost)
			wwwRedirects = append(wwwRedirects, wwwHost)
		}
	}

	// implicit www->non-www redirect(s)
	if len(wwwRedirects) > 0 {
		router.Pre(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				host := c.Request().Host

				if strings.HasPrefix(host, "www.") && list.ExistInSlice(host, wwwRedirects) {
					return c.Redirect(
						http.StatusTemporaryRedirect,
						(c.Scheme() + "://" + host[4:] + c.Request().RequestURI),
					)
				}

				return next(c)
			}
		})
	}

	certManager := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache(filepath.Join(app.DataDir(), ".autocert_cache")),
		HostPolicy: autocert.HostWhitelist(hostNames...),
	}

	// base request context used for cancelling long running requests
	// like the SSE connections
	baseCtx, cancelBaseCtx := context.WithCancel(context.Background())
	defer cancelBaseCtx()

	server := &http.Server{
		TLSConfig: &tls.Config{
			MinVersion:     tls.VersionTLS12,
			GetCertificate: certManager.GetCertificate,
			NextProtos:     []string{acme.ALPNProto},
		},
		ReadTimeout:       10 * time.Minute,
		ReadHeaderTimeout: 30 * time.Second,
		// WriteTimeout: 60 * time.Second, // breaks sse!
		Handler: router,
		Addr:    mainAddr,
		BaseContext: func(l net.Listener) context.Context {
			return baseCtx
		},
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
		addr := server.Addr

		if config.HttpsAddr != "" {
			schema = "https"

			if len(config.CertificateDomains) > 0 {
				addr = config.CertificateDomains[0]
			}
		}

		date := new(strings.Builder)
		log.New(date, "", log.LstdFlags).Print()

		bold := color.New(color.Bold).Add(color.FgGreen)
		bold.Printf(
			"%s Server started at %s\n",
			strings.TrimSpace(date.String()),
			color.CyanString("%s://%s", schema, addr),
		)

		regular := color.New()
		regular.Printf("├─ REST API: %s\n", color.CyanString("%s://%s/api/", schema, addr))
		regular.Printf("└─ Admin UI: %s\n", color.CyanString("%s://%s/_/", schema, addr))
	}

	// WaitGroup to block until server.ShutDown() returns because Serve and similar methods exit immediately.
	// Note that the WaitGroup would not do anything if the app.OnTerminate() hook isn't triggered.
	var wg sync.WaitGroup

	// try to gracefully shutdown the server on app termination
	app.OnTerminate().Add(func(e *core.TerminateEvent) error {
		cancelBaseCtx()

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		wg.Add(1)
		server.Shutdown(ctx)
		if e.IsRestart {
			// wait for execve and other handlers up to 5 seconds before exit
			time.AfterFunc(5*time.Second, func() {
				wg.Done()
			})
		} else {
			wg.Done()
		}

		return nil
	})

	// wait for the graceful shutdown to complete before exit
	defer wg.Wait()

	// ---
	// @todo consider removing the server return value because it is
	// not really useful when combined with the blocking serve calls
	// ---

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
