/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.annotate;

import java.lang.annotation.*;

/**
 * Uniform format adapter. This is not standard program interface, just for compatible.
 *
 * <pre>
 *     private static final DateTimeFormatter DATE = DateTimeFormatter.ofPattern("yyyy-MM-dd");
 *     private static final DateTimeFormatter TIME = DateTimeFormatter.ofPattern("HH:mm:ss");
 *     private static final DateTimeFormatter DATE_TIME = DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss");
 * </pre>
 *
 * @author coyzeng@gmail.com
 */
@Inherited
@Documented
@Target({ElementType.FIELD, ElementType.PARAMETER, ElementType.TYPE})
@Retention(RetentionPolicy.RUNTIME)
public @interface Format {

    /**
     * Format name.
     */
    String value() default "";

    /**
     * Format pattern.
     */
    String pattern() default "";

    /**
     * Formatter spi name.
     */
    String name() default "";

    /**
     * Format type kind.
     */
    Class<?> former() default Object.class;

}
