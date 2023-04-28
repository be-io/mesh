/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.struct;

import com.be.mesh.client.annotate.Index;
import com.be.mesh.client.mpc.*;
import com.be.mesh.client.tool.Tool;
import lombok.Data;

import java.io.Serializable;
import java.nio.ByteBuffer;
import java.util.Optional;

/**
 * @author coyzeng@gmail.com
 */
@Data
public class Entity implements Serializable {

    private static final long serialVersionUID = 1330865476976014949L;
    @Index(0)
    private String codec;
    @Index(5)
    private String schema;
    @Index(10)
    private byte[] buffer;

    /**
     * Create an empty entity.
     */
    public static Entity empty() {
        return new Entity();
    }

    /**
     * Wrap an object with entity.
     */
    public static Entity wrap(Object value) {
        if (null == value) {
            return empty();
        }
        ServiceLoader<Codec> codec = ServiceLoader.load(Codec.class);
        ParameterizedTypes pt = new ParameterizedTypes(value.getClass());
        Entity entity = new Entity();
        entity.setCodec(codec.defaultName());
        entity.setSchema(Tool.compress(codec.getDefault().encodeString(pt)));
        entity.setBuffer(codec.getDefault().encode(value).array());
        return entity;
    }

    /**
     * Is the entity is present.
     */
    public boolean present() {
        return Tool.optional(this.schema);
    }

    /**
     * Try to get object with default type.
     *
     * @param <T> generic type
     * @return instance
     */
    public <T> T readObject() {
        ParameterizedTypes pt = loadCodec(this.getCodec()).decodeString(Tool.decompress(this.getSchema()), Types.of(ParameterizedTypes.class));
        return tryReadObject(Types.of(PatternParameterizedType.make(pt)));
    }

    /**
     * Try to get object with the type.
     *
     * @param type parameterized types
     * @param <T>  generic type
     * @return instance
     */
    public <T> T tryReadObject(Types<T> type) {
        if (null == this.getBuffer()) {
            return null;
        }
        return loadCodec(this.getCodec()).decode(ByteBuffer.wrap(this.getBuffer()), type);
    }

    /**
     * Get the codec with name.
     *
     * @param name codec name
     * @return codec provider
     */
    private Codec loadCodec(String name) {
        return ServiceLoader.load(Codec.class).get(Optional.ofNullable(name).orElse(Codec.JSON));
    }
}
