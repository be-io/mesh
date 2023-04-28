/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.system;

import com.be.mesh.client.annotate.MPS;
import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.prsim.Endpoint;

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
