/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.annotate;

import java.lang.annotation.*;

/**
 * Index for protobuf or thrift etc.
 *
 * @author coyzeng@gmail.com
 */
@Inherited
@Documented
@Target({ElementType.FIELD, ElementType.PARAMETER})
@Retention(RetentionPolicy.RUNTIME)
public @interface Index {

    /**
     * Index position.
     */
    int value();

    /**
     * Parameter name.
     */
    String name() default "";

    /**
     * Transparent when serialize/deserialize
     */
    boolean transparent() default false;

    /**
     * Alias
     */
    String[] alias() default {};
}
