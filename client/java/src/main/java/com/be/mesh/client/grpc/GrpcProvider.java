/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.grpc;

import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.mpc.Eden;
import com.be.mesh.client.mpc.Provider;
import com.be.mesh.client.mpc.ProviderProperties;
import com.be.mesh.client.mpc.ServiceLoader;
import com.be.mesh.client.prsim.RuntimeHook;
import com.be.mesh.client.tool.Tool;
import com.google.common.net.InetAddresses;
import io.grpc.BindableService;
import io.grpc.Server;
import io.grpc.netty.NettyServerBuilder;
import io.netty.channel.epoll.EpollEventLoopGroup;
import io.netty.channel.epoll.EpollServerDomainSocketChannel;
import io.netty.channel.unix.DomainSocketAddress;
import lombok.extern.slf4j.Slf4j;
import org.springframework.util.unit.DataSize;

import java.net.InetSocketAddress;
import java.time.Duration;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicBoolean;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
@SPI("grpc")
public class GrpcProvider implements Provider, RuntimeHook {

    private final AtomicBoolean running = new AtomicBoolean(false);
    private static final Duration SHUTDOWN_DELAY = Duration.ofSeconds(3);
    private volatile Server server;

    /**
     * Start the gRPC {@link Server}.
     */
    public void start() throws Throwable {
        if (ServiceLoader.load(Eden.class).getDefault().inferTypes().isEmpty()) {
            log.debug("None mesh service has been annotated, dont start the grpc provider.");
            return;
        }
        if (running.get()) {
            return;
        }
        ProviderProperties properties = new ProviderProperties();
        NettyServerBuilder building = createServerBuilder(properties);
        if (properties.isEnableKeepAlive()) {
            building.keepAliveTime(properties.getKeepAliveTime().toNanos(), TimeUnit.NANOSECONDS);
            building.keepAliveTimeout(properties.getKeepAliveTimeout().toNanos(), TimeUnit.NANOSECONDS);
        }
        building.permitKeepAliveTime(properties.getPermitKeepAliveTime().toNanos(), TimeUnit.NANOSECONDS);
        building.permitKeepAliveWithoutCalls(properties.isPermitKeepAliveWithoutCalls());

        if (properties.getMaxConnectionIdle() != null) {
            building.maxConnectionIdle(properties.getMaxConnectionIdle().toNanos(), TimeUnit.NANOSECONDS);
        }
        if (properties.getMaxConnectionAge() != null) {
            building.maxConnectionAge(properties.getMaxConnectionAge().toNanos(), TimeUnit.NANOSECONDS);
        }
        if (properties.getMaxConnectionAgeGrace() != null) {
            building.maxConnectionAgeGrace(properties.getMaxConnectionAgeGrace().toNanos(), TimeUnit.NANOSECONDS);
        }
        building.maxInboundMessageSize((int) DataSize.ofMegabytes(512).toBytes());
        building.maxInboundMetadataSize((int) DataSize.ofMegabytes(64).toBytes());

        building.intercept(GrpcInterceptor.INTERCEPTOR);
        building.addService(new GrpcBindableService());
        for (BindableService service : ServiceLoader.load(BindableService.class).list()) {
            building.addService(service);
        }
        this.server = building.build();
        this.server.start();
        log.info("Mesh service grpc provider bootstrap with [{}]", properties.getAddress());
        Runtime.getRuntime().addShutdownHook(new Thread(() -> {
            try {
                close();
            } catch (Exception e) {
                log.error("", e);
                log.error("Shutdown grpc building failed.");
            }
        }));
        running.set(true);
    }

    @Override
    public void refresh() throws Throwable {
        start();
    }

    @Override
    public void stop() throws Throwable {
        close();
    }

    /**
     * Shutdown the gRPC {@link Server} when this object is closed.
     */
    @Override
    public void close() throws Exception {
        if (ServiceLoader.load(Eden.class).getDefault().inferTypes().isEmpty()) {
            log.debug("None mesh service has been annotated, dont stop the grpc provider.");
            return;
        }
        if (null == server || server.isShutdown()) {
            return;
        }
        try {
            server.shutdown();
            server.awaitTermination(SHUTDOWN_DELAY.toMillis(), TimeUnit.MILLISECONDS);
        } finally {
            server.shutdownNow();
        }
        log.info("Stop mesh service grpc provider finish.");
    }

    /**
     * Creates a new server builder.
     *
     * @return The newly created server builder.
     */
    private NettyServerBuilder createServerBuilder(ProviderProperties properties) {
        String address = Tool.anyone(properties.getAddress(), "*");
        int port = Tool.anyone(properties.getExpose().getPort(), 80);
        if (address.startsWith(GrpcUtils.DOMAIN_SOCKET_ADDRESS_PREFIX)) {
            String path = GrpcUtils.extractDomainSocketAddressPath(address);
            return NettyServerBuilder.forAddress(new DomainSocketAddress(path))
                    .channelType(EpollServerDomainSocketChannel.class)
                    .bossEventLoopGroup(new EpollEventLoopGroup(1))
                    .workerEventLoopGroup(new EpollEventLoopGroup());
        } else if (GrpcUtils.ANY_IP_ADDRESS.equals(address)) {
            return NettyServerBuilder.forPort(port);
        } else {
            return NettyServerBuilder.forAddress(new InetSocketAddress(InetAddresses.forString(address), port));
        }
    }

}
