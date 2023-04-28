/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.system;

import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.mpc.ServiceProxy;
import com.be.mesh.client.prsim.Queue;
import com.be.mesh.client.struct.Entity;

import java.time.Duration;

/**
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public class MeshQueue implements Queue<Object> {

    private final Queue<?> queue = ServiceProxy.proxy(Queue.class);

    @Override
    public Entity peek(String channel) {
        return queue.peek(channel);
    }

    @Override
    public Entity pop(String channel, Duration timeout) {
        return queue.pop(channel, timeout);
    }

    @Override
    public void push(String channel, Entity element) {
        queue.push(channel, element);
    }
}
