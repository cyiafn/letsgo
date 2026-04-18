package letsgo

import "testing"

func TestDumpToJSON(t *testing.T) {
	t.Run("nil returns empty string", func(t *testing.T) {
		if got := DumpToJSON(nil); got != "" {
			t.Fatalf("DumpToJSON(nil) = %q, want \"\"", got)
		}
	})

	t.Run("string", func(t *testing.T) {
		got := DumpToJSON("hello")
		want := `"hello"`
		if got != want {
			t.Fatalf("DumpToJSON(\"hello\") = %q, want %q", got, want)
		}
	})

	t.Run("int", func(t *testing.T) {
		if got := DumpToJSON(42); got != "42" {
			t.Fatalf("DumpToJSON(42) = %q, want \"42\"", got)
		}
	})

	t.Run("struct", func(t *testing.T) {
		type Point struct {
			X int `json:"x"`
			Y int `json:"y"`
		}
		got := DumpToJSON(Point{X: 1, Y: 2})
		want := `{"x":1,"y":2}`
		if got != want {
			t.Fatalf("DumpToJSON(Point{1,2}) = %q, want %q", got, want)
		}
	})

	t.Run("slice", func(t *testing.T) {
		got := DumpToJSON([]int{1, 2, 3})
		want := "[1,2,3]"
		if got != want {
			t.Fatalf("DumpToJSON([1,2,3]) = %q, want %q", got, want)
		}
	})

	t.Run("unmarshalable returns empty", func(t *testing.T) {
		got := DumpToJSON(make(chan int))
		if got != "" {
			t.Fatalf("DumpToJSON(chan) = %q, want \"\"", got)
		}
	})
}
