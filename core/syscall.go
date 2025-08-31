//go:build !(js && wasm)

package core

import "syscall"

// execve invokes the execve(2) system call.
func execve(argv0 string, argv []string, envv []string) error {
	return syscall.Exec(argv0, argv, envv)
}
