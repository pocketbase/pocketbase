package cmd

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/labstack/echo/v5/middleware"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
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
		Short: "Starts the web server (default to localhost:8090)",
		Run: func(command *cobra.Command, args []string) {
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

			// ensure that the latest migrations are applied before starting the server
			if err := runMigrations(app); err != nil {
				panic(err)
			}

			// reload app settings in case a new default value was set with a migration
			// (or if this is the first time the init migration was executed)
			if err := app.RefreshSettings(); err != nil {
				color.Yellow("=====================================")
				color.Yellow("WARNING - Settings load error! \n%v", err)
				color.Yellow("Fallback to the application defaults.")
				color.Yellow("=====================================")
			}

			// if no admins are found, create the first one
			totalAdmins, err := app.Dao().TotalAdmins()
			if err != nil {
				log.Fatalln(err)
				return
			}
			if totalAdmins == 0 {
				if err := promptCreateAdmin(app); err != nil {
					log.Fatalln(err)
					return
				}
			}

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
				ReadTimeout: 60 * time.Second,
				// WriteTimeout: 60 * time.Second, // breaks sse!
				Handler: router,
				Addr:    mainAddr,
			}

			if showStartBanner {
				schema := "http"
				if httpsAddr != "" {
					schema = "https"
				}
				bold := color.New(color.Bold).Add(color.FgGreen)
				bold.Printf("> Server started at: %s\n", color.CyanString("%s://%s", schema, serverConfig.Addr))
				fmt.Printf("  - REST API: %s\n", color.CyanString("%s://%s/api/", schema, serverConfig.Addr))
				fmt.Printf("  - Admin UI: %s\n", color.CyanString("%s://%s/_/", schema, serverConfig.Addr))
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
		"localhost:8090",
		"api HTTP server address",
	)

	command.PersistentFlags().StringVar(
		&httpsAddr,
		"https",
		"",
		"api HTTPS server address (auto TLS via Let's Encrypt)\nthe incomming --http address traffic also will be redirected to this address",
	)

	return command
}

func runMigrations(app core.App) error {
	connections := migrationsConnectionsMap(app)

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

func promptCreateAdmin(app core.App) error {
	color.White("-------------------------------------")
	color.Cyan("Lets create your first admin account:")
	color.White("-------------------------------------")

	prompts := []*survey.Question{
		{
			Name:   "Email",
			Prompt: &survey.Input{Message: "Email:"},
			Validate: func(val any) error {
				if err := survey.Required(val); err != nil {
					return err
				}
				if err := is.Email.Validate(val); err != nil {
					return err
				}

				return nil
			},
		},
		{
			Name:   "Password",
			Prompt: &survey.Password{Message: "Pass (min 10 chars):"},
			Validate: func(val any) error {
				if str, ok := val.(string); !ok || len(str) < 10 {
					return errors.New("The password must be at least 10 characters.")
				}
				return nil
			},
		},
	}

	result := struct {
		Email    string
		Password string
	}{}
	if err := survey.Ask(prompts, &result); err != nil {
		return err
	}

	form := forms.NewAdminUpsert(app, &models.Admin{})
	form.Email = result.Email
	form.Password = result.Password
	form.PasswordConfirm = result.Password

	if err := form.Submit(); err != nil {
		return err
	}

	color.Green("Successfully created admin %s!", result.Email)
	fmt.Println("")

	return nil
}
