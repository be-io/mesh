/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.struct;

import com.be.mesh.client.annotate.Index;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.io.Serializable;

/**
 * Any fixed information of principal.
 *
 * @author coyzeng@gmail.com
 */
@Data
@NoArgsConstructor
@AllArgsConstructor
public class Principal implements Serializable {

    private static final long serialVersionUID = 6758247984696267747L;

    /**
     * Event own node id.
     */
    @Index(value = 0, name = "node_id")
    private String nodeId;
    /**
     * Event own top institution id.
     */
    @Index(value = 5, name = "inst_id")
    private String instId;

    public Principal(String nodeId) {
        this.nodeId = nodeId;
        this.instId = nodeId;
    }

    /**
     * Of instId.
     *
     * @param instId inst id.
     * @return principal
     */
    public static Principal ofInstId(String instId) {
        return new Principal("", instId);
    }

}
