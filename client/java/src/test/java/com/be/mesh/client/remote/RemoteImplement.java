/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.remote;

import com.be.mesh.client.annotate.MPI;
import com.be.mesh.client.annotate.MPS;
import com.be.mesh.client.mpc.Mesh;
import com.be.mesh.client.prsim.Routable;
import com.be.mesh.client.struct.Principal;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.BeansException;
import org.springframework.context.ApplicationContext;
import org.springframework.context.ApplicationContextAware;

import java.util.Optional;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
@MPS
public class RemoteImplement implements RemoteService, ApplicationContextAware {

    @MPI
    private Routable<RemoteReference> reference;

    @MPI
    private RemoteReference reference0;

    @Override
    public String ping(String hi) {
        return String.format("i see, %s at %s", hi, Optional.ofNullable(Mesh.context().getPrincipals().peek()).map(Principal::getNodeId).orElse("unknown"));
    }

    @Override
    public String pong(String hei) {
        return reference.local().pong(hei);
    }

    @Override
    public void setApplicationContext(ApplicationContext context) throws BeansException {
        context.getBeansOfType(RemoteReference.class).forEach((k, v) -> {
            log.info("{}-{}", k, v.getClass().getName());
        });
    }
}
