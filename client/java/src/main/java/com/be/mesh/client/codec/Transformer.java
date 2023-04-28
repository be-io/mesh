/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.codec;

import com.be.mesh.client.annotate.SPI;

import java.io.IOException;
import java.lang.reflect.Field;

/**
 * Transformer program interface.
 *
 * @author coyzeng@gmail.com
 */
@SPI("")
public interface Transformer<I> {

    /**
     * Transform given value.
     */
    void form(Writer writer, Field field, I value) throws IOException;

    /**
     * Transform from given value.
     */
    I from(Reader reader, Field field) throws IOException;

    /**
     * Is transformer can transform given value.
     */
    boolean matches(Field field);

    /**
     * Transform write decorator.
     */
    interface Writer {
        void write(Number value) throws IOException;

        void write(String value) throws IOException;

        void write(Boolean value) throws IOException;

        void writeNull() throws IOException;
    }

    /**
     * Transform read decorator.
     */
    interface Reader {
        Number readNumber() throws IOException;

        String readString() throws IOException;

        Boolean readBoolean() throws IOException;
    }
}
