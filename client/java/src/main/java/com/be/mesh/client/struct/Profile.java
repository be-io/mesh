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
public class Profile implements Serializable {

    private static final long serialVersionUID = 4142940169613175696L;
    @Index(value = 0, name = "data_id")
    private String dataId;
    @Index(1)
    private String content;
}
