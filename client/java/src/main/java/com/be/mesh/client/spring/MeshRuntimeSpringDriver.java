/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.spring;

import com.be.mesh.client.mpc.ServiceLoader;
import com.be.mesh.client.prsim.RuntimeHook;
import com.be.mesh.client.tool.Mode;
import com.be.mesh.client.tool.Tool;
import lombok.extern.slf4j.Slf4j;
import org.springframework.context.event.*;

/**
 * Lifecycle bean that automatically starts and stops the grpc server.
 *
 * @author coyzeng@gmail.com
 */
@Slf4j
public class MeshRuntimeSpringDriver {

    @EventListener(ContextStartedEvent.class)
    public void onStarted(ContextStartedEvent event) {
        if (!Mode.DISABLE.match(Tool.MESH_MODE.get())) {
            ServiceLoader.load(RuntimeHook.class).list().forEach(hook -> Tool.uncheck(hook::start));
        }
    }

    @EventListener(ContextStoppedEvent.class)
    public void onStop(ContextStoppedEvent event) {
        if (!Mode.DISABLE.match(Tool.MESH_MODE.get())) {
            ServiceLoader.load(RuntimeHook.class).list().forEach(hook -> Tool.uncheck(hook::stop));
        }
    }

    @EventListener(ContextRefreshedEvent.class)
    public synchronized void onRefresh(ContextRefreshedEvent event) {
        if (!Mode.DISABLE.match(Tool.MESH_MODE.get())) {
            ServiceLoader.load(RuntimeHook.class).list().forEach(hook -> Tool.uncheck(hook::refresh));
        }
    }

    @EventListener(ContextClosedEvent.class)
    public void onClose(ContextClosedEvent event) {
        if (!Mode.DISABLE.match(Tool.MESH_MODE.get())) {
            ServiceLoader.load(RuntimeHook.class).list().forEach(hook -> Tool.uncheck(hook::stop));
        }
    }


}
