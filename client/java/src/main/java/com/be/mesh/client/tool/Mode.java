/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.tool;

import lombok.AllArgsConstructor;
import lombok.Getter;

import java.io.Serializable;

/**
 * @author coyzeng@gmail.com
 */
@Getter
@AllArgsConstructor
public class Mode implements Serializable {

    private static final long serialVersionUID = -3541296255932306873L;
    private final long code;

    public static final Mode DISABLE = new Mode(1L);
    public static final Mode FAILFAST = new Mode(2L);
    public static final Mode NOLOG = new Mode(4L);
    public static final Mode JSONLOG = new Mode(8L);
    public static final Mode RCache = new Mode(16L);
    public static final Mode PHeader = new Mode(32L);
    public static final Mode Metrics = new Mode(64L);
    public static final Mode RLog = new Mode(128L);
    public static final Mode MGrpc = new Mode(256L);
    public static final Mode PermitCirculate = new Mode(512);
    public static final Mode NoStdColor = new Mode(1024);
    public static final Mode NoCommunal = new Mode(2048);
    public static final Mode DisableTee = new Mode(4096);
    public static final Mode NoStaticFile = new Mode(8192);
    public static final Mode OpenTelemetry = new Mode(16384);

    public boolean match(Mode mode) {
        if (null == mode) {
            return false;
        }
        return (this.getCode() & mode.getCode()) == mode.getCode();
    }

    public static Mode from(String code) {
        if (!Tool.isNumeric(code)) {
            return FAILFAST;
        }
        return new Mode(Long.parseLong(code));
    }

    @Override
    public String toString() {
        return String.valueOf(this.code);
    }
}
