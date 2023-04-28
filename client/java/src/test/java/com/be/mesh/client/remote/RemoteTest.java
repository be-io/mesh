/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.remote;

import com.be.mesh.client.annotate.MPI;
import com.be.mesh.client.mpc.Mesh;
import com.be.mesh.client.spring.MeshConfiguration;
import com.be.mesh.client.struct.Principal;
import lombok.extern.slf4j.Slf4j;
import org.springframework.test.context.ContextConfiguration;
import org.springframework.test.context.testng.AbstractTestNGSpringContextTests;
import org.testng.annotations.Test;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
@ContextConfiguration(classes = MeshConfiguration.class)
public class RemoteTest extends AbstractTestNGSpringContextTests {

    @MPI
    private RemoteService remoteService;
    @MPI
    private RemoteReference reference;

    @Test
    public void pingOuterTest() {
        try {
            Mesh.context().getPrincipals().push(new Principal("LX0000000", "LX0000000"));
            log.info(remoteService.ping("hello earth"));
        } finally {
            Mesh.context().getPrincipals().poll();
        }
    }

    @Test
    public void pingInnerTest() {
        log.info(remoteService.ping("hello earth"));
    }

    @Test
    public void pongTest() throws Exception {
        log.info(remoteService.pong("hello earth"));
    }

}
