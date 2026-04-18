// Package collections provides useful collection data structures for Go.
package collections

import (
	"errors"
	"fmt"
)

// Set is a generic unordered collection of unique elements of type T.
//
// Semantics:
//   - Backed by a map[T]struct{}; every method has the same big-O behavior
//     as the equivalent map operation (O(1) amortized for Add/Has/Remove,
//     O(n) for bulk operations).
//   - NOT goroutine-safe. Synchronize externally if multiple goroutines
//     may mutate the same Set.
//   - Reference semantics: because the underlying type is a map, assigning
//     or passing a Set copies the header but shares the storage. Use Copy()
//     to get an independent Set.
//   - Iteration order is unspecified (inherited from Go maps).
//   - T must be comparable (the usual map-key constraint).
//
// Example:
//
//	s := collections.NewSet[int]().AddAll([]int{1, 2, 3})
//	s.Has(2)                 // -> true
//	s.Size()                 // -> 3
type Set[T comparable] map[T]struct{}

// NewSet constructs an empty Set[T]. The concrete type argument must be
// comparable (hashable).
//
// Semantics:
//   - Allocates a fresh underlying map; the returned Set is always non-nil
//     and has Size() == 0.
//   - NOT concurrent-safe — see the Set type doc.
//
// Example:
//
//	ints := collections.NewSet[int]()
//	names := collections.NewSet[string]()
func NewSet[T comparable]() Set[T] {
	return Set[T]{}
}

// Add inserts element into the Set. No-op if element is already present.
//
// Semantics:
//   - Mutates the receiver.
//   - Returns the receiver to allow method chaining.
//   - Idempotent: adding an existing element does not change Size().
//
// Example:
//
//	s := collections.NewSet[string]()
//	s.Add("a").Add("b").Add("a") // Size() == 2
func (s Set[T]) Add(element T) Set[T] {
	s[element] = struct{}{}
	return s
}

// AddAll inserts every element from the given slice into the Set.
//
// Semantics:
//   - Mutates the receiver.
//   - Returns the receiver to allow chaining.
//   - Duplicates within elements (and duplicates against existing members)
//     are silently ignored, consistent with set semantics.
//   - Nil or empty elements is a safe no-op.
//
// Example:
//
//	s := collections.NewSet[int]().AddAll([]int{1, 2, 3, 2})
//	s.Size() // -> 3
func (s Set[T]) AddAll(elements []T) Set[T] {
	for _, v := range elements {
		s.Add(v)
	}
	return s
}

// Has reports whether element is in the Set.
//
// Semantics:
//   - Pure: does not mutate the Set.
//   - O(1) average.
//
// Example:
//
//	s := collections.NewSet[string]().Add("a")
//	s.Has("a") // -> true
//	s.Has("z") // -> false
func (s Set[T]) Has(element T) bool {
	if _, ok := s[element]; ok {
		return true
	}
	return false
}

// HasAllOf reports whether every element in the given slice is present in
// the Set.
//
// Semantics:
//   - Pure: does not mutate the Set.
//   - Short-circuits on the first missing element.
//   - Vacuously true for a nil or empty elements slice.
//
// Example:
//
//	s := collections.NewSet[int]().AddAll([]int{1, 2, 3})
//	s.HasAllOf([]int{1, 2})   // -> true
//	s.HasAllOf([]int{1, 99})  // -> false
//	s.HasAllOf(nil)           // -> true (vacuously)
func (s Set[T]) HasAllOf(elements []T) bool {
	for _, v := range elements {
		if _, ok := s[v]; !ok {
			return false
		}
	}
	return true
}

// Remove deletes element from the Set.
//
// Semantics:
//   - Mutates the receiver.
//   - Returns a non-nil error if element is not in the Set; the Set is
//     left unchanged in that case.
//   - Unlike Add/AddAll, this does not return the Set, because the error
//     channel is already occupying the return slot.
//
// Example:
//
//	s := collections.NewSet[int]().AddAll([]int{1, 2})
//	s.Remove(1)  // nil, s now {2}
//	s.Remove(99) // non-nil error, s unchanged
func (s Set[T]) Remove(element T) error {
	if !s.Has(element) {
		return errors.New(fmt.Sprintf("element %v is not in the Set", element))
	}
	delete(s, element)
	return nil
}

// MustRemove is the panicking counterpart to Remove, for use when the
// element is known to be present and a missing element indicates a
// programmer error.
//
// Semantics:
//   - Mutates the receiver.
//   - Panics if element is not present (wrapping the error from Remove).
//   - Returns the receiver to allow chaining.
//
// Example:
//
//	s := collections.NewSet[int]().Add(1)
//	s.MustRemove(1) // ok, size 0
//	s.MustRemove(1) // panics: element 1 is not in the Set
func (s Set[T]) MustRemove(element T) Set[T] {
	if err := s.Remove(element); err != nil {
		panic(err)
	}
	return s
}

// RemoveAll deletes every element in the given slice from the Set,
// atomically.
//
// Semantics:
//   - Mutates the receiver.
//   - Atomic: if ANY element is missing, the Set is left UNCHANGED and a
//     non-nil error is returned. No partial removal occurs.
//   - Returns a non-nil error when elements is nil or empty — this method
//     treats "remove nothing" as a likely caller bug rather than a no-op.
//     Use Clear() if you want to empty the Set.
//
// Example:
//
//	s := collections.NewSet[int]().AddAll([]int{1, 2, 3})
//	s.RemoveAll([]int{1, 2})   // nil, s == {3}
//	s.RemoveAll([]int{3, 99})  // error, s still {3}
func (s Set[T]) RemoveAll(elements []T) error {
	if elements == nil || len(elements) == 0 {
		return errors.New("elements array is nil or empty")
	}
	if !s.HasAllOf(elements) {
		return errors.New("not all elements are in the Set")
	}

	for _, v := range elements {
		delete(s, v)
	}
	return nil
}

// MustRemoveAll is the panicking counterpart to RemoveAll.
//
// Semantics:
//   - Mutates the receiver.
//   - Atomic like RemoveAll: on a missing element the Set is unchanged
//     and the function panics.
//   - Returns the receiver to allow chaining.
//
// Example:
//
//	s := collections.NewSet[int]().AddAll([]int{1, 2, 3})
//	s.MustRemoveAll([]int{1, 2}) // ok
//	s.MustRemoveAll([]int{99})   // panics, s unchanged
func (s Set[T]) MustRemoveAll(elements []T) Set[T] {
	if err := s.RemoveAll(elements); err != nil {
		panic(err)
	}
	return s
}

// Clear removes every element from the Set.
//
// Semantics:
//   - Mutates the receiver in place (keeps the same underlying map;
//     other variables referring to the same Set will observe the change).
//   - Returns the receiver to allow chaining.
//
// Example:
//
//	s := collections.NewSet[int]().AddAll([]int{1, 2, 3}).Clear()
//	s.Size() // -> 0
func (s Set[T]) Clear() Set[T] {
	for k := range s {
		delete(s, k)
	}
	return s
}

// Copy returns a shallow, independent copy of the Set.
//
// Semantics:
//   - Allocates a new underlying map with the same elements.
//   - Mutations to the copy do NOT affect the original, and vice versa.
//   - The elements themselves are copied by assignment; for pointer or
//     reference types, the pointees are still shared.
//
// Example:
//
//	orig := collections.NewSet[int]().AddAll([]int{1, 2, 3})
//	dup := orig.Copy()
//	dup.Add(99)
//	orig.Has(99) // -> false
func (s Set[T]) Copy() Set[T] {
	newSet := make(Set[T], len(s))

	for key, value := range s {
		newSet[key] = value
	}
	return newSet
}

// Size returns the number of elements in the Set.
//
// Semantics:
//   - Pure: does not mutate.
//   - O(1), delegates to len on the underlying map.
//
// Example:
//
//	collections.NewSet[int]().AddAll([]int{1, 2, 3}).Size() // -> 3
func (s Set[T]) Size() int {
	return len(s)
}

// ToSlice returns the Set's elements as a slice.
//
// Semantics:
//   - Pure: the Set is not modified.
//   - UNORDERED: element order in the returned slice is not specified and
//     may differ between calls. Sort the result if you need determinism.
//   - Allocates a fresh slice of length Size().
//
// Example:
//
//	s := collections.NewSet[int]().AddAll([]int{3, 1, 2})
//	out := s.ToSlice()
//	sort.Ints(out) // -> [1 2 3]
func (s Set[T]) ToSlice() []T {
	slice := make([]T, len(s))
	i := 0
	for key := range s {
		slice[i] = key
		i++
	}
	return slice
}

// IsSuperSetOf reports whether s contains every element of other.
//
// Semantics:
//   - Pure: neither s nor other is modified.
//   - Treats a nil or empty other as trivially contained (returns true).
//   - Equal sets are supersets of each other by this definition.
//
// Example:
//
//	a := collections.NewSet[int]().AddAll([]int{1, 2, 3})
//	b := collections.NewSet[int]().AddAll([]int{1, 2})
//	a.IsSuperSetOf(b) // -> true
//	b.IsSuperSetOf(a) // -> false
func (s Set[T]) IsSuperSetOf(other Set[T]) bool {
	if other == nil || other.Size() == 0 {
		return true
	}

	for key := range other {
		if !s.Has(key) {
			return false
		}
	}

	return true
}

// IsSubSetOf reports whether every element of s is in other.
//
// Semantics:
//   - Pure: neither s nor other is modified.
//   - An empty (or nil) s is a subset of any Set, returning true.
//   - Equal sets are subsets of each other by this definition.
//
// Example:
//
//	a := collections.NewSet[int]().AddAll([]int{1, 2})
//	b := collections.NewSet[int]().AddAll([]int{1, 2, 3})
//	a.IsSubSetOf(b) // -> true
//	b.IsSubSetOf(a) // -> false
func (s Set[T]) IsSubSetOf(other Set[T]) bool {
	if s == nil || s.Size() == 0 {
		return true
	}

	for key := range s {
		if !other.Has(key) {
			return false
		}
	}

	return true
}

// Diff removes from s every element that also appears in other (in-place
// set difference: s := s \ other).
//
// Semantics:
//   - Mutates the receiver s; other is unchanged.
//   - Returns the receiver to allow chaining.
//   - Elements that are in other but not in s are silently ignored.
//
// Example:
//
//	s := collections.NewSet[int]().AddAll([]int{1, 2, 3, 4})
//	o := collections.NewSet[int]().AddAll([]int{2, 4, 99})
//	s.Diff(o)            // s becomes {1, 3}
func (s Set[T]) Diff(other Set[T]) Set[T] {
	for key := range other {
		_ = s.Remove(key)
	}
	return s
}

// NewDiff returns a new Set containing the elements of s that are not in
// other (non-mutating set difference: result = s \ other).
//
// Semantics:
//   - Pure: neither s nor other is modified.
//   - Allocates a fresh Set.
//
// Example:
//
//	s := collections.NewSet[int]().AddAll([]int{1, 2, 3, 4})
//	o := collections.NewSet[int]().AddAll([]int{2, 4})
//	s.NewDiff(o).ToSlice() // -> [1 3] (unordered)
//	s.Size()               // -> 4 (original untouched)
func (s Set[T]) NewDiff(other Set[T]) Set[T] {
	newSet := NewSet[T]()
	for key := range s {
		if !other.Has(key) {
			newSet.Add(key)
		}
	}
	return newSet
}

// Union adds every element of other to s in place (set union: s := s ∪ other).
//
// Semantics:
//   - Mutates the receiver s; other is unchanged.
//   - Returns the receiver to allow chaining.
//   - nil or empty other is a no-op.
//
// Example:
//
//	a := collections.NewSet[int]().AddAll([]int{1, 2})
//	b := collections.NewSet[int]().AddAll([]int{2, 3})
//	a.Union(b)   // a becomes {1, 2, 3}
func (s Set[T]) Union(other Set[T]) Set[T] {
	if other == nil || other.Size() == 0 {
		return s
	}
	for key := range other {
		s.Add(key)
	}
	return s
}

// NewUnion returns a new Set equal to s ∪ other.
//
// Semantics:
//   - Pure: neither s nor other is modified.
//   - Allocates a fresh Set.
//
// Example:
//
//	a := collections.NewSet[int]().AddAll([]int{1, 2})
//	b := collections.NewSet[int]().AddAll([]int{2, 3})
//	u := a.NewUnion(b)   // u == {1, 2, 3}
//	a.Size()             // -> 2 (original untouched)
func (s Set[T]) NewUnion(other Set[T]) Set[T] {
	newSet := s.Copy()
	return newSet.Union(other)
}

// Intersect restricts s to elements also in other (in-place intersection:
// s := s ∩ other).
//
// Semantics:
//   - Mutates the receiver s; other is unchanged.
//   - Returns the receiver to allow chaining.
//   - QUIRK: passing a nil or empty other is treated as a no-op (s is
//     returned unchanged) rather than clearing s. Callers who want
//     "intersect with empty produces empty" should call s.Clear() instead.
//
// Example:
//
//	a := collections.NewSet[int]().AddAll([]int{1, 2, 3})
//	b := collections.NewSet[int]().AddAll([]int{2, 3, 4})
//	a.Intersect(b)   // a becomes {2, 3}
func (s Set[T]) Intersect(other Set[T]) Set[T] {
	if other == nil || other.Size() == 0 {
		return s
	}
	for key := range s {
		if !other.Has(key) {
			_ = s.Remove(key)
		}
	}

	return s
}

// NewIntersect returns a new Set equal to s ∩ other.
//
// Semantics:
//   - Pure: neither s nor other is modified.
//   - Allocates a fresh Set.
//   - Inherits Intersect's quirk: a nil/empty other yields a copy of s,
//     not an empty Set.
//
// Example:
//
//	a := collections.NewSet[int]().AddAll([]int{1, 2, 3})
//	b := collections.NewSet[int]().AddAll([]int{2, 3, 4})
//	i := a.NewIntersect(b) // i == {2, 3}
//	a.Size()               // -> 3 (original untouched)
func (s Set[T]) NewIntersect(other Set[T]) Set[T] {
	newSet := s.Copy()
	return newSet.Intersect(other)
}
