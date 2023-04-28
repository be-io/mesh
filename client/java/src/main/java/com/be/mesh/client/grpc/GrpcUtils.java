/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.grpc;

import com.be.mesh.client.cause.CompatibleException;
import io.grpc.LoadBalancerRegistry;
import io.grpc.MethodDescriptor;

/**
 * Utility class that contains methods to extract some information from grpc classes.
 *
 * @author coyzeng@gmail.com
 */
public final class GrpcUtils {

    /**
     * Mesh invoke uniform service name.
     */
    public static final String MESH_SERVICE_NAME = "mesh-rpc";

    /**
     * Mesh invoke uniform method.
     */
    public static final String MESH_INVOKE_METHOD = "mesh-rpc/v1";

    /**
     * A constant that defines, that the server should listen to any IPv4 and IPv6 address.
     */
    public static final String ANY_IP_ADDRESS = "*";

    /**
     * A constant that defines, that the server should listen to any IPv4 address.
     */
    public static final String ANY_IPv4_ADDRESS = "0.0.0.0";

    /**
     * A constant that defines, that the server should listen to any IPv6 address.
     */
    public static final String ANY_IPv6_ADDRESS = "::";

    /**
     * Sets the default load balancing policy for this channel. This config might be overwritten by the service config
     * received from the target address. The names have to be resolvable from the {@link LoadBalancerRegistry}. By
     * default this the {@code round_robin} policy. Please note that this policy is different from the normal grpc-java
     * default policy {@code pick_first}.
     */
    public static final String DEFAULT_DEFAULT_LOAD_BALANCING_POLICY = "round_robin";

    /**
     * A constant that defines, the scheme of a Unix domain socket address.
     */
    public static final String DOMAIN_SOCKET_ADDRESS_SCHEME = "unix";

    /**
     * A constant that defines, the scheme prefix of a Unix domain socket address.
     */
    public static final String DOMAIN_SOCKET_ADDRESS_PREFIX = DOMAIN_SOCKET_ADDRESS_SCHEME + ":";

    /**
     * The cloud discovery metadata key used to identify the grpc port.
     */
    public static final String CLOUD_DISCOVERY_METADATA_PORT = "gRPC_port";

    /**
     * The constant for the grpc server port, -1 represents don't start an inter process server.
     */
    public static final int INTER_PROCESS_DISABLE = -1;

    /**
     * Extracts the domain socket address specific path from the given full address. The address must fulfill the
     * requirements as specified by <a href="https://grpc.github.io/grpc/cpp/md_doc_naming.html">grpc</a>.
     *
     * @param address The address to extract it from.
     * @return The extracted domain socket address specific path.
     * @throws IllegalArgumentException If the given address is not a valid address.
     */
    public static String extractDomainSocketAddressPath(final String address) {
        if (!address.startsWith(DOMAIN_SOCKET_ADDRESS_PREFIX)) {
            throw new CompatibleException(address + " is not a valid domain socket address.");
        }
        String path = address.substring(DOMAIN_SOCKET_ADDRESS_PREFIX.length());
        if (path.startsWith("//")) {
            path = path.substring(2);
            // We don't check this as there is no reliable way to check that it's an absolute path,
            // especially when Windows adds support for these in the future
            // if (!path.startsWith("/")) {
            // throw new IllegalArgumentException("If the path is prefixed with '//', then the path must be absolute");
            // }
        }
        return path;
    }

    /**
     * Extracts the service name from the given method.
     *
     * @param method The method to get the service name from.
     * @return The extracted service name.
     * @see MethodDescriptor#extractFullServiceName(String)
     * @see #extractMethodName(MethodDescriptor)
     */
    public static String extractServiceName(final MethodDescriptor<?, ?> method) {
        return MethodDescriptor.extractFullServiceName(method.getFullMethodName());
    }

    /**
     * Extracts the method name from the given method.
     *
     * @param method The method to get the method name from.
     * @return The extracted method name.
     * @see #extractServiceName(MethodDescriptor)
     */
    public static String extractMethodName(final MethodDescriptor<?, ?> method) {
        // This method is the equivalent of MethodDescriptor.extractFullServiceName
        final String fullMethodName = method.getFullMethodName();
        final int index = fullMethodName.lastIndexOf('/');
        if (index == -1) {
            return fullMethodName;
        }
        return fullMethodName.substring(index + 1);
    }

    private GrpcUtils() {
    }

}
