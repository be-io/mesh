/*
 * Copyright (c) 2000, 2023, ducesoft and/or its affiliates. All rights reserved.
 * DUCESOFT PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.schema.context;

import com.be.mesh.client.cause.MeshException;

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
