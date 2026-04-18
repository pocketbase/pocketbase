package apis

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/hook"
	"github.com/pocketbase/pocketbase/ui"
)

// bindUIExtensions binds the superuser UI extensions routes to the ServeEvent.Router.
//
// This method does nothing if the superuser UI is not bundled (aka. build with "no_ui" tag),
func bindUIExtensions(app core.App) {
	if ui.DistDirFS == nil {
		return
	}

	app.OnServe().Bind(&hook.Handler[*core.ServeEvent]{
		Priority: 9999, // execute as latest as possible
		Func: func(se *core.ServeEvent) error {
			uiGroup := se.Router.Group("/_").
				BindFunc(func(e *core.RequestEvent) error {
					if !e.App.IsDev() && e.Response.Header().Get("Cache-Control") == "" {
						e.Response.Header().Set("Cache-Control", "max-age=1209600, stale-while-revalidate=86400")
					}

					if e.Response.Header().Get("Content-Security-Policy") == "" {
						e.Response.Header().Set("Content-Security-Policy", defaultCSP)
					}

					return e.Next()
				}).
				Bind(Gzip())

			// register static extension routes
			for _, ext := range se.UIExtensions {
				if ext.Name == "" || ext.FS == nil {
					se.App.Logger().Debug("Invalid UI extension configuration", slog.Any("extension", ext))
					continue
				}

				uiGroup.GET("/extensions/"+ext.Name+"/{path...}", Static(ext.FS, false))
			}

			// combine all extensions main.js in one file
			//
			// note: don't cache in memory to allow previewing changes without restart
			uiGroup.GET("/extensions.js", func(re *core.RequestEvent) error {
				buf := new(bytes.Buffer)

				for _, ext := range se.UIExtensions {
					err := copyExtensionMainjs(buf, ext)
					if err != nil {
						return re.InternalServerError("An error occurred while generating the main.js extension file", err)
					}
				}

				return re.Stream(200, "text/javascript", buf)
			}).Bind(SkipSuccessActivityLog())

			return se.Next()
		},
	})
}

func copyExtensionMainjs(buf *bytes.Buffer, ext core.UIExtension) error {
	f, err := ext.FS.Open("main.js")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil // nothing to copy
		}

		return fmt.Errorf("[UI extension %q] main.js open error: %w", ext.Name, err)
	}
	defer f.Close()

	// wrap in a self-executing function to avoid scope and concatenation issues
	_, _ = buf.WriteString("(function(){")

	_, err = io.Copy(buf, f)
	if err != nil {
		return fmt.Errorf("[UI extension %q] main.js copy error: %w", ext.Name, err)
	}

	_, _ = buf.WriteString("})();")

	return nil
}
