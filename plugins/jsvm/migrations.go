package jsvm

import (
	"os"
	"path/filepath"

	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// MigrationsLoaderOptions defines optional struct to customize the default plugin behavior.
type MigrationsLoaderOptions struct {
	// Dir is the app migrations directory from where the js files will be loaded
	// (default to pb_data/migrations)
	Dir string
}

// migrationsLoader is the plugin definition.
// Usually it is instantiated via RegisterMigrationsLoader or MustRegisterMigrationsLoader.
type migrationsLoader struct {
	app     core.App
	options *MigrationsLoaderOptions
}

//
// MustRegisterMigrationsLoader registers the plugin to the provided
// app instance and panics if it fails.
//
// It it calls RegisterMigrationsLoader(app, options)
//
// If options is nil, by default the js files from pb_data/migrations are loaded.
// Set custom options.Dir if you want to change it to some other directory.
func MustRegisterMigrationsLoader(app core.App, options *MigrationsLoaderOptions) {
	if err := RegisterMigrationsLoader(app, options); err != nil {
		panic(err)
	}
}

// RegisterMigrationsLoader registers the plugin to the provided app instance.
//
// If options is nil, by default the js files from pb_data/migrations are loaded.
// Set custom options.Dir if you want to change it to some other directory.
func RegisterMigrationsLoader(app core.App, options *MigrationsLoaderOptions) error {
	l := &migrationsLoader{app: app}

	if options != nil {
		l.options = options
	} else {
		l.options = &MigrationsLoaderOptions{}
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
		vm := NewBaseVM(l.app)
		registry.Enable(vm)
		console.Enable(vm)

		vm.Set("migrate", func(up, down func(db dbx.Builder) error) {
			m.AppMigrations.Register(up, down, file)
		})

		_, err := vm.RunString(string(content))
		if err != nil {
			return err
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
		if f.IsDir() {
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
