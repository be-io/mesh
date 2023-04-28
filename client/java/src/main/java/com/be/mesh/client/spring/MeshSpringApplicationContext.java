/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.spring;

import org.springframework.beans.BeansException;
import org.springframework.context.ApplicationContext;
import org.springframework.context.ApplicationContextAware;
import org.springframework.context.ApplicationListener;
import org.springframework.context.ConfigurableApplicationContext;
import org.springframework.context.event.ApplicationContextEvent;
import org.springframework.lang.NonNull;

import java.util.ArrayList;
import java.util.List;

/**
 * @author coyzeng@gmail.com
 */
public class MeshSpringApplicationContext implements ApplicationContextAware, ApplicationListener<ApplicationContextEvent> {

    static final List<ConfigurableApplicationContext> CTX = new ArrayList<>();

    @Override
    public void setApplicationContext(@NonNull ApplicationContext ctx) throws BeansException {
        if (!CTX.contains(ctx)) {
            CTX.add((ConfigurableApplicationContext) ctx);
        }
    }

    @Override
    public void onApplicationEvent(ApplicationContextEvent event) {
        setApplicationContext(event.getApplicationContext());
    }

}
