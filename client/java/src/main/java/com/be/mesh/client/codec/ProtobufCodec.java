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
import io.protostuff.*;
import io.protostuff.runtime.RuntimeSchema;

import java.nio.ByteBuffer;
import java.util.ArrayList;
import java.util.Collection;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

/**
 * @author coyzeng@gmail.com
 */
@SPI(Codec.PROTOBUF)
public class ProtobufCodec implements Codec {

    private final Map<Types<?>, Schema<Object>> schemas = new ConcurrentHashMap<>();

    @SuppressWarnings("unchecked")
    @Override
    public ByteBuffer encode(Object value) {
        return this.encode0(value, object -> {
            // this is lazily created and cached by RuntimeSchema
            // so its safe to call RuntimeSchema.getSchema(Foo.class) over and over
            // The getSchema method is also thread-safe
            Schema<Object> schema = RuntimeSchema.getSchema((Class<Object>) object.getClass());
            // Re-use (manage) this buffer to avoid allocating on every serialization
            LinkedBuffer buffer = LinkedBuffer.allocate(512);
            try {
                return ByteBuffer.wrap(ProtostuffIOUtil.toByteArray(object, schema, buffer));
            } finally {
                buffer.clear();
            }
        });
    }

    @Override
    public <T> T decode(ByteBuffer buffer, Types<T> type) {
        return this.decode0(buffer, type, (x, y) -> {
            Schema<T> schema = createGenericTypeIfAbsent(y);
            T object = schema.newMessage();
            ProtostuffIOUtil.mergeFrom(x.array(), object, schema);
            return object;
        });
    }

    @SuppressWarnings("unchecked")
    private <T> Schema<T> createGenericTypeIfAbsent(Types<T> type) {
        return (Schema<T>) schemas.computeIfAbsent(type, kind -> (Schema<Object>) createSchema(kind));
    }

    private <T> Object createSchema(Types<T> kind) {
        if (kind.isAssignableFrom(Map.class)) {
            if (kind.getActualTypeArguments().length > 0) {
                Schema<?> ks = RuntimeSchema.getSchema((Class<?>) kind.getActualTypeArguments()[0]);
                Schema<?> vs = RuntimeSchema.getSchema((Class<?>) kind.getActualTypeArguments()[1]);
                return new MessageMapSchema<>(ks, vs);
            }
            return RuntimeSchema.getSchema(HashMap.class);
        }
        if (kind.isAssignableFrom(Collection.class)) {
            if (kind.getActualTypeArguments().length > 0) {
                Schema<?> vs = RuntimeSchema.getSchema((Class<?>) kind.getActualTypeArguments()[0]);
                return new MessageCollectionSchema<>(vs, false);
            }
            return RuntimeSchema.getSchema(ArrayList.class);
        }
        return RuntimeSchema.getSchema((Class<?>) kind.getRawType());
    }
}
