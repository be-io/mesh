/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.remote;

import com.be.mesh.client.annotate.MPS;

/**
 * @author coyzeng@gmail.com
 */
@MPS
public class RemoteReferenceImplement implements RemoteReference {

    @Override
    public String pong(String hei) {
        return "i am here";
    }
}
