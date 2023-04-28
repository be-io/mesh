/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import ch.qos.logback.classic.Level;
import ch.qos.logback.classic.Logger;
import com.be.mesh.client.annotate.MPI;
import com.be.mesh.client.prsim.Network;
import lombok.extern.slf4j.Slf4j;
import org.slf4j.LoggerFactory;
import org.testng.annotations.Test;

import java.util.Map;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.atomic.AtomicBoolean;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
public class ConcurrencyTest {

    @Test
    public void test() throws Exception {
        ((Logger) LoggerFactory.getLogger("root")).setLevel(Level.INFO);
        AtomicBoolean terminal = new AtomicBoolean();
        Network network = ServiceProxy.proxy(Network.class);
        System.setProperty("mesh.address", "gaia-mesh");
        network.getRoutes();
        ExecutorService executors = Executors.newFixedThreadPool(Runtime.getRuntime().availableProcessors());
        for (int index = 0; index < Runtime.getRuntime().availableProcessors() * 2; index++) {
            int fi = index;
            executors.submit(() -> {
                for (; ; ) {
                    if (terminal.get()) {
                        log.info("Thread {} terminal.", fi);
                        return;
                    }
                    try {
                        log.info("Thread {} success. {}", fi, network.getRoutes());
                    } catch (Throwable e) {
                        log.error(String.valueOf(fi), e);
                    }
                }
            });
        }
        new CountDownLatch(1).await();
    }

    private interface JanusFacade {
        @MPI("janus.open.invoke")
        Map<String, Object> invoke(Map<String, Object> input);
    }
}
