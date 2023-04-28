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
import java.util.List;

/**
 * @author coyzeng@gmail.com
 */
@Data
public class License implements Serializable {

    private static final long serialVersionUID = 4359645465378617503L;
    @Index(0)
    private String version;
    @Index(5)
    private long level;
    @Index(10)
    private String name;
    @Index(value = 15, name = "create_by")
    private String createBy;
    @Index(value = 20, name = "create_at")
    private long createAt;
    @Index(value = 25, name = "active_at")
    private long activeAt;
    @Index(30)
    private List<String> factors;
    @Index(35)
    private String signature;
    @Index(value = 40, name = "node_id")
    private String nodeId;
    @Index(value = 45, name = "inst_id")
    private String instId;
    @Index(50)
    private String server;
    @Index(55)
    private String crt;
    @Index(60)
    private List<String> group;
    @Index(65)
    private long replicas;
    @Index(value = 70, name = "max_cooperators")
    private long maxCooperators;
    @Index(value = 75, name = "max_tenants")
    private long maxTenants;
    @Index(value = 80, name = "max_users")
    private long maxUsers;
    @Index(value = 85, name = "max_mills")
    private long maxMills;
    @Index(value = 90, name = "white_urns")
    private List<String> whiteUrns;
    @Index(value = 95, name = "black_urns")
    private List<String> blackUrns;


}
