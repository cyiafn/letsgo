// Package letsgo provides generic useful functions for Go.
package letsgo

// TernaryOp returns a if condition is true, otherwise it returns b.
func TernaryOp[T any](condition bool, a, b T) T {
	if condition {
		return a
	}
	return b
}

// TernaryOpFunc executes function a if condition is true, otherwise it executes b.
func TernaryOpFunc[T any](condition bool, a, b func() T) T {
	if condition {
		return a()
	}
	return b()
}
