/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.struct;

import com.be.mesh.client.annotate.Index;
import lombok.Data;

import java.io.Serializable;

/**
 * @author coyzeng@gmail.com
 */
@Data
public class Timeout implements Serializable {

    private static final long serialVersionUID = 6375030576735413376L;

    /**
     * Timeout task id.
     */
    @Index(value = 0, name = "task_id")
    private String taskId;

    /**
     * The timeout binding keys.
     */
    @Index(5)
    private Topic binding;

    /**
     * Timeout event status. &1=1 expired, &2=2 removed, &4=4 stopped.
     * <p>
     * If the timeout is stopped. If stopped will be pushed into stopped queue.
     */
    @Index(10)
    private long status;

    /**
     * Timeout initial time.
     */
    @Index(value = 15, name = "create_at")
    private long createAt;

    /**
     * Timeout expect time.
     */
    @Index(value = 20, name = "invoke_at")
    private long invokeAt;

    /**
     * Timeout payload entity
     */
    @Index(25)
    private Entity entity;

}
