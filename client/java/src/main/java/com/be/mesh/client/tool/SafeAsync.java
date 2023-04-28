/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.tool;

import com.be.mesh.client.cause.MeshException;
import com.be.mesh.client.mpc.Mesh;
import com.be.mesh.client.prsim.Context;
import lombok.extern.slf4j.Slf4j;

import java.util.concurrent.Callable;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
public class SafeAsync<T> implements Runnable, Callable<T> {

    private static final String MSG = "Safe async with error";
    private final Context context;
    private Runnable runnable;
    private Callable<T> callable;

    public SafeAsync(Runnable runnable) {
        this.context = Mesh.context().resume();
        this.runnable = runnable;
    }

    public SafeAsync(Callable<T> callable) {
        this.context = Mesh.context().resume();
        this.callable = callable;
    }

    @Override
    public void run() {
        try {
            this.safeExec();
        } catch (Throwable e) {
            log.error(MSG, e);
        }
    }

    @Override
    public T call() throws Exception {
        try {
            return this.safeExec();
        } catch (Exception e) {
            log.error(MSG, e);
            throw e;
        } catch (Throwable e) {
            log.error(MSG, e);
            throw new MeshException(e);
        }
    }

    private T safeExec() throws Throwable {
        boolean contextDisable = Mesh.isEmpty();
        try {
            if (contextDisable) {
                Mesh.reset(context);
            }
            if (null != this.runnable) {
                this.runnable.run();
            }
            if (null != this.callable) {
                return this.callable.call();
            }
            return null;
        } finally {
            if (contextDisable) {
                Mesh.release();
            }
        }
    }

}
