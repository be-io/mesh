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
@SPI(Codec.THRIFT)
public class ThriftCodec implements Codec {

    @Override
    public ByteBuffer encode(Object value) {
        return null;
    }

    @Override
    public <T> T decode(ByteBuffer buffer, Types<T> type) {
        return null;
    }
}
