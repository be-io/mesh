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
public class TimeoutException extends MeshException {

    public TimeoutException(Throwable e) {
        super(MeshCode.TIMEOUT, e);
    }

    public TimeoutException(String message, Object... args) {
        super(MeshCode.TIMEOUT, message, args);
    }
}
