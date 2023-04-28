/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.system;

import io.be.mesh.macro.MPS;
import io.be.mesh.macro.SPI;
import io.be.mesh.prsim.Endpoint;

import java.nio.ByteBuffer;

/**
 * @author coyzeng@gmail.com
 */
@MPS
@SPI("mesh")
public class MeshEndpoint implements Endpoint, Endpoint.Sticker<ByteBuffer, ByteBuffer> {

    @Override
    public ByteBuffer fuzzy(ByteBuffer buff) {
        return null;
    }

    @Override
    public ByteBuffer stick(ByteBuffer varg) {
        return null;
    }
}
