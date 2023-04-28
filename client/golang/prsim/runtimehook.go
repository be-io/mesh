/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import "context"

var IRuntimeHook = (*RuntimeHook)(nil)
var IRuntimeAware = (*RuntimeAware)(nil)

type Runtime interface {

	// Submit the daemon to goroutine pool.
	Submit(routine func())
}

type RuntimeAware interface {
	// Init life cycle.
	Init() error
}

// RuntimeHook spi
// @SPI("mesh")
type RuntimeHook interface {

	// Start Trigger when mesh runtime is start.
	Start(ctx context.Context, runtime Runtime) error

	// Stop Trigger when mesh runtime is stop.
	Stop(ctx context.Context, runtime Runtime) error

	// Refresh Trigger then mesh runtime context is refresh or metadata is refresh.
	Refresh(ctx context.Context, runtime Runtime) error
}

type Waiter interface {
	Wait(ctx context.Context, runtime Runtime)
}

type MeshRuntime struct {
}

func (that *MeshRuntime) Submit(routine func()) {
	go routine()
}
