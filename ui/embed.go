// Package ui handles the PocketBase Admin frontend embedding.
package ui

import (
	"embed"

	"github.com/labstack/echo/v5"
)

//go:embed all:dist
var distDir embed.FS

//go:embed dist/index.html
var indexHTML embed.FS

// DistDirFS contains the embedded dist directory files (without the "dist" prefix)
var DistDirFS = echo.MustSubFS(distDir, "dist")

// DistIndexHTML contains the embedded dist/index.html file (without the "dist" prefix)
var DistIndexHTML = echo.MustSubFS(indexHTML, "dist")
