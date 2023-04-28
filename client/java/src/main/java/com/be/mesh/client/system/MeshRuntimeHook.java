/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.system;

import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.mpc.ServiceLoader;
import com.be.mesh.client.prsim.RuntimeHook;
import com.be.mesh.client.tool.Tool;
import lombok.extern.slf4j.Slf4j;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
@SPI("mesh")
public class MeshRuntimeHook implements RuntimeHook {

    @Override
    public void start() throws Throwable {
        refresh();
    }

    @Override
    public void stop() throws Throwable {
        //
    }

    @Override
    public void refresh() throws Throwable {
        Runtime.getRuntime().addShutdownHook(new Thread(() -> ServiceLoader.load(RuntimeHook.class).list().forEach(hook -> {
            try {
                Tool.uncheck(hook::stop);
            } catch (Throwable e) {
                log.error("Shutdown hook exec with error. ", e);
            }
        })));
    }
}
