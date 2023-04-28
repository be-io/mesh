/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.system;

import io.be.mesh.macro.SPI;
import io.be.mesh.mpc.ServiceProxy;
import io.be.mesh.prsim.Queue;
import io.be.mesh.struct.Entity;

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
