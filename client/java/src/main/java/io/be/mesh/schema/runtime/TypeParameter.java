/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.schema.runtime;

import lombok.Data;

import java.io.Serializable;
import java.util.List;

/**
 * @author coyzeng@gmail.com
 */
@Data
public class TypeParameter implements Serializable {

    private static final long serialVersionUID = 5582403502774610497L;

    private String fullName;

    private String aliasName;

    private String name;

    private String type;

    private String comment;

    private boolean required;

    private boolean array;

    private List<TypeParameter> parameters;

}
