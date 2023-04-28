/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.annotate;

import org.springframework.stereotype.Component;

import java.lang.annotation.*;

/**
 * @author coyzeng@gmail.com
 */
@Component
@Inherited
@Documented
@Target({ElementType.TYPE, ElementType.METHOD})
@Retention(RetentionPolicy.RUNTIME)
public @interface Bindings {

    /**
     * Subscribe bindings.
     */
    Binding[] value();
}
