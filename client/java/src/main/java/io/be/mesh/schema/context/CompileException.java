/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.schema.context;

import io.be.mesh.cause.MeshException;

/**
 * @author coyzeng@gmail.com
 */
public class CompileException extends MeshException {

    public CompileException(String message, Object... args) {
        super(message, args);
    }

    public CompileException(String message, Throwable e, Object... args) {
        super(message, e, args);
    }

    public CompileException(Throwable e) {
        super(e);
    }
}
