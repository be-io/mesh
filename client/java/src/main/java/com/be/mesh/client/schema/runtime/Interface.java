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
 * Service schema.
 *
 * @author coyzeng@gmail.com
 */
@Data
public class Interface implements Serializable {

    private static final long serialVersionUID = -3147616045858009908L;
    /**
     * Interface version.
     */
    @Index(0)
    private String version;
    /**
     * Interface kind.
     */
    @Index(1)
    private String kind;
    /**
     * Interface name.
     */
    @Index(2)
    private String name;
    /**
     * Interface alias.
     */
    @Index(3)
    private String alias;
    /**
     * Interface comment.
     */
    @Index(4)
    private String comment;
    /**
     * Interface attributes.
     */
    @Index(5)
    private List<Function> functions;

}
