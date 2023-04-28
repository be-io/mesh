/*
 * Copyright (c) 2000, 2023, ducesoft and/or its affiliates. All rights reserved.
 * DUCESOFT PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.schema.context;

import javax.tools.ForwardingJavaFileManager;
import javax.tools.JavaFileObject;
import javax.tools.JavaFileObject.Kind;
import javax.tools.StandardJavaFileManager;
import java.io.File;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.*;
import java.util.stream.Collectors;
import java.util.stream.Stream;
import java.util.stream.StreamSupport;

/**
 * @author coyzeng@gmail.com
 */
public class JavacFileManager extends ForwardingJavaFileManager<StandardJavaFileManager> implements StandardJavaFileManager {

    private static final String SOURCE_PATH = "src.main.java.";
    private static final String DOT = "\\.";

    protected JavacFileManager(StandardJavaFileManager fileManager) {
        super(fileManager);
    }

    @Override
    public Iterable<JavaFileObject> list(Location location, String pack, Set<Kind> kinds, boolean recurse) throws
            IOException {
        return Optional
                .ofNullable(super.list(location, pack, kinds, recurse))
                .map(fos -> StreamSupport.stream(fos.spliterator(), false))
                .orElseGet(() -> Optional
                        .ofNullable(kinds)
                        .orElseGet(() -> new HashSet<>(Arrays.asList(Kind.values())))
                        .stream()
                        .flatMap(kind -> {
                            Path root = Paths.get((SOURCE_PATH + pack).replaceAll(DOT, File.separator));
                            return findFiles(root, kind).stream().map(JavacFileObject::new);
                        }))
                .collect(Collectors.toList());
    }

    @Override
    public Iterable<? extends JavaFileObject> getJavaFileObjectsFromFiles(Iterable<? extends File> files) {
        return Optional
                .ofNullable(files)
                .map(v -> StreamSupport.stream(v.spliterator(), false))
                .orElseGet(Stream::empty)
                .map(v -> new JavacFileObject(v.toPath()))
                .collect(Collectors.toList());
    }

    @Override
    public Iterable<? extends JavaFileObject> getJavaFileObjects(File... files) {
        return Optional
                .ofNullable(files).map(Arrays::stream).orElseGet(Stream::empty)
                .map(v -> new JavacFileObject(v.toPath()))
                .collect(Collectors.toList());
    }

    @Override
    public Iterable<? extends JavaFileObject> getJavaFileObjectsFromStrings(Iterable<String> names) {
        return null;
    }

    @Override
    public Iterable<? extends JavaFileObject> getJavaFileObjects(String... names) {
        return null;
    }

    @Override
    public void setLocation(Location location, Iterable<? extends File> path) throws IOException {

    }

    @Override
    public Iterable<? extends File> getLocation(Location location) {
        return findFiles(Paths.get(System.getProperty(location.getName())), Kind.OTHER)
                .stream()
                .map(Path::toFile)
                .collect(Collectors.toList());
    }

    private List<Path> findFiles(Path root, Kind kind) {
        try {
            List<Path> paths = new ArrayList<>();
            Deque<Path> queue = new ArrayDeque<>();
            queue.push(root);
            while (!queue.isEmpty()) {
                root = queue.pop();
                if (root.toFile().isDirectory()) {
                    Files.list(root).forEach(queue::push);
                    continue;
                }
                if (root.toString().endsWith(kind.extension)) {
                    paths.add(root);
                }
            }
            return paths;
        } catch (IOException e) {
            throw new CompileException("Compile java sources failed.", e);
        }
    }
}
