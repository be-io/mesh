/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.spring;

import com.be.mesh.client.prsim.Routable;
import lombok.AllArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.FactoryBean;
import org.springframework.core.ResolvableType;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
@AllArgsConstructor
public class MeshRoutableFactoryBean implements FactoryBean<Routable<Object>> {

    private final Object mpi;

    @Override
    public Routable<Object> getObject() throws Exception {
        return Routable.of(mpi);
    }

    @Override
    public Class<?> getObjectType() {
        return ResolvableType.forClassWithGenerics(Routable.class, mpi.getClass()).resolve();
    }
}
