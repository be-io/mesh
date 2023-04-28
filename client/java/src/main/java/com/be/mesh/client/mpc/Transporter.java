/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.annotate.SPI;

import java.nio.ByteBuffer;

import static com.be.mesh.client.mpc.Transporter.PROVIDER;

/**
 * @author coyzeng@gmail.com
 */
@SPI(PROVIDER)
public interface Transporter {

    String PROVIDER = "provider";
    String CONSUMER = "consumer";

    /**
     * Transport the stream.
     *
     * @param buffer input stream
     * @return output stream
     * @throws Throwable e
     */
    ByteBuffer transport(String urn, ByteBuffer buffer) throws Throwable;
}
