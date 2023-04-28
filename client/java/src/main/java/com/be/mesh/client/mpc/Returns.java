/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.struct.Cause;

/**
 * @author coyzeng@gmail.com
 */
public interface Returns {

    /**
     * Return code.
     */
    String getCode();

    /**
     * Return code.
     */
    void setCode(String code);

    /**
     * Return message.
     */
    String getMessage();

    /**
     * Return cause.
     */
    Cause getCause();

    /**
     * Return cause.
     */
    void setCause(Cause cause);

    /**
     * Return message.
     */
    void setMessage(String message);

    /**
     * Return content.
     */
    Object getContent();

    /**
     * Return content.
     */
    void setContent(Object content);

}
