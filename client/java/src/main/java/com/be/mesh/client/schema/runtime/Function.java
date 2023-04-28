/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.schema.runtime;

import com.be.mesh.client.annotate.Index;
import lombok.Data;

import java.io.Serializable;
import java.util.List;

/**
 * @author coyzeng@gmail.com
 */
@Data
public class Function implements Serializable {

    private static final long serialVersionUID = -6806846808293887969L;
    /**
     * Function name.
     */
    @Index(0)
    private String name;
    /**
     * Function alias.
     */
    @Index(1)
    private String alias;
    /**
     * Function comment.
     */
    @Index(2)
    private String comment;
    /**
     * Function comment.
     */
    @Index(3)
    private List<Attribute> attributes;
}
