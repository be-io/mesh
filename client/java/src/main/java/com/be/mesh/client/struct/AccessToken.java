/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.struct;

import com.be.mesh.client.annotate.Index;
import lombok.Data;

import java.io.Serializable;

/**
 * @author coyzeng@gmail.com
 */
@Data
public class AccessToken implements Serializable {

    private static final long serialVersionUID = -1322503539021034512L;
    @Index(value = 0)
    private String token;
    @Index(value = 5)
    private String type;
    @Index(value = 10, name = "expires_at")
    private long expiresAt;
    @Index(15)
    private String scope;
    @Index(value = 20, name = "refresh_token")
    private String refreshToken;
}
