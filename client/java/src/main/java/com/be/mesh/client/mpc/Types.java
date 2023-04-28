/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.tool.Tool;
import lombok.AllArgsConstructor;

import java.io.Serializable;
import java.lang.reflect.Method;
import java.lang.reflect.ParameterizedType;
import java.lang.reflect.Type;
import java.util.List;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

/**
 * @author coyzeng@gmail.com
 */
public abstract class Types<T> implements ParameterizedType, Serializable {

    private static final long serialVersionUID = 6722418930517318263L;

    private static final Map<Type, Types<?>> kind = new ConcurrentHashMap<>();
    private static final Map<Type, Types<?>> list = new ConcurrentHashMap<>();
    private static final Map<Type, Map<Type, Types<?>>> dict = new ConcurrentHashMap<>();

    public static final Types<Map<String, String>> MapString = new Types<Map<String, String>>() {
    };
    public static final Types<Map<String, Object>> MapObject = new Types<Map<String, Object>>() {
    };
    public static final Types<List<String>> ListString = new Types<List<String>>() {
    };

    @SuppressWarnings("unchecked")
    public static <T> Types<T> of(Class<T> type) {
        return (Types<T>) kind.computeIfAbsent(type, key -> new ClassicTypes(type));
    }

    @SuppressWarnings("unchecked")
    public static <T> Types<T> of(Type type) {
        return (Types<T>) kind.computeIfAbsent(type, key -> new ParameterizedTypes(type));
    }

    @SuppressWarnings("unchecked")
    public static <T> Types<List<T>> list(Class<T> type) {
        return (Types<List<T>>) list.computeIfAbsent(type, key -> new ComplexTypes(List.class, new Type[]{type}));
    }

    @SuppressWarnings("unchecked")
    public static <T, K> Types<Map<T, K>> map(Class<T> type, Class<K> v) {
        return (Types<Map<T, K>>) dict.computeIfAbsent(type, key -> new ConcurrentHashMap<>()).computeIfAbsent(v, key -> new ComplexTypes(Map.class, new Type[]{type, v}));
    }

    @SuppressWarnings("unchecked")
    public static <T> Types<T> denest(Type type) {
        return (Types<T>) kind.computeIfAbsent(type, key -> new NestedTypes(type));
    }

    public static Type unbox(Method method, Class<?>... kinds) {
        if (method.getReturnType() == void.class || method.getReturnType() == Void.class) {
            return Object.class;
        }
        Class<?> rt = method.getReturnType();
        if (Tool.required(kinds)) {
            for (Class<?> kind : kinds) {
                if (kind.isAssignableFrom(rt)) {
                    return Types.denest(method.getGenericReturnType());
                }
            }
        }
        return Types.of(method.getGenericReturnType());
    }

    @Override
    public Type[] getActualTypeArguments() {
        Type type = getGenericType();
        if (type instanceof ParameterizedType) {
            return ((ParameterizedType) type).getActualTypeArguments();
        }
        return new Type[0];
    }

    @Override
    public Type getRawType() {
        Type type = getGenericType();
        if (type instanceof ParameterizedType) {
            return ((ParameterizedType) type).getRawType();
        }
        return type;
    }

    @Override
    public Type getOwnerType() {
        return null;
    }

    @Override
    public String toString() {
        if (null == getActualTypeArguments()) {
            return getRawType().getTypeName();
        }
        StringBuilder name = new StringBuilder(getRawType().getTypeName()).append('<');
        for (Type type : getActualTypeArguments()) {
            name.append(type.getTypeName()).append(',');
        }
        return name.append('>').toString();
    }

    private Type getGenericType() {
        Type sc = getClass().getGenericSuperclass();
        return ((ParameterizedType) sc).getActualTypeArguments()[0];
    }

    public boolean isAssignableFrom(Type type) {
        if (null == type) {
            return false;
        }
        if (this.getRawType() == type) {
            return true;
        }
        boolean iAmClass = this.getRawType() instanceof Class;
        boolean itIsClass = type instanceof Class;
        if (iAmClass && itIsClass) {
            return ((Class<?>) this.getRawType()).isAssignableFrom((Class<?>) type);
        }
        if (!iAmClass) {
            return of(this.getRawType()).isAssignableFrom(type);
        }
        if (type instanceof ParameterizedType) {
            return of(this.getRawType()).isAssignableFrom(((ParameterizedType) type).getRawType());
        }
        return false;
    }

    @AllArgsConstructor
    private static final class ClassicTypes extends Types<Object> {
        private final Class<?> type;

        @Override
        public Type getRawType() {
            return this.type;
        }

        @Override
        public Type[] getActualTypeArguments() {
            return new Type[0];
        }
    }

    @AllArgsConstructor
    private static final class ParameterizedTypes extends Types<Object> {
        private final transient Type type;

        @Override
        public Type getRawType() {
            if (this.type instanceof ParameterizedType) {
                return ((ParameterizedType) this.type).getRawType();
            }
            return this.type;
        }

        @Override
        public Type[] getActualTypeArguments() {
            if (this.type instanceof ParameterizedType) {
                return ((ParameterizedType) this.type).getActualTypeArguments();
            }
            return new Type[0];
        }
    }

    @AllArgsConstructor
    private static final class NestedTypes extends Types<Object> {
        private final transient Type type;

        @Override
        public Type getRawType() {
            Type nested = getNestedType();
            if (nested instanceof ParameterizedType) {
                return ((ParameterizedType) nested).getRawType();
            }
            return nested;
        }

        @Override
        public Type[] getActualTypeArguments() {
            Type nested = getNestedType();
            if (nested instanceof ParameterizedType) {
                return ((ParameterizedType) nested).getActualTypeArguments();
            }
            return new Type[0];
        }

        private Type getNestedType() {
            if (this.type instanceof ParameterizedType) {
                Type[] types = ((ParameterizedType) this.type).getActualTypeArguments();
                if (Tool.required(types)) {
                    return types[0];
                }
                return Object.class;
            }
            return Object.class;
        }

    }

    @AllArgsConstructor
    public static final class ComplexTypes extends Types<Object> {

        private final Type raw;
        private final Type[] args;

        @Override
        public Type getRawType() {
            return this.raw;
        }

        @Override
        public Type[] getActualTypeArguments() {
            if (Tool.required(this.args)) {
                return this.args;
            }
            return new Type[0];
        }

    }
}