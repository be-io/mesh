/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.struct.Reference;

import java.io.InputStream;
import java.util.concurrent.Future;

/**
 * @author coyzeng@gmail.com
 */
@SPI(Consumer.GRPC)
public interface Consumer extends AutoCloseable {

    String HTTP = "http";
    String GRPC = "grpc";
    String TCP = "tcp";
    String MQTT = "mqtt";

    /**
     * Start the mesh broker.
     *
     * @throws Throwable cause
     */
    void start() throws Throwable;

    /**
     * Consume the input payload.
     *
     * @param address   Remote address.
     * @param urn       Actual uniform resource domain name.
     * @param execution Service reference.
     * @param inbound   Input arguments.
     * @return Output payload
     */
    Future<InputStream> consume(String address, String urn, Execution<Reference> execution, InputStream inbound);
}
