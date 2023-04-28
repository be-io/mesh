/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.schema.macro;

import org.slf4j.event.Level;

import java.lang.annotation.*;

/**
 * @author coyzeng@gmail.com
 */
@Inherited
@Documented
@Target({ElementType.METHOD, ElementType.TYPE})
@Retention(RetentionPolicy.SOURCE)
public @interface Digest {

    /** digest name, default class name. */
    String value() default "";

    /** digest log level. */
    Level level() default Level.INFO;
}
