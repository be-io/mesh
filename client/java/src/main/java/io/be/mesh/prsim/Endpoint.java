/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.prsim;

import io.be.mesh.macro.MPI;
import io.be.mesh.macro.SPI;

import java.nio.ByteBuffer;

/**
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public interface Endpoint {

    /**
     * Fuzzy call with generic param
     * In multi returns, it's an array.
     */
    @MPI("${mesh.uname}")
    ByteBuffer fuzzy(ByteBuffer buff);

    /**
     * @param <I> is the input
     * @param <O> is the output
     */
    interface Sticker<I, O> {

        // Stick with generic param
        O stick(I varg);

    }

}
