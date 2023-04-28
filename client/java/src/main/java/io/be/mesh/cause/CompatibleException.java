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
public class CompatibleException extends MeshException {

    private static final long serialVersionUID = -8321043549481473497L;

    public CompatibleException(Throwable e) {
        super(MeshCode.COMPATIBLE, e);
    }

    public CompatibleException(String message, Object... args) {
        super(MeshCode.COMPATIBLE, message, args);
    }

}
