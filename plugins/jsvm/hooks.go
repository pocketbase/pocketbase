package jsvm

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/eventloop"
	"github.com/dop251/goja_nodejs/process"
	"github.com/dop251/goja_nodejs/require"
	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/jsvm/internal/docs/generated"
)

const (
	hooksExtension = ".pb.js"

	typesFileName = ".types.d.ts"

	typesReferenceDirective = `/// <reference path="./` + typesFileName + `" />`
)

// HooksConfig defines the config options of the JS app hooks plugin.
type HooksConfig struct {
	// Dir specifies the directory with the JS app hooks.
	//
	// If not set it fallbacks to a relative "pb_data/../pb_hooks" directory.
	Dir string

	// Watch enables auto app restarts when a JS app hook file changes.
	//
	// Note that currently the application cannot be automatically restarted on Windows
	// because the restart process relies on execve.
	Watch bool
}

// MustRegisterHooks registers the JS hooks plugin to
// the provided app instance and panics if it fails.
//
// Example usage:
//
//	jsvm.MustRegisterHooks(app, jsvm.HooksConfig{})
func MustRegisterHooks(app core.App, config HooksConfig) {
	if err := RegisterHooks(app, config); err != nil {
		panic(err)
	}
}

// RegisterHooks registers the JS hooks plugin to the provided app instance.
func RegisterHooks(app core.App, config HooksConfig) error {
	p := &hooks{app: app, config: config}

	if p.config.Dir == "" {
		p.config.Dir = filepath.Join(app.DataDir(), "../pb_hooks")
	}

	// fetch all js hooks sorted by their filename
	files, err := filesContent(p.config.Dir, `^.*`+regexp.QuoteMeta(hooksExtension)+`$`)
	if err != nil {
		return err
	}

	// prepend the types reference directive to empty files
	for name, content := range files {
		if len(content) != 0 {
			continue
		}
		path := filepath.Join(p.config.Dir, name)
		if err := prependToEmptyFile(path, typesReferenceDirective+"\n\n"); err != nil {
			color.Yellow("Unable to prepend the types reference: %v", err)
		}
	}

	registry := new(require.Registry) // this can be shared by multiple runtimes

	loop := eventloop.NewEventLoop()

	loop.Run(func(vm *goja.Runtime) {
		registry.Enable(vm)
		console.Enable(vm)
		process.Enable(vm)
		baseBinds(vm)
		dbxBinds(vm)
		filesystemBinds(vm)
		tokensBinds(vm)
		securityBinds(vm)
		formsBinds(vm)
		apisBinds(vm)

		vm.Set("$app", app)

		for file, content := range files {
			_, err := vm.RunString(string(content))
			if err != nil {
				if p.config.Watch {
					color.Red("Failed to execute %s: %v", file, err)
				} else {
					panic(err)
				}
			}
		}
	})

	loop.Start()

	app.OnAfterBootstrap().Add(func(e *core.BootstrapEvent) error {
		// always update the app types on start to ensure that
		// the user has the latest generated declarations
		if len(files) > 0 {
			if err := p.saveTypesFile(); err != nil {
				color.Yellow("Unable to save app types file: %v", err)
			}
		}

		return nil
	})

	app.OnTerminate().Add(func(e *core.TerminateEvent) error {
		loop.StopNoWait()

		return nil
	})

	if p.config.Watch {
		return p.watchFiles()
	}

	return nil
}

type hooks struct {
	app    core.App
	config HooksConfig
}

func (p *hooks) watchFiles() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	var debounceTimer *time.Timer

	stopDebounceTimer := func() {
		if debounceTimer != nil {
			debounceTimer.Stop()
		}
	}

	p.app.OnTerminate().Add(func(e *core.TerminateEvent) error {
		watcher.Close()

		stopDebounceTimer()

		return nil
	})

	// start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					stopDebounceTimer()
					return
				}

				// skip TS declaration files change
				if strings.HasSuffix(event.Name, ".d.ts") {
					continue
				}

				stopDebounceTimer()
				debounceTimer = time.AfterFunc(50*time.Millisecond, func() {
					// app restart is currently not supported on Windows
					if runtime.GOOS == "windows" {
						color.Yellow("File %s changed, please restart the app", event.Name)
					} else {
						color.Yellow("File %s changed, restarting...", event.Name)
						if err := p.app.Restart(); err != nil {
							color.Red("Failed to restart the app:", err)
						}
					}
				})
			case err, ok := <-watcher.Errors:
				if !ok {
					stopDebounceTimer()
					return
				}
				color.Red("Watch error:", err)
			}
		}
	}()

	// add the directory to watch
	err = watcher.Add(p.config.Dir)
	if err != nil {
		watcher.Close()
		return err
	}

	return nil
}

func (p *hooks) saveTypesFile() error {
	data, _ := generated.Types.ReadFile("types.d.ts")

	if err := os.WriteFile(filepath.Join(p.config.Dir, typesFileName), data, 0644); err != nil {
		return err
	}

	return nil
}

// prependToEmptyFile prepends the specified text to an empty file.
//
// If the file is not empty this method does nothing.
func prependToEmptyFile(path, text string) error {
	info, err := os.Stat(path)

	if err == nil && info.Size() == 0 {
		return os.WriteFile(path, []byte(text), 0644)
	}

	return err
}
