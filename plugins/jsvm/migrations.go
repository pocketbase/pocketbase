package jsvm

import (
	"fmt"
	"path/filepath"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/process"
	"github.com/dop251/goja_nodejs/require"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// MigrationsConfig defines the config options of the JS migrations loader plugin.
type MigrationsConfig struct {
	// Dir specifies the directory with the JS migrations.
	//
	// If not set it fallbacks to a relative "pb_data/../pb_migrations" directory.
	Dir string
}

// MustRegisterMigrations registers the JS migrations loader plugin to
// the provided app instance and panics if it fails.
//
// Example usage:
//
//	jsvm.MustRegisterMigrations(app, jsvm.MigrationsConfig{})
func MustRegisterMigrations(app core.App, config MigrationsConfig) {
	if err := RegisterMigrations(app, config); err != nil {
		panic(err)
	}
}

// RegisterMigrations registers the JS migrations loader hooks plugin
// to the provided app instance.
func RegisterMigrations(app core.App, config MigrationsConfig) error {
	l := &migrations{app: app, config: config}

	if l.config.Dir == "" {
		l.config.Dir = filepath.Join(app.DataDir(), "../pb_migrations")
	}

	files, err := filesContent(l.config.Dir, `^.*\.js$`)
	if err != nil {
		return err
	}

	registry := new(require.Registry) // this can be shared by multiple runtimes

	for file, content := range files {
		vm := goja.New()
		registry.Enable(vm)
		console.Enable(vm)
		process.Enable(vm)
		dbxBinds(vm)
		tokensBinds(vm)
		securityBinds(vm)

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

type migrations struct {
	app    core.App
	config MigrationsConfig
}
