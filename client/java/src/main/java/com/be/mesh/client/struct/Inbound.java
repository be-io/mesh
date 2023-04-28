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
import java.util.Map;

/**
 * @author coyzeng@gmail.com
 */
@Data
public class Inbound implements Serializable {

    private static final long serialVersionUID = -7034013150268079840L;
    @Index(0)
    private Object[] arguments;
    @Index(1)
    private Map<String, String> attachments;
}
