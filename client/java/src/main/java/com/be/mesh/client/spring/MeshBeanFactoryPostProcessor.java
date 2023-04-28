/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.spring;

import com.be.mesh.client.tool.Tool;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.BeansException;
import org.springframework.beans.factory.config.BeanFactoryPostProcessor;
import org.springframework.beans.factory.config.ConfigurableListableBeanFactory;
import org.springframework.beans.factory.support.AbstractBeanDefinition;
import org.springframework.beans.factory.support.BeanDefinitionBuilder;
import org.springframework.beans.factory.support.BeanDefinitionRegistry;
import org.springframework.beans.factory.support.BeanDefinitionRegistryPostProcessor;
import org.springframework.core.Ordered;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
public class MeshBeanFactoryPostProcessor implements BeanFactoryPostProcessor, Ordered, BeanDefinitionRegistryPostProcessor {

    @Override
    public void postProcessBeanFactory(ConfigurableListableBeanFactory factory) throws BeansException {
        //
    }

    @Override
    public int getOrder() {
        return Ordered.HIGHEST_PRECEDENCE;
    }

    @Override
    public void postProcessBeanDefinitionRegistry(BeanDefinitionRegistry registry) throws BeansException {
        //
    }

    public void registerBeanDefinition(Class<?> type, BeanDefinitionRegistry registry) {
        String name = Tool.formatObjectName(type);
        if (registry.containsBeanDefinition(name)) {
            log.info("The mesh bean definition [name : {}, class : {}] has been registered.", name, type.getName());
            return;
        }
        BeanDefinitionBuilder builder = BeanDefinitionBuilder.rootBeanDefinition(type);
        AbstractBeanDefinition beanDefinition = builder.getBeanDefinition();
        registry.registerBeanDefinition(name, beanDefinition);
    }

}
