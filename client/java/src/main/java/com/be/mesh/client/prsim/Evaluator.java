/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.prsim;

import com.be.mesh.client.annotate.MPI;
import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.struct.Page;
import com.be.mesh.client.struct.Paging;
import com.be.mesh.client.struct.Script;

import java.util.List;
import java.util.Map;

/**
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public interface Evaluator {

    /**
     * Compile the named rule.
     */
    @MPI("mesh.eval.compile")
    String compile(Script script);

    /**
     * Exec the script with code.
     */
    @MPI("mesh.eval.exec")
    String exec(String code, Map<String, String> args, String dft);

    /**
     * Dump the scripts.
     */
    @MPI("mesh.eval.dump")
    List<Script> dump(Map<String, String> feature);

    /**
     * Index the scripts.
     */
    @MPI("mesh.eval.index")
    Page<Script> index(Paging index);
}
