/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.struct.workflow;

import lombok.AllArgsConstructor;
import lombok.Getter;

/**
 * @author coyzeng@gmail.com
 */
@Getter
@AllArgsConstructor
public enum WorkVertexKind {
    START(1), FINISH(2), JOB(4), MAN(8), WAY(16), TIMER(32),
    ;
    private final int code;

    public boolean is(long kind) {
        return (this.code & kind) == this.code;
    }
}
