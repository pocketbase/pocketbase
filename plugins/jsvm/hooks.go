package jsvm

import (
	"path/filepath"
	"runtime"
	"time"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/eventloop"
	"github.com/dop251/goja_nodejs/process"
	"github.com/dop251/goja_nodejs/require"
	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
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
	files, err := filesContent(p.config.Dir, `^.*\.pb\.js$`)
	if err != nil {
		return err
	}

	dbx.HashExp{}.Build(app.DB(), nil)

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
					// return err
				}
			}
		}
	})

	loop.Start()

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

func (h *hooks) watchFiles() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	h.app.OnTerminate().Add(func(e *core.TerminateEvent) error {
		watcher.Close()

		return nil
	})

	var debounceTimer *time.Timer

	// start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if debounceTimer != nil {
					debounceTimer.Stop()
				}
				debounceTimer = time.AfterFunc(100*time.Millisecond, func() {
					// app restart is currently not supported on Windows
					if runtime.GOOS == "windows" {
						color.Yellow("File %s changed, please restart the app", event.Name)
					} else {
						color.Yellow("File %s changed, restarting...", event.Name)
						if err := h.app.Restart(); err != nil {
							color.Red("Failed to restart the app:", err)
						}
					}
				})
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				color.Red("Watch error:", err)
			}
		}
	}()

	// add the directory to watch
	err = watcher.Add(h.config.Dir)
	if err != nil {
		watcher.Close()
		return err
	}

	return nil
}
