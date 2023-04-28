/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.annotate.MPI;
import com.be.mesh.client.annotate.MPS;
import com.be.mesh.client.struct.Reference;
import com.be.mesh.client.struct.Service;
import com.be.mesh.client.tool.Tool;
import lombok.extern.slf4j.Slf4j;
import org.testng.annotations.Test;

import java.util.ArrayList;
import java.util.List;
import java.util.regex.Pattern;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
public class EdenTest {

    private interface X {
        @MPI("${mesh.name}.x.x")
        void x();
    }

    @MPS
    private static class Y implements X {

        @Override
        public void x() {

        }
    }

    @MPI
    private X x;

    @Test
    public void testReference() throws Exception {
        System.setProperty("mesh.name", "mesh");
        System.setProperty("mesh.address", "10.90.30.168");
        Eden eden = ServiceLoader.load(Eden.class).getDefault();
        MPI mpi = EdenTest.class.getDeclaredField("x").getAnnotation(MPI.class);
        Execution<Reference> execution = eden.refer(mpi, X.class, new MethodInspector(X.class, X.class.getDeclaredMethod("x")));
        log.info(execution.schema().getUrn());

        eden.store(Y.class, new Y());
        Execution<Service> service = eden.infer("x.x.mesh.0000000102030000000012700000001104123.lx000000000000000000x.trustbe.cn");
        log.info(service.schema().getUrn());
    }

    @Test
    public void testRegexSplit() throws Exception {
        Codec codec = ServiceLoader.load(Codec.class).getDefault();
        Pattern pattern = Pattern.compile("\\$\\{*?}");
        log.info(codec.encodeString(pattern.split("${mesh.address}.ns.y.z")));

        String[] names = "${mesh.address}.ns.y.z".split("\\.");
        List<String> ns = new ArrayList<>();
        StringBuilder buffer = new StringBuilder();
        for (String name : names) {
            if (Tool.startWith(name, "${")) {
                buffer.append(name).append('.');
                continue;
            }
            if (Tool.endsWith(name, "}")) {
                buffer.append(name);
                ns.add(buffer.toString());
                buffer.delete(0, buffer.length());
                continue;
            }
            ns.add(name);
        }
        log.info(codec.encodeString(ns));
    }
}
