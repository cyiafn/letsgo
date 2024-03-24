// Package letsgo provides generic useful functions for Go.
package letsgo

// IsNil returns true if the given pointer is nil, otherwise it returns false.
func IsNil[T any](v *T) bool {
	return v == nil
}

// Ptr returns a pointer to the given value.
// This allows more concise way to return a pointer to a primitive.
func Ptr[T any](v T) *T {
	return &v
}

// IfNilDefault returns the given pointer if it is not nil, otherwise it returns the default value.
func IfNilDefault[T any](v *T, def T) *T {
	if IsNil(v) {
		return &def
	}
	return v
}

// IfNilDefaultValue returns the value of the given pointer if it is not nil, otherwise it returns the default value.
func IfNilDefaultValue[T any](v *T, def T) T {
	if v == nil {
		return def
	}
	return *v
}