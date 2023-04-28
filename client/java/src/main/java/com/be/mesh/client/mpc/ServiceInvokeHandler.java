/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.struct.Service;
import com.be.mesh.client.tool.Tool;
import lombok.AllArgsConstructor;

import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.UndeclaredThrowableException;
import java.util.Optional;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.Future;
import java.util.concurrent.TimeUnit;

/**
 * @author coyzeng@gmail.com
 */
public class ServiceInvokeHandler implements Invoker<Object> {

    private final Invoker<?> invoker;

    public ServiceInvokeHandler(Object service) {
        this.invoker = Filter.composite(new ServiceInvoker(service), Filter.PROVIDER);
    }

    @Override
    public Object invoke(Invocation invocation) throws Throwable {
        try {
            return this.invoker.invoke(invocation);
        } catch (InvocationTargetException | UndeclaredThrowableException e) {
            throw Tool.destructor(e.getCause(), invocation.getInspector().getExceptionTypes());
        } catch (Throwable e) {
            throw Tool.destructor(e, invocation.getInspector().getExceptionTypes());
        }
    }

    @AllArgsConstructor
    private static final class ServiceInvoker implements Invoker<Object> {

        private final Object service;

        @Override
        public Object invoke(Invocation invocation) throws Throwable {
            Object result = invocation.getInspector().invoke(this.service, invocation.getParameters().arguments());
            if (!(result instanceof Future)) {
                return result;
            }
            try {
                Eden eden = ServiceLoader.load(Eden.class).getDefault();
                Execution<Service> execution = eden.infer(Mesh.context().getUrn());
                long timeout = Optional.ofNullable(execution).map(Execution::schema).map(Service::getTimeout).orElse(3000L);
                return ((Future<?>) result).get(Math.max(timeout, 3000L), TimeUnit.MILLISECONDS);
            } catch (ExecutionException e) {
                throw e.getCause();
            }
        }
    }

}
