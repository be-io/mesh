/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.tool;

/**
 * @author coyzeng@gmail.com
 */
public interface LDCable {

    String getInst();

    String getNode();

    String getAddress();

    String getProto();

    String getCodec();

    String getVersion();

    String getZone();

    String getCluster();

    String getCell();

    String getGroup();
}
