// Package collections provides useful collection data structures for Go.
package collections

import (
	"errors"
	"fmt"
)

type set[T comparable] map[T]struct{}

// NewSet returns a new set. This set is NOT concurrent safe.
// You must include a comparable (hashable) type in the type parameter.
// Example: NewSet[int]()
func NewSet[T comparable]() set[T] {
	return set[T]{}
}

// Add function simply adds an element to the set.
// Returns the set itself for chaining purposes.
func (s set[T]) Add(element T) set[T] {
	s[element] = struct{}{}
	return s
}

// AddAll adds all elements to the set
// Returns the set itself for chaining purposes.
func (s set[T]) AddAll(elements []T) set[T] {
	for _, v := range elements {
		s.Add(v)
	}
	return s
}

// Has returns true if the element is in the set.
func (s set[T]) Has(element T) bool {
	if _, ok := s[element]; ok {
		return true
	}
	return false
}

// HasAllOf returns true if all elements are in the set
func (s set[T]) HasAllOf(elements []T) bool {
	for _, v := range elements {
		if _, ok := s[v]; !ok {
			return false
		}
	}
	return true
}

// Remove removes an element from the set.
// Returns the set itself for chaining purposes.
// Returns an error if the element is not in the set.
func (s set[T]) Remove(element T) error {
	if !s.Has(element) {
		return errors.New(fmt.Sprintf("element %v is not in the set", element))
	}
	delete(s, element)
	return nil
}

// MustRemove removes an element from the set.
// Panics if the element is not in the set.
func (s set[T]) MustRemove(element T) set[T] {
	if err := s.Remove(element); err != nil {
		panic(err)
	}
	return s
}

// RemoveAll removes all elements in the array from the set
// Returns an error if any of the elements are not in the set
// This function is atomic.
func (s set[T]) RemoveAll(elements []T) error {
	if elements == nil || len(elements) == 0 {
		return errors.New("elements array is nil or empty")
	}
	if !s.HasAllOf(elements) {
		return errors.New("not all elements are in the set")
	}

	for _, v := range elements {
		delete(s, v)
	}
	return nil
}

// MustRemoveAll removes all elements in the array from the set.
// Panics if any of the elements are not in the set.
func (s set[T]) MustRemoveAll(elements []T) set[T] {
	if err := s.RemoveAll(elements); err != nil {
		panic(err)
	}
	return s
}

// Clear removes all elements in the set, might be useful in niche cases
// Returns the set itself for chaining purposes.
func (s set[T]) Clear() set[T] {
	for k := range s {
		delete(s, k)
	}
	return s
}

// Copy copies the set to a new set
func (s set[T]) Copy() set[T] {
	newSet := make(set[T], len(s))

	for key, value := range s {
		newSet[key] = value
	}
	return newSet
}

// Size gets the size of the set
func (s set[T]) Size() int {
	return len(s)
}

// ToSlice returns a slice of all elements in the set
// This is NOT ordered as the set does not guarantee order
func (s set[T]) ToSlice() []T {
	slice := make([]T, len(s))
	i := 0
	for key := range s {
		slice[i] = key
		i++
	}
	return slice
}

// IsSuperSetOf returns true if the set is a super set of the other set
func (s set[T]) IsSuperSetOf(other set[T]) bool {
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

// IsSubSetOf returns true if the set is a subset of the other set
func (s set[T]) IsSubSetOf(other set[T]) bool {
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

// Diff removes all the elements in s that the other set has too.
func (s set[T]) Diff(other set[T]) set[T] {
	for key := range other {
		_ = s.Remove(key)
	}
	return s
}

// NewDiff returns a new set with all the elements in s that the other set does not have.
func (s set[T]) NewDiff(other set[T]) set[T] {
	newSet := NewSet[T]()
	for key := range s {
		if !other.Has(key) {
			newSet.Add(key)
		}
	}
	return newSet
}

// Union unions the original set with another set
func (s set[T]) Union(other set[T]) set[T] {
	if other == nil || other.Size() == 0 {
		return s
	}
	for key := range other {
		s.Add(key)
	}
	return s
}

// NewUnion returns a new set with the union of the original set and another set
func (s set[T]) NewUnion(other set[T]) set[T] {
	newSet := s.Copy()
	return newSet.Union(other)
}

// Intersect intersects the original set with another set
func (s set[T]) Intersect(other set[T]) set[T] {
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

// NewIntersect returns a new set with the intersection of the original set and another set
func (s set[T]) NewIntersect(other set[T]) set[T] {
	newSet := s.Copy()
	return newSet.Intersect(other)
}
