/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.annotate.Index;
import com.be.mesh.client.tool.Tool;

import java.lang.reflect.Field;
import java.lang.reflect.Type;
import java.util.HashMap;
import java.util.Map;

/**
 * Without get and set.
 *
 * @author coyzeng@gmail.com
 */
public interface Parameters {

    /**
     * Convert parameters to map.
     */
    Map<String, Object> map();

    /**
     * Parameters declared type.
     */
    Class<?> type();

    /**
     * Generic arguments array.
     */
    Object[] arguments();

    /**
     * Generic arguments array.
     */
    void arguments(Object[] arguments);

    /**
     * Get the generic attachments. The attributes will be serialized. The attachments are mutable.
     *
     * @return attachments.
     */
    Map<String, String> attachments();

    /**
     * Attachment arguments.
     */
    void attachments(Map<String, String> attachments);

    /**
     * Arguments types.
     */
    default Map<Integer, Type> argumentTypes() {
        Map<Integer, Type> types = new HashMap<>();
        for (Field f : this.getClass().getDeclaredFields()) {
            if (f.isAnnotationPresent(Index.class) && !Tool.equals("attachments", f.getName())) {
                types.put(f.getAnnotation(Index.class).value(), f.getGenericType());
            }
        }
        return types;
    }
}
