/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.spring;

import com.be.mesh.client.mpc.ConsumerProperties;
import com.be.mesh.client.mpc.ProviderProperties;
import org.springframework.boot.autoconfigure.AutoConfiguration;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.ComponentScan;
import org.springframework.context.annotation.Configuration;

/**
 * @author coyzeng@gmail.com
 */
@Configuration
@AutoConfiguration
@ComponentScan("com.be.mesh.client")
public class MeshConfiguration {

    @Bean
    public static MeshBeanFactoryPostProcessor meshBeanFactoryPostProcessor() {
        return new MeshBeanFactoryPostProcessor();
    }

    @Bean
    public static MeshInstantiationAwareBeanPostProcessor meshInstantiationAwareBeanPostProcessor() {
        return new MeshInstantiationAwareBeanPostProcessor();
    }

    @Bean
    public static MeshConsumerBeanPostProcessor meshConsumerBeanPostProcessor() {
        return new MeshConsumerBeanPostProcessor();
    }

    @Bean
    public static MeshProviderBeanPostProcessor meshProviderBeanPostProcessor() {
        return new MeshProviderBeanPostProcessor();
    }

    @Bean
    public static MeshRuntimeSpringDriver meshRuntimeSpringDriver() {
        return new MeshRuntimeSpringDriver();
    }

    @Bean
    public static MeshSpringApplicationContext meshSpringApplicationContext() {
        return new MeshSpringApplicationContext();
    }

    @Bean
    public ConsumerProperties consumerProperties() {
        return new ConsumerProperties();
    }

    @Bean
    public ProviderProperties providerProperties() {
        return new ProviderProperties();
    }
}
