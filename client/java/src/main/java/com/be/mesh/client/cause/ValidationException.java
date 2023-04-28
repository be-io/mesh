/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.cause;

import com.be.mesh.client.mpc.MeshCode;
import com.be.mesh.client.prsim.Codeable;

/**
 * @author coyzeng@gmail.com
 */
public class ValidationException extends MeshException {

    public ValidationException(String message, Object... args) {
        super(MeshCode.VALIDATE.getCode(), message, args);
    }

    public ValidationException(Codeable code, String message, Object... args) {
        super(code, message, args);
    }

}
