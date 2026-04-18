package letsgo

import (
	"sync/atomic"
	"testing"
)

func TestOnPanicRun_NoPanicRunsCleanups(t *testing.T) {
	var count int32
	func() {
		defer OnPanicRun(true, func(r any) { t.Error("handleRecovery should not be called") },
			func() { atomic.AddInt32(&count, 1) },
			func() { atomic.AddInt32(&count, 1) },
		)
	}()
	if got := atomic.LoadInt32(&count); got != 2 {
		t.Fatalf("cleanups ran %d times, want 2", got)
	}
}

func TestOnPanicRun_PanicWithRecovery(t *testing.T) {
	var recovered any
	var count int32
	func() {
		defer OnPanicRun(true, func(r any) { recovered = r },
			func() { atomic.AddInt32(&count, 1) },
		)
		panic("boom")
	}()
	if recovered != "boom" {
		t.Fatalf("recovered = %v, want \"boom\"", recovered)
	}
	if atomic.LoadInt32(&count) != 1 {
		t.Fatalf("cleanup didn't run after panic")
	}
}

func TestOnPanicRun_PanicSuppressedWhenProcessRecoveryFalse(t *testing.T) {
	var handleCalled bool
	var count int32
	func() {
		defer OnPanicRun(false, func(r any) { handleCalled = true },
			func() { atomic.AddInt32(&count, 1) },
		)
		panic("x")
	}()
	if handleCalled {
		t.Fatal("handleRecovery should not be called when processRecovery=false")
	}
	if atomic.LoadInt32(&count) != 1 {
		t.Fatal("cleanups should still run after a recovered panic")
	}
}

func TestOnPanicRun_NoCleanupsProvided(t *testing.T) {
	func() {
		defer OnPanicRun(true, func(r any) {})
	}()
}

func TestCleanupWhenShutdown_DoesNotBlockOrRunImmediately(t *testing.T) {
	var ran int32
	CleanupWhenShutdown(nil, func() { atomic.AddInt32(&ran, 1) })
	if atomic.LoadInt32(&ran) != 0 {
		t.Fatal("cleanup ran before any signal was delivered")
	}
}

func TestCleanupWhenShutdown_NoCleanupsDoesNotPanic(t *testing.T) {
	CleanupWhenShutdown(nil)
}
