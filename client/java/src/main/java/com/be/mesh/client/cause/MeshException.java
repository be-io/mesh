/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.cause;

import com.be.mesh.client.mpc.MeshCode;
import com.be.mesh.client.prsim.Codeable;
import com.be.mesh.client.tool.Tool;
import lombok.Getter;

/**
 * @author coyzeng@gmail.com
 */
public class MeshException extends RuntimeException implements Codeable {

    private static final long serialVersionUID = -1375097411218335673L;

    @Getter
    private final String code;
    private final String message;

    public MeshException(Throwable e) {
        super(e);
        this.code = MeshCode.SYSTEM_ERROR.getCode();
        this.message = MeshCode.SYSTEM_ERROR.getMessage();
    }

    public MeshException(Codeable code, Throwable e) {
        super(e);
        this.code = code.getCode();
        this.message = code.getMessage();
    }

    public MeshException(Codeable code, String message, Object... args) {
        super(format(message, args));
        this.code = code.getCode();
        this.message = code.getMessage();
    }

    public MeshException(String code, String message, Object... args) {
        super(format(message, args));
        this.code = code;
        this.message = super.getMessage();
    }

    public MeshException(String message, Object... args) {
        super(format(message, args));
        this.code = MeshCode.SYSTEM_ERROR.getCode();
        this.message = super.getMessage();
    }

    public MeshException(Throwable e, String message, Object... args) {
        super(format(message, args), e);
        this.code = MeshCode.SYSTEM_ERROR.getCode();
        this.message = super.getMessage();
    }

    @Override
    public String getMessage() {
        return Tool.anyone(this.message, super.getMessage());
    }

    public String getRootMessage() {
        return super.getMessage();
    }

    private static String format(String message, Object... args) {
        if (null == args || args.length < 1) {
            return message;
        }
        return String.format(message, args);
    }

    public boolean unavailable() {
        return MeshCode.TIMEOUT.is(this) || MeshCode.NET_UNAVAILABLE.is(this);
    }
}
