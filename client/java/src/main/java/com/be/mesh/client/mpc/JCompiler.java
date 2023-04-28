/*
 * Copyright (c) 2000, 2023, ducesoft and/or its affiliates. All rights reserved.
 * DUCESOFT PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.schema.runtime.TypeStruct;
import com.be.mesh.client.struct.Reference;
import com.be.mesh.client.struct.Service;
import com.be.mesh.client.tool.Tool;
import com.be.mesh.client.tool.UUID;

import java.lang.reflect.Method;
import java.util.Arrays;
import java.util.List;
import java.util.stream.Collectors;

/**
 * @author coyzeng@gmail.com
 */
@SPI(JCompiler.JAVASSIST)
public interface JCompiler {

    String JAVAC = "javac";
    String JAVASSIST = "javassist";

    /**
     * Get the parameter class type.
     * <p>
     * (Ljava/util/Map<Ljava/lang/String;Ljava/util/List<Ljava/util/Map<Ljava/lang/String;Lcom/be/mesh/client/struct/Principal;>;>;>;)V
     * (Ljava/util/Map<Ljava/lang/String;Ljava/util/List<Ljava/util/Map<Ljava/lang/String;Lcom/be/mesh/client/struct/Principal;>;>;>;Ljava/util/Set<Ljava/lang/String;>;)Z
     *
     * @param method Service method.
     * @return Parameter type.
     */
    <T extends Parameters> Class<T> intype(Method method);

    /**
     * Get the parameter class type by reference.
     *
     * @param reference Service reference.
     * @return Parameter type.
     */
    <T extends Parameters> Class<T> intype(Reference reference);

    /**
     * Get the parameter class type.
     * <p>
     * (Ljava/util/Map<Ljava/lang/String;Ljava/util/List<Ljava/util/Map<Ljava/lang/String;Lcom/be/mesh/client/struct/Principal;>;>;>;)V
     * (Ljava/util/Map<Ljava/lang/String;Ljava/util/List<Ljava/util/Map<Ljava/lang/String;Lcom/be/mesh/client/struct/Principal;>;>;>;Ljava/util/Set<Ljava/lang/String;>;)Z
     *
     * @param method Service method.
     * @return Parameter type.
     */
    <T extends Returns> Class<T> retype(Method method);

    /**
     * Get the parameter class type by reference.
     *
     * @param service service.
     * @return Parameter type.
     */
    <T extends Returns> Class<T> retype(Service service);

    /**
     * List all service documents.
     *
     * @return Service documents.
     */
    List<TypeStruct> documents();

    /**
     * Compile the code of interfaces.
     *
     * @param interfaces Interfaces spec
     * @param implement  Interfaces implement
     * @param <T>        Generic interface type
     * @return Implement class
     */
    <T> Class<? extends T> compile(Class<T> interfaces, String implement);

    /**
     * Get the type name.
     */
    default String getTypeName(Method method, String pattern, boolean pack) {
        return String.format("%s%s%s%s",
                pack ? method.getDeclaringClass().getName() : method.getDeclaringClass().getSimpleName(),
                Tool.firstUpperCase(method.getName()),
                pattern,
                UUID.getInstance().shortUUID());
    }

    /**
     * Get the signature.
     */
    default String getSignature(Method method, String pattern) {
        return method.getDeclaringClass().getName() + '.' + Tool.firstUpperCase(method.getName()) + pattern + '(' + Arrays.stream(method.getParameterTypes()).map(Class::getCanonicalName).collect(Collectors.joining(",")) + ')';
    }
}
