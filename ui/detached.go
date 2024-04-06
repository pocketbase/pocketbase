//go:build detachableadmin

package ui

import (
	"os"

	"github.com/labstack/echo/v5"
)

func init() {
	DistDirFS = echo.MustSubFS(os.DirFS(os.Getenv("PB_ADMIN_UI_DIR")), "dist")
}
