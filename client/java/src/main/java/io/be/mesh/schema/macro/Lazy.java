/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.schema.macro;

import java.lang.annotation.*;
import java.lang.reflect.Modifier;

/**
 * @author coyzeng@gmail.com
 */
@Inherited
@Documented
@Target({ElementType.FIELD, ElementType.TYPE, ElementType.METHOD})
@Retention(RetentionPolicy.SOURCE)
public @interface Lazy {

    /**
     * is static sharable.
     */
    boolean share() default false;

    /**
     * the variable modifier level.
     */
    int modifier() default Modifier.PRIVATE;
}
