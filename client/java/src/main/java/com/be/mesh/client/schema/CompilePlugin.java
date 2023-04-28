/*
 * Copyright (c) 2000, 2023, ducesoft and/or its affiliates. All rights reserved.
 * DUCESOFT PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.schema;

import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.schema.context.CompileContext;
import com.be.mesh.client.schema.context.CompilePrinter;
import com.be.mesh.client.schema.runtime.Javadoc;

import java.lang.annotation.Annotation;
import java.util.Arrays;
import java.util.List;
import java.util.stream.Collectors;

/**
 * Compile plugin for javac compile.
 * <p>
 * See {@link Javadoc} with {@link Enum}.
 *
 * @author coyzeng@gmail.com
 */
public interface CompilePlugin extends CompilePrinter {

    /**
     * Compile plugin process with annotations.
     *
     * @param <A> annotation type
     * @return annotation collect
     */
    @SuppressWarnings({"unchecked"})
    default <A extends Annotation> List<Class<A>> with() {
        return Arrays
                .stream(this.getClass().getAnnotation(SPI.class).meta())
                .map(a -> (Class<A>) a)
                .distinct()
                .collect(Collectors.toList());
    }

    /**
     * Compile process.
     *
     * @param context javac compile context
     */
    void proceed(CompileContext context) throws Throwable;

}
