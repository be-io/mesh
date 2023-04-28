/*
 * Copyright (c) 2000, 2023, ducesoft and/or its affiliates. All rights reserved.
 * DUCESOFT PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.schema.macro;

import java.lang.annotation.*;

/**
 * @author coyzeng@gmail.com
 */
@Inherited
@Documented
@Target({ElementType.TYPE, ElementType.PACKAGE})
@Retention(RetentionPolicy.SOURCE)
public @interface License {

    /**
     * License text description.
     *
     * @return license
     */
    String value();

    /**
     * Active license plugin.
     *
     * @return true active
     */
    boolean active() default false;
}
