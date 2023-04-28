/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import (
	"context"
	"time"
)

var ILocker = (*Locker)(nil)

// Locker spi
// @MPI(macro.MeshMPI)
// @SPI(macro.MeshSPI)
type Locker interface {

	// Lock create write lock.
	// @MPI("mesh.locker.w.lock")
	Lock(ctx context.Context, rid string, timeout time.Duration) (bool, error)

	// Unlock release write lock.
	// @MPI("mesh.locker.w.unlock")
	Unlock(ctx context.Context, rid string) error

	// ReadLock create read lock.
	// @MPI("mesh.locker.r.lock")
	ReadLock(ctx context.Context, rid string, timeout time.Duration) (bool, error)

	// ReadUnlock release read lock.
	// @MPI("mesh.locker.r.unlock")
	ReadUnlock(ctx context.Context, rid string) error
}
