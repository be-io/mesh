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
public class Credential implements Serializable {

    private static final long serialVersionUID = -5830360176719400179L;
    @Index(value = 0, name = "client_id")
    private String clientId;
    @Index(value = 5, name = "client_key")
    private String clientKey;
    @Index(value = 10)
    private String username;
    @Index(value = 15)
    private String password;

}
