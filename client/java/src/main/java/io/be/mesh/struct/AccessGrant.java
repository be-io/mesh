/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.struct;

import io.be.mesh.macro.Index;
import lombok.Data;

import java.io.Serializable;

/**
 * @author coyzeng@gmail.com
 */
@Data
public class AccessGrant implements Serializable {

    private static final long serialVersionUID = 3234497595939081169L;
    @Index(0)
    private String code;

}
