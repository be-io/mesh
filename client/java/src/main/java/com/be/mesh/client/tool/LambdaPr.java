/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.tool;

/**
 * @author coyzeng@gmail.com
 */
@FunctionalInterface
public interface LambdaPr<T> {
    boolean test(T input) throws Throwable;
}