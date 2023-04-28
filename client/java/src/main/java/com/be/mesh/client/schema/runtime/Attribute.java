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
 * Service schema attributes.
 *
 * @author coyzeng@gmail.com
 */
@Data
public class Attribute implements Serializable {

    private static final long serialVersionUID = -6286387922204567692L;
    /**
     * Attribute kind.
     */
    @Index(0)
    private String kind;
    /**
     * Attribute index.
     */
    @Index(1)
    private int index;
    /**
     * Attribute name;
     */
    @Index(2)
    private String name;
    /**
     * Attribute alias.
     */
    @Index(3)
    private String alias;
    /**
     * Attribute comment.
     */
    @Index(4)
    private String comment;
    /**
     * Attribute is optional.
     */
    @Index(5)
    private boolean optional;
    /**
     * Attribute flags. input or output etc.
     */
    @Index(6)
    private int flags;
    /**
     * SubAttributes
     */
    @Index(7)
    private List<Attribute> attributes;

}
