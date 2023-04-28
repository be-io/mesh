/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.system;

import io.be.mesh.macro.SPI;
import io.be.mesh.mpc.ServiceProxy;
import io.be.mesh.prsim.Evaluator;
import io.be.mesh.struct.Page;
import io.be.mesh.struct.Paging;
import io.be.mesh.struct.Script;

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
