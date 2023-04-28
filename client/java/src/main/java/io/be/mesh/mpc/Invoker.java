/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.mpc;


import io.be.mesh.macro.SPI;

/**
 * @author coyzeng@gmail.com
 */
@SPI(pattern = "*")
public interface Invoker<T> {

    /**
     * Invoke the next invoker.
     *
     * @param invocation Invoke context.
     * @return result
     * @throws Throwable cause
     */
    Object invoke(Invocation invocation) throws Throwable;

}
