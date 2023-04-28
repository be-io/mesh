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
@SPI("grpc")
public interface Provider extends AutoCloseable {

    /**
     * Start the mesh broker.
     *
     * @throws Throwable cause
     */
    void start() throws Throwable;

}
