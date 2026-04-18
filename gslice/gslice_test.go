package gslice

import (
	"reflect"
	"sort"
	"testing"
)

func TestMap(t *testing.T) {
	got := Map([]int{1, 2, 3}, func(v int) int { return v * 2 })
	want := []int{2, 4, 6}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Map = %v, want %v", got, want)
	}

	got2 := Map([]int{1, 2, 3}, func(v int) string {
		if v == 2 {
			return "two"
		}
		return ""
	})
	if !reflect.DeepEqual(got2, []string{"", "two", ""}) {
		t.Fatalf("Map type change failed: %v", got2)
	}

	if out := Map([]int{}, func(v int) int { return v }); len(out) != 0 {
		t.Fatalf("Map of empty returned %v", out)
	}
}

func TestFilter(t *testing.T) {
	got := Filter([]int{1, 2, 3, 4}, func(v int) bool { return v%2 == 0 })
	want := []int{2, 4}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Filter = %v, want %v", got, want)
	}

	if out := Filter([]int{1, 2, 3}, func(int) bool { return false }); len(out) != 0 {
		t.Fatalf("Filter all-out returned %v", out)
	}

	if out := Filter([]int{}, func(int) bool { return true }); len(out) != 0 {
		t.Fatalf("Filter empty returned %v", out)
	}
}

func TestReduce(t *testing.T) {
	sum := Reduce([]int{1, 2, 3, 4}, 0, func(a, v int) int { return a + v })
	if sum != 10 {
		t.Fatalf("Reduce sum = %d, want 10", sum)
	}

	joined := Reduce([]string{"a", "b", "c"}, "", func(acc, v string) string { return acc + v })
	if joined != "abc" {
		t.Fatalf("Reduce join = %q, want \"abc\"", joined)
	}

	if got := Reduce([]int{}, 42, func(a, v int) int { return a + v }); got != 42 {
		t.Fatalf("Reduce empty returned %d, want 42", got)
	}
}

func TestForEach(t *testing.T) {
	seen := []int{}
	ForEach([]int{1, 2, 3}, func(v int) { seen = append(seen, v) })
	if !reflect.DeepEqual(seen, []int{1, 2, 3}) {
		t.Fatalf("ForEach saw %v", seen)
	}

	ForEach([]int{}, func(int) { t.Fatal("should not run on empty") })
}

func TestAny(t *testing.T) {
	if !Any([]int{1, 2, 3}, func(v int) bool { return v == 2 }) {
		t.Fatal("Any found nothing, want true")
	}
	if Any([]int{1, 2, 3}, func(v int) bool { return v > 10 }) {
		t.Fatal("Any returned true, want false")
	}
	if Any([]int{}, func(int) bool { return true }) {
		t.Fatal("Any over empty should return false")
	}
}

func TestAll(t *testing.T) {
	if !All([]int{2, 4, 6}, func(v int) bool { return v%2 == 0 }) {
		t.Fatal("All even = false, want true")
	}
	if All([]int{2, 3, 4}, func(v int) bool { return v%2 == 0 }) {
		t.Fatal("All even = true with odd present")
	}
	if !All([]int{}, func(int) bool { return false }) {
		t.Fatal("All over empty should be true (vacuously)")
	}
}

func TestContains(t *testing.T) {
	if !Contains([]int{1, 2, 3}, 2) {
		t.Fatal("Contains 2 = false")
	}
	if Contains([]int{1, 2, 3}, 9) {
		t.Fatal("Contains 9 = true")
	}
	if Contains([]string{}, "x") {
		t.Fatal("Contains on empty = true")
	}
}

func TestIndexOf(t *testing.T) {
	if got := IndexOf([]int{10, 20, 30}, 20); got != 1 {
		t.Fatalf("IndexOf = %d, want 1", got)
	}
	if got := IndexOf([]int{10, 20, 30}, 99); got != -1 {
		t.Fatalf("IndexOf missing = %d, want -1", got)
	}
	if got := IndexOf([]int{}, 1); got != -1 {
		t.Fatalf("IndexOf empty = %d, want -1", got)
	}
}

func TestFind(t *testing.T) {
	v, ok := Find([]int{1, 2, 3, 4}, func(v int) bool { return v > 2 })
	if !ok || v != 3 {
		t.Fatalf("Find = %d/%v, want 3/true", v, ok)
	}
	_, ok = Find([]int{1, 2}, func(v int) bool { return v > 100 })
	if ok {
		t.Fatal("Find expected miss")
	}
	v2, ok := Find([]string{}, func(string) bool { return true })
	if ok || v2 != "" {
		t.Fatalf("Find on empty = %q/%v, want \"\"/false", v2, ok)
	}
}

func TestReverse(t *testing.T) {
	got := Reverse([]int{1, 2, 3})
	if !reflect.DeepEqual(got, []int{3, 2, 1}) {
		t.Fatalf("Reverse = %v, want [3 2 1]", got)
	}
	if out := Reverse([]int{}); len(out) != 0 {
		t.Fatalf("Reverse empty = %v", out)
	}

	src := []int{1, 2, 3}
	_ = Reverse(src)
	if !reflect.DeepEqual(src, []int{1, 2, 3}) {
		t.Fatal("Reverse mutated input")
	}
}

func TestUnique(t *testing.T) {
	got := Unique([]int{1, 2, 2, 3, 1, 4})
	if !reflect.DeepEqual(got, []int{1, 2, 3, 4}) {
		t.Fatalf("Unique = %v, want [1 2 3 4]", got)
	}
	if out := Unique([]string{}); len(out) != 0 {
		t.Fatalf("Unique empty = %v", out)
	}
	if got := Unique([]string{"a", "a", "a"}); !reflect.DeepEqual(got, []string{"a"}) {
		t.Fatalf("Unique all-same = %v", got)
	}
}

func TestFlatten(t *testing.T) {
	got := Flatten([][]int{{1, 2}, {3}, {}, {4, 5}})
	if !reflect.DeepEqual(got, []int{1, 2, 3, 4, 5}) {
		t.Fatalf("Flatten = %v", got)
	}
	if out := Flatten([][]int{}); len(out) != 0 {
		t.Fatalf("Flatten empty outer = %v", out)
	}
	if out := Flatten([][]int{{}, {}}); len(out) != 0 {
		t.Fatalf("Flatten all-empty-inner = %v", out)
	}
}

func TestFlatMap(t *testing.T) {
	got := FlatMap([]int{1, 2, 3}, func(v int) []int { return []int{v, v * 10} })
	if !reflect.DeepEqual(got, []int{1, 10, 2, 20, 3, 30}) {
		t.Fatalf("FlatMap = %v", got)
	}
	got2 := FlatMap([]int{1, 2, 3}, func(v int) []int {
		if v == 2 {
			return nil
		}
		return []int{v}
	})
	if !reflect.DeepEqual(got2, []int{1, 3}) {
		t.Fatalf("FlatMap with nils = %v", got2)
	}
}

func TestGroupBy(t *testing.T) {
	got := GroupBy([]int{1, 2, 3, 4, 5, 6}, func(v int) string {
		if v%2 == 0 {
			return "even"
		}
		return "odd"
	})
	sort.Ints(got["even"])
	sort.Ints(got["odd"])
	if !reflect.DeepEqual(got["even"], []int{2, 4, 6}) {
		t.Fatalf("even = %v", got["even"])
	}
	if !reflect.DeepEqual(got["odd"], []int{1, 3, 5}) {
		t.Fatalf("odd = %v", got["odd"])
	}

	empty := GroupBy([]int{}, func(int) int { return 0 })
	if len(empty) != 0 {
		t.Fatalf("GroupBy empty = %v", empty)
	}
}

func TestChunk(t *testing.T) {
	got := Chunk([]int{1, 2, 3, 4, 5}, 2)
	want := [][]int{{1, 2}, {3, 4}, {5}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Chunk 2 = %v, want %v", got, want)
	}

	exact := Chunk([]int{1, 2, 3, 4}, 2)
	if !reflect.DeepEqual(exact, [][]int{{1, 2}, {3, 4}}) {
		t.Fatalf("Chunk exact = %v", exact)
	}

	big := Chunk([]int{1, 2}, 10)
	if !reflect.DeepEqual(big, [][]int{{1, 2}}) {
		t.Fatalf("Chunk oversized = %v", big)
	}

	if out := Chunk([]int{1, 2, 3}, 0); out != nil {
		t.Fatalf("Chunk 0 = %v, want nil", out)
	}
	if out := Chunk([]int{1, 2, 3}, -1); out != nil {
		t.Fatalf("Chunk negative = %v, want nil", out)
	}

	if out := Chunk([]int{}, 3); len(out) != 0 {
		t.Fatalf("Chunk empty = %v", out)
	}
}
