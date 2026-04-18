// Package choose provides ternary-style conditional helpers for Go.
// Go lacks a `cond ? a : b` operator; these functions offer the same
// expression-level shape with explicit control over evaluation strategy.
package choose

// If returns trueValue when condition is true, otherwise falseValue.
//
// Semantics:
//   - Eager: BOTH trueValue and falseValue are evaluated by the caller
//     before If is invoked, as is normal in Go function calls. If either
//     argument is expensive to produce or has side effects you want to
//     avoid, use IfLazyL, IfLazyR, or IfLazy instead.
//   - Pure: no state changes; simply returns one of the two values.
//   - Type-parameterized: trueValue and falseValue must share the same
//     type T (use `any` at the call site if you need to mix types).
//
// Example:
//
//	label := If(count == 1, "item", "items")
//	limit := If(isAdmin, 1000, 10)
func If[T any](condition bool, trueValue, falseValue T) T {
	if condition {
		return trueValue
	}
	return falseValue
}

// IfLazyL (Lazy-Left) returns trueValue() when condition is true, otherwise
// returns the pre-computed falseValue.
//
// Semantics:
//   - The true-branch thunk is ONLY invoked when the condition holds; it is
//     not called when the condition is false.
//   - falseValue is evaluated eagerly by the caller before IfLazyL is invoked.
//   - Use this when the true branch is expensive or has side effects but
//     the false branch is cheap/constant.
//
// Example:
//
//	// Only hit the DB when the cache misses.
//	user := IfLazyL(cached == nil, func() User { return db.Load(id) }, *cached)
func IfLazyL[T any](condition bool, trueValue func() T, falseValue T) T {
	if condition {
		return trueValue()
	}
	return falseValue
}

// IfLazyR (Lazy-Right) returns the pre-computed trueValue when condition is
// true, otherwise returns falseValue().
//
// Semantics:
//   - The false-branch thunk is ONLY invoked when the condition is false.
//   - trueValue is evaluated eagerly by the caller.
//   - Use this when the false branch is expensive or has side effects but
//     the true branch is cheap/constant.
//
// Example:
//
//	// Cheap happy path, expensive fallback.
//	cfg := IfLazyR(loaded != nil, *loaded, func() Config { return loadFromDisk() })
func IfLazyR[T any](condition bool, trueValue T, falseValue func() T) T {
	if condition {
		return trueValue
	}
	return falseValue()
}

// IfLazy returns trueValue() when condition is true, otherwise falseValue().
// Both branches are thunks, so EXACTLY ONE of them is evaluated — the
// closest equivalent in Go to a true short-circuiting ternary.
//
// Semantics:
//   - Lazy on both sides: whichever branch is not selected is never called,
//     so side effects and expensive computations in the unused branch do
//     not happen.
//   - Use when both branches are expensive, or when evaluating the wrong
//     branch would be incorrect (e.g. would panic on nil, hit the DB, etc.).
//
// Example:
//
//	value := IfLazy(useCache,
//	    func() V { return cache.Get(key) },
//	    func() V { return fetchFromOrigin(key) },
//	)
//	// Exactly one of cache.Get or fetchFromOrigin runs.
func IfLazy[T any](condition bool, trueValue, falseValue func() T) T {
	if condition {
		return trueValue()
	}
	return falseValue()
}
