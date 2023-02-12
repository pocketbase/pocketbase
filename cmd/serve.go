package cmd

import (
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
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/migrations/logs"
	"github.com/pocketbase/pocketbase/tools/migrate"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

// NewServeCommand creates and returns new command responsible for
// starting the default PocketBase web server.
func NewServeCommand(app core.App, showStartBanner bool) *cobra.Command {
	var allowedOrigins []string
	var httpAddr string
	var httpsAddr string

	command := &cobra.Command{
		Use:   "serve",
		Short: "Starts the web server (default to 127.0.0.1:8090)",
		Run: func(command *cobra.Command, args []string) {
			// ensure that the latest migrations are applied before starting the server
			if err := runMigrations(app); err != nil {
				panic(err)
			}

			// reload app settings in case a new default value was set with a migration
			// (or if this is the first time the init migration was executed)
			if err := app.RefreshSettings(); err != nil {
				color.Yellow("=====================================")
				color.Yellow("WARNING: Settings load error! \n%v", err)
				color.Yellow("Fallback to the application defaults.")
				color.Yellow("=====================================")
			}

			router, err := apis.InitApi(app)
			if err != nil {
				panic(err)
			}

			// configure cors
			router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
				Skipper:      middleware.DefaultSkipper,
				AllowOrigins: allowedOrigins,
				AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
			}))

			// start http server
			// ---
			mainAddr := httpAddr
			if httpsAddr != "" {
				mainAddr = httpsAddr
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
				ReadTimeout:       5 * time.Minute,
				ReadHeaderTimeout: 30 * time.Second,
				// WriteTimeout: 60 * time.Second, // breaks sse!
				Handler: router,
				Addr:    mainAddr,
			}

			if showStartBanner {
				schema := "http"
				if httpsAddr != "" {
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

			var serveErr error
			if httpsAddr != "" {
				// if httpAddr is set, start an HTTP server to redirect the traffic to the HTTPS version
				if httpAddr != "" {
					go http.ListenAndServe(httpAddr, certManager.HTTPHandler(nil))
				}

				// start HTTPS server
				serveErr = serverConfig.ListenAndServeTLS("", "")
			} else {
				// start HTTP server
				serveErr = serverConfig.ListenAndServe()
			}

			if serveErr != http.ErrServerClosed {
				log.Fatalln(serveErr)
			}
		},
	}

	command.PersistentFlags().StringSliceVar(
		&allowedOrigins,
		"origins",
		[]string{"*"},
		"CORS allowed domain origins list",
	)

	command.PersistentFlags().StringVar(
		&httpAddr,
		"http",
		"127.0.0.1:8090",
		"api HTTP server address",
	)

	command.PersistentFlags().StringVar(
		&httpsAddr,
		"https",
		"",
		"api HTTPS server address (auto TLS via Let's Encrypt)\nthe incoming --http address traffic also will be redirected to this address",
	)

	return command
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
