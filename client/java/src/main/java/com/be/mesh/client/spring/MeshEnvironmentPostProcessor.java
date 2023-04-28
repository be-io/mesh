/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.spring;

import com.be.mesh.client.mpc.Codec;
import com.be.mesh.client.mpc.ServiceLoader;
import com.be.mesh.client.prsim.Network;
import lombok.extern.slf4j.Slf4j;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.env.EnvironmentPostProcessor;
import org.springframework.boot.env.PropertySourceLoader;
import org.springframework.boot.env.RandomValuePropertySource;
import org.springframework.boot.env.YamlPropertySourceLoader;
import org.springframework.core.env.ConfigurableEnvironment;
import org.springframework.core.env.MutablePropertySources;
import org.springframework.core.env.PropertySource;
import org.springframework.core.io.ByteArrayResource;

import java.io.IOException;
import java.util.List;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
public class MeshEnvironmentPostProcessor implements EnvironmentPostProcessor {

    private final PropertySourceLoader loader = new YamlPropertySourceLoader();

    @Override
    public void postProcessEnvironment(ConfigurableEnvironment environment, SpringApplication application) {
        // The binder is used to obtain the address of an object in the remote configuration file
//        Binder binder = Binder.get(environment);
//        String[] addrs = binder.bind("mesh.address", String[].class).orElse(new String[]{});
//        MutablePropertySources propertySources = environment.getPropertySources();
//        for (int index = 0; index < addrs.length; index++) {
//            // The name of the configuration file cannot be consistent. If it is consistent, it will be overwritten
//            try {
//                loadProperties("" + index, addrs[index], propertySources);
//            } catch (IOException e) {
//                log.error("Load fail, url is: " + addrs[index], e);
//            }
//        }
    }

    private void loadProperties(String name, String url, MutablePropertySources destination) throws IOException {
        Codec codec = ServiceLoader.load(Codec.class).getDefault();
        Network network = ServiceLoader.load(Network.class).getDefault();
        ByteArrayResource resource = new ByteArrayResource(codec.encode(network.getEnviron()).array());
        if (resource.exists()) {
            // If the resource exists, use PropertySourceLoader to load the configuration file
            List<PropertySource<?>> load = loader.load(name, resource);
            // Put the corresponding resources in front of RandomValuePropertySource to ensure that the loaded remote resources will take precedence over the system configuration
            load.forEach(it -> destination.addBefore(RandomValuePropertySource.RANDOM_PROPERTY_SOURCE_NAME, it));
            log.info("Load configuration success from " + url);
        } else {
            log.error("Load configuration fail from " + url + ", don't load this.");
        }
    }
}
