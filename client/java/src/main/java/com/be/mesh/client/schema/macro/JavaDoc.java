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
@Target(ElementType.TYPE)
@Retention(RetentionPolicy.RUNTIME)
public @interface JavaDoc {

    /**
     * javadoc package name, default current class.
     */
    String[] value() default "";

    /**
     * module name.
     */
    String name() default "";

    /**
     * exclude package.
     */
    String[] exclude() default {};

    /**
     * Ignore package.
     */
    String[] ignore() default {"java", "javax", "org"};

    /**
     * Javadoc version.
     */
    String version() default "1.0.0";

}
