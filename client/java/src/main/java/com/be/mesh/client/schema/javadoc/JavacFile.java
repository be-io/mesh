/*
 * Copyright (c) 2000, 2023, ducesoft and/or its affiliates. All rights reserved.
 * DUCESOFT PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.schema.javadoc;

import com.be.mesh.client.cause.MeshException;
import com.be.mesh.client.schema.context.CompileContext;
import com.be.mesh.client.schema.macro.JavaDoc;
import com.be.mesh.client.tool.Tool;

import javax.lang.model.element.TypeElement;
import java.io.File;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.*;
import java.util.stream.Stream;

/**
 * @author coyzeng@gmail.com
 */
public class JavacFile {

    private static final String SOURCE_PATH = "src.main.java.";
    private static final String JAVA_SUFFIX = ".java";
    private static final String DOT_REGEX = ".";

    public static List<Path> findSources(CompileContext context) {
        JavaDoc javadoc = context.annotatedElement().getAnnotation(JavaDoc.class);
        String fullName = ((TypeElement) context.annotatedElement()).getQualifiedName().toString();
        String pack = Tool.anyone(Tool.anyone(javadoc.value()), fullName.substring(0, fullName.lastIndexOf(".")));
        Path path = Paths.get((SOURCE_PATH + pack).replace(DOT_REGEX, File.separator));
        if (!Files.exists(path)) {
            return Collections.emptyList();
        }
        return findSource(path);
    }

    private static List<Path> findSource(Path root) {
        try {
            List<Path> paths = new ArrayList<>();
            Deque<Path> queue = new ArrayDeque<>();
            queue.push(root);
            while (!queue.isEmpty()) {
                root = queue.pop();
                if (root.toFile().isDirectory()) {
                    try (Stream<Path> files = Files.list(root)) {
                        files.forEach(queue::push);
                    }
                    continue;
                }
                if (root.toString().endsWith(JAVA_SUFFIX)) {
                    paths.add(root);
                }
            }
            return paths;
        } catch (IOException e) {
            throw new MeshException("Compile java sources failed.", e);
        }
    }
}
