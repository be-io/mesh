/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.annotate.MPI;

import java.lang.annotation.Annotation;
import java.lang.reflect.Method;
import java.lang.reflect.Parameter;
import java.lang.reflect.Type;

/**
 * @author coyzeng@gmail.com
 */
public class MethodInspector implements Inspector {

    private MPI mpi;
    private final Class<?> type;
    private final Method method;

    public MethodInspector(MPI mpi, Class<?> type, Method method) {
        this.mpi = mpi;
        this.type = type;
        this.method = method;
    }

    public MethodInspector(Class<?> type, Method method) {
        this.type = type;
        this.method = method;
    }

    @Override
    public String getSignature() {
        StringBuilder signature = new StringBuilder(method.getName());
        for (Parameter parameter : method.getParameters()) {
            signature.append(parameter.getType().getName());
        }
        return signature.toString();
    }

    @Override
    public Class<?> getType() {
        return this.type;
    }

    @Override
    public String getName() {
        return this.method.getName();
    }

    @SuppressWarnings("unchecked")
    @Override
    public <T extends Annotation> T getAnnotation(Class<T> kind) {
        if (null != mpi && kind == MPI.class) {
            return (T) mpi;
        }
        if (method.isAnnotationPresent(kind)) {
            return method.getAnnotation(kind);
        }
        return type.getAnnotation(kind);
    }

    @Override
    public Class<?> getReturnType() {
        return method.getReturnType();
    }

    @Override
    public Type getReturnGenericType() {
        return method.getGenericReturnType();
    }

    @Override
    public Class<?>[] getExceptionTypes() {
        return method.getExceptionTypes();
    }

    @Override
    public Object invoke(Object obj, Object... args) throws ReflectiveOperationException {
        return method.invoke(obj, args);
    }
}
