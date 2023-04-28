/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.boost;

import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.prsim.Locker;

import java.time.Duration;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReadWriteLock;
import java.util.concurrent.locks.ReentrantReadWriteLock;

/**
 * @author coyzeng@gmail.com
 */
@SPI(value = "hazelcast", prototype = true)
public class HazelcastDistributionLocker implements Locker {

    @Override
    public boolean lock(String rid, Duration timeout) {
        return HazelcastReference.INST.get().getCPSubsystem().getLock(rid).tryLock(timeout.getNano(), TimeUnit.NANOSECONDS);
    }

    @Override
    public void unlock(String rid) {
        HazelcastReference.INST.get().getCPSubsystem().getLock(rid).unlock();
    }

    @Override
    public boolean readLock(String rid, Duration timeout) {
        return HazelcastReference.INST.get().getCPSubsystem().getLock(rid).tryLock(timeout.getNano(), TimeUnit.NANOSECONDS);
    }

    @Override
    public void readUnlock(String rid) {
        HazelcastReference.INST.get().getCPSubsystem().getLock(rid).unlock();
    }

    @Override
    public ReadWriteLock getReadWriteLock(String rid) {
        return new ReentrantReadWriteLock();
    }

    @Override
    public Lock getLock(String rid) {
        return HazelcastReference.INST.get().getCPSubsystem().getLock(rid);
    }

}
