// Package jsvm implements pluggable utilities for binding a JS goja runtime
// to the PocketBase instance (loading migrations, attaching to app hooks, etc.).
//
// Example:
//
//	jsvm.MustRegister(app, jsvm.Config{
//		WatchHooks: true,
//	})
package jsvm

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/process"
	"github.com/dop251/goja_nodejs/require"
	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/plugins/jsvm/internal/types/generated"
	"github.com/pocketbase/pocketbase/tools/template"
)

const (
	typesFileName = "types.d.ts"
)

// Config defines the config options of the jsvm plugin.
type Config struct {
	// HooksWatch enables auto app restarts when a JS app hook file changes.
	//
	// Note that currently the application cannot be automatically restarted on Windows
	// because the restart process relies on execve.
	HooksWatch bool

	// HooksDir specifies the JS app hooks directory.
	//
	// If not set it fallbacks to a relative "pb_data/../pb_hooks" directory.
	HooksDir string

	// HooksFilesPattern specifies a regular expression pattern that
	// identify which file to load by the hook vm(s).
	//
	// If not set it fallbacks to `^.*(\.pb\.js|\.pb\.ts)$`, aka. any
	// HookdsDir file ending in ".pb.js" or ".pb.ts" (the last one is to enforce IDE linters).
	HooksFilesPattern string

	// HooksPoolSize specifies how many goja.Runtime instances to prewarm
	// and keep for the JS app hooks gorotines execution.
	//
	// Zero or negative value means that it will create a new goja.Runtime
	// on every fired goroutine.
	HooksPoolSize int

	// MigrationsDir specifies the JS migrations directory.
	//
	// If not set it fallbacks to a relative "pb_data/../pb_migrations" directory.
	MigrationsDir string

	// If not set it fallbacks to `^.*(\.js|\.ts)$`, aka. any MigrationDir file
	// ending in ".js" or ".ts" (the last one is to enforce IDE linters).
	MigrationsFilesPattern string

	// TypesDir specifies the directory where to store the embedded
	// TypeScript declarations file.
	//
	// If not set it fallbacks to "pb_data".
	TypesDir string
}

// MustRegister registers the jsvm plugin in the provided app instance
// and panics if it fails.
//
// Example usage:
//
//	jsvm.MustRegister(app, jsvm.Config{})
func MustRegister(app core.App, config Config) {
	if err := Register(app, config); err != nil {
		panic(err)
	}
}

// Register registers the jsvm plugin in the provided app instance.
func Register(app core.App, config Config) error {
	p := &plugin{app: app, config: config}

	if p.config.HooksDir == "" {
		p.config.HooksDir = filepath.Join(app.DataDir(), "../pb_hooks")
	}

	if p.config.MigrationsDir == "" {
		p.config.MigrationsDir = filepath.Join(app.DataDir(), "../pb_migrations")
	}

	if p.config.HooksFilesPattern == "" {
		p.config.HooksFilesPattern = `^.*(\.pb\.js|\.pb\.ts)$`
	}

	if p.config.MigrationsFilesPattern == "" {
		p.config.MigrationsFilesPattern = `^.*(\.js|\.ts)$`
	}

	if p.config.TypesDir == "" {
		p.config.TypesDir = app.DataDir()
	}

	p.app.OnAfterBootstrap().Add(func(e *core.BootstrapEvent) error {
		// always update the app types on start to ensure that
		// the user has the latest generated declarations
		if err := p.saveTypesFile(); err != nil {
			color.Yellow("Unable to save app types file: %v", err)
		}

		return nil
	})

	if err := p.registerMigrations(); err != nil {
		return fmt.Errorf("registerMigrations: %w", err)
	}

	if err := p.registerHooks(); err != nil {
		return fmt.Errorf("registerHooks: %w", err)
	}

	return nil
}

type plugin struct {
	app    core.App
	config Config
}

// registerMigrations registers the JS migrations loader.
func (p *plugin) registerMigrations() error {
	// fetch all js migrations sorted by their filename
	files, err := filesContent(p.config.MigrationsDir, p.config.MigrationsFilesPattern)
	if err != nil {
		return err
	}

	registry := new(require.Registry) // this can be shared by multiple runtimes

	for file, content := range files {
		vm := goja.New()
		registry.Enable(vm)
		console.Enable(vm)
		process.Enable(vm)
		baseBinds(vm)
		dbxBinds(vm)
		tokensBinds(vm)
		securityBinds(vm)
		// note: disallow for now and give the authors of custom SaaS offerings
		// 		 some time to adjust their code to avoid eventual security issues
		//
		// osBinds(vm)
		// filepathBinds(vm)
		// httpClientBinds(vm)

		vm.Set("migrate", func(up, down func(db dbx.Builder) error) {
			m.AppMigrations.Register(up, down, file)
		})

		_, err := vm.RunString(string(content))
		if err != nil {
			return fmt.Errorf("failed to run migration %s: %w", file, err)
		}
	}

	return nil
}

// registerHooks registers the JS app hooks loader.
func (p *plugin) registerHooks() error {
	// fetch all js hooks sorted by their filename
	files, err := filesContent(p.config.HooksDir, p.config.HooksFilesPattern)
	if err != nil {
		return err
	}

	// prepend the types reference directive
	//
	// note: it is loaded during startup to handle conveniently also
	// the case when the HooksWatch option is enabled and the application
	// restart on newly created file
	for name, content := range files {
		if len(content) != 0 {
			// skip non-empty files for now to prevent accidental overwrite
			continue
		}
		path := filepath.Join(p.config.HooksDir, name)
		directive := `/// <reference path="` + p.relativeTypesPath(p.config.HooksDir) + `" />`
		if err := prependToEmptyFile(path, directive+"\n\n"); err != nil {
			color.Yellow("Unable to prepend the types reference: %v", err)
		}
	}

	// initialize the hooks dir watcher
	if p.config.HooksWatch {
		if err := p.watchHooks(); err != nil {
			return err
		}
	}

	if len(files) == 0 {
		// no need to register the vms since there are no entrypoint files anyway
		return nil
	}

	absHooksDir, err := filepath.Abs(p.config.HooksDir)
	if err != nil {
		return err
	}

	p.app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.HTTPErrorHandler = p.normalizeServeExceptions(e.Router.HTTPErrorHandler)
		return nil
	})

	// safe to be shared across multiple vms
	requireRegistry := new(require.Registry)
	templateRegistry := template.NewRegistry()

	sharedBinds := func(vm *goja.Runtime) {
		requireRegistry.Enable(vm)
		console.Enable(vm)
		process.Enable(vm)
		baseBinds(vm)
		dbxBinds(vm)
		filesystemBinds(vm)
		tokensBinds(vm)
		securityBinds(vm)
		osBinds(vm)
		filepathBinds(vm)
		httpClientBinds(vm)
		formsBinds(vm)
		apisBinds(vm)
		vm.Set("$app", p.app)
		vm.Set("$template", templateRegistry)
		vm.Set("__hooks", absHooksDir)
	}

	// initiliaze the executor vms
	executors := newPool(p.config.HooksPoolSize, func() *goja.Runtime {
		executor := goja.New()
		sharedBinds(executor)
		return executor
	})

	// initialize the loader vm
	loader := goja.New()
	sharedBinds(loader)
	hooksBinds(p.app, loader, executors)
	cronBinds(p.app, loader, executors)
	routerBinds(p.app, loader, executors)

	for file, content := range files {
		func() {
			defer func() {
				if err := recover(); err != nil {
					fmtErr := fmt.Errorf("Failed to execute %s:\n - %v", file, err)

					if p.config.HooksWatch {
						color.Red("%v", fmtErr)
					} else {
						panic(fmtErr)
					}
				}
			}()

			_, err := loader.RunString(string(content))
			if err != nil {
				panic(err)
			}
		}()
	}

	return nil
}

// normalizeExceptions wraps the provided error handler and returns a new one
// with extracted goja exception error value for consistency when throwing or returning errors.
func (p *plugin) normalizeServeExceptions(oldErrorHandler echo.HTTPErrorHandler) echo.HTTPErrorHandler {
	return func(c echo.Context, err error) {
		defer func() {
			oldErrorHandler(c, err)
		}()

		if err == nil || c.Response().Committed {
			return // no error or already committed
		}

		jsException, ok := err.(*goja.Exception)
		if !ok {
			return // no exception
		}

		switch v := jsException.Value().Export().(type) {
		case error:
			err = v
		case map[string]any: // goja.GoError
			if vErr, ok := v["value"].(error); ok {
				err = vErr
			}
		}
	}
}

// watchHooks initializes a hooks file watcher that will restart the
// application (*if possible) in case of a change in the hooks directory.
//
// This method does nothing if the hooks directory is missing.
func (p *plugin) watchHooks() error {
	if _, err := os.Stat(p.config.HooksDir); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil // no hooks dir to watch
		}
		return err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	var debounceTimer *time.Timer

	stopDebounceTimer := func() {
		if debounceTimer != nil {
			debounceTimer.Stop()
			debounceTimer = nil
		}
	}

	p.app.OnTerminate().Add(func(e *core.TerminateEvent) error {
		watcher.Close()

		stopDebounceTimer()

		return nil
	})

	// start listening for events.
	go func() {
		defer stopDebounceTimer()

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
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
					return
				}
				color.Red("Watch error:", err)
			}
		}
	}()

	// add directories to watch
	//
	// @todo replace once recursive watcher is added (https://github.com/fsnotify/fsnotify/issues/18)
	dirsErr := filepath.Walk(p.config.HooksDir, func(path string, info fs.FileInfo, err error) error {
		// ignore hidden directories and node_modules
		if !info.IsDir() || info.Name() == "node_modules" || strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		return watcher.Add(path)
	})
	if dirsErr != nil {
		watcher.Close()
	}

	return dirsErr
}

// fullTypesPathReturns returns the full path to the generated TS file.
func (p *plugin) fullTypesPath() string {
	return filepath.Join(p.config.TypesDir, typesFileName)
}

// relativeTypesPath returns a path to the generated TS file relative
// to the specified basepath.
//
// It fallbacks to the full path if generating the relative path fails.
func (p *plugin) relativeTypesPath(basepath string) string {
	fullPath := p.fullTypesPath()

	rel, err := filepath.Rel(basepath, fullPath)
	if err != nil {
		// fallback to the full path
		rel = fullPath
	}

	return rel
}

// saveTypesFile saves the embedded TS declarations as a file on the disk.
func (p *plugin) saveTypesFile() error {
	fullPath := p.fullTypesPath()

	// ensure that the types directory exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	// retrieve the types data to write
	data, err := generated.Types.ReadFile(typesFileName)
	if err != nil {
		return err
	}

	if err := os.WriteFile(fullPath, data, 0644); err != nil {
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

// filesContent returns a map with all direct files within the specified dir and their content.
//
// If directory with dirPath is missing or no files matching the pattern were found,
// it returns an empty map and no error.
//
// If pattern is empty string it matches all root files.
func filesContent(dirPath string, pattern string) (map[string][]byte, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return map[string][]byte{}, nil
		}
		return nil, err
	}

	var exp *regexp.Regexp
	if pattern != "" {
		var err error
		if exp, err = regexp.Compile(pattern); err != nil {
			return nil, err
		}
	}

	result := map[string][]byte{}

	for _, f := range files {
		if f.IsDir() || (exp != nil && !exp.MatchString(f.Name())) {
			continue
		}

		raw, err := os.ReadFile(filepath.Join(dirPath, f.Name()))
		if err != nil {
			return nil, err
		}

		result[f.Name()] = raw
	}

	return result, nil
}
