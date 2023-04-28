/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package plugin

import (
	"testing"

	"time"
)

func TestSchedule(t *testing.T) {
	size := 5
	pool := NewWorkerPool(size)
	p := pool.(*workerQueue)
	for i := 0; i < 5; i++ {
		pool.Schedule(func() {
			time.Sleep(100 * time.Millisecond)
		})
	}
	now := time.Now()
	pool.Schedule(func() {})
	if time.Now().Before(now.Add(50 * time.Millisecond)) {
		t.Errorf("Test Schedule() error, should wait for 20 millisecond")
	}
	if len(p.sem) != size {
		t.Errorf("Test Schedule() error, should be 5")
	}
}

func TestScheduleAlways(t *testing.T) {
	size := 5
	pool := NewWorkerPool(size)
	p := pool.(*workerQueue)
	for i := 0; i < 20; i++ {
		pool.ScheduleAlways(func() {
		})
		time.Sleep(10 * time.Millisecond)
	}

	t.Logf("Test ScheduleAlways p.sem is %d", len(p.sem))
	if len(p.sem) == 1 {
		t.Errorf("Test ScheduleAlways() error, should not be 1")
	}

	for i := 0; i < 5; i++ {
		pool.ScheduleAlways(func() {
			time.Sleep(100 * time.Millisecond)
		})
	}
	now := time.Now()
	pool.ScheduleAlways(func() {})
	if time.Now().After(now.Add(50 * time.Millisecond)) {
		t.Errorf("Test Schedule() error, should run it now")
	}
	if len(p.sem) != size {
		t.Errorf("Test Schedule() error, should be 5")
	}
}

func TestScheduleAuto(t *testing.T) {
	size := 5
	pool := NewWorkerPool(size)
	p := pool.(*workerQueue)
	for i := 0; i < 3; i++ {
		pool.ScheduleAuto(func() {
			time.Sleep(time.Millisecond)
		})
		time.Sleep(50 * time.Millisecond)
	}
	if len(p.sem) != 1 {
		t.Errorf("Test ScheduleAuto() error, should be 1, but get %d", len(p.sem))
	}

	for i := 0; i < 3; i++ {
		pool.ScheduleAuto(func() {
			time.Sleep(10 * time.Millisecond)
		})
	}
	time.Sleep(50 * time.Millisecond)
	if len(p.sem) != 3 {
		t.Errorf("Test ScheduleAuto() error, should be 3, but get %d", len(p.sem))
	}

	for i := 0; i < 10; i++ {
		pool.ScheduleAuto(func() {
			time.Sleep(time.Millisecond)
		})
	}
	time.Sleep(10 * time.Millisecond)
	if len(p.sem) != size {
		t.Errorf("Test ScheduleAuto() error, should be %d, but get %d", size, len(p.sem))
	}
}

func TestPanic(t *testing.T) {
	pool := NewWorkerPool(10)
	p := pool.(*workerQueue)
	pool.ScheduleAlways(func() {
		panic("hello")
	})
	time.Sleep(10 * time.Millisecond)
	if len(p.sem) != 0 {
		t.Errorf("Test ScheduleAuto() error, should be 0")
	}
}
