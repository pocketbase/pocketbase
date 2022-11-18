// Package types implements some commonly used db serializable types
// like datetime, json, etc.
package types

// Pointer is a generic helper that returns val as *T.
func Pointer[T any](val T) *T {
	return &val
}
