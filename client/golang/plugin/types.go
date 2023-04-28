/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package plugin

import "context"

type Factory interface {

	// Ptt is the plugin's attributes
	Ptt() *Ptt
}

type Plugin interface {

	// Start will be thread safe invoked after install. Panic if init error.
	Start(ctx context.Context, runtime Runtime)

	// Stop will thread safe release the plugin.
	Stop(ctx context.Context, runtime Runtime)
}

type Container interface {

	// Start the parent.
	Start(ctx context.Context, kvs ...string)

	// Stop the parent.
	Stop(ctx context.Context)

	// Wait the system signal to callback.
	Wait(ctx context.Context)

	// WaitAny will wait any lock hold, it must invoke before Start.
	WaitAny(ctx context.Context, locker Locker)

	// Plugins load the plugins in container.
	Plugins(plugins ...Name)

	// Hook for the parent.
	Hook
}

type Runtime interface {

	// Parse the environment profile as struct.
	Parse(ptr interface{}) error

	// Parse0 the conf with runtime parser.
	Parse0(ptr interface{}, conf string) error

	// Load the plugin in parent.
	Load(plugin Name)

	// Submit the daemon to goroutine pool.
	Submit(routine func())

	// Hook will be called around plugin cycle.
	Hook
}

type Hook interface {

	// StartHook add hook function when the daemon process is started.
	StartHook(hook func())

	// ShutdownHook add hook function when the daemon process is shutdown.
	ShutdownHook(hook func())
}

type Locker interface {

	// TryRLock is the preferred function for taking a shared file lock. This
	// function takes an RW-mutex lock before it tries to lock the file, so there is
	// the possibility that this function may block for a short time if another
	// goroutine is trying to take any action.
	//
	// The actual file lock is non-blocking. If we are unable to get the shared file
	// lock, the function will return false instead of waiting for the lock. If we
	// get the lock, we also set the *Flock instance as being share-locked.
	TryRLock() (bool, error)

	// TryLock is the preferred function for taking an exclusive file lock. This
	// function takes an RW-mutex lock before it tries to lock the file, so there is
	// the possibility that this function may block for a short time if another
	// goroutine is trying to take any action.
	//
	// The actual file lock is non-blocking. If we are unable to get the exclusive
	// file lock, the function will return false instead of waiting for the lock. If
	// we get the lock, we also set the *Flock instance as being exclusive-locked.
	TryLock() (bool, error)

	// Unlock is a function to unlock the file. This file takes a RW-mutex lock, so
	// while it is running the Locked() and RLocked() functions will be blocked.
	//
	// This function short-circuits if we are unlocked already. If not, it calls
	// syscall.LOCK_UN on the file and closes the file descriptor. It does not
	// remove the file from disk. It's up to your application to do.
	//
	// Please note, if your shared lock became an exclusive lock this may
	// unintentionally drop the exclusive lock if called by the consumer that
	// believes they have a shared lock. Please see Lock() for more details.
	Unlock() error
}

type Ptt struct {
	// Name is the plugin name.
	Name Name

	// Flags must return the plugin bootstrap flags struct.
	Flags interface{}

	// Create a plugin instance.
	Create func() Plugin

	// Plugin run priority
	Priority int

	// WaitAny will wait the lock release.
	WaitAny bool
}

type Informer func(value string) (interface{}, error)
