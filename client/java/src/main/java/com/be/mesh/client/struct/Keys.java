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
public class Keys implements Serializable {

    private static final long serialVersionUID = 7786786600610178192L;

    @Index(value = 0, name = "private_key")
    private String privateKey;

    @Index(value = 5, name = "public_key")
    private String publicKey;

    @Index(10)
    private String algorithm;

    @Index(15)
    private String csr;

    @Index(20)
    private String cert;
}
