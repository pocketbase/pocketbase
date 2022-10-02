package pocketbase

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/pocketbase/pocketbase/cmd"
	"github.com/pocketbase/pocketbase/core"
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

	debugFlag         bool
	dataDirFlag       string
	encryptionEnvFlag string
	hideStartBanner   bool

	// RootCmd is the main console command
	RootCmd *cobra.Command
}

// Config is the PocketBase initialization config struct.
type Config struct {
	// optional default values for the console flags
	DefaultDebug         bool
	DefaultDataDir       string // if not set, it will fallback to "./pb_data"
	DefaultEncryptionEnv string

	// hide the default console server info on app startup
	HideStartBanner bool
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
		DefaultDebug: isUsingGoRun,
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
			Use:     "pocketbase",
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
		debugFlag:         config.DefaultDebug,
		dataDirFlag:       config.DefaultDataDir,
		encryptionEnvFlag: config.DefaultEncryptionEnv,
		hideStartBanner:   config.HideStartBanner,
	}

	// parse base flags
	// (errors are ignored, since the full flags parsing happens on Execute())
	pb.eagerParseFlags(config)

	// initialize the app instance
	pb.appWrapper = &appWrapper{core.NewBaseApp(
		pb.dataDirFlag,
		pb.encryptionEnvFlag,
		pb.debugFlag,
	)}

	// hide the default help command (allow only `--help` flag)
	pb.RootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	// hook the bootstrap process
	pb.RootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		return pb.Bootstrap()
	}

	return pb
}

// Start starts the application, aka. registers the default system
// commands (serve, migrate, version) and executes pb.RootCmd.
func (pb *PocketBase) Start() error {
	// register system commands
	pb.RootCmd.AddCommand(cmd.NewServeCommand(pb, !pb.hideStartBanner))
	pb.RootCmd.AddCommand(cmd.NewMigrateCommand(pb))

	return pb.Execute()
}

// Execute initializes the application (if not already) and executes
// the pb.RootCmd with graceful shutdown support.
//
// This method differs from pb.Start() by not registering the default
// system commands!
func (pb *PocketBase) Execute() error {
	var wg sync.WaitGroup

	wg.Add(1)

	// wait for interrupt signal to gracefully shutdown the application
	go func() {
		defer wg.Done()
		quit := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit
	}()

	// execute the root command
	go func() {
		defer wg.Done()
		if err := pb.RootCmd.Execute(); err != nil {
			log.Println(err)
		}
	}()

	wg.Wait()

	// cleanup
	return pb.onTerminate()
}

// onTerminate tries to release the app resources on app termination.
func (pb *PocketBase) onTerminate() error {
	return pb.ResetBootstrapState()
}

// eagerParseFlags parses the global app flags before calling pb.RootCmd.Execute().
// so we can have all PocketBase flags ready for use on initialization.
func (pb *PocketBase) eagerParseFlags(config Config) error {
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
		&pb.debugFlag,
		"debug",
		config.DefaultDebug,
		"enable debug mode, aka. showing more detailed logs",
	)

	return pb.RootCmd.ParseFlags(os.Args[1:])
}

// tries to find the base executable directory and how it was run
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
