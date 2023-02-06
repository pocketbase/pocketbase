package jsvm

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/process"
	"github.com/dop251/goja_nodejs/require"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// MigrationsOptions defines optional struct to customize the default migrations loader behavior.
type MigrationsOptions struct {
	// Dir specifies the directory with the JS migrations.
	//
	// If not set it fallbacks to a relative "pb_data/../pb_migrations" directory.
	Dir string
}

// migrations is the migrations loader plugin definition.
// Usually it is instantiated via RegisterMigrations or MustRegisterMigrations.
type migrations struct {
	app     core.App
	options *MigrationsOptions
}

// MustRegisterMigrations registers the migrations loader plugin to
// the provided app instance and panics if it fails.
//
// Internally it calls RegisterMigrations(app, options).
//
// If options is nil, by default the js files from pb_data/migrations are loaded.
// Set custom options.Dir if you want to change it to some other directory.
func MustRegisterMigrations(app core.App, options *MigrationsOptions) {
	if err := RegisterMigrations(app, options); err != nil {
		panic(err)
	}
}

// RegisterMigrations registers the plugin to the provided app instance.
//
// If options is nil, by default the js files from pb_data/migrations are loaded.
// Set custom options.Dir if you want to change it to some other directory.
func RegisterMigrations(app core.App, options *MigrationsOptions) error {
	l := &migrations{app: app}

	if options != nil {
		l.options = options
	} else {
		l.options = &MigrationsOptions{}
	}

	if l.options.Dir == "" {
		l.options.Dir = filepath.Join(app.DataDir(), "../pb_migrations")
	}

	files, err := readDirFiles(l.options.Dir)
	if err != nil {
		return err
	}

	registry := new(require.Registry) // this can be shared by multiple runtimes

	for file, content := range files {
		vm := NewBaseVM()
		registry.Enable(vm)
		console.Enable(vm)
		process.Enable(vm)

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

// readDirFiles returns a map with all directory files and their content.
//
// If directory with dirPath is missing, it returns an empty map and no error.
func readDirFiles(dirPath string) (map[string][]byte, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string][]byte{}, nil
		}
		return nil, err
	}

	result := map[string][]byte{}

	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".js") {
			continue // not a .js file
		}
		raw, err := os.ReadFile(filepath.Join(dirPath, f.Name()))
		if err != nil {
			return nil, err
		}
		result[f.Name()] = raw
	}

	return result, nil
}
