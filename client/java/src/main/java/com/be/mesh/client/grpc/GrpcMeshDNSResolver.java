/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.grpc;

import com.be.mesh.client.tool.Tool;
import io.grpc.EquivalentAddressGroup;
import io.grpc.NameResolver;
import io.grpc.NameResolverProvider;
import io.grpc.Status;
import io.grpc.internal.DnsNameResolver;
import io.grpc.internal.DnsNameResolverProvider;
import lombok.extern.slf4j.Slf4j;

import java.net.InetSocketAddress;
import java.net.URI;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

/**
 * A provider for {@link DnsNameResolver}.
 *
 * <p>It resolves a target URI whose scheme is {@code "dns"}. The (optional) authority of the target
 * URI is reserved for the address of alternative DNS server (not implemented yet). The path of the
 * target URI, excluding the leading slash {@code '/'}, is treated as the host name and the optional
 * port to be resolved by DNS. Example target URIs:
 *
 * <ul>
 *   <li>{@code "dns:///foo.googleapis.com:8080"} (using default DNS)</li>
 *   <li>{@code "dns://8.8.8.8/foo.googleapis.com:8080"} (using alternative DNS (not implemented
 *   yet))</li>
 *   <li>{@code "dns:///foo.googleapis.com"} (without port)</li>
 * </ul>
 *
 * @author coyzeng@gmail.com
 */
public class GrpcMeshDNSResolver extends NameResolverProvider {

    public static final GrpcMeshDNSResolver INSTANCE = new GrpcMeshDNSResolver();
    private static final List<String> SCHEME = Arrays.asList("dns", "http", "https");

    @Override
    protected boolean isAvailable() {
        return false;
    }

    @Override
    protected int priority() {
        return new DnsNameResolverProvider().priority() + 1;
    }

    @Override
    public NameResolver newNameResolver(URI uri, NameResolver.Args args) {
        NameResolver resolver = new MeshDNSResolver(uri, args);
        resolver.refresh();
        return resolver;
    }

    @Override
    public String getDefaultScheme() {
        return SCHEME.get(0);
    }


    /**
     * {@link DnsNameResolver}
     */
    @Slf4j
    private static final class MeshDNSResolver extends NameResolver {
        private final URI uri;
        private final NameResolver.Args args;

        private MeshDNSResolver(URI uri, Args args) {
            this.uri = uri;
            this.args = args;
        }

        private InetSocketAddress resolve() {
            String[] hosts = Tool.split(Tool.MESH_ADDRESS.get().any(), ":");
            if (hosts.length < 2) {
                return new InetSocketAddress(hosts[0], Tool.DEFAULT_MESH_PORT);
            }
            return new InetSocketAddress(hosts[0], Integer.parseInt(hosts[1]));
        }

        @Override
        public String getServiceAuthority() {
            return this.uri.getAuthority();
        }

        @Override
        public void start(Listener2 listener) {
            try {
                // MPC-MESH
                InetSocketAddress address = resolve();
                List<EquivalentAddressGroup> groups = Collections.singletonList(new EquivalentAddressGroup(address));
                ResolutionResult result = ResolutionResult.newBuilder().setAddresses(groups).build();
                listener.onResult(result);
            } catch (Throwable e) {
                log.error("Resolve domain failed", e);
                listener.onError(Status.NOT_FOUND);
            }
        }

        @Override
        public void shutdown() {
            //
        }

    }
}
