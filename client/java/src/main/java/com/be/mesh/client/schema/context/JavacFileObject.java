/*
 * Copyright (c) 2000, 2023, ducesoft and/or its affiliates. All rights reserved.
 * DUCESOFT PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.schema.context;

import lombok.AllArgsConstructor;

import javax.lang.model.element.Modifier;
import javax.lang.model.element.NestingKind;
import javax.tools.JavaFileObject;
import java.io.*;
import java.net.URI;
import java.nio.charset.Charset;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.Arrays;

/**
 * @author coyzeng@gmail.com
 */
@AllArgsConstructor
public class JavacFileObject implements JavaFileObject {

    private final Path path;

    @Override
    public Kind getKind() {
        return Arrays
                .stream(Kind.values())
                .filter(kind -> this.path.endsWith(Paths.get(kind.extension)))
                .findFirst()
                .orElse(Kind.OTHER);
    }

    @Override
    public boolean isNameCompatible(String simpleName, Kind kind) {
        return true;
    }

    @Override
    public NestingKind getNestingKind() {
        return NestingKind.TOP_LEVEL;
    }

    @Override
    public Modifier getAccessLevel() {
        return Modifier.PUBLIC;
    }

    @Override
    public URI toUri() {
        return this.path.toUri();
    }

    @Override
    public String getName() {
        return this.path.getFileName().toString();
    }

    @Override
    public InputStream openInputStream() throws IOException {
        return this.path.toUri().toURL().openStream();
    }

    @Override
    public OutputStream openOutputStream() throws IOException {
        return this.path.toUri().toURL().openConnection().getOutputStream();
    }

    @Override
    public Reader openReader(boolean ignoreEncodingErrors) throws IOException {
        return new InputStreamReader(openInputStream());
    }

    @Override
    public CharSequence getCharContent(boolean ignoreEncodingErrors) throws IOException {
        return new String(Files.readAllBytes(this.path), Charset.defaultCharset());
    }

    @Override
    public Writer openWriter() throws IOException {
        return new OutputStreamWriter(openOutputStream());
    }

    @Override
    public long getLastModified() {
        return this.path.toFile().lastModified();
    }

    @Override
    public boolean delete() {
        return this.path.toFile().delete();
    }
}
