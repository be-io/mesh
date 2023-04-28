/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.annotate.Profiles;
import com.be.mesh.client.grpc.GrpcUtils;
import com.be.mesh.client.tool.Tool;
import io.grpc.internal.GrpcUtil;
import lombok.Data;
import org.springframework.util.unit.DataSize;

import java.io.Serializable;
import java.net.URI;
import java.time.Duration;
import java.time.temporal.ChronoUnit;

/**
 * @author coyzeng@gmail.com
 */
@Data
@Profiles("mesh.provider")
public class ProviderProperties implements Serializable {

    /**
     * Bind address for the server. Defaults to ANY_IP_ADDRESS "*". Alternatively you can restrict this to
     * ANY_IPv4_ADDRESS "0.0.0.0" or ANY_IPv6_ADDRESS "::". Or restrict it to exactly one IP address.
     * On unix systems it is also possible to prefix it with {@link GrpcUtils#DOMAIN_SOCKET_ADDRESS_PREFIX unix:} to use
     * domain socket addresses/paths (Additional dependencies may be required).
     */
    private String address = "0.0.0.0";

    /**
     * Server port to listen on. Defaults to {@code 9090}. If set to {@code 0} a random available port will be selected
     * and used. Use {@code -1} to disable the inter-process server (for example if you only want to use the in-process
     * server).
     */
    private URI expose = Tool.MESH_RUNTIME.get();

    /**
     * The name of the in-process server. If not set, then the in process server won't be started.
     */
    private String inProcessName;

    /**
     * The time to wait for the server to gracefully shutdown (completing all requests after the server started to
     * shutdown). If set to a negative value, the server waits forever. If set to {@code 0} the server will force
     * shutdown immediately. Defaults to {@code 30s}.
     */
    private Duration shutdownGracePeriod = Duration.of(30, ChronoUnit.SECONDS);

    /**
     * Setting to enable keepAlive. Default to {@code false}.
     */
    private boolean enableKeepAlive = false;

    /**
     * The default delay before we send a keepAlives. Defaults to {@code 2h}. Default unit {@link ChronoUnit#SECONDS
     * SECONDS}.
     */
    private Duration keepAliveTime = Duration.of(2, ChronoUnit.HOURS);

    /**
     * The default timeout for a keepAlives ping request. Defaults to {@code 20s}. Default unit
     * {@link ChronoUnit#SECONDS SECONDS}.
     */
    private Duration keepAliveTimeout = Duration.of(20, ChronoUnit.SECONDS);

    /**
     * Specify the most aggressive keep-alive time clients are permitted to configure. Defaults to {@code 5min}. Default
     * unit {@link ChronoUnit#SECONDS SECONDS}.
     */
    private Duration permitKeepAliveTime = Duration.of(5, ChronoUnit.MINUTES);

    /**
     * Whether clients are allowed to send keep-alive HTTP/2 PINGs even if there are no outstanding RPCs on the
     * connection. Defaults to {@code false}.
     */
    private boolean permitKeepAliveWithoutCalls = false;

    /**
     * Specify a max connection idle time. Defaults to disabled. Default unit {@link ChronoUnit#SECONDS SECONDS}.
     */
    private Duration maxConnectionIdle = null;

    /**
     * Specify a max connection age. Defaults to disabled. Default unit {@link ChronoUnit#SECONDS SECONDS}.
     */
    private Duration maxConnectionAge = null;

    /**
     * Specify a grace time for the graceful max connection age termination. Defaults to disabled. Default unit
     * {@link ChronoUnit#SECONDS SECONDS}.
     */
    private Duration maxConnectionAgeGrace = null;

    /**
     * The maximum message size allowed to be received by the server. If not set ({@code null}) then
     * {@link GrpcUtil#DEFAULT_MAX_MESSAGE_SIZE gRPC's default} should be used.
     */
    private DataSize maxInboundMessageSize = null;

    /**
     * The maximum size of metadata allowed to be received. If not set ({@code null}) then
     * {@link GrpcUtil#DEFAULT_MAX_HEADER_LIST_SIZE gRPC's default} should be used.
     */
    private DataSize maxInboundMetadataSize = null;

    /**
     * Whether gRPC health service is enabled or not. Defaults to {@code true}.
     */
    private boolean healthServiceEnabled = true;

    /**
     * Whether proto reflection service is enabled or not. Defaults to {@code true}.
     */
    private boolean reflectionServiceEnabled = true;
}
