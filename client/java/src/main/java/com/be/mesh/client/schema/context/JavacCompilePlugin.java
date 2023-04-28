/*
 * Copyright (c) 2000, 2023, ducesoft and/or its affiliates. All rights reserved.
 * DUCESOFT PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.schema.context;

import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.mpc.ServiceLoader;
import com.be.mesh.client.schema.CompilePlugin;
import com.be.mesh.client.tool.Once;

import javax.tools.*;
import java.io.StringWriter;
import java.io.Writer;
import java.lang.annotation.Annotation;
import java.util.*;
import java.util.stream.Collectors;
import java.util.stream.StreamSupport;

/**
 * @author coyzeng@gmail.com
 */
@SPI(value = "*", meta = Object.class)
public class JavacCompilePlugin implements CompilePlugin {

    private final Once<Map<String, CompilePlugin>> plugins = Once.with(() -> ServiceLoader.load(CompilePlugin.class).map());

    @Override
    public <A extends Annotation> List<Class<A>> with() {
        return plugins.get().values().stream().flatMap(plugin -> {
            List<Class<A>> ns = plugin.with();
            return ns.stream();
        }).distinct().collect(Collectors.toList());
    }

    @Override
    public void proceed(CompileContext context) {
        plugins.get()
                .values()
                .stream()
                .filter(plugin -> plugin.with().stream().anyMatch(a -> a.getName().equals(context.annotationElement().
                        getQualifiedName().toString()))).forEach(plugin -> {
                    try {
                        plugin.proceed(context);
                    } catch (Throwable e) {
                        error(e, String.format("Compiler plugin %s executed failed!", plugin.getClass().getName()));
                        context.error(e, e.getMessage());
                    }
                });
    }

    public static void compile(CompileContext context) {
        try {
            Iterable<JavaFileObject> fileObjects = context.fileManager().list(StandardLocation.SOURCE_PATH,
                    "*", new HashSet<>(Collections.singletonList(JavaFileObject.Kind.SOURCE)), true);
            if (null == fileObjects) {
                return;
            }
            Iterable<? extends JavaFileObject> compilationUnits = StreamSupport
                    .stream(fileObjects.spliterator(), false)
                    .filter(fileObject -> context
                            .roundEnvironment()
                            .getRootElements()
                            .stream()
                            .anyMatch(e -> fileObject.getName().contentEquals(e.getSimpleName())))
                    .collect(Collectors.toList());
            List<String> options = context
                    .processingEnvironment()
                    .getOptions()
                    .entrySet()
                    .stream()
                    .map(kv -> Arrays.asList(kv.getKey(), kv.getValue()))
                    .flatMap(Collection::stream)
                    .collect(Collectors.toList());
            Writer out = new StringWriter();
            DiagnosticCollector<JavaFileObject> diagnostics = new DiagnosticCollector<>();
            JavaCompiler.CompilationTask task = ToolProvider.getSystemJavaCompiler().getTask(out,
                    context.fileManager(),
                    diagnostics,
                    options,
                    Collections.emptyList(),
                    compilationUnits);
            task.setLocale(Locale.getDefault());
            task.call();
            for (Diagnostic<? extends JavaFileObject> diagnostic : diagnostics.getDiagnostics()) {
                context.error(String.format("Error on line %d in %s with %s",
                        diagnostic.getLineNumber(),
                        diagnostic.getSource().toUri(),
                        diagnostic.getMessage(Locale.getDefault())));
            }
            context.info(out.toString());
        } catch (Exception e) {
            context.error(e, e.getMessage());
        }
    }
}
