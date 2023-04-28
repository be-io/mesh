/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.system;

import io.be.mesh.macro.SPI;
import io.be.mesh.mpc.ServiceProxy;
import io.be.mesh.prsim.Scheduler;
import io.be.mesh.struct.Timeout;
import io.be.mesh.struct.Topic;

import java.time.Duration;
import java.util.Set;

/**
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public class MeshScheduler implements Scheduler {

    private final Scheduler scheduler = ServiceProxy.proxy(Scheduler.class);

    @Override
    public String timeout(Timeout timeout, Duration duration) {
        return scheduler.timeout(timeout, duration);
    }

    @Override
    public String cron(String cron, Topic binding) {
        return scheduler.cron(cron, binding);
    }

    @Override
    public String period(Duration duration, Topic binding) {
        return scheduler.period(duration, binding);
    }

    @Override
    public Set<String> dump() {
        return scheduler.dump();
    }

    @Override
    public boolean cancel(String taskId) {
        return scheduler.cancel(taskId);
    }

    @Override
    public boolean stop(String taskId) {
        return scheduler.stop(taskId);
    }

    @Override
    public void emit(Topic topic) {
        scheduler.emit(topic);
    }

    @Override
    public boolean isShutdown() {
        return scheduler.isShutdown();
    }

    @Override
    public void shutdown(Duration duration) throws InterruptedException {
        //
    }
}
