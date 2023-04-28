/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.prsim;


import com.be.mesh.client.remote.RemoteImplement;
import com.be.mesh.client.remote.RemoteService;
import com.be.mesh.client.struct.Principal;
import lombok.extern.slf4j.Slf4j;
import org.testng.annotations.Test;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
public class StreamTest {

    private final Routable<RemoteService> service = Routable.of(new RemoteImplement());

    @Test
    public void mapTest() {
        log.info(service.local().ping("local"));
        log.info(service.any(new Principal("2")).ping("any"));
        service.many(new Principal("3"), new Principal("4")).parallel().forEach(x -> log.info(x.ping("many")));
    }

}
