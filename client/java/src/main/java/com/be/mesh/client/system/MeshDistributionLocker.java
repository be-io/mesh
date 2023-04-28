/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.system;

import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.mpc.ServiceProxy;
import com.be.mesh.client.prsim.Locker;

import java.time.Duration;

/**
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public class MeshDistributionLocker implements Locker {

    private final Locker locker = ServiceProxy.proxy(Locker.class);

    @Override
    public boolean lock(String rid, Duration timeout) {
        return locker.lock(rid, timeout);
    }

    @Override
    public void unlock(String rid) {
        locker.unlock(rid);
    }

    @Override
    public boolean readLock(String rid, Duration timeout) {
        return locker.readLock(rid, timeout);
    }

    @Override
    public void readUnlock(String rid) {
        locker.readUnlock(rid);
    }

}
