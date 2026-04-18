// Package letsgo provides generic useful functions for Go.
package letsgo

// IsNil reports whether the given typed pointer is nil.
//
// Semantics:
//   - Pure: purely a comparison, no allocation, no side effects.
//   - Only meaningful for typed pointers (*T). It does NOT unwrap
//     interface-typed values that wrap a nil pointer — an interface
//     containing a typed nil pointer is itself non-nil and should be
//     checked with reflect, not this helper.
//
// Example:
//
//	var p *int
//	IsNil(p)            // -> true
//	x := 5
//	IsNil(&x)           // -> false
func IsNil[T any](v *T) bool {
	return v == nil
}

// Ptr returns a pointer to a copy of v. Useful for taking the address of
// a literal or a function-return value inline (which Go does not allow
// syntactically).
//
// Semantics:
//   - Allocates: v is copied onto the heap (in practice — escape analysis
//     may optimize small cases) and its address is returned. Mutating the
//     returned pointer does NOT affect the caller's variable.
//   - Non-nil: the returned pointer is never nil, even when T's zero value
//     was passed in.
//
// Example:
//
//	p := Ptr(42)         // *int pointing to 42
//	q := Ptr("hi")       // *string pointing to "hi"
//	type Req struct{ Name *string }
//	r := Req{Name: Ptr("alice")} // inline struct field
func Ptr[T any](v T) *T {
	return &v
}

// IfNilDefault returns v if v is non-nil, otherwise returns a pointer to
// a copy of def.
//
// Semantics:
//   - Allocates only on the nil branch (to take the address of def).
//   - The returned pointer is never nil.
//   - When v is non-nil, the SAME pointer is returned (not a copy), so
//     mutations through it are visible to the caller.
//   - When v is nil, def is copied; mutations via the returned pointer
//     do NOT affect the caller's def variable.
//
// Example:
//
//	var missing *string
//	IfNilDefault(missing, "fallback") // -> *string pointing to "fallback"
//	s := "given"
//	IfNilDefault(&s, "fallback")      // -> &s (same pointer)
func IfNilDefault[T any](v *T, def T) *T {
	if IsNil(v) {
		return &def
	}
	return v
}

// IfNilDefaultValue returns *v if v is non-nil, otherwise returns def.
// This is the value-returning counterpart of IfNilDefault for callers
// that want a T rather than a *T.
//
// Semantics:
//   - Pure: no allocation, no mutation.
//   - Dereferences v when non-nil — make sure v points to a valid T
//     (the normal non-nil invariant; this is NOT a safety wrapper that
//     survives pointing into freed memory).
//
// Example:
//
//	var missing *int
//	IfNilDefaultValue(missing, 10)  // -> 10
//	x := 7
//	IfNilDefaultValue(&x, 10)       // -> 7
func IfNilDefaultValue[T any](v *T, def T) T {
	if v == nil {
		return def
	}
	return *v
}
