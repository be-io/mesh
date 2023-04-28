/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.spring;

import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.mpc.Factory;
import com.be.mesh.client.tool.Tool;

import java.io.ByteArrayInputStream;
import java.io.InputStream;
import java.nio.charset.StandardCharsets;
import java.util.Objects;
import java.util.stream.Stream;

/**
 * Spring application extension factory.
 *
 * @author coyzeng@gmail.com
 */
@SPI("spring")
public class MeshSpringExtensionFactory implements Factory {

    private boolean notSpring() {
        return !Tool.isClassPresent("org.springframework.context.ApplicationContext");
    }

    @Override
    public <T> T getProvider(String key, Class<T> type) {
        if (notSpring()) {
            return null;
        }
        return MeshSpringApplicationContext.CTX.stream().map(c -> {
            try {
                return c.getBean(key, type);
            } catch (Exception e) {
                return null;
            }
        }).filter(Objects::nonNull).findFirst().orElse(null);
    }

    @Override
    public <T> Stream<T> getProvider(Class<T> spi) {
        if (notSpring()) {
            return null;
        }
        return MeshSpringApplicationContext.CTX.stream().map(c -> {
            try {
                return c.getBeansOfType(spi);
            } catch (Exception e) {
                return null;
            }
        }).filter(Objects::nonNull).map(x -> x.values().stream()).findFirst().orElseGet(Stream::empty);
    }

    @Override
    public InputStream getResource(String name) {
        if (notSpring()) {
            return new ByteArrayInputStream("".getBytes(StandardCharsets.UTF_8));
        }
        String value = MeshSpringApplicationContext.CTX.stream().
                map(org.springframework.context.ConfigurableApplicationContext::getEnvironment).
                map(x -> x.getProperty(name)).filter(Tool::required).findFirst().orElse("");
        return new ByteArrayInputStream(value.getBytes(StandardCharsets.UTF_8));
    }

}
