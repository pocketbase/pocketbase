//go:build js && wasm

package core

import "errors"

func execve(argv0 string, argv []string, envv []string) error {
	return errors.ErrUnsupported
}
