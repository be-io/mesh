/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.annotate.SPI;
import lombok.AllArgsConstructor;

import java.io.InputStream;
import java.nio.ByteBuffer;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.Future;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;

import static com.be.mesh.client.mpc.Transporter.CONSUMER;

/**
 * @author coyzeng@gmail.com
 */
@SPI(CONSUMER)
public class ConsumerTransporter implements Transporter {

    @Override
    public ByteBuffer transport(String urn, ByteBuffer buffer) throws Throwable {
        return null;
    }

    @AllArgsConstructor
    public static class ConsumerFuture<T> implements Future<T> {

        private final Future<InputStream> future;

        @Override
        public boolean cancel(boolean mayInterruptIfRunning) {
            return this.future.cancel(mayInterruptIfRunning);
        }

        @Override
        public boolean isCancelled() {
            return this.future.isCancelled();
        }

        @Override
        public boolean isDone() {
            return this.future.isDone();
        }

        @Override
        public T get() throws InterruptedException, ExecutionException {
            return (T) this.future.get();
        }

        @Override
        public T get(long timeout, TimeUnit unit) throws InterruptedException, ExecutionException, TimeoutException {
            return (T) this.future.get(timeout, unit);
        }
    }
}
