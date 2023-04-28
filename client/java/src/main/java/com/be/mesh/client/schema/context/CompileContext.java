/*
 * Copyright (c) 2000, 2023, ducesoft and/or its affiliates. All rights reserved.
 * DUCESOFT PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.schema.context;

import com.be.mesh.client.schema.compiler.JavaCompiler;

import javax.annotation.processing.ProcessingEnvironment;
import javax.annotation.processing.RoundEnvironment;
import javax.lang.model.element.Element;
import javax.lang.model.element.TypeElement;
import javax.lang.model.type.TypeMirror;
import javax.lang.model.util.Elements;
import javax.tools.DiagnosticCollector;
import javax.tools.JavaFileManager.Location;
import javax.tools.JavaFileObject;
import javax.tools.StandardJavaFileManager;
import java.util.Optional;

/**
 * @author coyzeng@gmail.com
 */
public interface CompileContext extends AutoCloseable, CompilePrinter {

    JavaCompiler compiler();

    Elements utilities();

    <A extends Element> A annotatedElement();

    <A extends TypeElement> A annotationElement();

    <A extends RoundEnvironment> A roundEnvironment();

    ProcessingEnvironment processingEnvironment();

    StandardJavaFileManager fileManager();

    DiagnosticCollector<JavaFileObject> diagnostics();

    boolean isAssignable(TypeMirror inter, TypeMirror mirror);

    /**
     * <code>
     * JavaFileObject file       = processingEnvironment().getFiler().createSourceFile("ducer");
     * Path           sourcePath = Paths.get(file.toUri()).resolve("src");
     * </code>
     */
    Optional<JavaFileObject> getObject(Location location, String relativePath);

    String getQualifiedName(Element element);

    JavaCompiler getCompiler();
}
