// Package slices provides useful collection data structures for Go.
package slices

// Has returns true if the given slice contains the given element, otherwise it returns false.
func Has[T comparable](s []T, a T) bool {
	for _, v := range s {
		if v == a {
			return true
		}
	}
	return false
}

// AllThatSatisfies returns a new slice that contains all elements that satisfy the given function.
func AllThatSatisfies[T any](s []T, f func(T) bool) []T {
	if s == nil || len(s) == 0 {
		return make([]T, 0)
	}
	result := make([]T, 0, len(s))
	for _, v := range s {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}
