/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.schema;

import io.be.mesh.macro.SPI;
import io.be.mesh.schema.context.CompileContext;
import io.be.mesh.schema.context.CompilePrinter;
import io.be.mesh.schema.runtime.Javadoc;

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
