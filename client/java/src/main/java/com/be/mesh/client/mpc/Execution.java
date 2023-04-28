/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

/**
 * @author coyzeng@gmail.com
 */
public interface Execution<T> extends Invoker<T> {

    /**
     * Execution schema.
     */
    T schema();

    /**
     * Inspect execution.
     */
    Inspector inspect();

    /**
     * Execution input type.
     */
    <I extends Parameters> Types<I> intype();

    /**
     * Execution output return type.
     */
    <O extends Returns> Types<O> retype();

    /**
     * Reflect input type.
     */
    Parameters inflect();

    /**
     * Reflect output return type.
     */
    Returns reflect();
}
