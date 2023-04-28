/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.remote;

/**
 * @author coyzeng@gmail.com
 */
public interface RemoteService {

    String ping(String hi);

    String pong(String hei);
}
