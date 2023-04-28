/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.metric;

import io.be.mesh.tool.Once;
import io.prometheus.client.CollectorRegistry;

/**
 * @author coyzeng@gmail.com
 */
public class Registry {

    public static final Once<CollectorRegistry> REGISTRY = Once.with(() -> new CollectorRegistry(false));
}
