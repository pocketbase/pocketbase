//go:build !detachableadmin

package ui

import (
	"embed"

	"github.com/labstack/echo/v5"
)

//go:embed all:dist
var distDir embed.FS

func init() {
	// DistDirFS contains the embedded dist directory files (without the "dist" prefix)
	DistDirFS = echo.MustSubFS(distDir, "dist")
}
