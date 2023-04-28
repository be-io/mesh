/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import lombok.Data;

import java.util.Map;

/**
 * @author coyzeng@gmail.com
 */
@Data
public class ServiceInvocation implements Invocation {

    private static final long serialVersionUID = -176271661042439336L;
    private transient Invoker<?> proxy;
    private transient Execution<?> execution;
    private transient Inspector inspector;
    private transient Parameters parameters;
    private transient URN urn;

    @Override
    public Object[] getArguments() {
        return this.parameters.arguments();
    }

    @Override
    public Map<String, String> getAttachments() {
        return this.parameters.attachments();
    }

    @SuppressWarnings("unchecked")
    @Override
    public <T> Execution<T> getExecution() {
        return (Execution<T>) this.execution;
    }

}
