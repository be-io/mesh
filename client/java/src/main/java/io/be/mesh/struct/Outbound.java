/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.struct;

import io.be.mesh.macro.Index;
import io.be.mesh.mpc.Returns;
import lombok.Data;

import java.io.Serializable;

/**
 * @author coyzeng@gmail.com
 */
@Data
public class Outbound implements Serializable, Returns {

    private static final long serialVersionUID = 4142652099532835971L;

    @Index(0)
    private String code;

    @Index(5)
    private String message;

    @Index(10)
    private Cause cause;

    @Index(15)
    private Object content;

}
