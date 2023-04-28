/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.mpc;

import io.be.mesh.macro.MPI;
import io.be.mesh.tool.Tool;

/**
 * @author coyzeng@gmail.com
 */
public class AnnotationHandler {

    public static final AnnotationHandler REF = new AnnotationHandler();

    @MPI
    private String hack;

    public MPI getDefaultMPI() {
        return Tool.uncheck(() -> AnnotationHandler.class.getDeclaredField("hack").getAnnotation(MPI.class));
    }

}
