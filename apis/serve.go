package apis

import (
	"context"
	"crypto/tls"
	"errors"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/hook"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/routine"
	"github.com/pocketbase/pocketbase/ui"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

// ServeConfig defines a configuration struct for apis.Serve().
type ServeConfig struct {
	// ShowStartBanner indicates whether to show or hide the server start console message.
	ShowStartBanner bool

	// HttpAddr is the TCP address to listen for the HTTP server (eg. "127.0.0.1:80").
	HttpAddr string

	// HttpsAddr is the TCP address to listen for the HTTPS server (eg. "127.0.0.1:443").
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
func Serve(app core.App, config ServeConfig) error {
	if len(config.AllowedOrigins) == 0 {
		config.AllowedOrigins = []string{"*"}
	}

	// ensure that the latest migrations are applied before starting the server
	err := app.RunAllMigrations()
	if err != nil {
		return err
	}

	pbRouter, err := NewRouter(app)
	if err != nil {
		return err
	}

	pbRouter.Bind(CORS(CORSConfig{
		AllowOrigins: config.AllowedOrigins,
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	pbRouter.GET("/_/{path...}", Static(ui.DistDirFS, false)).
		BindFunc(func(e *core.RequestEvent) error {
			// ignore root path
			if e.Request.PathValue(StaticWildcardParam) != "" {
				e.Response.Header().Set("Cache-Control", "max-age=1209600, stale-while-revalidate=86400")
			}

			// add a default CSP
			if e.Response.Header().Get("Content-Security-Policy") == "" {
				e.Response.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' http://127.0.0.1:* https://tile.openstreetmap.org data: blob:; connect-src 'self' http://127.0.0.1:* https://nominatim.openstreetmap.org; script-src 'self' 'sha256-GRUzBA7PzKYug7pqxv5rJaec5bwDCw1Vo6/IXwvD3Tc='")
			}

			return e.Next()
		}).
		Bind(Gzip())

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
		pbRouter.Bind(wwwRedirect(wwwRedirects))
	}

	certManager := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache(filepath.Join(app.DataDir(), core.LocalAutocertCacheDirName)),
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
		// higher defaults to accommodate large file uploads/downloads
		WriteTimeout:      5 * time.Minute,
		ReadTimeout:       5 * time.Minute,
		ReadHeaderTimeout: 1 * time.Minute,
		Addr:              mainAddr,
		BaseContext: func(l net.Listener) context.Context {
			return baseCtx
		},
		ErrorLog: log.New(&serverErrorLogWriter{app: app}, "", 0),
	}

	serveEvent := new(core.ServeEvent)
	serveEvent.App = app
	serveEvent.Router = pbRouter
	serveEvent.Server = server
	serveEvent.CertManager = certManager
	serveEvent.InstallerFunc = DefaultInstallerFunc

	var listener net.Listener

	// graceful shutdown
	// ---------------------------------------------------------------
	// WaitGroup to block until server.ShutDown() returns because Serve and similar methods exit immediately.
	// Note that the WaitGroup would do nothing if the app.OnTerminate() hook isn't triggered.
	var wg sync.WaitGroup

	// try to gracefully shutdown the server on app termination
	app.OnTerminate().Bind(&hook.Handler[*core.TerminateEvent]{
		Id: "pbGracefulShutdown",
		Func: func(te *core.TerminateEvent) error {
			cancelBaseCtx()

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			wg.Add(1)

			_ = server.Shutdown(ctx)

			if te.IsRestart {
				// wait for execve and other handlers up to 3 seconds before exit
				time.AfterFunc(3*time.Second, func() {
					wg.Done()
				})
			} else {
				wg.Done()
			}

			return te.Next()
		},
		Priority: -9999,
	})

	// wait for the graceful shutdown to complete before exit
	defer func() {
		wg.Wait()

		if listener != nil {
			_ = listener.Close()
		}
	}()
	// ---------------------------------------------------------------

	var baseURL string

	// trigger the OnServe hook and start the tcp listener
	serveHookErr := app.OnServe().Trigger(serveEvent, func(e *core.ServeEvent) error {
		handler, err := e.Router.BuildMux()
		if err != nil {
			return err
		}

		e.Server.Handler = handler

		if config.HttpsAddr == "" {
			baseURL = "http://" + serverAddrToHost(serveEvent.Server.Addr)
		} else {
			baseURL = "https://"
			if len(config.CertificateDomains) > 0 {
				baseURL += config.CertificateDomains[0]
			} else {
				baseURL += serverAddrToHost(serveEvent.Server.Addr)
			}
		}

		addr := e.Server.Addr
		if addr == "" {
			// fallback similar to the std Server.ListenAndServe/ListenAndServeTLS
			if config.HttpsAddr != "" {
				addr = ":https"
			} else {
				addr = ":http"
			}
		}

		listener, err = net.Listen("tcp", addr)
		if err != nil {
			return err
		}

		if e.InstallerFunc != nil {
			app := e.App
			installerFunc := e.InstallerFunc
			routine.FireAndForget(func() {
				if err := loadInstaller(app, baseURL, installerFunc); err != nil {
					app.Logger().Warn("Failed to initialize installer", "error", err)
				}
			})
		}

		return nil
	})
	if serveHookErr != nil {
		return serveHookErr
	}

	if listener == nil {
		//nolint:staticcheck
		return errors.New("The OnServe finalizer wasn't invoked. Did you forget to call the ServeEvent.Next() method?")
	}

	if config.ShowStartBanner {
		date := new(strings.Builder)
		log.New(date, "", log.LstdFlags).Print()

		bold := color.New(color.Bold).Add(color.FgGreen)
		bold.Printf(
			"%s Server started at %s\n",
			strings.TrimSpace(date.String()),
			color.CyanString("%s", baseURL),
		)

		regular := color.New()
		regular.Printf("├─ REST API:  %s\n", color.CyanString("%s/api/", baseURL))
		regular.Printf("└─ Dashboard: %s\n", color.CyanString("%s/_/", baseURL))
	}

	var serveErr error
	if config.HttpsAddr != "" {
		if config.HttpAddr != "" {
			// start an additional HTTP server for redirecting the traffic to the HTTPS version
			go http.ListenAndServe(config.HttpAddr, certManager.HTTPHandler(nil))
		}

		// start HTTPS server
		serveErr = serveEvent.Server.ServeTLS(listener, "", "")
	} else {
		// OR start HTTP server
		serveErr = serveEvent.Server.Serve(listener)
	}
	if serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
		return serveErr
	}

	return nil
}

// serverAddrToHost loosely converts http.Server.Addr string into a host to print.
func serverAddrToHost(addr string) string {
	if addr == "" || strings.HasSuffix(addr, ":http") || strings.HasSuffix(addr, ":https") {
		return "127.0.0.1"
	}
	return addr
}

type serverErrorLogWriter struct {
	app core.App
}

func (s *serverErrorLogWriter) Write(p []byte) (int, error) {
	s.app.Logger().Debug(strings.TrimSpace(string(p)))

	return len(p), nil
}
