/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package log

import (
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"
)

func init() {
	catchSignals()

	onProcessExit = append(onProcessExit, func() {
		if pidFile != "" {
			Catch(os.Remove(pidFile))
		}
	})
}

var (
	pidFile               string
	onProcessExit         []func()
	shutdownCallbacksOnce sync.Once
	shutdownCallbacks     []func() error
	signalCallback        = make(map[syscall.Signal][]func())
)

func catchSignals() {
	catchSignalsCrossPlatform()
	catchSignalsPosix()
}

func catchSignalsCrossPlatform() {
	go func() {
		defer func() {
			if err := recover(); nil != err {
				Error0("panic %v\n%s", err, string(debug.Stack()))
			}
		}()
		channel := make(chan os.Signal, 1)
		signal.Notify(channel, syscall.SIGTERM, syscall.SIGHUP,
			syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGINT)

		for sig := range channel {
			Debug0("signal %s received!", sig)
			switch sig {
			case syscall.SIGQUIT:
				// quit
				for _, hook := range onProcessExit {
					hook() // only perform important cleanup actions
				}
				os.Exit(0)
			case syscall.SIGTERM:
				// stop to quit
				exitCode := ExecuteShutdownCallbacks("SIGTERM")
				for _, hook := range onProcessExit {
					hook() // only perform important cleanup actions
				}
				//Stop()
				executeSignalCallback(syscall.SIGTERM)
				os.Exit(exitCode)
			case syscall.SIGUSR1:
				// reopen
			case syscall.SIGUSR2:
				// do nothing
			case syscall.SIGHUP:
				executeSignalCallback(syscall.SIGHUP)
			case syscall.SIGINT:
				executeSignalCallback(syscall.SIGINT)
			}
		}
	}()
}

func executeSignalCallback(sig syscall.Signal) {
	if hooks, ok := signalCallback[sig]; ok {
		for _, hook := range hooks {
			hook()
		}
	}
}

func catchSignalsPosix() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				Error0("panic %v\n%s", r, string(debug.Stack()))
			}
		}()
		shutdown := make(chan os.Signal, 1)
		signal.Notify(shutdown, os.Interrupt)

		for index := 0; true; index++ {
			<-shutdown

			if index > 0 {
				for _, hook := range onProcessExit {
					hook() // important cleanup actions only
				}
				os.Exit(2)
			}

			// important cleanup actions before shutdown callbacks
			for _, hook := range onProcessExit {
				hook()
			}

			go func() {
				defer func() {
					if err := recover(); err != nil {
						Error0("panic %v\n%s", err, string(debug.Stack()))
					}
				}()
				os.Exit(ExecuteShutdownCallbacks("SIGINT"))
			}()
		}
	}()
}

func ExecuteShutdownCallbacks(signalName string) (exitCode int) {
	shutdownCallbacksOnce.Do(func() {
		var errs []error

		for _, cb := range shutdownCallbacks {
			// If the callback is performing normally,
			// err does not need to be saved to prevent
			// the exit code from being non-zero
			if err := cb(); nil != err {
				errs = append(errs, err)
			}
		}

		if len(errs) > 0 {
			for _, err := range errs {
				Error0(" %s shutdown: %v", signalName, err)
			}
			exitCode = 4
		}
	})

	return
}

func AddProcessExitHook(hook func()) {
	onProcessExit = append(onProcessExit, hook)
}

func AddProcessShutdownHook(hook func() error) {
	shutdownCallbacks = append(shutdownCallbacks, hook)
}

// AddProcessShutDownFirstHook insert the callback func into the header
func AddProcessShutDownFirstHook(hook func() error) {
	var firstCallbacks []func() error
	firstCallbacks = append(firstCallbacks, hook)
	firstCallbacks = append(firstCallbacks, shutdownCallbacks...)
	// replace current firstCallbacks
	shutdownCallbacks = firstCallbacks
}

func AddSignalCallbackHook(cb func(), signals ...syscall.Signal) {
	for _, sig := range signals {
		signalCallback[sig] = append(signalCallback[sig], cb)
	}
}
