/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.codec;

import io.be.mesh.macro.SPI;
import io.be.mesh.mpc.Codec;
import io.be.mesh.mpc.Types;

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
