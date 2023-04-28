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
import java.util.Map;

/**
 * @author coyzeng@gmail.com
 */
@Data
public class Registration implements Serializable {

    private static final long serialVersionUID = -7483417635501048566L;
    public static final String METADATA = "metadata";
    public static final String PROXY = "proxy";

    @Index(value = 0, name = "instance_id")
    private String instanceId;
    @Index(5)
    private String name;
    @Index(10)
    private String kind;
    @Index(15)
    private String address;
    @Index(20)
    private Object content;
    @Index(25)
    private long timestamp;
    @Index(30)
    private Map<String, String> attachments;
}
