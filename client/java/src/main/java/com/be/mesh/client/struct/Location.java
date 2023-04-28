/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.struct;

import com.be.mesh.client.annotate.Index;
import lombok.Data;
import lombok.EqualsAndHashCode;

/**
 * @author coyzeng@gmail.com
 */
@Data
@EqualsAndHashCode(callSuper = true)
public class Location extends Principal {

    private static final long serialVersionUID = 4059350185355979983L;

    @Index(10)
    private String ip;
    @Index(15)
    private String port;
    @Index(20)
    private String host;
    @Index(25)
    private String name;
}
