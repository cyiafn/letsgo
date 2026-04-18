// Package gslice provides generic, allocation-friendly functional helpers
// for Go slices: map/filter/reduce plus common combinators.
//
// All functions take the input slice as their first argument and return a
// fresh slice (never mutating the input) unless explicitly documented
// otherwise.
package gslice

// Map applies f to every element of s and returns a new slice containing
// the results, in the same order.
//
// Semantics:
//   - Pure w.r.t. s: the input slice is never mutated (but f is free to
//     have side effects — those are the caller's concern).
//   - Always returns a non-nil slice of length len(s); for an empty input
//     the result is an empty (not nil) slice.
//   - Supports type transformation: input type T and output type U may differ.
//
// Example:
//
//	Map([]int{1, 2, 3}, func(v int) int { return v * 2 })     // -> [2 4 6]
//	Map([]int{1, 2, 3}, func(v int) string { ... })            // -> []string{...}
//	Map([]int{}, func(v int) int { return v })                 // -> [] (len 0)
func Map[T, U any](s []T, f func(T) U) []U {
	r := make([]U, len(s))
	for i, v := range s {
		r[i] = f(v)
	}
	return r
}

// Filter returns a new slice containing only the elements of s for which
// f returns true, in the original order.
//
// Semantics:
//   - Pure w.r.t. s: the input is never mutated.
//   - Always returns a non-nil slice (possibly of length 0).
//   - Preserves order and duplicates.
//
// Example:
//
//	Filter([]int{1, 2, 3, 4}, func(v int) bool { return v%2 == 0 }) // -> [2 4]
//	Filter([]int{1, 2, 3}, func(int) bool { return false })         // -> [] (len 0)
func Filter[T any](s []T, f func(T) bool) []T {
	r := make([]T, 0)
	for _, v := range s {
		v := v
		if f(v) {
			r = append(r, v)
		}
	}

	return r
}

// Reduce folds s into a single accumulator by repeatedly applying f to
// the running accumulator and each element, starting from initial.
//
// Semantics:
//   - Pure w.r.t. s: the input slice is never mutated.
//   - Left fold: elements are processed front-to-back. For a non-associative
//     f, order matters.
//   - If s is empty, initial is returned unchanged.
//
// Example:
//
//	Reduce([]int{1, 2, 3, 4}, 0, func(a, v int) int { return a + v })                // -> 10
//	Reduce([]string{"a","b","c"}, "", func(acc, v string) string { return acc + v }) // -> "abc"
//	Reduce([]int{}, 42, func(a, v int) int { return a + v })                         // -> 42
func Reduce[T, U any](s []T, initial U, f func(U, T) U) U {
	acc := initial
	for _, v := range s {
		acc = f(acc, v)
	}
	return acc
}

// ForEach invokes f on every element of s in order, purely for its side
// effects.
//
// Semantics:
//   - Does not return anything — use Map/Filter/Reduce when you want a
//     value back.
//   - Pure w.r.t. s: input is not mutated (but f may mutate captured state).
//   - Safe on empty/nil slices (simply does nothing).
//
// Example:
//
//	seen := []int{}
//	ForEach([]int{1, 2, 3}, func(v int) { seen = append(seen, v) })
//	// seen == [1 2 3]
func ForEach[T any](s []T, f func(T)) {
	for _, v := range s {
		f(v)
	}
}

// Any reports whether at least one element of s satisfies f.
//
// Semantics:
//   - Pure w.r.t. s.
//   - Short-circuits on the first match.
//   - Returns false for an empty/nil slice (existential over empty is false).
//
// Example:
//
//	Any([]int{1, 2, 3}, func(v int) bool { return v == 2 })   // -> true
//	Any([]int{1, 2, 3}, func(v int) bool { return v > 10 })   // -> false
//	Any([]int{}, func(int) bool { return true })              // -> false
func Any[T any](s []T, f func(T) bool) bool {
	for _, v := range s {
		if f(v) {
			return true
		}
	}
	return false
}

// All reports whether every element of s satisfies f.
//
// Semantics:
//   - Pure w.r.t. s.
//   - Short-circuits on the first element that fails f.
//   - Returns true for an empty/nil slice (vacuous truth — universal over
//     empty is true).
//
// Example:
//
//	All([]int{2, 4, 6}, func(v int) bool { return v%2 == 0 }) // -> true
//	All([]int{2, 3, 4}, func(v int) bool { return v%2 == 0 }) // -> false
//	All([]int{}, func(int) bool { return false })             // -> true (vacuous)
func All[T any](s []T, f func(T) bool) bool {
	for _, v := range s {
		if !f(v) {
			return false
		}
	}
	return true
}

// Contains reports whether target appears in s.
//
// Semantics:
//   - Pure w.r.t. s.
//   - Uses == for equality, which is why T must be comparable.
//   - Short-circuits on the first match; O(n) worst case.
//   - Returns false for an empty/nil slice.
//
// Example:
//
//	Contains([]int{1, 2, 3}, 2)     // -> true
//	Contains([]int{1, 2, 3}, 9)     // -> false
//	Contains([]string{}, "x")       // -> false
func Contains[T comparable](s []T, target T) bool {
	for _, v := range s {
		if v == target {
			return true
		}
	}
	return false
}

// IndexOf returns the zero-based index of the first element in s equal to
// target, or -1 if no such element exists.
//
// Semantics:
//   - Pure w.r.t. s.
//   - Uses == for equality (T must be comparable).
//   - Returns -1 for an empty/nil slice.
//
// Example:
//
//	IndexOf([]int{10, 20, 30}, 20) // -> 1
//	IndexOf([]int{10, 20, 30}, 99) // -> -1
//	IndexOf([]int{}, 1)            // -> -1
func IndexOf[T comparable](s []T, target T) int {
	for i, v := range s {
		if v == target {
			return i
		}
	}
	return -1
}

// Find returns the first element of s for which f returns true, along with
// a boolean indicating whether a match was found.
//
// Semantics:
//   - Pure w.r.t. s.
//   - On no match, returns T's zero value and false — the zero value is
//     NOT a valid "found" result, always check the bool first (same idiom
//     as map-lookup's comma-ok).
//   - Short-circuits on the first match.
//
// Example:
//
//	v, ok := Find([]int{1, 2, 3, 4}, func(v int) bool { return v > 2 }) // -> 3, true
//	_, ok := Find([]int{1, 2}, func(v int) bool { return v > 100 })     // -> _, false
func Find[T any](s []T, f func(T) bool) (T, bool) {
	for _, v := range s {
		if f(v) {
			return v, true
		}
	}
	var zero T
	return zero, false
}

// Reverse returns a new slice with the elements of s in reverse order.
//
// Semantics:
//   - Pure: the input slice is NOT mutated — a fresh slice is allocated.
//     Use slices.Reverse from the standard library if you want in-place.
//   - Always returns a non-nil slice of length len(s).
//
// Example:
//
//	Reverse([]int{1, 2, 3}) // -> [3 2 1]
//	src := []int{1, 2, 3}
//	_ = Reverse(src)
//	// src is still [1 2 3]
func Reverse[T any](s []T) []T {
	r := make([]T, len(s))
	for i, v := range s {
		r[len(s)-1-i] = v
	}
	return r
}

// Unique returns a new slice containing the distinct elements of s, in the
// order of their first occurrence.
//
// Semantics:
//   - Pure w.r.t. s.
//   - Uses == / map keys for equality (T must be comparable).
//   - Stable: preserves the order of first appearance.
//   - Allocates an auxiliary map of size up to len(s).
//
// Example:
//
//	Unique([]int{1, 2, 2, 3, 1, 4})  // -> [1 2 3 4]
//	Unique([]string{"a","a","a"})    // -> ["a"]
//	Unique([]string{})               // -> [] (len 0)
func Unique[T comparable](s []T) []T {
	seen := make(map[T]struct{}, len(s))
	r := make([]T, 0, len(s))
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		r = append(r, v)
	}
	return r
}

// Flatten concatenates a slice-of-slices into a single slice, preserving
// order.
//
// Semantics:
//   - Pure: neither the outer nor the inner slices are mutated.
//   - Pre-sizes the result to avoid repeated reallocation.
//   - Empty inner slices are tolerated and contribute nothing.
//
// Example:
//
//	Flatten([][]int{{1, 2}, {3}, {}, {4, 5}}) // -> [1 2 3 4 5]
//	Flatten([][]int{})                         // -> [] (len 0)
func Flatten[T any](s [][]T) []T {
	n := 0
	for _, inner := range s {
		n += len(inner)
	}
	r := make([]T, 0, n)
	for _, inner := range s {
		r = append(r, inner...)
	}
	return r
}

// FlatMap applies f to each element of s and concatenates the resulting
// slices into one, preserving order (equivalent to Map followed by Flatten).
//
// Semantics:
//   - Pure w.r.t. s.
//   - f may return a nil or empty slice for an element — those contribute
//     nothing to the output, which makes FlatMap a convenient
//     "map-and-filter-out" primitive.
//   - Supports type transformation: input T, output U.
//
// Example:
//
//	FlatMap([]int{1, 2, 3}, func(v int) []int { return []int{v, v*10} })
//	// -> [1 10 2 20 3 30]
//
//	FlatMap([]int{1, 2, 3}, func(v int) []int {
//	    if v == 2 { return nil }
//	    return []int{v}
//	}) // -> [1 3]
func FlatMap[T, U any](s []T, f func(T) []U) []U {
	r := make([]U, 0, len(s))
	for _, v := range s {
		r = append(r, f(v)...)
	}
	return r
}

// GroupBy partitions s into a map keyed by f(element), where each value
// is the slice of elements that produced that key, in original order.
//
// Semantics:
//   - Pure w.r.t. s.
//   - K must be comparable (map-key constraint).
//   - For an empty/nil input, returns an empty (non-nil) map.
//   - Preserves intra-group order.
//
// Example:
//
//	GroupBy([]int{1, 2, 3, 4, 5, 6}, func(v int) string {
//	    if v%2 == 0 { return "even" }
//	    return "odd"
//	})
//	// -> {"even": [2 4 6], "odd": [1 3 5]}
func GroupBy[T any, K comparable](s []T, f func(T) K) map[K][]T {
	r := make(map[K][]T)
	for _, v := range s {
		k := f(v)
		r[k] = append(r[k], v)
	}
	return r
}

// Chunk splits s into consecutive sub-slices of at most size elements.
//
// Semantics:
//   - Pure w.r.t. s.
//   - Returns nil when size <= 0 (invalid request) — distinguish this from
//     the empty [][]T{} you get for an empty input with a valid size.
//   - The last chunk may be shorter than size when len(s) is not a
//     multiple of size.
//   - The returned sub-slices SHARE the backing array of s (standard Go
//     slicing behavior). Mutating an element through a chunk mutates s,
//     and extending a chunk with append may overlap neighboring chunks.
//     Copy each chunk if you need independence.
//
// Example:
//
//	Chunk([]int{1, 2, 3, 4, 5}, 2) // -> [[1 2] [3 4] [5]]
//	Chunk([]int{1, 2}, 10)         // -> [[1 2]]
//	Chunk([]int{1, 2, 3}, 0)       // -> nil
//	Chunk([]int{}, 3)              // -> [] (len 0)
func Chunk[T any](s []T, size int) [][]T {
	if size <= 0 {
		return nil
	}
	r := make([][]T, 0, (len(s)+size-1)/size)
	for i := 0; i < len(s); i += size {
		end := i + size
		if end > len(s) {
			end = len(s)
		}
		r = append(r, s[i:end])
	}
	return r
}
