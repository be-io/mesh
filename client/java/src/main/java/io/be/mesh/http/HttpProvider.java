/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.http;

import io.be.mesh.macro.SPI;
import io.be.mesh.mpc.Provider;
import lombok.extern.slf4j.Slf4j;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
@SPI("http")
public class HttpProvider implements Provider {

    @Override
    public void start() throws Throwable {

    }

    @Override
    public void close() throws Exception {

    }
}
