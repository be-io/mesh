/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.prsim;

import com.be.mesh.client.annotate.Index;
import com.be.mesh.client.annotate.MPI;
import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.struct.Timeout;
import com.be.mesh.client.struct.Topic;

import java.time.Duration;
import java.util.Set;
import java.util.concurrent.RejectedExecutionException;

/**
 * Schedules {@link Timeout}s for one-time future execution in a background
 * thread.
 *
 * @author coyzeng@gmail.com
 */
@SPI("jvm")
public interface Scheduler {

    /**
     * Schedules the specified {@link Timeout} for one-time execution after
     * the specified delay.
     *
     * @return a handle which is associated with the specified task
     * @throws IllegalStateException      if this timer has been {@linkplain #stop(String) stopped} already
     * @throws RejectedExecutionException if the pending timeouts are too many and creating new timeout
     *                                    can cause instability in the system.
     */
    @MPI("mesh.schedule.timeout")
    String timeout(@Index(0) Timeout timeout, @Index(1) Duration duration);

    /**
     * Schedules with the cron expression. "0 * * 1-3 * ? *"
     *
     * @param cron    cron expression
     * @param binding listener bindings
     * @return task id
     */
    @MPI("mesh.schedule.cron")
    String cron(@Index(0) String cron, @Index(1) Topic binding);

    /**
     * Period schedule with fixed duration.
     *
     * @param duration Delay duration.
     * @param binding  Emit topic.
     * @return task id
     */
    @MPI("mesh.schedule.period")
    String period(@Index(0) Duration duration, @Index(1) Topic binding);

    /**
     * Releases all resources acquired by this {@link Scheduler} and cancels all
     * tasks which were scheduled but not executed yet.
     *
     * @return the handles associated with the tasks which were canceled by
     * this method
     */
    @MPI("mesh.schedule.dump")
    Set<String> dump();

    /**
     * Attempts to cancel the {@link com.be.mesh.client.struct.Timeout} associated with this handle.
     * If the task has been executed or cancelled already, it will return with
     * no side effect.
     *
     * @return True if the cancellation completed successfully, otherwise false
     */
    @MPI("mesh.schedule.cancel")
    boolean cancel(@Index(value = 0, name = "task_id") String taskId);

    /**
     * Attempts to cancel the {@link com.be.mesh.client.struct.Timeout} associated with this handle.
     * If the task has been executed or cancelled already, it will return with
     * no side effect.
     *
     * @return True if the cancellation completed successfully, otherwise false
     */
    @MPI("mesh.schedule.stop")
    boolean stop(@Index(value = 0, name = "task_id") String taskId);

    /**
     * Emit the scheduler topic
     */
    @MPI("mesh.schedule.emit")
    void emit(Topic topic);

    /**
     * the timer is shutdown.
     *
     * @return true for stop
     */
    boolean isShutdown();

    /**
     * Shutdown with terminal await time.
     *
     * @param duration Await await
     */
    void shutdown(Duration duration) throws InterruptedException;
}