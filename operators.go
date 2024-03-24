// Package letsgo provides generic useful functions for Go.
package letsgo

// If returns a if condition is true, otherwise it returns b.
func If[T any](condition bool, a, b T) T {
	if condition {
		return a
	}
	return b
}

// IfFunc executes function a if condition is true, otherwise it executes b.
func IfFunc[T any](condition bool, a, b func() T) T {
	if condition {
		return a()
	}
	return b()
}