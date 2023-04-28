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
public class CryptException extends MeshException {

    public CryptException(Throwable e) {
        super(MeshCode.CRYPT_ERROR, e);
    }

    public CryptException(Codeable codeable, Throwable e) {
        super(codeable, e);
    }

    public CryptException(String message, Object... args) {
        super(MeshCode.CRYPT_ERROR, message, args);
    }

}
