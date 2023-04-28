/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.prsim;

import com.be.mesh.client.annotate.Index;
import com.be.mesh.client.annotate.MPI;
import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.struct.*;

import java.time.Duration;

/**
 * <pre>
 *      +--------+                               +---------------+
 *      |        |--(A)- Authorization Request ->|   Resource    |
 *      |        |                               |     Owner     |
 *      |        |<-(B)-- Authorization Grant ---|               |
 *      |        |                               +---------------+
 *      |        |
 *      |        |                               +---------------+
 *      |        |--(C)-- Authorization Grant -->| Authorization |
 *      | Client |                               |     Server    |
 *      |        |<-(D)----- Access Token -------|               |
 *      |        |                               +---------------+
 *      |        |
 *      |        |                               +---------------+
 *      |        |--(E)----- Access Token ------>|    Resource   |
 *      |        |                               |     Server    |
 *      |        |<-(F)--- Protected Resource ---|               |
 *      +--------+                               +---------------+
 * </pre>
 *
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public interface Tokenizer {

    /**
     * Apply a node token.
     *
     * @param kind     token kind
     * @param duration max life cycle
     * @return node token
     */
    @MPI("mesh.trust.apply")
    String apply(@Index(0) String kind, @Index(1) Duration duration);

    /**
     * Verify some token verifiable.
     *
     * @param token node token
     * @return true if verify passed
     */
    @MPI("mesh.trust.verify")
    boolean verify(@Index(0) String token);

    /**
     * OAuth2 quick authorize, contains grant code and code authorize.
     */
    @MPI("mesh.oauth2.quickauth")
    AccessToken quickauth(@Index(0) Credential credential);

    /**
     * OAuth2 code grant.
     */
    @MPI("mesh.oauth2.grant")
    AccessGrant grant(@Index(0) Credential credential);

    /**
     * OAuth2 accept grant code.
     */
    @MPI("mesh.oauth2.accept")
    AccessCode accept(@Index(0) String code);

    /**
     * OAuth2 reject grant code.
     */
    @MPI("mesh.oauth2.accept")
    void reject(@Index(0) String code);

    /**
     * OAuth2 code authorize.
     */
    @MPI("mesh.oauth2.authorize")
    AccessToken authorize(@Index(0) String code);

    /**
     * OAuth2 authenticate.
     */
    @MPI("mesh.oauth2.authenticate")
    AccessID authenticate(@Index(0) String token);

    /**
     * OAuth2 auth token refresh.
     */
    @MPI("mesh.oauth2.refresh")
    AccessToken refresh(@Index(0) String token);
}
