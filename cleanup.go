// Package letsgo provides generic useful functions for Go.
package letsgo

import (
	"os"
	"os/signal"
	"syscall"
)

// OnPanicRun is a panic-recovery helper intended to be installed with defer.
// If the surrounding goroutine panics, it recovers the panic, optionally
// hands the recovered value to handleRecovery, and then invokes every
// function in cleanups in order.
//
// Semantics:
//   - Not pure: calls recover() and therefore MUST be used via `defer` in
//     the same goroutine whose panic you want to catch. Using it in a
//     non-deferred call is a no-op for panic handling.
//   - Always swallows the panic. Even when processRecovery is false, the
//     panic does not propagate — it is silently discarded and cleanups
//     still run. If you need to re-panic, do it yourself inside
//     handleRecovery.
//   - handleRecovery is only invoked when processRecovery is true AND a
//     panic occurred. Passing a nil handleRecovery with processRecovery
//     set to true will itself panic.
//   - cleanups always run (on the panic path AND on a normal return),
//     in the order they were passed. They run sequentially in the calling
//     goroutine.
//   - Does not call os.Exit; control returns to the caller's deferred chain.
//
// Example:
//
//	import (
//	    "fmt"
//	    "runtime/debug"
//	    "github.com/cyiafn/letsgo"
//	)
//
//	func main() {
//	    defer letsgo.OnPanicRun(true, func(r any) {
//	        fmt.Println("panic:", r)
//	        debug.PrintStack()
//	    }, closeDB, flushLogs)
//
//	    doWork() // if this panics: prints the panic, runs closeDB, then flushLogs.
//	}
func OnPanicRun(processRecovery bool, handleRecovery func(r any), cleanups ...func()) {
	if r := recover(); r != nil {
		if processRecovery {
			handleRecovery(r)
		}
	}

	if len(cleanups) == 0 {
		return
	}

	for _, cleanup := range cleanups {
		cleanup()
	}
}

// CleanupWhenShutdown installs an asynchronous OS-signal handler that runs
// the provided cleanup functions exactly once when the process receives a
// shutdown signal, then calls os.Exit(0).
//
// Semantics:
//   - Non-blocking: returns immediately. The signal handler runs in its own
//     goroutine.
//   - Installs a process-wide signal handler via signal.Notify. Calling it
//     multiple times will register additional (redundant) handlers — it is
//     intended to be called once at startup.
//   - If interceptSignals is nil or empty, it defaults to SIGINT and SIGTERM.
//     Pass a custom slice to listen for other signals (e.g. SIGHUP for reload).
//   - cleanups run sequentially, in the order they were passed, from within
//     the handler goroutine. If any cleanup panics, subsequent cleanups will
//     NOT run and the process will terminate via the default panic path
//     instead of os.Exit(0). Recover inside each cleanup if you need
//     best-effort execution of all of them.
//   - Terminates via os.Exit(0), bypassing any further deferred functions.
//
// Example:
//
//	import (
//	    "os"
//	    "syscall"
//	    "github.com/cyiafn/letsgo"
//	)
//
//	func main() {
//	    letsgo.CleanupWhenShutdown(
//	        []os.Signal{syscall.SIGINT, syscall.SIGTERM},
//	        closeDB,
//	        flushMetrics,
//	    )
//	    runServer() // on Ctrl-C: closeDB -> flushMetrics -> os.Exit(0).
//	}
func CleanupWhenShutdown(interceptSignals []os.Signal, cleanups ...func()) {
	go func() {
		sigs := make(chan os.Signal, 1)

		if len(interceptSignals) == 0 {
			interceptSignals = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
		}
		signal.Notify(sigs, interceptSignals...)
		go func() {
			<-sigs
			for _, cleanup := range cleanups {
				cleanup()
			}
			os.Exit(0)
		}()
	}()
}
