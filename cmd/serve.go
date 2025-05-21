package cmd

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/spf13/cobra"
)

// NewServeCommand creates and returns new command responsible for
// starting the default PocketBase web server.
func NewServeCommand(app core.App, showStartBanner bool) *cobra.Command {
	var allowedOrigins []string
	var httpAddr string
	var httpsAddr string
	var pathPrefix string

	command := &cobra.Command{
		Use:          "serve [domain(s)]",
		Args:         cobra.ArbitraryArgs,
		Short:        "Starts the web server (default to 127.0.0.1:8090 if no domain is specified)",
		SilenceUsage: true,
		RunE: func(command *cobra.Command, args []string) error {
			// set default listener addresses if at least one domain is specified
			if len(args) > 0 {
				if httpAddr == "" {
					httpAddr = "0.0.0.0:80"
				}
				if httpsAddr == "" {
					httpsAddr = "0.0.0.0:443"
				}
			} else {
				if httpAddr == "" {
					httpAddr = "127.0.0.1:8090"
				}
			}

			err := apis.Serve(app, apis.ServeConfig{
				HttpAddr:           httpAddr,
				HttpsAddr:          httpsAddr,
				PathPrefix:         pathPrefix,
				ShowStartBanner:    showStartBanner,
				AllowedOrigins:     allowedOrigins,
				CertificateDomains: args,
			})

			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}

			return err
		},
	}

	var origins []string
	if os.Getenv("PB_ALLOWED_ORIGINS") != "" {
		origins = strings.Split(os.Getenv("PB_ALLOWED_ORIGINS"), ",")
	} else {
		origins = []string{"*"}
	}

	command.PersistentFlags().StringSliceVar(
		&allowedOrigins,
		"origins",
		origins,
		"CORS allowed domain origins list",
	)

	command.PersistentFlags().StringVar(
		&httpAddr,
		"http",
		os.Getenv("PB_HTTP_ADDR"),
		"TCP address to listen for the HTTP server\n(if domain args are specified - default to 0.0.0.0:80, otherwise - default to 127.0.0.1:8090)",
	)

	command.PersistentFlags().StringVar(
		&httpsAddr,
		"https",
		os.Getenv("PB_HTTPS_ADDR"),
		"TCP address to listen for the HTTPS server\n(if domain args are specified - default to 0.0.0.0:443, otherwise - default to empty string, aka. no TLS)\nThe incoming HTTP traffic also will be auto redirected to the HTTPS version",
	)

	command.PersistentFlags().StringVar(
		&pathPrefix,
		"pathPrefix",
		os.Getenv("PB_PATH_PREFIX"),
		"URL path prefix for the server routes. Must start with a slash (eg. /pb). Useful when reuse same domain for diffrent sites behind nginx",
	)

	return command
}
