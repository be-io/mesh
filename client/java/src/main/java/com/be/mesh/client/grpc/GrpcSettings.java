/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.grpc;

import io.netty.channel.pool.ChannelPool;
import lombok.Data;

import java.time.Duration;

/**
 * @author coyzeng@gmail.com
 */
@Data
public class GrpcSettings {

    /**
     * How often to check and possibly resize the {@link ChannelPool}.
     */
    public static final Duration RESIZE_INTERVAL = Duration.ofMinutes(1);
    /**
     * The maximum number of channels that can be added or removed at a time.
     */
    public static final int MAX_RESIZE_DELTA = 2;
    /**
     * Threshold to start scaling down the channel pool.
     *
     * <p>When the average of the maximum number of outstanding RPCs in a single minute drop below
     * this threshold, channels will be removed from the pool.
     */
    private int minRpcPerChannel;

    /**
     * Threshold to start scaling up the channel pool.
     *
     * <p>When the average of the maximum number of outstanding RPCs in a single minute surpass this
     * threshold, channels will be added to the pool. For google services, gRPC channels will start
     * locally queuing RPC when there are 100 concurrent RPCs.
     */
    private int maxRpcPerChannel;

    /**
     * The absolute minimum size of the channel pool.
     *
     * <p>Regardless of the current throughput, the number of channels will not drop below this limit
     */
    private int minChannels;

    /**
     * The absolute maximum size of the channel pool.
     *
     * <p>Regardless of the current throughput, the number of channels will not exceed this limit
     */
    private int maxChannels;

    /**
     * If all of the channels should be replaced on an hourly basis.
     *
     * <p>The GFE will forcibly disconnect active channels after an hour. To minimize the cost of
     * reconnects, this will create a new channel asynchronuously, prime it and then swap it with an
     * old channel.
     */
    private boolean preemptiveRefreshEnabled;

    /**
     * Helper to check if the {@link ChannelPool} implementation can skip dynamic size logic
     */
    static boolean isStaticSize(GrpcSettings settings) {
        // When range is restricted to a single size
        if (settings.getMinChannels() == settings.getMaxChannels()) {
            return true;
        }
        // When the scaling threshold are not set
        return settings.getMinRpcPerChannel() == 0 && settings.getMaxRpcPerChannel() == Integer.MAX_VALUE;
    }

}
