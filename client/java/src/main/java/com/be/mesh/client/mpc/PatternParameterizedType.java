/*
 * Copyright (c) 2000, 2023, ducesoft and/or its affiliates. All rights reserved.
 * DUCESOFT PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.cause.MeshException;
import com.be.mesh.client.tool.Tool;

import java.io.Serializable;
import java.lang.reflect.ParameterizedType;
import java.lang.reflect.Type;
import java.util.Arrays;

/**
 * List&lt;List&lt;Map&lt;Map&lt;String, String&gt;, Map&lt;String, Map&lt;Object, String&gt;&gt;&gt;&gt;&gt;
 *
 * @author coyzeng@gmail.com
 */
public class PatternParameterizedType implements ParameterizedType,
        Serializable {

    private static final long serialVersionUID = -3695969395497282211L;
    private final transient Type ownerType;
    private final transient Type rawType;
    private final transient Type[] typeArguments;

    private PatternParameterizedType(Class<?> rawType, Type[] typeArguments) {
        this.rawType = rawType;
        this.ownerType = ClassOwnership.detectJvmBehavior().getOwnerType(rawType);
        this.typeArguments = typeArguments;
    }

    public static PatternParameterizedType make(ParameterizedTypes struct) {
        try {
            if (Tool.optional(struct.getArgs())) {
                return new PatternParameterizedType(Class.forName(struct.getRaw()), new Type[0]);
            }
            Type[] types = struct.getArgs().stream().map(PatternParameterizedType::make).toArray(Type[]::new);
            return new PatternParameterizedType(Class.forName(struct.getRaw()), types);
        } catch (Exception e) {
            throw new MeshException("Make type {} failed.", e, struct.toString());
        }
    }

    @Override
    public Type getOwnerType() {
        return this.ownerType;
    }

    @Override
    public Type getRawType() {
        return this.rawType;
    }

    @Override
    public Type[] getActualTypeArguments() {
        return this.typeArguments;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) {
            return true;
        }
        if (!(other instanceof ParameterizedType)) {
            return false;
        }
        ParameterizedType otherType = (ParameterizedType) other;
        return (otherType.getOwnerType() == null &&
                this.rawType.equals(otherType.getRawType()) &&
                Arrays.equals(this.typeArguments, otherType.getActualTypeArguments()));
    }

    @Override
    public int hashCode() {
        return Arrays.hashCode(typeArguments) ^
                (ownerType == null ? 0 : ownerType.hashCode()) ^
                (rawType == null ? 0 : rawType.hashCode());
    }

    @Override
    public String toString() {
        StringBuilder st = new StringBuilder();

        if (null != ownerType) {
            if (ownerType instanceof Class) {
                st.append(((Class<?>) ownerType).getName());
            } else {
                st.append(ownerType.toString());
            }

            st.append(".");

            if (ownerType instanceof PatternParameterizedType) {
                // Find simple name of nested type by removing the
                // shared prefix with owner.
                st.append(rawType.getTypeName().replace(((PatternParameterizedType) ownerType).rawType.getTypeName() + "$", ""));
            } else {
                st.append(rawType.getTypeName());
            }
        } else {
            st.append(rawType.getTypeName());
        }

        if (null != typeArguments && typeArguments.length > 0) {
            st.append("<");
            boolean first = true;
            for (Type t : typeArguments) {
                if (!first) {
                    st.append(", ");
                }
                if (t instanceof Class) {
                    st.append(((Class<?>) t).getName());
                } else {
                    st.append(t.toString());
                }
                first = false;
            }
            st.append(">");
        }

        return st.toString();
    }

    enum ClassOwnership {
        OWNED_BY_ENCLOSING_CLASS {
            @Override
            Class<?> getOwnerType(Class<?> rawType) {
                return rawType.getEnclosingClass();
            }
        },
        LOCAL_CLASS_HAS_NO_OWNER {
            @Override
            Class<?> getOwnerType(Class<?> rawType) {
                if (rawType.isLocalClass()) {
                    return null;
                } else {
                    return rawType.getEnclosingClass();
                }
            }
        };

        abstract Class<?> getOwnerType(Class<?> rawType);

        static final ClassOwnership JVM_BEHAVIOR = detectJvmBehavior();

        private static ClassOwnership detectJvmBehavior() {
            class LocalClass<T> {
            }
            Class<?> subclass = new LocalClass<String>() {
            }.getClass();
            ParameterizedType parameterizedType = (ParameterizedType) subclass.getGenericSuperclass();
            for (ClassOwnership behavior : ClassOwnership.values()) {
                if (behavior.getOwnerType(LocalClass.class) == parameterizedType.getOwnerType()) {
                    return behavior;
                }
            }
            throw new AssertionError();
        }
    }

}

