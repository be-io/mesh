/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.struct;

import lombok.Data;

import java.io.Serializable;

/**
 * * Standard output payload.
 *
 * @author coyzeng@gmail.com
 */
@Data
public class Response<T> implements Serializable {

    private static final long serialVersionUID = 7784141366669360075L;
    private String code;
    private T content;
    private String message;
    private String cause;
}
