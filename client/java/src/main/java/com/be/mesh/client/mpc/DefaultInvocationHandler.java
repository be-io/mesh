/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.tool.Tool;
import lombok.extern.slf4j.Slf4j;

import java.lang.invoke.MethodHandles;
import java.lang.reflect.Field;
import java.lang.reflect.Method;
import java.lang.reflect.Modifier;

/**
 * JDK8.
 *
 * @author coyzeng@gmail.com
 */
@Slf4j
public class DefaultInvocationHandler implements DefaultHandler {

    private static final MethodHandles.Lookup lookup = MethodHandles.lookup();

    static {
        try {
            Field allowedModes = MethodHandles.Lookup.class.getDeclaredField("allowedModes");
            if (Modifier.isFinal(allowedModes.getModifiers())) {
                final Field modifiersField = Field.class.getDeclaredField("modifiers");
                Tool.setFieldInt(allowedModes, modifiersField, allowedModes.getModifiers() & ~Modifier.FINAL);
                Tool.setField(lookup, allowedModes, -1);
            }
        } catch (Exception e) {
            log.warn("{}", e.getMessage());
        }
    }

    @Override
    public Object invoke(Object proxy, Method method, Object[] args) throws Throwable {
        Class<?> c = method.getDeclaringClass();
        return lookup.in(c).unreflectSpecial(method, c).bindTo(proxy).invokeWithArguments(args);
    }
}
