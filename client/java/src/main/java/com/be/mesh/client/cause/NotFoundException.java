/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.cause;

import com.be.mesh.client.mpc.MeshCode;

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
