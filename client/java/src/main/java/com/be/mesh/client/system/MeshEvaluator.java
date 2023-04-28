/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.system;

import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.mpc.ServiceProxy;
import com.be.mesh.client.prsim.Evaluator;
import com.be.mesh.client.struct.Page;
import com.be.mesh.client.struct.Paging;
import com.be.mesh.client.struct.Script;

import java.util.List;
import java.util.Map;

/**
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public class MeshEvaluator implements Evaluator {

    private final Evaluator evaluator = ServiceProxy.proxy(Evaluator.class);

    @Override
    public String compile(Script script) {
        return evaluator.compile(script);
    }

    @Override
    public String exec(String code, Map<String, String> args, String dft) {
        return evaluator.exec(code, args, dft);
    }

    @Override
    public List<Script> dump(Map<String, String> feature) {
        return evaluator.dump(feature);
    }

    @Override
    public Page<Script> index(Paging index) {
        return evaluator.index(index);
    }
}
