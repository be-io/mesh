/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import java.lang.reflect.Method;

/**
 * @author coyzeng@gmail.com
 */
public interface DefaultHandler {

    /**
     * Invoke the method in default mode.
     *
     * @param proxy  proxy
     * @param method method
     * @param args   arguments
     * @return result
     * @throws Throwable e
     */
    Object invoke(Object proxy, Method method, Object[] args) throws Throwable;

}
