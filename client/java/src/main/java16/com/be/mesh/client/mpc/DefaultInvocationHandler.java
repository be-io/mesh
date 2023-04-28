/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import java.lang.reflect.Method;

/**
 * JDK16.
 *
 * @author coyzeng@gmail.com
 */
public class DefaultInvocationHandler implements DefaultHandler {

    @Override
    public Object invoke(Object proxy, Method method, Object[] args) throws Throwable {
        return java.lang.reflect.InvocationHandler.invokeDefault(proxy, method, args);
    }

}
