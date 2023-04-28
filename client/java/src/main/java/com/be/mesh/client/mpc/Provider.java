/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.annotate.SPI;

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
