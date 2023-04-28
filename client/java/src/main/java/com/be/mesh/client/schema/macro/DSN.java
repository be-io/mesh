/*
 * Copyright (c) 2000, 2023, ducesoft and/or its affiliates. All rights reserved.
 * DUCESOFT PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.schema.macro;

import java.lang.annotation.Inherited;
import java.lang.annotation.Retention;
import java.lang.annotation.Target;

import static java.lang.annotation.ElementType.*;
import static java.lang.annotation.RetentionPolicy.RUNTIME;

/**
 * multi data source support.<br/>
 * warn: for transactional. you should ensure one data source one transactional.
 * otherwise see @TransactionalX
 * <p>
 * Example:
 * <pre>
 * &#064;Primary
 * &#064;Bean
 * public DataSource a(){...}
 * &#064;Bean
 * public DataSource b(){...}
 * &#064;Bean
 * public DataSource c(){
 *     bean.a(); // support different data source method nested call
 * }
 * &#064;Bean
 * public TransactionManager tm(IDataSource dataSource) {...}
 *
 * &#064;DataSource("a")
 * public void insert(){...}
 * &#064;DataSource("b")
 * public void update(){...}
 * &#064;DataSource("c")
 * public void select(){...}
 * &#064;DataSource() // will use the primary data source a as default
 * public void insert(){...}
 *
 * also can use on field inject:
 * &#064;Resource
 * &#064;DataSource("a")
 * private IDao dao; // this dao will use the data source a
 *
 * much see at EnableMultiDataSource
 *
 * @author coyzeng@gmail.com
 */

@Inherited
@Target({TYPE, METHOD, FIELD, TYPE_USE})
@Retention(RUNTIME)
public @interface DSN {

    /**
     * data source id or name
     */
    String value() default "";
}
