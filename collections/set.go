// Package collections provides useful collection data structures for Go.
package collections

import (
	"errors"
	"fmt"
)

type Set[T comparable] map[T]struct{}

// NewSet returns a new Set. This Set is NOT concurrent safe.
// You must include a comparable (hashable) type in the type parameter.
// Example: NewSet[int]()
func NewSet[T comparable]() Set[T] {
	return Set[T]{}
}

// Add function simply adds an element to the Set.
// Returns the Set itself for chaining purposes.
func (s Set[T]) Add(element T) Set[T] {
	s[element] = struct{}{}
	return s
}

// AddAll adds all elements to the Set
// Returns the Set itself for chaining purposes.
func (s Set[T]) AddAll(elements []T) Set[T] {
	for _, v := range elements {
		s.Add(v)
	}
	return s
}

// Has returns true if the element is in the Set.
func (s Set[T]) Has(element T) bool {
	if _, ok := s[element]; ok {
		return true
	}
	return false
}

// HasAllOf returns true if all elements are in the Set
func (s Set[T]) HasAllOf(elements []T) bool {
	for _, v := range elements {
		if _, ok := s[v]; !ok {
			return false
		}
	}
	return true
}

// Remove removes an element from the Set.
// Returns the Set itself for chaining purposes.
// Returns an error if the element is not in the Set.
func (s Set[T]) Remove(element T) error {
	if !s.Has(element) {
		return errors.New(fmt.Sprintf("element %v is not in the Set", element))
	}
	delete(s, element)
	return nil
}

// MustRemove removes an element from the Set.
// Panics if the element is not in the Set.
func (s Set[T]) MustRemove(element T) Set[T] {
	if err := s.Remove(element); err != nil {
		panic(err)
	}
	return s
}

// RemoveAll removes all elements in the array from the Set
// Returns an error if any of the elements are not in the Set
// This function is atomic.
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

// MustRemoveAll removes all elements in the array from the Set.
// Panics if any of the elements are not in the Set.
func (s Set[T]) MustRemoveAll(elements []T) Set[T] {
	if err := s.RemoveAll(elements); err != nil {
		panic(err)
	}
	return s
}

// Clear removes all elements in the Set, might be useful in niche cases
// Returns the Set itself for chaining purposes.
func (s Set[T]) Clear() Set[T] {
	for k := range s {
		delete(s, k)
	}
	return s
}

// Copy copies the Set to a new Set
func (s Set[T]) Copy() Set[T] {
	newSet := make(Set[T], len(s))

	for key, value := range s {
		newSet[key] = value
	}
	return newSet
}

// Size gets the size of the Set
func (s Set[T]) Size() int {
	return len(s)
}

// ToSlice returns a slice of all elements in the Set
// This is NOT ordered as the Set does not guarantee order
func (s Set[T]) ToSlice() []T {
	slice := make([]T, len(s))
	i := 0
	for key := range s {
		slice[i] = key
		i++
	}
	return slice
}

// IsSuperSetOf returns true if the Set is a super Set of the other Set
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

// IsSubSetOf returns true if the Set is a subSet of the other Set
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

// Diff removes all the elements in s that the other Set has too.
func (s Set[T]) Diff(other Set[T]) Set[T] {
	for key := range other {
		_ = s.Remove(key)
	}
	return s
}

// NewDiff returns a new Set with all the elements in s that the other Set does not have.
func (s Set[T]) NewDiff(other Set[T]) Set[T] {
	newSet := NewSet[T]()
	for key := range s {
		if !other.Has(key) {
			newSet.Add(key)
		}
	}
	return newSet
}

// Union unions the original Set with another Set
func (s Set[T]) Union(other Set[T]) Set[T] {
	if other == nil || other.Size() == 0 {
		return s
	}
	for key := range other {
		s.Add(key)
	}
	return s
}

// NewUnion returns a new Set with the union of the original Set and another Set
func (s Set[T]) NewUnion(other Set[T]) Set[T] {
	newSet := s.Copy()
	return newSet.Union(other)
}

// Intersect intersects the original Set with another Set
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

// NewIntersect returns a new Set with the intersection of the original Set and another Set
func (s Set[T]) NewIntersect(other Set[T]) Set[T] {
	newSet := s.Copy()
	return newSet.Intersect(other)
}
