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
public class Page<T> implements Serializable {

    private static final long serialVersionUID = -2185748030617065721L;
    @Index(0)
    private String sid;
    @Index(5)
    private long index;
    @Index(10)
    private long limit;
    @Index(15)
    private long total;
    @Index(20)
    private boolean next;
    @Index(25)
    private T data;

}
