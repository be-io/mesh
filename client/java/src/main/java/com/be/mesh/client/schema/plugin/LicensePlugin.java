/*
 * Copyright (c) 2000, 2023, ducesoft and/or its affiliates. All rights reserved.
 * DUCESOFT PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.schema.plugin;

import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.cause.MeshException;
import com.be.mesh.client.schema.CompilePlugin;
import com.be.mesh.client.schema.context.CompileContext;
import com.be.mesh.client.schema.macro.License;
import com.be.mesh.client.tool.Tool;

import javax.lang.model.element.Element;
import javax.lang.model.element.ElementKind;
import javax.lang.model.element.QualifiedNameable;
import javax.tools.JavaFileObject;
import javax.tools.StandardLocation;
import java.io.File;
import java.io.InputStream;
import java.util.ArrayList;
import java.util.List;
import java.util.Optional;

/**
 * @author coyzeng@gmail.com
 */
@SPI(value = "license", meta = License.class)
public class LicensePlugin implements CompilePlugin {

    private final List<String> licensed = new ArrayList<>();

    @Override
    public void proceed(CompileContext context) {
        if (!(context.annotatedElement() instanceof QualifiedNameable)) {
            return;
        }
        String qualifiedName = ((QualifiedNameable) context.annotatedElement()).getQualifiedName().toString();
        if (licensed.stream().anyMatch(name -> name.contains(qualifiedName))) {
            return;
        }
        licensed.add(qualifiedName);
        String license = context.annotatedElement().getAnnotation(License.class).value();
        if (Tool.optional(license)) {
            return;
        }
        context.
                roundEnvironment()
                .getRootElements()
                .parallelStream()
                .filter(this::canLicense)
                .filter(e -> e instanceof QualifiedNameable)
                .map(e -> (QualifiedNameable) e)
                .forEach(element -> {
                    String path = context.getQualifiedName(element).replace(".", File.separator);
                    //
                    Optional<JavaFileObject> object = context.getObject(StandardLocation.SOURCE_PATH, path);
                    object.ifPresent(that -> {
                        try {
                            try (InputStream stream = that.openInputStream()) {
                                Tool.read(stream);
                            }
                        } catch (RuntimeException e) {
                            throw e;
                        } catch (Exception e) {
                            throw new MeshException(e);
                        }
                    });
                });
    }

    private boolean canLicense(Element e) {
        return e.getKind() == ElementKind.CLASS ||
                e.getKind() == ElementKind.INTERFACE ||
                e.getKind() == ElementKind.ENUM ||
                //e.getKind() == ElementKind.MODULE ||
                e.getKind() == ElementKind.PACKAGE ||
                e.getKind() == ElementKind.ANNOTATION_TYPE;
    }
}
