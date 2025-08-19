//go:build js && wasm

package core

import "errors"

func exec(argv0 string, argv []string, envv []string) error {
	return errors.ErrUnsupported
}
