/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import lombok.extern.slf4j.Slf4j;
import org.testng.annotations.Test;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
public class ContextTest {

    @Test
    public void testSpan() throws Exception {
        Mesh.contextSafeUncheck(() -> {
            log.info("{},{}", Mesh.context().getTraceId(), Mesh.context().getSpanId());
            Mesh.contextSafeUncheck(() -> {
                log.info("{},{}", Mesh.context().getTraceId(), Mesh.context().getSpanId());
                Mesh.contextSafeUncheck(() -> {
                    log.info("{},{}", Mesh.context().getTraceId(), Mesh.context().getSpanId());
                    Mesh.contextSafeUncheck(() -> {
                        log.info("{},{}", Mesh.context().getTraceId(), Mesh.context().getSpanId());
                    });
                    Mesh.contextSafeUncheck(() -> {
                        log.info("{},{}", Mesh.context().getTraceId(), Mesh.context().getSpanId());
                    });
                    Mesh.contextSafeUncheck(() -> {
                        log.info("{},{}", Mesh.context().getTraceId(), Mesh.context().getSpanId());
                    });
                });
            });
            Mesh.contextSafeUncheck(() -> {
                log.info("{},{}", Mesh.context().getTraceId(), Mesh.context().getSpanId());
            });
            Mesh.contextSafeUncheck(() -> {
                log.info("{},{}", Mesh.context().getTraceId(), Mesh.context().getSpanId());
            });
        });
    }
}
