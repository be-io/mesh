/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.struct.workflow;

import com.be.mesh.client.annotate.Index;
import lombok.Data;

import java.io.Serializable;
import java.util.Map;

/**
 * @author coyzeng@gmail.com
 */
@Data
public class WorkIntent implements Serializable {

    private static final long serialVersionUID = -821661586102661015L;
    @Index(0)
    private String bno;
    @Index(1)
    private String cno;
    @Index(2)
    private Map<String, String> context;
    @Index(3)
    private Worker applier;
}
