/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.tools;

import lombok.extern.slf4j.Slf4j;
import org.testng.annotations.Test;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
public class LogTest {

    @Test
    public void testEscape() {
        log.info("{}-{}-{}", 1, 2, 3);
        log.info("{\"x\":\"y\"}");
        log.error("123456", new RuntimeException());
    }
}
