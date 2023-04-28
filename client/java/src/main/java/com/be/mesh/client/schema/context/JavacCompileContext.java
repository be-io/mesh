/*
 * Copyright (c) 2000, 2023, ducesoft and/or its affiliates. All rights reserved.
 * DUCESOFT PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.schema.context;

import com.be.mesh.client.schema.compiler.JavaCompiler;
import com.be.mesh.client.schema.compiler.JdkCompiler;
import com.be.mesh.client.tool.Tool;

import javax.annotation.processing.ProcessingEnvironment;
import javax.annotation.processing.RoundEnvironment;
import javax.lang.model.element.Element;
import javax.lang.model.element.QualifiedNameable;
import javax.lang.model.element.TypeElement;
import javax.lang.model.type.TypeMirror;
import javax.lang.model.util.Elements;
import javax.tools.DiagnosticCollector;
import javax.tools.JavaFileManager.Location;
import javax.tools.JavaFileObject;
import javax.tools.StandardJavaFileManager;
import javax.tools.ToolProvider;
import java.lang.reflect.Method;
import java.nio.charset.Charset;
import java.util.Locale;
import java.util.Map;
import java.util.Optional;
import java.util.Set;
import java.util.concurrent.ConcurrentHashMap;

/**
 * @author coyzeng@gmail.com
 */
public class JavacCompileContext implements CompileContext {

    private final Map<Class<?>, Object> cached = new ConcurrentHashMap<>();
    private ProcessingEnvironment processingEnvironment;
    private RoundEnvironment environment;
    private TypeElement annotationElement;
    private Element annotatedElement;
    private JavaCompiler compiler;
    private StandardJavaFileManager fileManager;
    private DiagnosticCollector<JavaFileObject> diagnostics;

    public JavacCompileContext setProcessingEnvironment(ProcessingEnvironment processingEnvironment) {
        this.processingEnvironment = processingEnvironment;
        this.compiler = new JdkCompiler();
        this.diagnostics = new DiagnosticCollector<>();
        this.fileManager = ToolProvider.getSystemJavaCompiler().getStandardFileManager(this.diagnostics,
                Locale.getDefault(),
                Charset.defaultCharset());
        this.fileManager = new JavacFileManager(this.fileManager);
        return this;
    }

    public JavacCompileContext setEnvironment(RoundEnvironment environment) {
        this.environment = environment;
        return this;
    }

    public JavacCompileContext setAnnotationElement(TypeElement annotationElement) {
        this.annotationElement = annotationElement;
        return this;
    }

    public JavacCompileContext setAnnotatedElement(Element annotatedElement) {
        this.annotatedElement = annotatedElement;
        return this;
    }

    @Override
    public JavaCompiler compiler() {
        return this.compiler;
    }

    @Override
    public Elements utilities() {
        return processingEnvironment.getElementUtils();
    }

    @SuppressWarnings({"unchecked"})
    @Override
    public <A extends Element> A annotatedElement() {
        return (A) this.annotatedElement;
    }

    @SuppressWarnings({"unchecked"})
    @Override
    public <A extends TypeElement> A annotationElement() {
        return (A) this.annotationElement;
    }

    @SuppressWarnings({"unchecked"})
    @Override
    public <A extends RoundEnvironment> A roundEnvironment() {
        return (A) this.environment;
    }

    @Override
    public ProcessingEnvironment processingEnvironment() {
        return this.processingEnvironment;
    }

    @Override
    public StandardJavaFileManager fileManager() {
        return this.fileManager;
    }

    @Override
    public DiagnosticCollector<JavaFileObject> diagnostics() {
        return this.diagnostics;
    }

    @Override
    public boolean isAssignable(TypeMirror inter, TypeMirror mirror) {
        return this.processingEnvironment.getTypeUtils().isAssignable(inter, mirror);
    }

    @SuppressWarnings({"unchecked"})
    @Override
    public Optional<JavaFileObject> getObject(Location location, String relativePath) {
        Set<JavaFileObject> objects = (Set<JavaFileObject>) cached.computeIfAbsent(JavaFileObject.class, key -> {
            try {
                Class<?> argumentsClass = Class.forName("com.sun.tools.javac.main.Arguments");
                Class<?> contextClass = Class.forName("com.sun.tools.javac.util.Context");
                // JavacProcessingEnvironment
                Method getContext = Tool.getMethod(processingEnvironment, "getContext");
                Object context = getContext.invoke(processingEnvironment);
                Method instance = Tool.getMethod(argumentsClass, "instance", contextClass);
                Object arguments = instance.invoke(null, context);
                Method getFileObjects = Tool.getMethod(arguments, "getFileObjects");
                return getFileObjects.invoke(arguments);
            } catch (Exception e) {
                warn("Unable to determine source file path!");
                throw new CompileException(e);
            }
        });
        return objects.stream().filter(object -> object.toUri().getPath().contains(relativePath)).findFirst();
    }

    @Override
    public String getQualifiedName(Element element) {
        if (element instanceof QualifiedNameable) {
            return getPackageName((QualifiedNameable) element) + "." + getClassName((QualifiedNameable) element);
        }
        return element.getSimpleName().toString();
    }

    @Override
    public JavaCompiler getCompiler() {
        return this.compiler;
    }

    @Override
    public void close() throws Exception {
        fileManager.close();
    }

    private String getPackageName(QualifiedNameable element) {
        if (element.getKind().isClass() || element.getKind().isInterface()) {
            return element
                    .getQualifiedName()
                    .subSequence(0, element.getQualifiedName().length() - element.getSimpleName().length() - 1)
                    .toString();
        }
        return element.getQualifiedName().toString();
    }

    private String getClassName(QualifiedNameable element) {
        if (element.getKind().isClass() || element.getKind().isInterface()) {
            return element.getSimpleName().toString();
        }
        return "package-info";
    }

}
