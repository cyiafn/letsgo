package collections

import (
	"sort"
	"testing"
)

func TestNewSet(t *testing.T) {
	s := NewSet[int]()
	if s == nil {
		t.Fatal("NewSet returned nil")
	}
	if s.Size() != 0 {
		t.Fatalf("new set size = %d, want 0", s.Size())
	}
}

func TestAddAndHas(t *testing.T) {
	s := NewSet[string]()
	s.Add("a").Add("b")
	if !s.Has("a") || !s.Has("b") {
		t.Fatal("Has missing expected elements")
	}
	if s.Has("c") {
		t.Fatal("Has reported absent element as present")
	}
	s.Add("a")
	if s.Size() != 2 {
		t.Fatalf("duplicate add changed size to %d", s.Size())
	}
}

func TestAddAll(t *testing.T) {
	s := NewSet[int]().AddAll([]int{1, 2, 3, 2})
	if s.Size() != 3 {
		t.Fatalf("size after AddAll = %d, want 3", s.Size())
	}
	if !s.HasAllOf([]int{1, 2, 3}) {
		t.Fatal("HasAllOf failed after AddAll")
	}
}

func TestHasAllOf(t *testing.T) {
	s := NewSet[int]().AddAll([]int{1, 2, 3})
	if !s.HasAllOf([]int{1, 2}) {
		t.Fatal("subset check failed")
	}
	if s.HasAllOf([]int{1, 4}) {
		t.Fatal("should return false when missing elements")
	}
	if !s.HasAllOf([]int{}) {
		t.Fatal("empty should be vacuously true")
	}
}

func TestRemove(t *testing.T) {
	s := NewSet[int]().AddAll([]int{1, 2, 3})
	if err := s.Remove(2); err != nil {
		t.Fatalf("Remove present returned err: %v", err)
	}
	if s.Has(2) {
		t.Fatal("Remove didn't delete")
	}
	if err := s.Remove(99); err == nil {
		t.Fatal("Remove absent should return error")
	}
}

func TestMustRemove(t *testing.T) {
	s := NewSet[int]().Add(1)
	s.MustRemove(1)
	if s.Size() != 0 {
		t.Fatal("MustRemove didn't delete")
	}
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("MustRemove absent did not panic")
		}
	}()
	s.MustRemove(1)
}

func TestRemoveAll(t *testing.T) {
	s := NewSet[int]().AddAll([]int{1, 2, 3})
	if err := s.RemoveAll([]int{1, 2}); err != nil {
		t.Fatalf("RemoveAll returned err: %v", err)
	}
	if s.Size() != 1 || !s.Has(3) {
		t.Fatalf("unexpected state after RemoveAll: %v", s)
	}

	s2 := NewSet[int]().AddAll([]int{1, 2, 3})
	if err := s2.RemoveAll([]int{1, 99}); err == nil {
		t.Fatal("RemoveAll should error when any element is missing")
	}
	if s2.Size() != 3 {
		t.Fatal("RemoveAll should be atomic — no partial removal")
	}

	if err := s2.RemoveAll(nil); err == nil {
		t.Fatal("RemoveAll nil should error")
	}
	if err := s2.RemoveAll([]int{}); err == nil {
		t.Fatal("RemoveAll empty should error")
	}
}

func TestMustRemoveAll(t *testing.T) {
	s := NewSet[int]().AddAll([]int{1, 2, 3})
	s.MustRemoveAll([]int{1, 2})
	if s.Size() != 1 {
		t.Fatal("MustRemoveAll didn't delete")
	}
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("MustRemoveAll with missing should panic")
		}
	}()
	s.MustRemoveAll([]int{99})
}

func TestClear(t *testing.T) {
	s := NewSet[int]().AddAll([]int{1, 2, 3}).Clear()
	if s.Size() != 0 {
		t.Fatalf("Clear left size %d", s.Size())
	}
}

func TestCopy(t *testing.T) {
	s := NewSet[int]().AddAll([]int{1, 2, 3})
	c := s.Copy()
	if c.Size() != s.Size() || !c.HasAllOf([]int{1, 2, 3}) {
		t.Fatal("Copy missing elements")
	}
	c.Add(99)
	if s.Has(99) {
		t.Fatal("Copy not independent of original")
	}
}

func TestSize(t *testing.T) {
	s := NewSet[int]()
	if s.Size() != 0 {
		t.Fatal("empty size != 0")
	}
	s.Add(1).Add(2)
	if s.Size() != 2 {
		t.Fatalf("size = %d, want 2", s.Size())
	}
}

func TestToSlice(t *testing.T) {
	s := NewSet[int]().AddAll([]int{3, 1, 2})
	slice := s.ToSlice()
	sort.Ints(slice)
	if len(slice) != 3 || slice[0] != 1 || slice[1] != 2 || slice[2] != 3 {
		t.Fatalf("ToSlice = %v", slice)
	}

	empty := NewSet[int]().ToSlice()
	if len(empty) != 0 {
		t.Fatalf("empty ToSlice = %v", empty)
	}
}

func TestIsSuperSetOf(t *testing.T) {
	s := NewSet[int]().AddAll([]int{1, 2, 3})
	if !s.IsSuperSetOf(NewSet[int]().AddAll([]int{1, 2})) {
		t.Fatal("superset check failed")
	}
	if s.IsSuperSetOf(NewSet[int]().AddAll([]int{1, 99})) {
		t.Fatal("should not be superset")
	}
	if !s.IsSuperSetOf(NewSet[int]()) {
		t.Fatal("any set should be a superset of empty")
	}
	if !s.IsSuperSetOf(nil) {
		t.Fatal("any set should be a superset of nil")
	}
}

func TestIsSubSetOf(t *testing.T) {
	s := NewSet[int]().AddAll([]int{1, 2})
	if !s.IsSubSetOf(NewSet[int]().AddAll([]int{1, 2, 3})) {
		t.Fatal("subset check failed")
	}
	if s.IsSubSetOf(NewSet[int]().AddAll([]int{1, 99})) {
		t.Fatal("should not be subset")
	}
	if !NewSet[int]().IsSubSetOf(NewSet[int]().AddAll([]int{1, 2})) {
		t.Fatal("empty set should be subset of anything")
	}
}

func TestDiff(t *testing.T) {
	s := NewSet[int]().AddAll([]int{1, 2, 3, 4})
	other := NewSet[int]().AddAll([]int{2, 4, 99})
	s.Diff(other)
	slice := s.ToSlice()
	sort.Ints(slice)
	if len(slice) != 2 || slice[0] != 1 || slice[1] != 3 {
		t.Fatalf("Diff = %v, want [1 3]", slice)
	}
}

func TestNewDiff(t *testing.T) {
	s := NewSet[int]().AddAll([]int{1, 2, 3, 4})
	other := NewSet[int]().AddAll([]int{2, 4})
	diff := s.NewDiff(other)
	slice := diff.ToSlice()
	sort.Ints(slice)
	if len(slice) != 2 || slice[0] != 1 || slice[1] != 3 {
		t.Fatalf("NewDiff = %v", slice)
	}
	if s.Size() != 4 {
		t.Fatal("NewDiff mutated original")
	}
}

func TestUnion(t *testing.T) {
	s := NewSet[int]().AddAll([]int{1, 2})
	other := NewSet[int]().AddAll([]int{2, 3})
	s.Union(other)
	slice := s.ToSlice()
	sort.Ints(slice)
	if len(slice) != 3 || slice[0] != 1 || slice[1] != 2 || slice[2] != 3 {
		t.Fatalf("Union = %v", slice)
	}

	s2 := NewSet[int]().Add(1)
	s2.Union(NewSet[int]())
	if s2.Size() != 1 {
		t.Fatalf("Union with empty changed size: %v", s2.ToSlice())
	}
	s2.Union(nil)
	if s2.Size() != 1 {
		t.Fatalf("Union with nil changed size")
	}
}

func TestNewUnion(t *testing.T) {
	s := NewSet[int]().AddAll([]int{1, 2})
	other := NewSet[int]().AddAll([]int{2, 3})
	u := s.NewUnion(other)
	slice := u.ToSlice()
	sort.Ints(slice)
	if len(slice) != 3 {
		t.Fatalf("NewUnion = %v", slice)
	}
	if s.Size() != 2 {
		t.Fatal("NewUnion mutated original")
	}
}

func TestIntersect(t *testing.T) {
	s := NewSet[int]().AddAll([]int{1, 2, 3})
	s.Intersect(NewSet[int]().AddAll([]int{2, 3, 4}))
	slice := s.ToSlice()
	sort.Ints(slice)
	if len(slice) != 2 || slice[0] != 2 || slice[1] != 3 {
		t.Fatalf("Intersect = %v", slice)
	}

	s2 := NewSet[int]().Add(1)
	s2.Intersect(NewSet[int]())
	if s2.Size() != 1 {
		t.Fatal("Intersect with empty unexpectedly cleared set (per current impl)")
	}
}

func TestNewIntersect(t *testing.T) {
	s := NewSet[int]().AddAll([]int{1, 2, 3})
	other := NewSet[int]().AddAll([]int{2, 3, 4})
	inter := s.NewIntersect(other)
	slice := inter.ToSlice()
	sort.Ints(slice)
	if len(slice) != 2 || slice[0] != 2 || slice[1] != 3 {
		t.Fatalf("NewIntersect = %v", slice)
	}
	if s.Size() != 3 {
		t.Fatal("NewIntersect mutated original")
	}
}
