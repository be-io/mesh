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
public class Script implements Serializable {

    private static final long serialVersionUID = 4149025267197383315L;
    public static final String EXPRESSION = "EXPRESSION";
    public static final String VALUE = "VALUE";
    public static final String SCRIPT = "SCRIPT";
    @Index(0)
    private String code;
    @Index(5)
    private String name;
    @Index(10)
    private String desc;
    /**
     * @see #EXPRESSION
     * @see #VALUE
     */
    @Index(15)
    private String kind;
    @Index(20)
    private String expr;
    @Index(25)
    private Map<String, String> attachment;
}
