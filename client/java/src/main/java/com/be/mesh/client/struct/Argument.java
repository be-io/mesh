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
public class Argument implements Serializable {

    private static final long serialVersionUID = 1038399896960140230L;
    /**
     * Argument name.
     */
    @Index(0)
    private String name;
    /**
     * Argument type. type is a keywords for some language.
     */
    @Index(1)
    private String kind;
}
