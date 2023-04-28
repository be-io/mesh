/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.annotate.SPI;

import java.nio.ByteBuffer;
import java.nio.charset.StandardCharsets;
import java.util.Optional;
import java.util.function.BiFunction;
import java.util.function.Function;

/**
 * @author coyzeng@gmail.com
 */
@SPI(Codec.JSON)
public interface Codec {

    String JSON = "json";
    String JACKSON = "jackson";
    String PROTOBUF = "protobuf";
    String PROTOBUF4 = "protobuf4";
    String XML = "xml";
    String YAML = "yaml";
    String THRIFT = "thrift";
    String MESH = "mesh";

    /**
     * Encode object to binary array.
     *
     * @param value object
     * @return binary array
     */
    ByteBuffer encode(Object value);

    /**
     * Decode to binary array to object.
     *
     * @param buffer binary array
     * @param type   object type
     * @param <T>    type
     * @return typed object
     */
    <T> T decode(ByteBuffer buffer, Types<T> type);

    /**
     * JSON编码.
     *
     * @param value 数据
     * @return JSON
     */
    default String encodeString(Object value) {
        ByteBuffer buffer = this.encode(value);
        return new String(buffer.array(), StandardCharsets.UTF_8);
    }

    /**
     * JSON解码.
     *
     * @param value 数据
     * @param <T>   类型引用
     * @return Object
     */
    default <T> T decodeString(String value, Types<T> type) {
        ByteBuffer buffer = ByteBuffer.wrap(Optional.ofNullable(value).orElse("").getBytes(StandardCharsets.UTF_8));
        return this.decode(buffer, type);
    }

    /**
     * Encode object to binary array.
     *
     * @param value object
     * @return binary array
     */
    default ByteBuffer encode0(Object value, Function<Object, ByteBuffer> fn) {
        if (value instanceof ByteBuffer) {
            return (ByteBuffer) value;
        }
        if (value instanceof String) {
            return ByteBuffer.wrap(((String) value).getBytes(StandardCharsets.UTF_8));
        }
        if (value instanceof byte[]) {
            return ByteBuffer.wrap((byte[]) value);
        }
        return fn.apply(value);
    }

    /**
     * Decode to binary array to object.
     *
     * @param buffer binary array
     * @param type   object type
     * @param <T>    type
     * @return typed object
     */
    @SuppressWarnings("unchecked")
    default <T> T decode0(ByteBuffer buffer, Types<T> type, BiFunction<ByteBuffer, Types<T>, T> fn) {
        if (type.getRawType() == String.class) {
            return (T) new String(buffer.array(), StandardCharsets.UTF_8);
        }
        if (type.getRawType() == byte[].class) {
            return (T) buffer.array();
        }
        if (type.getRawType() == ByteBuffer.class) {
            return (T) buffer.array();
        }
        return fn.apply(buffer, type);
    }

}
