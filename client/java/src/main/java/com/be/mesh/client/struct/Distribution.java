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
 * Mesh Compute Center.
 *
 * @author coyzeng@gmail.com
 */
@Data
public class Distribution implements Serializable {

    private static final long serialVersionUID = 5581101152726626187L;
    /**
     * zone.
     */
    @Index(0)
    private String zone;
    /**
     * cluster.
     */
    @Index(5)
    private String cluster;
    /**
     * cell.
     */
    @Index(10)
    private String cell;
    /**
     * group.
     */
    @Index(15)
    private String group;
    /**
     * address.
     */
    @Index(20)
    private String address;
}
