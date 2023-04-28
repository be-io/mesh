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
public class Institution implements Serializable {

    private static final long serialVersionUID = -2794109194935386598L;
    @Index(value = 0, name = "node_id")
    private String nodeId;
    @Index(value = 5, name = "inst_id")
    private String instId;
    @Index(value = 10, name = "inst_name")
    private String instName;
    @Index(value = 15)
    private int status;

}
