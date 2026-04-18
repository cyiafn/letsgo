package choose

import "testing"

func TestIf(t *testing.T) {
	if got := If(true, "yes", "no"); got != "yes" {
		t.Fatalf("If(true) = %q, want \"yes\"", got)
	}
	if got := If(false, "yes", "no"); got != "no" {
		t.Fatalf("If(false) = %q, want \"no\"", got)
	}
	if got := If(true, 1, 2); got != 1 {
		t.Fatalf("If(true,1,2) = %d, want 1", got)
	}
}

func TestIfLazyL(t *testing.T) {
	trueCalled := false
	got := IfLazyL(true, func() int { trueCalled = true; return 1 }, 2)
	if got != 1 || !trueCalled {
		t.Fatalf("true branch: got=%d called=%v", got, trueCalled)
	}

	trueCalled = false
	got = IfLazyL(false, func() int { trueCalled = true; return 1 }, 2)
	if got != 2 {
		t.Fatalf("false branch: got=%d want 2", got)
	}
	if trueCalled {
		t.Fatal("true branch thunk should not be invoked when condition is false")
	}
}

func TestIfLazyR(t *testing.T) {
	falseCalled := false
	got := IfLazyR(false, 1, func() int { falseCalled = true; return 2 })
	if got != 2 || !falseCalled {
		t.Fatalf("false branch: got=%d called=%v", got, falseCalled)
	}

	falseCalled = false
	got = IfLazyR(true, 1, func() int { falseCalled = true; return 2 })
	if got != 1 {
		t.Fatalf("true branch: got=%d want 1", got)
	}
	if falseCalled {
		t.Fatal("false branch thunk should not be invoked when condition is true")
	}
}

func TestIfLazy(t *testing.T) {
	tCalled, fCalled := 0, 0
	got := IfLazy(true, func() int { tCalled++; return 1 }, func() int { fCalled++; return 2 })
	if got != 1 || tCalled != 1 || fCalled != 0 {
		t.Fatalf("true branch: got=%d tCalled=%d fCalled=%d", got, tCalled, fCalled)
	}

	tCalled, fCalled = 0, 0
	got = IfLazy(false, func() int { tCalled++; return 1 }, func() int { fCalled++; return 2 })
	if got != 2 || tCalled != 0 || fCalled != 1 {
		t.Fatalf("false branch: got=%d tCalled=%d fCalled=%d", got, tCalled, fCalled)
	}
}
