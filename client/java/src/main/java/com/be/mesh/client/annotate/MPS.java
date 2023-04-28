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
 * Multi Provider Service. Mesh provider service.
 *
 * @author coyzeng@gmail.com
 */
@Component
@Inherited
@Documented
@Target({ElementType.TYPE})
@Retention(RetentionPolicy.RUNTIME)
public @interface MPS {

    /**
     * Service name. As alias topic.
     */
    String name() default "";

    /**
     * Service version.
     */
    String version() default "*";

    /**
     * Service net/io protocol.
     */
    String proto() default MPI.GRPC;

    /**
     * Service codec.
     */
    String codec() default MPI.JSON;

    /**
     * Service invoke timeout. millions.
     */
    long timeout() default 3000;

    /**
     * Service flag 1 asyncable 2 encrypt 4 communal.
     */
    long flags() default 0;

}
