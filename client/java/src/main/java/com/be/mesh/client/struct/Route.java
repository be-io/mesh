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
public class Route implements Serializable {

    private static final long serialVersionUID = 9210017214247532109L;
    /**
     * 节点编号
     */
    @Index(value = 0, name = "node_id")
    private String nodeId;
    /**
     * 机构编号
     */
    @Index(value = 5, name = "inst_id")
    private String instId;
    /**
     * 节点名称
     */
    @Index(10)
    private String name;
    /**
     * 机构名称
     */
    @Index(value = 15, name = "inst_name")
    private String instName;
    /**
     * 节点地址
     */
    @Index(20)
    private String address;
    /**
     * 节点描述
     */
    @Index(25)
    private String describe;
    /**
     * Host root certifications.
     */
    @Index(value = 30, name = "host_root")
    private String hostRoot;
    /**
     * Host private certifications key
     */
    @Index(value = 35, name = "host_key")
    private String hostKey;
    /**
     * Host certification
     */
    @Index(value = 40, name = "host_crt")
    private String hostCrt;
    /**
     * Guest root certifications
     */
    @Index(value = 45, name = "guest_root")
    private String guestRoot;
    /**
     * Guest private certifications key.
     */
    @Index(value = 50, name = "guest_key")
    private String guestKey;
    /**
     * Guest certifications key.
     */
    @Index(value = 55, name = "guest_crt")
    private String guestCrt;
    /**
     * 状态
     */
    @Index(60)
    private int status;
    /**
     * 版本
     */
    @Index(65)
    private int version;
    /**
     * Auth code
     */
    @Index(value = 70, name = "auth_code")
    private String authCode;
    /**
     * Expire time
     */
    @Index(value = 75, name = "expire_at")
    private Long expireAt;
    /**
     * Extra info
     */
    @Index(80)
    private String extra;
    /**
     * Node status
     */
    @Index(value = 85, name = "create_at")
    private String createAt;
    /**
     * Node status
     */
    @Index(value = 90, name = "create_by")
    private String createBy;
    /**
     * Node status
     */
    @Index(value = 95, name = "update_at")
    private String updateAt;
    /**
     * Node status
     */
    @Index(value = 100, name = "update_by")
    private String updateBy;
    /**
     * Network group.
     */
    @Index(value = 105)
    private String group;
    /**
     * Upstream bandwidth.
     */
    @Index(value = 110)
    private long upstream;
    /**
     * Downstream bandwidth.
     */
    @Index(value = 115)
    private long downstream;
    /**
     * Static public ip address.
     */
    @Index(value = 120, name = "static_ip")
    private String staticIP;
    /**
     * Proxy endpoint in transport
     */
    @Index(125)
    private String proxy;
    /**
     * MPC concurrency
     */
    @Index(130)
    private long concurrency;
}
