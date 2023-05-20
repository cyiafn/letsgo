// Package letsgo provides generic useful functions for Go.
package letsgo

import (
	"os"
	"os/signal"
	"syscall"
)

// OnPanicRun executes the given functions f if a panic occurs.
// If processRecovery is enabled, your handeRecovery function will be executed.
// The handleRecovery function should accept a recovery value and execute your custom code such as logging
// the traceStack or anything you might desire.
//
// Sample usecase:
//
// import (
//
//		"fmt"
//		"runtime/debug"
//		"github.com/cyiafn/letsgo"
//	)
//
//	func main() {
//		defer letsgo.OnPanicRun(true, func(r interface{}) {
//			fmt.Println(r)
//			debug.PrintStack()
//		}, closeConnToDB())
//
// This should be ran as a deferred function at the start of your application.
func OnPanicRun(processRecovery bool, handleRecovery func(r any), cleanups ...func()) {
	if r := recover(); r != nil {
		if processRecovery {
			handleRecovery(r)
		}
	}

	if cleanups == nil || len(cleanups) == 0 {
		return
	}

	for _, cleanup := range cleanups {
		cleanup()
	}
}

// CleanupWhenShutdown executes the given functions cleanups when the process is about to shut down.
// It is useful for cleaning up your instance before shutting down such as freeing up connections to the DB, etc.
// It detects shutdowns using os.Signals. Most common shutdown signals are syscall.SIGTERM and syscall.SIGINT.
// This function defaults to SIGINT if no signals are passed in.
// Alternative, you can pass in your custom interceptSignals that you want it to detect.
//
// Sample Usage:
// import (
//
//	"github.com/cyiafn/letsgo"
//	"os"
//	"syscall"
//
// )
//
//	func main() {
//		letsgo.CleanupWhenShutdown([]os.Signal{syscall.SIGINT, syscall.SIGTERM}, closeConnToDB())
//	}
//
// This should be ran at the start of your application.
func CleanupWhenShutdown(interceptSignals []os.Signal, cleanups ...func()) {
	go func() {
		sigs := make(chan os.Signal, 1)

		if interceptSignals == nil || len(interceptSignals) == 0 {
			interceptSignals = []os.Signal{syscall.SIGINT}
		}
		signal.Notify(sigs, interceptSignals...)
		go func() {
			<-sigs
			for _, cleanup := range cleanups {
				cleanup()
			}
		}()
		os.Exit(0)
	}()
}
