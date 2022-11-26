// Example
//
// 	publicdir.MustRegister(app, &publicdir.Options{
// 		FlagsCmd:      app.RootCmd,
// 		IndexFallback: false,
// 	})
package publicdir

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type Options struct {
	Dir           string
	IndexFallback bool
}

type plugin struct {
	app     core.App
	options *Options
}

func MustRegister(app core.App, options *Options) {
	if err := Register(app, options); err != nil {
		panic(err)
	}
}

func Register(app core.App, options *Options) error {
	p := &plugin{app: app}

	if options != nil {
		p.options = options
	} else {
		p.options = &Options{}
	}

	if options.Dir == "" {
		options.Dir = defaultPublicDir()
	}

	p.app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// serves static files from the provided public dir (if exists)
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS(options.Dir), options.IndexFallback))

		return nil
	})

	return nil
}

func defaultPublicDir() string {
	if strings.HasPrefix(os.Args[0], os.TempDir()) {
		// most likely ran with go run
		return "./pb_public"
	}
	return filepath.Join(os.Args[0], "../pb_public")
}
