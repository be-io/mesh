/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.schema.runtime;

import com.be.mesh.client.annotate.*;
import com.be.mesh.client.mpc.Codec;
import com.be.mesh.client.mpc.Schema;
import com.be.mesh.client.mpc.ServiceLoader;
import com.be.mesh.client.tool.Tool;
import com.thoughtworks.qdox.JavaProjectBuilder;
import lombok.extern.slf4j.Slf4j;

import java.lang.reflect.Method;
import java.lang.reflect.Parameter;
import java.lang.reflect.Type;
import java.util.*;
import java.util.concurrent.ConcurrentHashMap;
import java.util.stream.Collectors;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
@SPI("dynamic")
public class DynamicSchema implements Schema {

    private final List<Interface> interfaces = new ArrayList<>();
    private final Map<String, Class<?>> types = new ConcurrentHashMap<>();
    private final Map<Class<?>, Interface> references = new ConcurrentHashMap<>();

    @Override
    public void imports(List<Class<?>> types) {
        JavaProjectBuilder project = new JavaProjectBuilder();
        for (Class<?> type : types) {
            Javadoc javadoc = Javadoc.getJavadoc(project, type);
            List<Method> methods = getDeclaredMethods(type);
            List<Function> functions = new ArrayList<>(methods.size());
            for (Method method : methods) {
                List<Attribute> attributes = new ArrayList<>();
                for (int index = 0; index < method.getParameters().length; index++) {
                    Parameter parameter = method.getParameters()[0];
                    Attribute attribute = new Attribute();
                    attribute.setKind(parameter.getType().getName());
                    attribute.setName(parameter.getName());
                    attribute.setAlias("");
                    attribute.setIndex(Optional.ofNullable(parameter.getAnnotation(Index.class)).map(Index::value).orElse(index));
                    attribute.setFlags(1);
                    attribute.setOptional(true);
                    attribute.setComment(javadoc.getMethodParameterComment(method.getName(), parameter.getName()));
                    attribute.setAttributes(resolveAttribute(parameter.getParameterizedType()));
                    attributes.add(attribute);
                }
                Attribute attribute = new Attribute();
                attribute.setKind(method.getGenericReturnType().getTypeName());
                attribute.setName("");
                attribute.setAlias("");
                attribute.setIndex(0);
                attribute.setFlags(2);
                attribute.setOptional(true);
                attribute.setComment(javadoc.getMethodReturnComment(method.getName()));
                attribute.setAttributes(resolveAttribute(method.getGenericReturnType()));
                attribute.setAttributes(attributes);

                Function function = new Function();
                function.setName(method.getName());
                function.setAlias(Optional.ofNullable(method.getAnnotation(MPS.class)).map(MPS::name).orElse(""));
                function.setComment(javadoc.getMethodComment(method.getName()));
                function.setAttributes(attributes);
                functions.add(function);
            }

            Interface kind = new Interface();
            kind.setVersion(resolveVersion(type));
            kind.setKind(type.getName());
            kind.setName(Tool.formatObjectName(type));
            kind.setAlias(resolveAlias(type));
            kind.setComment(javadoc.getComment());
            kind.setFunctions(functions);
            interfaces.add(kind);
        }
    }

    private List<Attribute> resolveAttribute(Type type) {
        return new ArrayList<>();
    }

    @Override
    public void imports(String schema) {

    }

    @Override
    public String exports() {
        Codec codec = ServiceLoader.load(Codec.class).getDefault();
        return codec.encodeString(this.interfaces);
    }

    @Override
    public Class<?> search(String urn) {
        return null;
    }

    private String resolveVersion(Class<?> type) {
        if (type.isAnnotationPresent(MPS.class)) {
            return type.getAnnotation(MPS.class).version();
        }
        if (type.isAnnotationPresent(Binding.class)) {
            return type.getAnnotation(Binding.class).version();
        }
        if (type.isAnnotationPresent(Bindings.class)) {
            return type.getAnnotation(Bindings.class).value()[0].version();
        }
        return "1.0.0";
    }

    private String resolveAlias(Class<?> type) {
        if (type.isAnnotationPresent(MPS.class)) {
            return type.getAnnotation(MPS.class).name();
        }
        if (type.isAnnotationPresent(Binding.class)) {
            return type.getAnnotation(Binding.class).topic();
        }
        if (type.isAnnotationPresent(Bindings.class)) {
            return type.getAnnotation(Bindings.class).value()[0].topic();
        }
        return "1.0.0";
    }

    private List<Method> getDeclaredMethods(Class<?> type) {
        List<Class<?>> kinds = new ArrayList<>(type.getInterfaces().length + 1);
        kinds.addAll(Arrays.asList(type.getInterfaces()));
        if (type.isInterface()) {
            kinds.add(type);
        }
        return kinds.stream().map(Class::getDeclaredMethods).flatMap(Arrays::stream).filter(Tool::canService).collect(Collectors.toList());
    }

}
