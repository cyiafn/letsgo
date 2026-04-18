package letsgo

import "testing"

func TestIsNil(t *testing.T) {
	var nilPtr *int
	if !IsNil(nilPtr) {
		t.Fatal("IsNil(nil) = false, want true")
	}
	v := 42
	if IsNil(&v) {
		t.Fatal("IsNil(&v) = true, want false")
	}
}

func TestPtr(t *testing.T) {
	p := Ptr(7)
	if p == nil {
		t.Fatal("Ptr(7) returned nil")
	}
	if *p != 7 {
		t.Fatalf("*Ptr(7) = %d, want 7", *p)
	}

	s := Ptr("hello")
	if *s != "hello" {
		t.Fatalf("*Ptr(\"hello\") = %q, want \"hello\"", *s)
	}
}

func TestIfNilDefault(t *testing.T) {
	var nilPtr *int
	got := IfNilDefault(nilPtr, 5)
	if got == nil || *got != 5 {
		t.Fatalf("IfNilDefault(nil, 5) = %v, want *5", got)
	}

	v := 10
	got = IfNilDefault(&v, 5)
	if got == nil || *got != 10 {
		t.Fatalf("IfNilDefault(&10, 5) = %v, want *10", got)
	}
}

func TestIfNilDefaultValue(t *testing.T) {
	var nilPtr *int
	if got := IfNilDefaultValue(nilPtr, 5); got != 5 {
		t.Fatalf("IfNilDefaultValue(nil, 5) = %d, want 5", got)
	}
	v := 10
	if got := IfNilDefaultValue(&v, 5); got != 10 {
		t.Fatalf("IfNilDefaultValue(&10, 5) = %d, want 10", got)
	}
}
