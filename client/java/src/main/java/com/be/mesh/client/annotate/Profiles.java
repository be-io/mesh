/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.annotate;

import java.lang.annotation.*;

/**
 * Uniform format adapter.
 *
 * @author coyzeng@gmail.com
 */
@Inherited
@Documented
@Target({ElementType.TYPE})
@Retention(RetentionPolicy.RUNTIME)
public @interface Profiles {

    /**
     * Format name.
     */
    String value() default "";
}
