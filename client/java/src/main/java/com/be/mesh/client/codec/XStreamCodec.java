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
import com.thoughtworks.xstream.XStream;
import com.thoughtworks.xstream.security.AnyTypePermission;

import java.io.ByteArrayInputStream;
import java.nio.ByteBuffer;
import java.nio.charset.StandardCharsets;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

/**
 * @author coyzeng@gmail.com
 */
@SPI(Codec.XML)
public class XStreamCodec implements Codec {

    private final XStream xStream = new XStream();
    private final Map<Class<?>, Boolean> cache = new ConcurrentHashMap<>();

    @Override
    public ByteBuffer encode(Object value) {
        return this.encode0(value, object -> {
            keep(object.getClass());
            return ByteBuffer.wrap(xStream.toXML(object).getBytes(StandardCharsets.UTF_8));
        });
    }

    @SuppressWarnings("unchecked")
    @Override
    public <T> T decode(ByteBuffer buffer, Types<T> type) {
        return this.decode0(buffer, type, (x, y) -> {
            keep((Class<?>) y.getRawType());
            return (T) xStream.fromXML(new ByteArrayInputStream(x.array()));
        });
    }

    private void keep(Class<?> type) {
        cache.computeIfAbsent(type, x -> {
            xStream.autodetectAnnotations(true);
            xStream.addImmutableType(x, true);
            xStream.processAnnotations(x);
            xStream.allowTypes(cache.keySet().toArray(new Class<?>[0]));
            xStream.addPermission(AnyTypePermission.ANY);
            xStream.ignoreUnknownElements();
            return true;
        });
    }
}
