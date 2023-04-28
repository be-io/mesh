/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package io.be.mesh.system;

import io.be.mesh.macro.SPI;
import io.be.mesh.mpc.ServiceProxy;
import io.be.mesh.prsim.Tokenizer;
import io.be.mesh.struct.*;

import java.time.Duration;

/**
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public class MeshTokenizer implements Tokenizer {

    private final Tokenizer tokenizer = ServiceProxy.proxy(Tokenizer.class);

    @Override
    public String apply(String kind, Duration duration) {
        return tokenizer.apply(kind, duration);
    }

    @Override
    public boolean verify(String token) {
        return tokenizer.verify(token);
    }

    @Override
    public AccessToken quickauth(Credential credential) {
        return tokenizer.quickauth(credential);
    }

    @Override
    public AccessGrant grant(Credential credential) {
        return tokenizer.grant(credential);
    }

    @Override
    public AccessCode accept(String code) {
        return tokenizer.accept(code);
    }

    @Override
    public void reject(String code) {
        tokenizer.reject(code);
    }

    @Override
    public AccessToken authorize(String code) {
        return tokenizer.authorize(code);
    }

    @Override
    public AccessID authenticate(String token) {
        return tokenizer.authenticate(token);
    }

    @Override
    public AccessToken refresh(String token) {
        return tokenizer.refresh(token);
    }

}
