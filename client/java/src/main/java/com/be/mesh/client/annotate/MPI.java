/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.annotate;

import java.lang.annotation.*;

/**
 * Multi Provider Interface. Mesh provider interface.
 *
 * @author coyzeng@gmail.com
 */
@Inherited
@Documented
@Target({ElementType.FIELD, ElementType.METHOD, ElementType.TYPE})
@Retention(RetentionPolicy.RUNTIME)
public @interface MPI {

    String GRPC = "grpc";
    String HTTP = "http";
    String PROTOBUF = "protobuf";
    String PROTOBUF4 = "protobuf4";
    String JSON = "json";

    /**
     * Service name. {@link MPS#name()}
     */
    String value() default "";

    /**
     * Service name. {@link MPS#name()}
     */
    String name() default "";

    /**
     * Service version. {@link MPS#version()}
     */
    String version() default "*";

    /**
     * Service net/io protocol.
     */
    String proto() default GRPC;

    /**
     * Service codec.
     */
    String codec() default JSON;

    /**
     * Service flag 1 asyncable 2 encrypt 4 communal.
     */
    long flags() default 0;

    /**
     * Service invoke timeout. millions.
     */
    long timeout() default 10000;

    /**
     * Invoke retry times.
     */
    int retries() default 3;

    /**
     * Service node identity.
     */
    String node() default "";

    /**
     * Service inst identity.
     */
    String inst() default "";

    /**
     * Service zone.
     */
    String zone() default "";

    /**
     * Service cluster.
     */
    String cluster() default "";

    /**
     * Service cell.
     */
    String cell() default "";

    /**
     * Service group.
     */
    String group() default "";

    /**
     * Service address.
     */
    String address() default "";
}
