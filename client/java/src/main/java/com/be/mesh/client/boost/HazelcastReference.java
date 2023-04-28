/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.boost;

import com.be.mesh.client.tool.Once;
import com.hazelcast.core.HazelcastInstance;

/**
 * @author coyzeng@gmail.com
 */
public class HazelcastReference {

    public static final Once<HazelcastInstance> INST = Once.with(com.hazelcast.core.Hazelcast::newHazelcastInstance);

}
