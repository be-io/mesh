/*
 * Copyright (c) 2000, 2023, ducesoft and/or its affiliates. All rights reserved.
 * DUCESOFT PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.annotate;

import java.lang.annotation.Documented;
import java.lang.annotation.Inherited;
import java.lang.annotation.Retention;
import java.lang.annotation.Target;
import java.util.ServiceLoader;

import static java.lang.annotation.ElementType.*;
import static java.lang.annotation.RetentionPolicy.RUNTIME;

/**
 * Metadata annotation for Serial Peripheral Interface. Can be used with {@link ServiceLoader#load(Class)} or
 * dependency injection at compile time and runtime time.
 * <p/>
 * Annotated on class or method will provider service as SPI.
 * <p/>
 * Example:
 * <pre>
 * &#064;SPI(value="name")
 * public class SPIService {
 *
 *     &#064;SPI(value="name")
 *     public SPIService spi(){...}
 * }
 * </pre>
 * Annotated on field or parameter will reference SPI service as reference for dependency injection.
 * <p/>
 * Example:
 * <pre>
 * public class Demo {
 *      &#064;SPI(value="name")
 *      private SPIService spi;
 *
 *      public void call(&#064;SPI(value="name") SPIService spi) {...}
 * }
 * </pre>
 *
 * @author coyzeng@gmail.com
 */
@Inherited
@Documented
@Retention(RUNTIME)
@Target({TYPE, METHOD, FIELD, PARAMETER})
public @interface SPI {

    /**
     * xxx.xxx.xxx
     */
    String value() default "";

    /**
     * Pattern match
     */
    String pattern() default "";

    /**
     * spi priority order by asc, highest is {@link Integer#MIN_VALUE}
     */
    int priority() default 0;

    /**
     * metadata for the spi.
     */
    boolean prototype() default false;

    /**
     * SPI alias.
     */
    String[] alias() default {};

    /**
     * metadata for the spi.
     */
    Class<?>[] meta() default {};

    /**
     * Exclude metadata.
     */
    String[] exclude() default {};

    /**
     * Enable the spi provider default.
     */
    boolean enable() default true;
}
