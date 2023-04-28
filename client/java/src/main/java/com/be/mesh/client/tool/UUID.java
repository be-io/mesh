/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.tool;

import com.be.mesh.client.cause.CompatibleException;
import lombok.extern.slf4j.Slf4j;

import java.util.Random;

/**
 * UUIDUtils.
 *
 * @author coyzeng@gmail.com
 */
@Slf4j
public final class UUID {

    private static final Random RANDOM = new Random();

    private static final long WORKER_ID_BITS = 5L;

    private static final long DATACENTERIDBITS = 5L;

    private static final long MAX_WORKER_ID = ~(-1L << WORKER_ID_BITS);

    private static final long MAX_DATACENTER_ID = ~(-1L << DATACENTERIDBITS);

    private static final long SEQUENCE_BITS = 12L;

    private static final long WORKER_ID_SHIFT = SEQUENCE_BITS;

    private static final long DATACENTER_ID_SHIFT = SEQUENCE_BITS + WORKER_ID_BITS;

    private static final long TIMESTAMP_LEFT_SHIFT = SEQUENCE_BITS + WORKER_ID_BITS + DATACENTERIDBITS;

    private static final long SEQUENCE_MASK = ~(-1L << SEQUENCE_BITS);

    private static final UUID ID_WORKER_UTILS = new UUID();

    private final long workerId;

    private final long datacenterId;

    private final long idepoch;

    private long sequence = '0';

    private long lastTimestamp = -1L;

    private UUID() {
        this(RANDOM.nextInt((int) MAX_WORKER_ID), RANDOM.nextInt((int) MAX_DATACENTER_ID), 1288834974657L);
    }

    private UUID(final long workerId, final long datacenterId, final long idepoch) {
        if (workerId > MAX_WORKER_ID || workerId < 0) {
            throw new CompatibleException(String.format("worker Id can't be greater than %d or less than 0", MAX_WORKER_ID));
        }
        if (datacenterId > MAX_DATACENTER_ID || datacenterId < 0) {
            throw new CompatibleException(String.format("datacenter Id can't be greater than %d or less than 0", MAX_DATACENTER_ID));
        }
        this.workerId = workerId;
        this.datacenterId = datacenterId;
        this.idepoch = idepoch;
    }

    /**
     * Gets instance.
     *
     * @return the instance
     */
    public static UUID getInstance() {
        return ID_WORKER_UTILS;
    }

    private synchronized long nextId() {
        long timestamp = timeGen();
        if (timestamp < lastTimestamp) {
            log.error(String.format("Clock moved backwards.  Refusing to generate id for %d milliseconds", lastTimestamp - timestamp));
        }
        if (lastTimestamp == timestamp) {
            sequence = (sequence + 1) & SEQUENCE_MASK;
            if (sequence == 0) {
                timestamp = tilNextMillis(lastTimestamp);
            }
        } else {
            sequence = 0L;
        }

        lastTimestamp = timestamp;

        return ((timestamp - idepoch) << TIMESTAMP_LEFT_SHIFT)
                | (datacenterId << DATACENTER_ID_SHIFT)
                | (workerId << WORKER_ID_SHIFT) | sequence;
    }

    private long tilNextMillis(final long lastTimestamp) {
        long timestamp = timeGen();
        while (timestamp <= lastTimestamp) {
            timestamp = timeGen();
        }
        return timestamp;
    }

    private long timeGen() {
        return System.currentTimeMillis();
    }

    /**
     * generate short uuid.
     *
     * @return short uuid.
     */
    public String shortUUID() {
        return String.valueOf(ID_WORKER_UTILS.nextId());
    }
}
