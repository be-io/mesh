/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.tool.FitMap;

import java.io.Serializable;
import java.lang.reflect.Type;
import java.util.*;

/**
 * @author coyzeng@gmail.com
 */
public class GenericParameters extends TreeMap<String, Object> implements Serializable, Parameters {

    private static final long serialVersionUID = -4338158729657717760L;
    private static final String ATTACHMENTS = "attachments";

    @Override
    public Map<String, Object> map() {
        return this;
    }

    @Override
    public Class<?> type() {
        return this.getClass();
    }

    @Override
    public Object[] arguments() {
        List<Object> args = new ArrayList<>(this.size());
        this.forEach((key, value) -> {
            if (ATTACHMENTS.equals(key)) {
                return;
            }
            args.add(value);
        });
        return args.toArray();
    }

    @Override
    public void arguments(Object[] arguments) {
        //
    }

    @Override
    public Map<Integer, Type> argumentTypes() {
        return new HashMap<>();
    }

    @SuppressWarnings("unchecked")
    @Override
    public Map<String, String> attachments() {
        Object attachments = this.computeIfAbsent(ATTACHMENTS, key -> new HashMap<>());
        if (attachments instanceof Map) {
            return new FitMap<>((Map<String, String>) attachments);
        }
        Codec codec = ServiceLoader.load(Codec.class).get(Codec.JSON);
        Map<String, String> des = codec.decode(codec.encode(attachments), Types.MapString);
        this.put(ATTACHMENTS, des);
        return new FitMap<>(des);
    }

    @Override
    public void attachments(Map<String, String> attachments) {
        this.put(ATTACHMENTS, attachments);
    }
}
