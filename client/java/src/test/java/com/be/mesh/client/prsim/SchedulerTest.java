/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.prsim;

import com.be.mesh.client.mpc.ServiceLoader;
import com.be.mesh.client.struct.Topic;
import com.be.mesh.client.struct.Timeout;
import com.be.mesh.client.tool.UUID;
import org.testng.Assert;
import org.testng.annotations.Test;

import java.time.Duration;

/**
 * @author coyzeng@gmail.com
 */
public class SchedulerTest {

    @Test
    public void timeoutTest() throws Exception {
        Topic tuple = new Topic();
        tuple.setTopic("net.mesh.registry");
        tuple.setCode("*");
        Timeout timeout = new Timeout();
        timeout.setTaskId(UUID.getInstance().shortUUID());
        timeout.setBinding(tuple);
        Scheduler scheduler = ServiceLoader.load(Scheduler.class).getDefault();
        String taskId = scheduler.timeout(timeout, Duration.ofSeconds(1));
        Assert.assertEquals(1, 1);
        scheduler.shutdown(Duration.ofSeconds(3));
    }
}
