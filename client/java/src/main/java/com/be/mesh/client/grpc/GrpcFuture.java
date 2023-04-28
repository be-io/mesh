/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.grpc;

import com.be.mesh.client.cause.MeshException;
import com.be.mesh.client.cause.NotFoundException;
import com.be.mesh.client.mpc.MeshCode;
import com.be.mesh.client.tool.Tool;
import io.grpc.StatusRuntimeException;
import lombok.extern.slf4j.Slf4j;

import java.util.concurrent.ExecutionException;
import java.util.concurrent.Future;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;
import java.util.function.Consumer;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
public class GrpcFuture<T> implements Future<T> {

    private final Future<T> future;
    private final Consumer<Throwable>[] hooks;

    @SafeVarargs
    public GrpcFuture(Future<T> future, Consumer<Throwable>... hooks) {
        this.future = future;
        this.hooks = hooks;
    }

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
        try {
            return this.future.get();
        } catch (ExecutionException e) {
            if (e.getCause() instanceof StatusRuntimeException) {
                RuntimeException ce = toException(((StatusRuntimeException) e.getCause()));
                runHook(ce);
                throw ce;
            }
            throw e;
        }
    }

    @Override
    public T get(long timeout, TimeUnit unit) throws InterruptedException, ExecutionException, TimeoutException {
        try {
            return this.future.get(timeout, unit);
        } catch (ExecutionException e) {
            if (e.getCause() instanceof StatusRuntimeException) {
                throw toException(((StatusRuntimeException) e.getCause()));
            }
            throw e;
        }
    }

    private void runHook(RuntimeException re) {
        if (Tool.optional(this.hooks)) {
            return;
        }
        for (Consumer<Throwable> hook : this.hooks) {
            try {
                hook.accept(re);
            } catch (Throwable e) {
                log.warn("Execute grpc consumer hook with cause. ", e);
            }
        }
    }

    private RuntimeException toException(StatusRuntimeException e) {
        switch (e.getStatus().getCode()) {
            case OK:
                break;
            case ABORTED:
            case CANCELLED:
            case UNKNOWN:
            case INTERNAL:
            case ALREADY_EXISTS:
                return new MeshException(MeshCode.SYSTEM_ERROR, e);
            case INVALID_ARGUMENT:
            case DATA_LOSS:
                return new MeshException(MeshCode.VALIDATE, e);
            case DEADLINE_EXCEEDED:
                return new com.be.mesh.client.cause.TimeoutException(e);
            case UNAVAILABLE:
                return new MeshException(MeshCode.NET_UNAVAILABLE, e);
            case NOT_FOUND:
            case UNIMPLEMENTED:
                return new NotFoundException(e);
            case OUT_OF_RANGE:
            case RESOURCE_EXHAUSTED:
            case UNAUTHENTICATED:
            case PERMISSION_DENIED:
            case FAILED_PRECONDITION:
                return new MeshException(MeshCode.UNAUTHORIZED, e);
        }
        return e;
    }
}
