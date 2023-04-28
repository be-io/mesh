/*
 * Copyright (c) 2000, 2023, ducesoft and/or its affiliates. All rights reserved.
 * DUCESOFT PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.schema.runtime;

import com.be.mesh.client.mpc.Factory;
import com.be.mesh.client.mpc.ServiceLoader;
import com.thoughtworks.qdox.JavaProjectBuilder;
import com.thoughtworks.qdox.model.DocletTag;
import com.thoughtworks.qdox.model.JavaAnnotatedElement;
import com.thoughtworks.qdox.model.JavaClass;
import lombok.AllArgsConstructor;
import lombok.extern.slf4j.Slf4j;

import java.io.InputStream;
import java.io.InputStreamReader;
import java.util.Collection;
import java.util.List;
import java.util.Optional;
import java.util.stream.Stream;


/**
 * @author coyzeng@gmail.com
 */
@Slf4j
@AllArgsConstructor
public final class Javadoc {

    private final JavaClass source;

    public static Javadoc getJavadoc(JavaProjectBuilder project, Class<?> type) {
        String uri = String.join("/", type.getName().split("\\."));
        String source = String.format("/%s.java", uri);
        try {
            InputStream input = ServiceLoader.load(Factory.class).getDefault().getResource(source);
            InputStreamReader reader = new InputStreamReader(input);
            project.addSource(reader);
        } catch (Exception e) {
            log.error(String.format("Cant found source of %s.", type.getName()), e);
        }
        JavaClass javaClass = project.getClassByName(source.substring(1));
        return Optional.ofNullable(javaClass).map(Javadoc::new).orElseGet(() -> new Javadoc(null));
    }

    public String getComment() {
        return Optional.ofNullable(source).map(JavaAnnotatedElement::getComment).map(String::trim).orElse("");
    }

    public String getMethodParameterComment(String methodName, String parameter) {
        return Optional.ofNullable(source.getMethods()).
                map(Collection::stream).orElseGet(Stream::empty).
                filter(method -> method.getName().equals(methodName)).findFirst().
                map(method -> {
                    List<DocletTag> tags = method.getTagsByName("param");
                    for (DocletTag tag : tags) {
                        List<String> parameters = tag.getParameters();
                        if (parameters.isEmpty()) {
                            continue;
                        }
                        String name = parameters.remove(0);
                        if (parameter.equals(name)) {
                            return String.join("", parameters);
                        }
                    }
                    return "";
                }).map(String::trim).orElse("");
    }

    public String getMethodReturnComment(String methodName) {
        return Optional.ofNullable(source.getMethods()).
                map(Collection::stream).orElseGet(Stream::empty).
                filter(method -> method.getName().equals(methodName)).findFirst().
                map(method -> method.getReturns().getComment()).map(String::trim).orElse("");
    }

    public String getFieldComment(String fieldName) {
        return Optional.ofNullable(source).map(x -> x.getFieldByName(fieldName)).
                map(JavaAnnotatedElement::getComment).
                map(String::trim).orElse("");
    }

    public String getMethodComment(String name) {
        return Optional.ofNullable(source.getMethods()).
                map(Collection::stream).orElseGet(Stream::empty).
                filter(method -> method.getName().equals(name)).findFirst().
                map(JavaAnnotatedElement::getComment).map(String::trim).orElse("");
    }

}
