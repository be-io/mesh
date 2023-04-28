/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.cause.NotFoundException;
import lombok.AllArgsConstructor;

import java.lang.annotation.Annotation;
import java.lang.reflect.Type;
import java.util.Map;

/**
 * @author coyzeng@gmail.com
 */
@AllArgsConstructor
public class GenericInspector implements Inspector {

    private final String name;

    @Override
    public String getSignature() {
        return this.name;
    }

    @Override
    public Class<?> getType() {
        return Object.class;
    }

    @Override
    public String getName() {
        return this.name;
    }

    @Override
    public <T extends Annotation> T getAnnotation(Class<T> kind) {
        return null;
    }

    @Override
    public Class<?> getReturnType() {
        return Map.class;
    }

    @Override
    public Type getReturnGenericType() {
        return Map.class;
    }

    @Override
    public Class<?>[] getExceptionTypes() {
        return new Class[0];
    }

    @Override
    public Object invoke(Object obj, Object... args) throws ReflectiveOperationException {
        throw new NotFoundException("Generic execution cant serve %s as mps.", Mesh.context().getUrn());
    }
}
