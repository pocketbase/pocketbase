package pocketbase

import (
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/fatih/color"
	"github.com/thinkonmay/pocketbase/cmd"
	"github.com/thinkonmay/pocketbase/core"
	"github.com/thinkonmay/pocketbase/tools/list"
	"github.com/spf13/cobra"
)

var _ core.App = (*PocketBase)(nil)

// Version of PocketBase
var Version = "(untracked)"

// appWrapper serves as a private core.App instance wrapper.
type appWrapper struct {
	core.App
}

// PocketBase defines a PocketBase app launcher.
//
// It implements [core.App] via embedding and all of the app interface methods
// could be accessed directly through the instance (eg. PocketBase.DataDir()).
type PocketBase struct {
	*appWrapper

	devFlag           bool
	dataDirFlag       string
	encryptionEnvFlag string
	hideStartBanner   bool

	// RootCmd is the main console command
	RootCmd *cobra.Command
}

// Config is the PocketBase initialization config struct.
type Config struct {
	// optional default values for the console flags
	DefaultDev           bool
	DefaultDataDir       string // if not set, it will fallback to "./pb_data"
	DefaultEncryptionEnv string

	// hide the default console server info on app startup
	HideStartBanner bool

	// optional DB configurations
	DataMaxOpenConns int // default to core.DefaultDataMaxOpenConns
	DataMaxIdleConns int // default to core.DefaultDataMaxIdleConns
	LogsMaxOpenConns int // default to core.DefaultLogsMaxOpenConns
	LogsMaxIdleConns int // default to core.DefaultLogsMaxIdleConns
}

// New creates a new PocketBase instance with the default configuration.
// Use [NewWithConfig()] if you want to provide a custom configuration.
//
// Note that the application will not be initialized/bootstrapped yet,
// aka. DB connections, migrations, app settings, etc. will not be accessible.
// Everything will be initialized when [Start()] is executed.
// If you want to initialize the application before calling [Start()],
// then you'll have to manually call [Bootstrap()].
func New() *PocketBase {
	_, isUsingGoRun := inspectRuntime()

	return NewWithConfig(Config{
		DefaultDev: isUsingGoRun,
	})
}

// NewWithConfig creates a new PocketBase instance with the provided config.
//
// Note that the application will not be initialized/bootstrapped yet,
// aka. DB connections, migrations, app settings, etc. will not be accessible.
// Everything will be initialized when [Start()] is executed.
// If you want to initialize the application before calling [Start()],
// then you'll have to manually call [Bootstrap()].
func NewWithConfig(config Config) *PocketBase {
	// initialize a default data directory based on the executable baseDir
	if config.DefaultDataDir == "" {
		baseDir, _ := inspectRuntime()
		config.DefaultDataDir = filepath.Join(baseDir, "pb_data")
	}

	pb := &PocketBase{
		RootCmd: &cobra.Command{
			Use:     filepath.Base(os.Args[0]),
			Short:   "PocketBase CLI",
			Version: Version,
			FParseErrWhitelist: cobra.FParseErrWhitelist{
				UnknownFlags: true,
			},
			// no need to provide the default cobra completion command
			CompletionOptions: cobra.CompletionOptions{
				DisableDefaultCmd: true,
			},
		},
		devFlag:           config.DefaultDev,
		dataDirFlag:       config.DefaultDataDir,
		encryptionEnvFlag: config.DefaultEncryptionEnv,
		hideStartBanner:   config.HideStartBanner,
	}

	// replace with a colored stderr writer
	pb.RootCmd.SetErr(newErrWriter())

	// parse base flags
	// (errors are ignored, since the full flags parsing happens on Execute())
	pb.eagerParseFlags(&config)

	// initialize the app instance
	pb.appWrapper = &appWrapper{core.NewBaseApp(core.BaseAppConfig{
		IsDev:            pb.devFlag,
		DataDir:          pb.dataDirFlag,
		EncryptionEnv:    pb.encryptionEnvFlag,
		DataMaxOpenConns: config.DataMaxOpenConns,
		DataMaxIdleConns: config.DataMaxIdleConns,
		LogsMaxOpenConns: config.LogsMaxOpenConns,
		LogsMaxIdleConns: config.LogsMaxIdleConns,
	})}

	// hide the default help command (allow only `--help` flag)
	pb.RootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	return pb
}

// Start starts the application, aka. registers the default system
// commands (serve, migrate, version) and executes pb.RootCmd.
func (pb *PocketBase) Start() error {
	// register system commands
	pb.RootCmd.AddCommand(cmd.NewAdminCommand(pb))
	pb.RootCmd.AddCommand(cmd.NewServeCommand(pb, !pb.hideStartBanner))

	return pb.Execute()
}

// Execute initializes the application (if not already) and executes
// the pb.RootCmd with graceful shutdown support.
//
// This method differs from pb.Start() by not registering the default
// system commands!
func (pb *PocketBase) Execute() error {
	if !pb.skipBootstrap() {
		if err := pb.Bootstrap(); err != nil {
			return err
		}
	}

	done := make(chan bool, 1)

	// listen for interrupt signal to gracefully shutdown the application
	go func() {
		sigch := make(chan os.Signal, 1)
		signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)
		<-sigch

		done <- true
	}()

	// execute the root command
	go func() {
		// note: leave to the commands to decide whether to print their error
		pb.RootCmd.Execute()

		done <- true
	}()

	<-done

	// trigger cleanups
	return pb.OnTerminate().Trigger(&core.TerminateEvent{
		App: pb,
	}, func(e *core.TerminateEvent) error {
		return e.App.ResetBootstrapState()
	})
}

// eagerParseFlags parses the global app flags before calling pb.RootCmd.Execute().
// so we can have all PocketBase flags ready for use on initialization.
func (pb *PocketBase) eagerParseFlags(config *Config) error {
	pb.RootCmd.PersistentFlags().StringVar(
		&pb.dataDirFlag,
		"dir",
		config.DefaultDataDir,
		"the PocketBase data directory",
	)

	pb.RootCmd.PersistentFlags().StringVar(
		&pb.encryptionEnvFlag,
		"encryptionEnv",
		config.DefaultEncryptionEnv,
		"the env variable whose value of 32 characters will be used \nas encryption key for the app settings (default none)",
	)

	pb.RootCmd.PersistentFlags().BoolVar(
		&pb.devFlag,
		"dev",
		config.DefaultDev,
		"enable dev mode, aka. printing logs and sql statements to the console",
	)

	return pb.RootCmd.ParseFlags(os.Args[1:])
}

// skipBootstrap eagerly checks if the app should skip the bootstrap process:
// - already bootstrapped
// - is unknown command
// - is the default help command
// - is the default version command
//
// https://github.com/thinkonmay/pocketbase/issues/404
// https://github.com/thinkonmay/pocketbase/discussions/1267
func (pb *PocketBase) skipBootstrap() bool {
	flags := []string{
		"-h",
		"--help",
		"-v",
		"--version",
	}

	if pb.IsBootstrapped() {
		return true // already bootstrapped
	}

	cmd, _, err := pb.RootCmd.Find(os.Args[1:])
	if err != nil {
		return true // unknown command
	}

	for _, arg := range os.Args {
		if !list.ExistInSlice(arg, flags) {
			continue
		}

		// ensure that there is no user defined flag with the same name/shorthand
		trimmed := strings.TrimLeft(arg, "-")
		if len(trimmed) > 1 && cmd.Flags().Lookup(trimmed) == nil {
			return true
		}
		if len(trimmed) == 1 && cmd.Flags().ShorthandLookup(trimmed) == nil {
			return true
		}
	}

	return false
}

// inspectRuntime tries to find the base executable directory and how it was run.
func inspectRuntime() (baseDir string, withGoRun bool) {
	if strings.HasPrefix(os.Args[0], os.TempDir()) {
		// probably ran with go run
		withGoRun = true
		baseDir, _ = os.Getwd()
	} else {
		// probably ran with go build
		withGoRun = false
		baseDir = filepath.Dir(os.Args[0])
	}
	return
}

// newErrWriter returns a red colored stderr writter.
func newErrWriter() *coloredWriter {
	return &coloredWriter{
		w: os.Stderr,
		c: color.New(color.FgRed),
	}
}

// coloredWriter is a small wrapper struct to construct a [color.Color] writter.
type coloredWriter struct {
	w io.Writer
	c *color.Color
}

// Write writes the p bytes using the colored writer.
func (colored *coloredWriter) Write(p []byte) (n int, err error) {
	colored.c.SetWriter(colored.w)
	defer colored.c.UnsetWriter(colored.w)

	return colored.c.Print(string(p))
}
