/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import (
	"context"
	"github.com/opendatav/mesh/client/golang/macro"
	"github.com/opendatav/mesh/client/golang/types"
	"time"
)

const SCSEQ = "PRSIM_SC_TASK"

var SystemClock = &macro.Btt{Topic: "mesh.schedule.system.clock", Code: "*"}

var IScheduler = (*Scheduler)(nil)

// Scheduler
// @MPI(macro.MeshMPI)
// @SPI(macro.MeshSPI)
// Schedules {@link Timeout}s for one-time future execution in a background thread.
type Scheduler interface {

	// Timeout
	// Schedules the specified {@link Timeout} for one-time execution after the specified delay.
	// @MPI("mesh.schedule.timeout")
	Timeout(ctx context.Context, timeout *types.Timeout, duration time.Duration) (string, error)

	// Cron
	// Schedules with the cron expression. "0 * * 1-3 * ? *"
	// @MPI("mesh.schedule.cron")
	Cron(ctx context.Context, cron string, topic *types.Topic) (string, error)

	// Period
	// schedule with fixed duration.
	// @MPI("mesh.schedule.period")
	Period(ctx context.Context, duration time.Duration, topic *types.Topic) (string, error)

	// Dump max is 1000 items
	// Releases all resources acquired by this {@link Scheduler} and cancels all
	// tasks which were scheduled but not executed yet.
	// @MPI("mesh.schedule.dump")
	Dump(ctx context.Context) ([]string, error)

	// Cancel
	// Attempts to cancel the {@link com.be.mesh.client.struct.Timeout} associated with this handle.
	// If the task has been executed or cancelled already, it will return with
	// no side effect.
	// @MPI("mesh.schedule.cancel")
	Cancel(ctx context.Context, taskId string) (bool, error)

	// Stop
	// Attempts to stop the {@link com.be.mesh.client.struct.Timeout} associated with this handle.
	// If the task has been executed or cancelled already, it will return with
	// no side effect.
	// @MPI("mesh.schedule.stop")
	Stop(ctx context.Context, taskId string) (bool, error)

	// Emit the scheduler topic
	// @MPI("mesh.schedule.emit")
	Emit(ctx context.Context, topic *types.Topic) error

	// Shutdown the scheduler
	// @MPI("mesh.schedule.shutdown")
	Shutdown(ctx context.Context, duration time.Duration) error
}
