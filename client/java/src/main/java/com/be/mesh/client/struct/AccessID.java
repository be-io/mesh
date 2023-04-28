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
public class AccessID implements Serializable {

    private static final long serialVersionUID = 4523013429808590134L;
    @Index(value = 0, name = "client_id")
    private String clientId;
    @Index(value = 10, name = "user_id")
    private String userId;
    @Index(value = 15, name = "expire_at")
    private long expireAt;

}
