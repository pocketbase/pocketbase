//go:build !(js && wasm)

package core

import "syscall"

func exec(argv0 string, argv []string, envv []string) error {
	return syscall.Exec(argv0, argv, envv)
}
