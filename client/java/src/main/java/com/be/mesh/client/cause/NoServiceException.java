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
public class NoServiceException extends MeshException {

    public NoServiceException(Throwable e) {
        super(MeshCode.NO_SERVICE, e);
    }

    public NoServiceException(String message, Object... args) {
        super(MeshCode.NO_SERVICE, message, args);
    }
}
