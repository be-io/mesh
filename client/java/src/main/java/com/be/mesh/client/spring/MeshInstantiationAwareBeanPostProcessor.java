/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.spring;

import com.be.mesh.client.annotate.MPI;
import com.be.mesh.client.mpc.Eden;
import com.be.mesh.client.mpc.ServiceLoader;
import com.be.mesh.client.prsim.Routable;
import com.be.mesh.client.tool.Tool;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.BeansException;
import org.springframework.beans.factory.config.BeanFactoryPostProcessor;
import org.springframework.beans.factory.config.ConfigurableListableBeanFactory;
import org.springframework.beans.factory.config.ConstructorArgumentValues;
import org.springframework.beans.factory.config.InstantiationAwareBeanPostProcessor;
import org.springframework.beans.factory.support.AbstractBeanDefinition;
import org.springframework.beans.factory.support.DefaultListableBeanFactory;
import org.springframework.beans.factory.support.RootBeanDefinition;
import org.springframework.core.Ordered;
import org.springframework.core.ResolvableType;

import java.lang.reflect.Field;
import java.util.ArrayDeque;
import java.util.Deque;

/**
 * <pre>
 *     if (isRoutable && factory instanceof DefaultListableBeanFactory) {
 *          ResolvableType resolvableType = ResolvableType.forClassWithGenerics(Routable.class, referenceType);
 *          RootBeanDefinition routableDefinition = new RootBeanDefinition();
 *          routableDefinition.setTargetType(resolvableType);
 *          routableDefinition.setAutowireMode(AbstractBeanDefinition.AUTOWIRE_BY_TYPE);
 *          routableDefinition.setAutowireCandidate(true);
 *          ((DefaultListableBeanFactory) factory).registerBeanDefinition(String.format("%sRoutable", beanName), routableDefinition);
 *     }
 * </pre>
 *
 * @author coyzeng@gmail.com
 */
@Slf4j
public class MeshInstantiationAwareBeanPostProcessor implements InstantiationAwareBeanPostProcessor, Ordered, BeanFactoryPostProcessor {

    private ConfigurableListableBeanFactory factory;

    @Override
    public Object postProcessBeforeInstantiation(Class<?> type, String name) throws BeansException {
        Deque<Class<?>> types = new ArrayDeque<>();
        types.push(type);
        while (!types.isEmpty()) {
            Class<?> kind = types.pop();
            if (Tool.contains(kind.getName(), "org.springframework")) {
                continue;
            }
            if (null != kind.getSuperclass() && Object.class != kind.getSuperclass()) {
                types.push(kind.getSuperclass());
            }
            for (Field variable : kind.getDeclaredFields()) {
                if (!variable.isAnnotationPresent(MPI.class)) {
                    continue;
                }
                Class<?> referenceType = Tool.detectReferenceType(variable);
                Eden eden = ServiceLoader.load(Eden.class).getDefault();
                Object bean = eden.define(variable.getAnnotation(MPI.class), referenceType);
                String beanName = Tool.formatObjectName(referenceType);
                if (factory.containsBean(beanName) && referenceType.isAssignableFrom(factory.getBean(beanName).getClass())) {
                    continue;
                }
                if (factory.containsBean(beanName)) {
                    log.debug("Bean has been initialized already, please check the bean definition conflict with name {}.", beanName);
                    beanName = String.format("%s%d", beanName, referenceType.hashCode());
                    if (factory.containsBean(beanName)) {
                        continue;
                    }
                }
                factory.registerSingleton(beanName, bean);
                factory.registerResolvableDependency(referenceType, bean);

                if (factory instanceof DefaultListableBeanFactory) {
                    ConstructorArgumentValues arguments = new ConstructorArgumentValues();
                    arguments.addIndexedArgumentValue(0, bean);
                    ResolvableType resolvableType = ResolvableType.forClassWithGenerics(Routable.class, referenceType);
                    RootBeanDefinition definition = new RootBeanDefinition();
                    definition.setTargetType(resolvableType);
                    definition.setAutowireMode(AbstractBeanDefinition.AUTOWIRE_BY_TYPE);
                    definition.setAutowireCandidate(true);
                    definition.setBeanClass(MeshRoutableFactoryBean.class);
                    definition.setConstructorArgumentValues(arguments);
                    ((DefaultListableBeanFactory) factory).registerBeanDefinition(String.format("%sRoutable", beanName), definition);
                }
            }
        }
        return InstantiationAwareBeanPostProcessor.super.postProcessBeforeInstantiation(type, name);
    }

    @Override
    public Object postProcessBeforeInitialization(Object bean, String name) throws BeansException {
        // Any bean create manual must process by this processor callback
        this.postProcessBeforeInstantiation(bean.getClass(), name);
        return InstantiationAwareBeanPostProcessor.super.postProcessBeforeInitialization(bean, name);
    }

    @Override
    public int getOrder() {
        return Ordered.HIGHEST_PRECEDENCE;
    }

    @Override
    public void postProcessBeanFactory(ConfigurableListableBeanFactory factory) throws BeansException {
        this.factory = factory;
    }

}
