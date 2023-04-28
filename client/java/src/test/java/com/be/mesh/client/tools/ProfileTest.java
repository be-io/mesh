/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.tools;

import com.be.mesh.client.tool.Tool;
import lombok.extern.slf4j.Slf4j;
import org.testng.annotations.Test;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
public class ProfileTest {

    @Test
    public void testRuntime() {
        log.info(Tool.MESH_RUNTIME.get().toString());
        Tool.MESH_RUNTIME.release();
        System.setProperty("mesh.runtime", "omega:6604");
        log.info(Tool.MESH_RUNTIME.get().toString());
        Tool.MESH_RUNTIME.release();
        System.setProperty("mesh.runtime", "omega");
        log.info(Tool.MESH_RUNTIME.get().toString());
        Tool.MESH_RUNTIME.release();
        System.setProperty("mesh.runtime", "127.0.0.1");
        log.info(Tool.MESH_RUNTIME.get().toString());
        Tool.MESH_RUNTIME.release();
        System.setProperty("mesh.runtime", "127.0.0.1:80");
        log.info(Tool.MESH_RUNTIME.get().toString());
        Tool.MESH_RUNTIME.release();
        System.setProperty("mesh.runtime", "https://127.0.0.1");
        log.info(Tool.MESH_RUNTIME.get().toString());
        Tool.MESH_RUNTIME.release();
        System.setProperty("mesh.runtime", "https://127.0.0.1:80");
        log.info(Tool.MESH_RUNTIME.get().toString());
    }
}
