/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.struct.Cause;

import java.io.Serializable;
import java.util.HashMap;

/**
 * @author coyzeng@gmail.com
 */
public class GenericReturns extends HashMap<String, Object> implements Returns, Serializable {

    private static final String CODE = "code";
    private static final String MESSAGE = "message";
    private static final String CONTENT = "content";
    private static final String CAUSE = "cause";

    @Override
    public String getCode() {
        if (null == this.get(CODE)) {
            return "";
        }
        return String.valueOf(this.get(CODE));
    }

    @Override
    public void setCode(String code) {
        this.put(CODE, code);
    }

    @Override
    public String getMessage() {
        if (null == this.get(MESSAGE)) {
            return "";
        }
        return String.valueOf(this.get(MESSAGE));
    }

    @Override
    public void setMessage(String message) {
        this.put(MESSAGE, message);
    }

    @Override
    public Object getContent() {
        return this.get(CONTENT);
    }

    @Override
    public void setContent(Object content) {
        this.put(CONTENT, content);
    }

    @Override
    public Cause getCause() {
        if (null == this.get(CAUSE) || this.get(CAUSE) instanceof Cause) {
            return (Cause) this.get(CAUSE);
        }
        Codec codec = ServiceLoader.load(Codec.class).get(Codec.JSON);
        return codec.decode(codec.encode(this.get(CAUSE)), Types.of(Cause.class));
    }

    @Override
    public void setCause(Cause cause) {
        this.put(CAUSE, cause);
    }
}
