/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.annotate.MPI;

import java.lang.reflect.Method;
import java.lang.reflect.Proxy;
import java.util.*;

/**
 * @author coyzeng@gmail.com
 */
public final class ServiceProxy {

    private ServiceProxy() {
    }

    @SuppressWarnings("unchecked")
    public static <T> T proxy(MPI mpi, Class<T> type) {
        List<Class<?>> types = new ArrayList<>(type.getInterfaces().length + 2);
        types.add(Meshable.class);
        types.addAll(Arrays.asList(type.getInterfaces()));
        if (type.isInterface()) {
            types.add(type);
        }
        Map<Method, Inspector> inspectors = new HashMap<>();
        for (Class<?> kind : types) {
            for (Method method : kind.getDeclaredMethods()) {
                inspectors.put(method, new MethodInspector(mpi, type, method));
            }
        }
        return (T) Proxy.newProxyInstance(type.getClassLoader(), types.toArray(new Class<?>[0]), new ReferenceInvokeHandler(inspectors));
    }

    public static <T> T proxy(Class<T> type) {
        return proxy(AnnotationHandler.REF.getDefaultMPI(), type);
    }
}
