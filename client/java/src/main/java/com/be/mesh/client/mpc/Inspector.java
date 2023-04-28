/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import java.lang.annotation.Annotation;
import java.lang.reflect.Type;

/**
 * @author coyzeng@gmail.com
 */
public interface Inspector {

    /**
     * Get the inspector signature.
     *
     * @return signature
     */
    String getSignature();

    /**
     * Get the declared type of inspect object declared.
     *
     * @return declared type.
     */
    Class<?> getType();

    /**
     * Get the name of inspector.
     *
     * @return name
     */
    String getName();

    /**
     * Get the annotations of inspector.
     *
     * @param kind annotation type
     * @param <T>  generic type
     * @return annotation
     */
    <T extends Annotation> T getAnnotation(Class<T> kind);

    /**
     * Get the return type of inspector.
     *
     * @return return type.
     */
    Class<?> getReturnType();

    /**
     * Get the generic return type of inspector.
     *
     * @return return type.
     */
    Type getReturnGenericType();

    /**
     * Get the exception types.
     *
     * @return exception types.
     */
    Class<?>[] getExceptionTypes();

    /**
     * Invoke the inspector object.
     *
     * @param obj  owner
     * @param args arguments
     * @return return value
     */
    Object invoke(Object obj, Object... args) throws ReflectiveOperationException;
}
