/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.system;

import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.mpc.ServiceProxy;
import com.be.mesh.client.prsim.Graphics;
import com.be.mesh.client.struct.Captcha;

import java.util.Map;

/**
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public class MeshGraphics implements Graphics {

    private final Graphics graphics = ServiceProxy.proxy(Graphics.class);

    @Override
    public Captcha captcha(String kind, Map<String, String> features) {
        return graphics.captcha(kind, features);
    }

    @Override
    public boolean verify(String mno, String value) {
        return graphics.verify(mno, value);
    }
}
