/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.prsim;

import com.be.mesh.client.annotate.Index;
import com.be.mesh.client.annotate.MPI;
import com.be.mesh.client.mpc.Mesh;
import com.be.mesh.client.mpc.ServiceLoader;
import com.be.mesh.client.mpc.ServiceProxy;
import com.be.mesh.client.struct.*;
import lombok.Data;
import lombok.extern.slf4j.Slf4j;
import org.testng.annotations.Test;

import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.*;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
public class NetworkTest {

    @Test
    public void versionTest() {
        Mesh.contextSafeUncheck(() -> {
            Mesh.context().getAttachments().put(Context.Metadata.MESH_SUBSET.getKey(), "LX0000010000000");
            log.info("{}", ServiceProxy.proxy(NetworkFacade.class).version("LX0000010000000"));
        });
    }

    @Test
    public void healthCheckTest() {
        System.setProperty("mesh.direct", "omega.=127.0.0.1:999");
        log.info("{}", ServiceProxy.proxy(NetworkFacade.class).findAll());
        log.info("{}", ServiceProxy.proxy(NetworkFacade.class).published());
        log.info("{}", ServiceProxy.proxy(NetworkFacade.class).healthCheck(new HashMap<>()));
        log.info("{}", ServiceProxy.proxy(NetworkFacade.class).findDataAssetDetail(2L));
    }

    @Test
    public void getEnvironTest() throws Exception {
        System.setProperty("mesh.address", "127.0.0.1:570");
        log.info("{}", ServiceProxy.proxy(Network.class).getEnviron());
        log.info("{}", ServiceProxy.proxy(Network.class).getEnviron());
        log.info("{}", ServiceProxy.proxy(Network.class).getEnviron());
        log.info("{}", ServiceProxy.proxy(Network.class).getEnviron());
        log.info("{}", ServiceProxy.proxy(NetworkFacade.class).getEnviron().get());
        Network network = ServiceLoader.load(Network.class).getDefault();
        log.info("{}", network.getEnviron());
        log.info("{}", network.getRoutes());
        log.info("{}", network.version("lx0000010000020"));
    }

    @Test
    public void weaveTest() {
        log.info("{}", ServiceProxy.proxy(Network.class).getRoutes());
        log.info("{}", ServiceProxy.proxy(Network.class).getEnviron());
        Route route = new Route();
        route.setNodeId("LX0000010000030");
        ServiceLoader.load(Network.class).getDefault().weave(route);
    }

    @Test
    public void accessibleTest() {
        // System.setProperty("mesh.host", "10.99.73.33");
        Network network = ServiceLoader.load(Network.class).getDefault();
        log.info("{}", network.getEnviron());
        log.info("{}", network.accessible(new Route()));
        log.info("{}", ServiceProxy.proxy(NetworkFacade.class).testAccessibility(new MeshEdge()));
    }

    @Test
    public void http2rpcTest() {
        System.setProperty("mesh.address", "127.0.0.1");
        Routable<NetworkFacade> network = Routable.of(ServiceProxy.proxy(NetworkFacade.class));
        List<Map<String, String>> partners = network.with("omega-token", "0af0f416-3b85-4160-a460-6f80383cf7cb").any(new Principal("LX0000010000030")).findAll();
        log.info("{}", partners);
        MeshEdge edge = new MeshEdge();
        network.any(new Principal("lx1101011100010", "lx1101011100010")).test(new HashMap<>());
    }

    @Test
    public void sequenceTest() throws Exception {
        Map<String, String> sequences = new ConcurrentHashMap<>();
        System.setProperty("mesh.address", "10.12.0.35");
        NetworkFacade facade = ServiceProxy.proxy(NetworkFacade.class);
        ExecutorService executor = Executors.newFixedThreadPool(Runtime.getRuntime().availableProcessors());
        for (int index = 0; index < 50; index++) {
            executor.submit(() -> {
                String seq = facade.nextNo("TASK");
                sequences.put(seq, seq);
                log.info(seq);
            });
        }
        executor.shutdown();
        log.info("{}-{}", executor.awaitTermination(10, TimeUnit.SECONDS), sequences.size());
    }

    @Test
    public void middlewarePluginTest() throws Exception {
        NetworkFacade facade = ServiceProxy.proxy(NetworkFacade.class);
        Map<String, String> x = facade.subscribe(Event.newInstance(new HashMap<>(), new Topic("com.trustbe.janus.sze.asset.status", "")));
        log.info("{}", x);
    }

    public interface NetworkFacade {

        /**
         * 生成序列号.
         */
        @MPI("omega.system.generator.sequence")
        String nextNo(@Index(0) String biz);

        @MPI(value = "mesh.net.version")
        Versions version(@Index(value = 0, name = "node_id") String nodeId);

        @MPI("mesh.net.environ")
        Future<Map<String, Object>> getEnviron();

        /**
         * 测试网络可达性.
         *
         * @param route 组网数据
         * @return true可达
         */
        @MPI("mesh.net.accessible")
        boolean testAccessibility(MeshEdge route);

        @MPI(name = "tensor.route.grpc")
        void test(Map<String, String> param);

        @MPI(value = "omega.cooperation.find.all", node = "LX0000010000010", address = "10.99.31.33:570")
        List<Map<String, String>> findAll();

        @MPI(value = "theta.cooperator.health.check", node = "LX0000010000010")
        Object healthCheck(Map<String, String> args);

        @MPI(value = "theta.asset.published.list", node = "JG0100000500000000")
        List<Map<String, Object>> published();

        @MPI(value = "theta.asset.query.byId", node = "JG0100000500000000")
        Map<String, Object> findDataAssetDetail(@Index(0) Long dataAssetId);

        @MPI("mesh.x.y.status")
        Map<String, String> subscribe(Event event);
    }

    @Data
    public static final class MeshEdge {
        private String address = "127.0.0.1:573";
    }
}
