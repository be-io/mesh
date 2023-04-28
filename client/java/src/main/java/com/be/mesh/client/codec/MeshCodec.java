/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.codec;

import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.mpc.Codec;
import com.be.mesh.client.mpc.Types;

import java.nio.ByteBuffer;

/**
 * @author coyzeng@gmail.com
 */
@SPI(Codec.MESH)
public class MeshCodec implements Codec {

    @Override
    public ByteBuffer encode(Object value) {
        return this.encode0(value, this::encodeTo);
    }

    @Override
    public <T> T decode(ByteBuffer buffer, Types<T> type) {
        return this.decode0(buffer, type, this::decodeFrom);
    }

    private ByteBuffer encodeTo(Object value) {
        return ByteBuffer.wrap(new byte[0]);
    }

    private <T> T decodeFrom(ByteBuffer buffer, Types<T> type) {
        return null;
    }

}
