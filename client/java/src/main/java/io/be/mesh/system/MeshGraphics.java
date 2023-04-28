/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.system;

import io.be.mesh.macro.SPI;
import io.be.mesh.mpc.ServiceProxy;
import io.be.mesh.prsim.Graphics;
import io.be.mesh.struct.Captcha;

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
