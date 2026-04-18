//go:build no_ui

package ui

import "io/fs"

// DistDirFS is deliberately not set to prevent bundling the UI with the binary.
var DistDirFS fs.FS
