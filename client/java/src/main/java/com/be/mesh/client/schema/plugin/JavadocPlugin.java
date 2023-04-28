/*
 * Copyright (c) 2000, 2023, ducesoft and/or its affiliates. All rights reserved.
 * DUCESOFT PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.schema.plugin;

import com.be.mesh.client.annotate.Exclude;
import com.be.mesh.client.annotate.MPI;
import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.mpc.Codec;
import com.be.mesh.client.mpc.ServiceLoader;
import com.be.mesh.client.schema.CompilePlugin;
import com.be.mesh.client.schema.context.CompileContext;
import com.be.mesh.client.schema.javadoc.Kind;
import com.be.mesh.client.schema.javadoc.Parser;
import com.be.mesh.client.schema.macro.JavaDoc;
import com.be.mesh.client.tool.Tool;
import lombok.Data;

import javax.lang.model.element.Element;
import javax.lang.model.element.ElementKind;
import javax.lang.model.element.TypeElement;
import javax.tools.StandardLocation;
import java.io.File;
import java.io.IOException;
import java.io.OutputStream;
import java.io.Serializable;
import java.util.*;
import java.util.stream.Collectors;
import java.util.stream.Stream;

/**
 * <pre>
 *     &lt;dependency&gt;
 *        &lt;groupId&gt;com.thoughtworks.qdox&lt;/groupId&gt;
 *        &lt;artifactId&gt;qdox&lt;/artifactId&gt;
 *        &lt;scope&gt;provided&lt;/scope&gt;
 *     &lt;/dependency&gt;
 * </pre>
 *
 * @author coyzeng@gmail.com
 */
@SPI(value = "javadoc", meta = JavaDoc.class)
public class JavadocPlugin implements CompilePlugin {

    private final List<String> processed = new ArrayList<>();

    @Override
    public void proceed(CompileContext context) throws Throwable {
        JavaDoc javadoc = context.annotatedElement().getAnnotation(JavaDoc.class);
        try (OutputStream stream = openSchema(context, javadoc)) {
            Codec codec = ServiceLoader.load(Codec.class).get(Codec.JSON);
            List<Kind> documents = scanAST(context).flatMap(e -> {
                        List<Element> paths = new ArrayList<>();
                        Deque<Element> queue = new ArrayDeque<>();
                        queue.push(e);
                        while (!queue.isEmpty()) {
                            e = queue.pop();
                            if (!this.canJavadoc(e)) {
                                continue;
                            }
                            if (e.getEnclosedElements().stream().anyMatch(this::canJavadoc)) {
                                e.getEnclosedElements().forEach(queue::push);
                            }
                            processed.add(e.toString());
                            paths.add(e);
                        }
                        return paths.stream();
                    })
                    .map(e -> (TypeElement) e)
                    .filter(this::shouldJavadoc)
                    .flatMap(doc -> Parser.parse(context, javadoc, doc).stream())
                    .distinct()
                    .collect(Collectors.toList());
            Map<String, Object> schema = new HashMap<>(2);
            schema.put("v", javadoc.version());
            schema.put("d", documents);
            schema.put("p", javadoc.value());
            stream.write(codec.encode(schema).array());
        } catch (Throwable e) {
            error(e, e.getMessage());
        }
    }

    private Stream<Element> scanAST(CompileContext context) {
        List<String> excludes = Arrays.asList(context.annotatedElement().getAnnotation(JavaDoc.class).exclude());
        List<String> includes = Arrays.asList(context.annotatedElement().getAnnotation(JavaDoc.class).value());
        return context.roundEnvironment().getRootElements().stream()
                .filter(this::canJavadoc).map(e -> (Element) e).filter(e -> {
                    if (includes.isEmpty() && excludes.isEmpty()) {
                        return true;
                    }
                    for (String exclude : excludes) {
                        if (Tool.startWith(e.asType().toString(), exclude)) {
                            return false;
                        }
                    }
                    for (String include : includes) {
                        if (Tool.startWith(e.asType().toString(), include)) {
                            return true;
                        }
                    }
                    return includes.isEmpty();
                });
    }

    private OutputStream openSchema(CompileContext context, JavaDoc doc) throws IOException {
        String location = String.format("META-INF%smesh%s%s.schema", File.separator, File.separator, Tool.anyone(doc.name(), Tool.MESH_NAME.get()));
        return context.processingEnvironment().getFiler().createResource(StandardLocation.CLASS_OUTPUT, "", location).openOutputStream();
    }

    private boolean canJavadoc(Element e) {
        return !processed.contains(e.toString()) && (e.getKind() == ElementKind.CLASS || e.getKind() == ElementKind.INTERFACE || e.getKind() == ElementKind.ENUM || e.getKind() == ElementKind.ANNOTATION_TYPE);
    }

    private boolean shouldJavadoc(TypeElement x) {
        if (null != x.getAnnotation(Exclude.class)) {
            return false;
        }
        return x.getSuperclass().toString().equals(Object.class.getCanonicalName()) ||
                x.getInterfaces().stream().anyMatch(v -> v.toString().equals(Serializable.class.getCanonicalName())) ||
                null != x.getAnnotation(Data.class) ||
                x.getEnclosedElements().stream().anyMatch(v -> null != v.getAnnotation(MPI.class));
    }
}
