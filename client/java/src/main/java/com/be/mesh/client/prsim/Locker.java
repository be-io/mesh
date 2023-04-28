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

import java.time.Duration;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReadWriteLock;
import java.util.concurrent.locks.ReentrantLock;
import java.util.concurrent.locks.ReentrantReadWriteLock;

/**
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public interface Locker {

    /**
     * Create write lock.
     *
     * @param rid     resource id
     * @param timeout lock timeout
     * @return Read write lock
     */
    @MPI("mesh.locker.w.lock")
    boolean lock(@Index(0) String rid, @Index(1) Duration timeout);

    /**
     * Release write lock.
     *
     * @param rid resource id
     */
    @MPI("mesh.locker.w.unlock")
    void unlock(@Index(0) String rid);

    /**
     * Create read lock.
     *
     * @param rid     resource id
     * @param timeout lock timeout
     * @return lock
     */
    @MPI("mesh.locker.r.lock")
    boolean readLock(@Index(0) String rid, @Index(1) Duration timeout);

    /**
     * Release read lock.
     *
     * @param rid resource id
     */
    @MPI("mesh.locker.r.unlock")
    void readUnlock(@Index(0) String rid);

    default ReadWriteLock getReadWriteLock(@Index(0) String rid) {
        return new ReentrantReadWriteLock();
    }

    default Lock getLock(@Index(0) String rid) {
        return new ReentrantLock();
    }

}
