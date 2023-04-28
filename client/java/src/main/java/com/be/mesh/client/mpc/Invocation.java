/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import java.io.Serializable;
import java.util.Map;

/**
 * @author coyzeng@gmail.com
 */
public interface Invocation extends Serializable {

    /**
     * Get the delegate target object.
     */
    Invoker<?> getProxy();

    /**
     * Get the invocation inspector.
     *
     * @return inspector
     */
    Inspector getInspector();

    /**
     * Invoke parameters. include arguments and parameters.
     */
    Parameters getParameters();

    /**
     * Invoke parameters.
     */
    Object[] getArguments();

    /**
     * Get the attachments. The attributes will be serialized.
     *
     * @return attachments.
     */
    Map<String, String> getAttachments();

    /**
     * Get the invocation execution.
     */
    <T> Execution<T> getExecution();

    /**
     * Get the invoked urn.
     */
    URN getUrn();
}
