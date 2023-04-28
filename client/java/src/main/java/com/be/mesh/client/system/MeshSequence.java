/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.system;

import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.cause.MeshException;
import com.be.mesh.client.mpc.ServiceProxy;
import com.be.mesh.client.prsim.Sequence;
import com.be.mesh.client.tool.Tool;
import lombok.extern.slf4j.Slf4j;

import java.util.Map;
import java.util.concurrent.BlockingDeque;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.LinkedBlockingDeque;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.locks.ReadWriteLock;
import java.util.concurrent.locks.ReentrantReadWriteLock;

/**
 * 提供近端池化能力.
 *
 * @author coyzeng@gmail.com
 */
@Slf4j
@SPI("mesh")
public class MeshSequence implements Sequence {

    private final Map<String, SequenceQueue> sections = new ConcurrentHashMap<>();
    private static final Sequence sequence = ServiceProxy.proxy(Sequence.class);

    @Override
    public String next(String kind, int length) {
        try {
            SequenceQueue queue = sections.computeIfAbsent(kind, key -> new SequenceQueue());
            return queue.getNext(6, kind, length);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new MeshException(e);
        }
    }

    @Override
    public String[] section(String kind, int size, int length) {
        return sequence.section(kind, size, length);
    }

    public static final class SequenceQueue {

        private final BlockingDeque<String> queue = new LinkedBlockingDeque<>();
        private final ReadWriteLock lock = new ReentrantReadWriteLock();

        public String getNext(int retry, String kind, int length) throws InterruptedException {
            for (int index = 0; index < retry; index++) {
                if (!lock.readLock().tryLock(500, TimeUnit.MILLISECONDS)) {
                    log.info("Get {} sequence timeout, retry {} times", kind, index + 1);
                    continue;
                }
                try {
                    String value = queue.poll();
                    if (Tool.required(value)) {
                        return value;
                    }
                } finally {
                    lock.readLock().unlock();
                }
                if (!lock.writeLock().tryLock(500, TimeUnit.MILLISECONDS)) {
                    log.info("Set {} sequence timeout, retry {} times", kind, index + 1);
                    continue;
                }
                try {
                    String[] values = sequence.section(kind, 10, length);
                    if (Tool.optional(values)) {
                        throw new MeshException("Get sequence empty sections for %s", kind);
                    }
                    for (int offset = 1; offset < values.length; offset++) {
                        queue.push(values[offset]);
                    }
                    return values[0];
                } finally {
                    lock.writeLock().unlock();
                }
            }
            throw new MeshException("Get sequence timeout for %s", kind);
        }

    }
}
