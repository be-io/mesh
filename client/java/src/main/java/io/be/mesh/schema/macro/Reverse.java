/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.schema.macro;

import java.lang.annotation.*;

/**
 * @author coyzeng@gmail.com
 */
@Inherited
@Documented
@Target(ElementType.TYPE)
@Retention(RetentionPolicy.RUNTIME)
public @interface Reverse {

    /** table name */
    String value();

    /** package name */
    String pack();

    /** is reverse override. */
    boolean override() default false;
}
