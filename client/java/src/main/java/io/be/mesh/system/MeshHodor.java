/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.system;

import io.be.mesh.macro.SPI;
import io.be.mesh.prsim.Hodor;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public class MeshHodor implements Hodor {

    @Override
    public Map<String, String> stats(List<String> features) {
        Map<String, String> stats = new HashMap<>(1);
        stats.put("status", "true");
        return stats;
    }

    @Override
    public void debug(List<String> features) {
        //
    }
}
