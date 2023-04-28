/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.cause;

import io.be.mesh.mpc.MeshCode;

/**
 * @author coyzeng@gmail.com
 */
public class NotFoundException extends MeshException {

    public NotFoundException(Throwable e) {
        super(MeshCode.NOT_FOUND, e);
    }

    public NotFoundException(String message, Object... args) {
        super(MeshCode.NOT_FOUND, message, args);
    }
}
