/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.spring;

import com.be.mesh.client.annotate.Binding;
import com.be.mesh.client.annotate.Bindings;
import com.be.mesh.client.annotate.MPS;
import com.be.mesh.client.mpc.Eden;
import com.be.mesh.client.mpc.ServiceLoader;
import org.springframework.beans.BeansException;
import org.springframework.beans.factory.config.SmartInstantiationAwareBeanPostProcessor;

import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

/**
 * @author coyzeng@gmail.com
 */
public class MeshProviderBeanPostProcessor implements SmartInstantiationAwareBeanPostProcessor {

    private final Map<String, Class<?>> services = new ConcurrentHashMap<>();

    @Override
    public Object postProcessBeforeInstantiation(Class<?> type, String name) throws BeansException {
        if (type.isAnnotationPresent(MPS.class) || type.isAnnotationPresent(Bindings.class) || type.isAnnotationPresent(Binding.class)) {
            services.put(name, type);
        }
        return SmartInstantiationAwareBeanPostProcessor.super.postProcessBeforeInstantiation(type, name);
    }

    @Override
    public Object postProcessAfterInitialization(Object bean, String name) throws BeansException {
        if (null != services.get(name)) {
            Eden eden = ServiceLoader.load(Eden.class).getDefault();
            eden.store(services.get(name), bean);
            services.remove(name);
        }
        return bean;
    }

}
