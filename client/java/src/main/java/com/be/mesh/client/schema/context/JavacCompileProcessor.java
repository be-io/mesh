/*
 * Copyright (c) 2000, 2023, ducesoft and/or its affiliates. All rights reserved.
 * DUCESOFT PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.schema.context;

import com.be.mesh.client.schema.CompilePlugin;
import com.be.mesh.client.tool.Tool;

import javax.annotation.processing.*;
import javax.lang.model.SourceVersion;
import javax.lang.model.element.TypeElement;
import java.util.Set;
import java.util.stream.Collectors;

/**
 * @author coyzeng@gmail.com
 */
@SupportedOptions({"mesh.schema.enable"})
@SupportedAnnotationTypes("*")
@SupportedSourceVersion(SourceVersion.RELEASE_8)
public class JavacCompileProcessor extends AbstractProcessor implements CompilePrinter {

    private CompilePlugin plugin;

    @Override
    public synchronized void init(ProcessingEnvironment processingEnv) {
        if (super.isInitialized()) {
            return;
        }
        plugin = new JavacCompilePlugin();
        super.init(processingEnv);
        Printers.WRITERS.add(new JavacPrinter(processingEnv));
        info("Mesh schema processor init!");
    }

    @Override
    public Set<String> getSupportedAnnotationTypes() {
        Set<String> types = plugin.with().stream().map(Class::getName).collect(Collectors.toSet());
        info("Mesh schema process %s", String.join(",", types));
        return types;
    }

    @Override
    public boolean process(Set<? extends TypeElement> annotations, RoundEnvironment roundEnv) {
        try {
            if ("false".equals(processingEnv.getOptions().getOrDefault("mesh.schema.enable", "true"))) {
                info("Mesh schema processor disabled!");
                return true;
            }
            if (roundEnv.processingOver()) {
                info("Mesh schema processor processed[{}]", roundEnv.processingOver());
                return true;
            }
            info("Mesh schema processor executing!");
            annotations.forEach(annotation -> {
                info("Mesh schema processor annotation[{}]", annotation.getQualifiedName());
                roundEnv.getElementsAnnotatedWith(annotation).forEach(annotated -> {
                    try (CompileContext context = new JavacCompileContext()
                            .setProcessingEnvironment(processingEnv)
                            .setEnvironment(roundEnv)
                            .setAnnotatedElement(annotated)
                            .setAnnotationElement(annotation)) {
                        info("Mesh schema processor annotated[{}]", annotated.asType().toString());
                        plugin.proceed(context);
                    } catch (Throwable e) {
                        error(Tool.getStackTrace(e));
                        error(e, String.format("Mesh schema processor annotated[%s]", annotated.asType().toString()));
                        error("Mesh schema processor " + e.getMessage() + "!");
                    }
                });
            });
            info("Mesh schema processor executed!");
            return true;
        } catch (Throwable e) {
            error(e, String.format("Mesh schema processor failed[%s]", e.getMessage()));
            error("Mesh schema processor " + e.getMessage() + "!");
            return false;
        }
    }

}
