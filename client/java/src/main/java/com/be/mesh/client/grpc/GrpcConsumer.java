/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.grpc;

import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.cause.TimeoutException;
import com.be.mesh.client.mpc.Consumer;
import com.be.mesh.client.mpc.Execution;
import com.be.mesh.client.struct.Reference;
import com.be.mesh.client.tool.Tool;
import io.grpc.*;
import io.grpc.netty.NegotiationType;
import io.grpc.netty.NettyChannelBuilder;
import io.grpc.stub.ClientCalls;
import io.netty.channel.ChannelOption;
import lombok.extern.slf4j.Slf4j;
import org.springframework.util.unit.DataSize;

import java.io.InputStream;
import java.net.InetSocketAddress;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.Future;
import java.util.concurrent.TimeUnit;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
@SPI(Consumer.GRPC)
public class GrpcConsumer implements AutoCloseable, Consumer {

    private final Map<String, GrpcChannelPool> channels = new ConcurrentHashMap<>(3);
    private final MethodDescriptor<InputStream, InputStream> descriptor = newDescriptor();

    @Override
    public void start() throws Throwable {
        //
    }

    @Override
    public Future<InputStream> consume(String address, String urn, Execution<Reference> execution, InputStream inbound) {
        if (isProtobufNative(inbound)) {
            log.warn("MPI service is protobuf native protocol.");
        }
        ManagedChannel channel = channels.computeIfAbsent(address, this::newChannelPool);
        CallOptions options = CallOptions.DEFAULT.withWaitForReady().withAuthority(urn).withDeadline(Deadline.after(execution.schema().getTimeout(), TimeUnit.MILLISECONDS));
        ClientCall<InputStream, InputStream> clientCall = channel.newCall(descriptor, options);
        return new GrpcFuture<>(ClientCalls.futureUnaryCall(clientCall, inbound), e -> resetIfTimeout(e, channel));
    }

    @Override
    public void close() throws Exception {
        for (ManagedChannel channel : channels.values()) {
            channel.awaitTermination(6, TimeUnit.SECONDS);
        }
    }

    /**
     * Reset transport if channel timeout.
     */
    private void resetIfTimeout(Throwable e, ManagedChannel channel) {
        if (e instanceof TimeoutException) {
            channel.resetConnectBackoff();
        }
    }

    /**
     * Resolve the mesh address.
     */
    private InetSocketAddress resolve(String address) {
        String[] hosts = Tool.split(address, ":");
        if (hosts.length < 2) {
            return new InetSocketAddress(hosts[0], Tool.DEFAULT_MESH_PORT);
        }
        return new InetSocketAddress(hosts[0], Integer.parseInt(hosts[1]));
    }

    /**
     * Initial channel pool.
     */
    private GrpcChannelPool newChannelPool(String address) {
        GrpcSettings settings = new GrpcSettings();
        settings.setMinRpcPerChannel(Tool.getProperty(100, "mesh.grpc.channel.rpc.min"));
        settings.setMaxRpcPerChannel(Tool.getProperty(200, "mesh.grpc.channel.rpc.max"));
        settings.setMinChannels(Tool.getProperty(Runtime.getRuntime().availableProcessors() * 2, "mesh.grpc.channel.min"));
        settings.setMaxChannels(Tool.getProperty(200, "mesh.grpc.channel.max"));
        settings.setPreemptiveRefreshEnabled(true);
        return GrpcChannelPool.create(settings, () -> this.newChannel(address));
    }

    /**
     * <pre>
     * if (Epoll.isAvailable()) {
     *     nc.channelType(EpollDomainSocketChannel.class);
     *     nc.eventLoopGroup(new EpollEventLoopGroup());
     *  }
     * </pre>
     */
    private ManagedChannel newChannel(String address) {
        NettyChannelBuilder nc = NettyChannelBuilder.forAddress(resolve(address));
        nc.overrideAuthority(Tool.MESH_ADDRESS.get().any());
        nc.negotiationType(NegotiationType.PLAINTEXT);
        nc.withOption(ChannelOption.CONNECT_TIMEOUT_MILLIS, (int) TimeUnit.SECONDS.toMillis(12));
        nc.keepAliveTime(60, TimeUnit.SECONDS);
        nc.keepAliveTimeout(60, TimeUnit.SECONDS);
        nc.idleTimeout(12, TimeUnit.SECONDS);
        nc.keepAliveWithoutCalls(true);
        nc.maxInboundMessageSize((int) DataSize.ofMegabytes(512).toBytes());
        nc.maxInboundMetadataSize((int) DataSize.ofMegabytes(65).toBytes());
        nc.enableFullStreamDecompression();
        nc.enableRetry();
        nc.maxRetryAttempts(3);
        nc.intercept(GrpcInterceptor.INTERCEPTOR);
        ManagedChannel channel = nc.build();
        channel.notifyWhenStateChanged(ConnectivityState.SHUTDOWN, channel::resetConnectBackoff);
        channel.notifyWhenStateChanged(ConnectivityState.TRANSIENT_FAILURE, channel::resetConnectBackoff);
        return channel;
    }

    /**
     * ProtoUtils.marshaller(new Inbound())
     */
    private MethodDescriptor<InputStream, InputStream> newDescriptor() {
        MethodDescriptor.Builder<InputStream, InputStream> db = MethodDescriptor.newBuilder();
        db.setType(MethodDescriptor.MethodType.UNARY);
        db.setFullMethodName(GrpcUtils.MESH_INVOKE_METHOD);
        db.setSampledToLocalTracing(true);
        db.setRequestMarshaller(new GrpcMarshaller());
        db.setResponseMarshaller(new GrpcMarshaller());
        db.setSchemaDescriptor(null);
        return db.build();
    }

    /**
     * Is the payload is protobuf message.
     */
    private boolean isProtobufNative(Object argument) {
        return argument instanceof com.google.protobuf.Message;
    }
}
