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
// It implements [core.App] via embedding and all of interface methods
// could be accessed directly through the instance (eg. PocketBase.DataDir()).
type PocketBase struct {
	*appWrapper

	// RootCmd is the main cli command
	RootCmd *cobra.Command

	// console flags
	debugFlag     bool
	dataDirFlag   string
	encryptionEnv string

	// console flag fallback values
	defaultDebug         bool
	defaultDataDir       string
	defaultEncryptionEnv string

	// serve start banner
	showStartBanner bool
}

// New creates a new PocketBase instance.
//
// Note that the application will not be initialized/bootstrapped yet,
// aka. DB connections, migrations, app settings, etc. will not be accessible.
// Everything will be initialized when Start() is executed.
// If you want to initialize the application before calling Start(),
// then you'll have to manually call Bootstrap().
func New() *PocketBase {
	// try to find the base executable directory and how it was run
	var withGoRun bool
	var baseDir string
	if strings.HasPrefix(os.Args[0], os.TempDir()) {
		// probably ran with go run...
		withGoRun = true
		baseDir, _ = os.Getwd()
	} else {
		// probably ran with go build...
		withGoRun = false
		baseDir = filepath.Dir(os.Args[0])
	}

	defaultDir := filepath.Join(baseDir, "pb_data")

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
		defaultDebug:         withGoRun,
		defaultDataDir:       defaultDir,
		defaultEncryptionEnv: "",
		showStartBanner:      true,
	}

	// parse base flags
	// (errors are ignored, since the full flags parsing happens on Execute())
	pb.eagerParseFlags()

	pb.appWrapper = &appWrapper{core.NewBaseApp(
		pb.dataDirFlag,
		pb.encryptionEnv,
		pb.debugFlag,
	)}

	return pb
}

// DefaultDebug sets the default --debug flag value.
func (pb *PocketBase) DefaultDebug(val bool) *PocketBase {
	pb.defaultDebug = val
	return pb
}

// DefaultDataDir sets the default --dir flag value.
func (pb *PocketBase) DefaultDataDir(val string) *PocketBase {
	pb.defaultDataDir = val
	return pb
}

// DefaultEncryptionEnv sets the default --encryptionEnv flag value.
func (pb *PocketBase) DefaultEncryptionEnv(val string) *PocketBase {
	pb.defaultEncryptionEnv = val
	return pb
}

// ShowStartBanner shows/hides the web server start banner.
func (pb *PocketBase) ShowStartBanner(val bool) *PocketBase {
	pb.showStartBanner = val
	return pb
}

// Start starts the application, aka. registers the default system
// commands (serve, migrate, version) and executes pb.RootCmd.
func (pb *PocketBase) Start() error {
	// register system commands
	pb.RootCmd.AddCommand(cmd.NewServeCommand(pb, pb.showStartBanner))
	pb.RootCmd.AddCommand(cmd.NewMigrateCommand(pb))

	return pb.Execute()
}

// Execute initializes the application (if not already) and executes
// the pb.RootCmd with graceful shutdown support.
//
// This method differs from pb.Start() by not registering the default
// system commands!
func (pb *PocketBase) Execute() error {
	if err := pb.Bootstrap(); err != nil {
		return err
	}

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
func (pb *PocketBase) eagerParseFlags() error {
	pb.RootCmd.PersistentFlags().StringVar(
		&pb.dataDirFlag,
		"dir",
		pb.defaultDataDir,
		"the PocketBase data directory",
	)

	pb.RootCmd.PersistentFlags().StringVar(
		&pb.encryptionEnv,
		"encryptionEnv",
		pb.defaultEncryptionEnv,
		"the env variable whose value of 32 chars will be used \nas encryption key for the app settings (default none)",
	)

	pb.RootCmd.PersistentFlags().BoolVar(
		&pb.debugFlag,
		"debug",
		pb.defaultDebug,
		"enable debug mode, aka. showing more detailed logs",
	)

	return pb.RootCmd.ParseFlags(os.Args[1:])
}
