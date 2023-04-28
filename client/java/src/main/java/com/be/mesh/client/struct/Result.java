/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.struct;

import com.be.mesh.client.mpc.MeshCode;
import com.be.mesh.client.prsim.Codeable;
import lombok.Data;

import java.io.Serializable;

/**
 * @author coyzeng@gmail.com
 */
@Data
public class Result implements Serializable {

    private static final long serialVersionUID = 5902320483315348521L;
    public static final Result SUCCESS = new Result();

    static {
        SUCCESS.setCode(MeshCode.SUCCESS.getCode());
        SUCCESS.setMessage(MeshCode.SUCCESS.getMessage());
    }

    /**
     * Event consumer principal information.
     */
    private Principal source;
    /**
     *
     */
    private Principal target;
    /**
     * Subscribe process code.
     */
    private String code;
    /**
     * Subscribe process message.
     */
    private String message;
    /**
     * Subscribe return value if synchronized.
     */
    private Object payload;

    /**
     * Fault result construct function.
     *
     * @param e cause
     * @return result
     */
    public static Result fault(Throwable e) {
        if (e instanceof Codeable) {
            return fault(((Codeable) e).getCode(), e.getMessage());
        }
        return fault(MeshCode.SYSTEM_ERROR.getCode(), e.getMessage());
    }

    /**
     * Fault result construct function.
     *
     * @param code    error code
     * @param message error message
     * @return result
     */
    public static Result fault(String code, String message) {
        Result result = new Result();
        result.setCode(code);
        result.setMessage(message);
        return result;
    }

    /**
     * Is pub ack.
     */
    public boolean ack() {
        return MeshCode.SUCCESS.is(getCode());
    }
}
