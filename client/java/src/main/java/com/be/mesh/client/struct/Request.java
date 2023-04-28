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
 * Standard input payload.
 *
 * @author coyzeng@gmail.com
 */
@Data
public class Request implements Serializable {

    private static final long serialVersionUID = -1157070419360040418L;
    private String method;
    private Object content;
    private String version;

}
