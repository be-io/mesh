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
 * Environ is mesh fixed environ information.
 *
 * @author coyzeng@gmail.com
 */
@Data
public class Environ implements Serializable {

    private static final long serialVersionUID = 8552412835278472975L;
    /**
     * 节点ID，所有节点按照标准固定分配一个全网唯一nodeId.
     */
    @Index(value = 0, name = "node_id")
    private String nodeId;
    /**
     * 节点证书，每一个节点有一个统一的蓝象入网证书私钥，私钥用来作为节点内服务入网凭证.
     */
    @Index(value = 5, name = "node_key")
    private String nodeKey;
    /**
     * 每一个节点有一个初始的机构ID作为该节点的拥有者.
     */
    @Index(value = 10, name = "inst_id")
    private String instId;
    /**
     * 一级机构名称.
     */
    @Index(value = 15, name = "inst_name")
    private String instName;
    /**
     * 每一个节点内拥有一副证书用于通信，该证书可以被动态替换.
     */
    @Index(value = 20, name = "private_key")
    private String privateKey;
    /**
     * 每一个节点内拥有一副证书用于通信，该证书可以被动态替换.
     */
    @Index(value = 25, name = "public_key")
    private String publicKey;
    /**
     * Mesh data center.
     */
    @Index(30)
    private Distribution distribution;

}
