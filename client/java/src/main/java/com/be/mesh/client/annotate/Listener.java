/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.annotate;

import java.lang.annotation.*;

/**
 * @author coyzeng@gmail.com
 */
@Inherited
@Documented
@Target({ElementType.TYPE, ElementType.METHOD})
@Retention(RetentionPolicy.RUNTIME)
public @interface Listener {

    /**
     * Topic No.
     */
    String topic();

    /**
     * Event code.
     */
    String code() default "*";

    /**
     * Event version.
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
     * Event subscribe asyncable.
     */
    boolean async() default false;

    /**
     * Service invoke timeout. millions.
     */
    long timeout() default 3000;

    /**
     * Is publish as a mesh service.
     */
    boolean meshable() default true;

}
