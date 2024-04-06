// Package ui handles the PocketBase Admin frontend, which can either be embedded
// or detached from the main binary. This logic is handled by build tags, the
// embedded version is enabled by default. To detach the Admin UI, build the
// binary with the "detachableadmin" build tag. e.g. `go build -tags
// detachableadmin`
package ui

import (
	"io/fs"
)

// DistDirFS is designed to point to the dist directory of the Admin UI. This can
// point to an embedded or detached directory depending on the build tags. The
// embedded version is enabled by default, if providing a detached directory,
// which is discoverable at runtime, use the PB_ADMIN_DIR environment variable to
// specify the directory path. Note that the directory must contain a "dist"
// directory containing the Admin UI files.
var DistDirFS fs.FS
