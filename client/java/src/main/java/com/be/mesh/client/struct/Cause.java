/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.struct;

import com.be.mesh.client.annotate.Index;
import com.be.mesh.client.cause.MeshException;
import com.be.mesh.client.mpc.Mesh;
import com.be.mesh.client.mpc.Types;
import com.be.mesh.client.prsim.Context;
import com.be.mesh.client.tool.Tool;
import lombok.Data;
import lombok.extern.slf4j.Slf4j;

import java.io.Serializable;

/**
 * @author coyzeng@gmail.com
 */
@Data
@Slf4j
public class Cause implements Serializable {

    private static final long serialVersionUID = -4907560275981296064L;
    public static final Context.Key<String> POS = new Context.Key<>("mesh.mpc.cause", Types.of(String.class));

    @Index(0)
    private String name;
    @Index(5)
    private String pos;
    @Index(10)
    private String text;
    @Index(15)
    private byte[] buff;

    public static Cause of(Throwable e) {
        Cause cause = new Cause();
        String name = e.getClass().getName();
        cause.setName(name);
        if (Tool.required(Mesh.context().getAttribute(POS))) {
            cause.setPos(Mesh.context().getAttribute(POS));
        } else {
            cause.setPos(String.format("%s(%s)(%s)", Tool.IP.get(), Tool.HOST_NAME.get(), Tool.MESH_NAME.get()));
        }
        if (!name.startsWith("java.") && !name.startsWith("javax.") && !(e instanceof MeshException)) {
            cause.setText(Tool.getStackTrace(e));
        }
        try {
            cause.setBuff(Tool.serializeCause(e));
        } catch (Throwable ioe) {
            log.warn("Serialize cause with error.", ioe);
        }
        return cause;
    }

    public static Throwable of(String code, String message, Cause cause) {
        try {
            Mesh.context().setAttribute(POS, cause.getPos());
            Class.forName(cause.getName());
            return Tool.deserializeCause(cause.getBuff());
        } catch (Throwable e) {
            log.warn(cause.getText());
            log.warn("Deserialize cause with error.", e);
            return new MeshException(code, message);
        }
    }
}
