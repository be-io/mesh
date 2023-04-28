/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package plugin

import (
	"github.com/be-io/mesh/client/golang/log"
	"runtime/debug"
)

// WorkerQueue provides a pool for goroutines
type WorkerQueue interface {

	// Schedule try to acquire pooled worker goroutine to execute the specified task,
	// this method would block if no worker goroutine is available
	Schedule(task func())

	// ScheduleAlways try to acquire pooled worker goroutine to execute the specified task first,
	// but would not block if no worker goroutine is available. A temp goroutine will be created for task execution.
	ScheduleAlways(task func())

	// ScheduleAuto auto
	ScheduleAuto(task func())
}

type workerQueue struct {
	work chan func()
	sem  chan struct{}
}

// NewWorkerPool create a worker pool
func NewWorkerPool(size int) WorkerQueue {
	return &workerQueue{
		work: make(chan func()),
		sem:  make(chan struct{}, size),
	}
}

func (p *workerQueue) Schedule(task func()) {
	select {
	case p.work <- task:
	case p.sem <- struct{}{}:
		go p.spawnWorker(task)
	}
}

func (p *workerQueue) ScheduleAlways(task func()) {
	select {
	case p.work <- task:
	case p.sem <- struct{}{}:
		go p.spawnWorker(task)
	default:
		// new temp goroutine for task execution
		log.Debug0("[syncpool] workerpool new goroutine")
		go task()
	}
}

func (p *workerQueue) ScheduleAuto(task func()) {
	select {
	case p.work <- task:
		return
	default:
	}
	select {
	case p.work <- task:
	case p.sem <- struct{}{}:
		go p.spawnWorker(task)
	default:
		// new temp goroutine for task execution
		log.Debug0("[syncpool] workerpool new goroutine")
		go task()
	}
}

func (p *workerQueue) spawnWorker(task func()) {
	defer func() {
		if r := recover(); r != nil {
			log.Warn0("syncpool", "[syncpool] panic %v\n%s", p, string(debug.Stack()))
		}
		<-p.sem
	}()
	for {
		task()
		task = <-p.work
	}
}
