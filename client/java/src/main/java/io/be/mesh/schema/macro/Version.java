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
public @interface Version {

    /**
     * Application name, like name()
     */
    String value() default "";

    /**
     * Application name, only one words.
     */
    String name() default "";

    /**
     * Application version, default determine from git.
     */
    String version() default "";

}
