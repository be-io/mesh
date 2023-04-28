/*
 * Copyright (c) 2000, 2023, ducesoft and/or its affiliates. All rights reserved.
 * DUCESOFT PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.schema.runtime;

import lombok.Data;

import java.io.Serializable;
import java.util.List;

/**
 * @author coyzeng@gmail.com
 */
@Data
public class TypeStruct implements Serializable {

    private static final long serialVersionUID = 1562967895064798196L;

    private String command;

    private String fullName;

    private String aliasName;

    private String method;

    private String version;

    private String comment;

    private String classComment;

    private String author;

    private List<TypeParameter> input;

    private List<TypeParameter> output;
}
