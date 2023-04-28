/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.cause.MeshException;

import java.io.ByteArrayInputStream;
import java.io.InputStream;
import java.nio.ByteBuffer;
import java.util.stream.Stream;

/**
 * SpiExtensionFactory.
 *
 * @author coyzeng@gmail.com
 */
@SPI("java")
public class ManifestFactory implements Factory {

    @Override
    public <T> T getProvider(String key, Class<T> spi) {
        if (spi.isInterface() && spi.isAnnotationPresent(SPI.class)) {
            ServiceLoader<T> serviceLoader = ServiceLoader.load(spi);
            return serviceLoader.getDefault();
        }
        String format = String.format(WARN, spi.getName());
        throw new MeshException(format);
    }

    @Override
    public <T> Stream<T> getProvider(Class<T> spi) {
        return ServiceLoader.load(spi).list().stream();
    }

    @Override
    public InputStream getResource(String name) {
        ByteBuffer buffer = ServiceLoader.resource(name);
        return new ByteArrayInputStream(buffer.array());
    }

}
