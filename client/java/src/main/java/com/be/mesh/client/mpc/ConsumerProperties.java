/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.annotate.Profiles;
import lombok.Data;
import org.springframework.util.unit.DataSize;

import java.io.Serializable;
import java.time.Duration;

/**
 * @author coyzeng@gmail.com
 */
@Data
@Profiles("mesh.consumer")
public class ConsumerProperties implements Serializable {


    private static final long serialVersionUID = 976224325055467957L;
    /**
     * Sets the target address uri for the channel. The target uri must be in the format:
     * {@code schema:[//[authority]][/path]}. If nothing is configured then the name of the client will be used along
     * with the default scheme. We recommend explicitly configuring the scheme used for the address resolutions such as
     * {@code dns:/}.
     *
     * <p>
     * <b>Examples</b>
     * </p>
     *
     * <ul>
     * <li>{@code static://localhost:9090} (refers to exactly one IPv4 or IPv6 address, dependent on the jre
     * configuration, it does not check whether there is actually someone listening on that network interface)</li>
     * <li>{@code static://10.0.0.10}</li>
     * <li>{@code static://10.0.0.10,10.11.12.11}</li>
     * <li>{@code static://10.0.0.10:9090,10.0.0.11:80,10.0.0.12:1234,[::1]:8080}</li>
     * <li>{@code dns:/localhost (might refer to the IPv4 or the IPv6 address or both, dependent on the system
     * configuration, it does not check whether there is actually someone listening on that network interface)}</li>
     * <li>{@code dns:/example.com}</li>
     * <li>{@code dns:/example.com:9090}</li>
     * <li>{@code dns:///example.com:9090}</li>
     * <li>{@code discovery:/foo-service}</li>
     * <li>{@code discovery:///foo-service}</li>
     * <li>{@code unix:<relative-path>} (Additional dependencies may be required)</li>
     * <li>{@code unix://</absolute-path>} (Additional dependencies may be required)</li>
     * </ul>
     * <p>
     * The string representation of an uri to use as target address or null to use a fallback.
     *
     * @see <a href="https://github.com/grpc/grpc/blob/master/doc/naming.md">gRPC Name Resolution</a>
     */
    private String address;
    private boolean enableKeepAlive;
    private Duration keepAliveTime = Duration.ofMinutes(12);
    private Duration keepAliveTimeout = Duration.ofMinutes(12);
    private boolean keepAliveWithoutCalls;
    private Duration shutdownGracePeriod = Duration.ofSeconds(5);
    private DataSize maxInboundMessageSize = DataSize.ofBytes(1024 << 16);  // 64MB
    private boolean fullStreamDecompression;
    private Duration immediateConnectTimeout = Duration.ofSeconds(5);

}
