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
import com.be.mesh.client.struct.Institution;
import lombok.Data;

import java.io.Serializable;
import java.util.List;

/**
 * 互联互通.
 *
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public interface Interconnection {

    /**
     * 获取合约关联的合作机构。会包含一个本方，目前默认只有两个，保留未来多个.
     */
    @MPI("mesh.inc.instx")
    List<Institution> instx(@Index(value = 0, name = "contract_id") String contractId);

    /**
     * 获取多个合作方之间的互联互通合约，为空表示没有合约，合作方大于2个时会返回多个合约
     */
    @MPI("mesh.inc.contract.ids")
    List<String> contractIds(@Index(value = 0, name = "inst_ids") List<String> instIds);

    /**
     * /v1/interconn/node/contract/query
     * APPLIED已申请
     * APPROVED已授权
     * REJECTED已拒绝
     * TERMINATED已解除
     */
    @MPI("mesh.inc.contract")
    IncState contract(IncContractID req);

    /**
     * /v1/interconn/node/query  /v1/platform/node/query
     */
    @MPI(name = "mesh.inc.describe", flags = 8)
    IncNode describe(IncNodeID req);

    /**
     * Weave /v1/interconn/node/contract/apply
     */
    @MPI(name = "mesh.inc.weave", flags = 8)
    IncContractID weave(IncNode req);

    /**
     * Ack /v1/interconn/node/contract/confirm
     */
    @MPI(name = "mesh.inc.ack", flags = 8)
    IncOption ack(IncAck req);

    /**
     * Abort /v1/interconn/node/contract/terminate
     */
    @MPI(name = "mesh.inc.abort", flags = 8)
    IncOption abort(IncContractID req);

    /**
     * Refresh /v1/interconn/node/update /v1/platform/node/update
     */
    @MPI(name = "mesh.inc.refresh", flags = 8)
    IncOption refresh(IncNode req);

    /**
     * Probe /v1/interconn/node/health
     * 节点健康状态。直接返回ok
     */
    @MPI(name = "mesh.inc.probe", flags = 8)
    IncOption probe(IncState req);


    @Data
    class IncContractID implements Serializable {

        private static final long serialVersionUID = 4576648697896713353L;
        @Index(value = 0, name = "contract_id")
        private String contractId;

    }

    @Data
    class IncState implements Serializable {

        private static final long serialVersionUID = -3333767381417774432L;
        @Index(value = 0, name = "status")
        private String status; // APPLIED已申请 APPROVED已授权 REJECTED已拒绝 TERMINATED已 解除

    }

    @Data
    class IncNodeID implements Serializable {

        private static final long serialVersionUID = 8262814780729011659L;
        @Index(value = 0, name = "node_id")
        private String nodeId;

    }

    @Data
    class IncNode implements Serializable {

        private static final long serialVersionUID = 7715869947106526402L;
        @Index(value = 1, name = "node_id")
        private String nodeId; // 合作方的节点ID
        @Index(value = 2, name = "name")
        private String name; // 节点名称
        @Index(value = 3, name = "institution")
        private String institution; // 节点所属机构
        @Index(value = 4, name = "inst_id")
        private String instId; // 机构Id
        @Index(value = 5, name = "system")
        private String system; // 技术服务提供系统
        @Index(value = 6, name = "system_version")
        private String systemVersion; // 系统版本
        @Index(value = 7, name = "address")
        private String address; // 节点服务地址
        @Index(value = 8, name = "description")
        private String description; // 节点说明 optional
        @Index(value = 9, name = "auth_type")
        private String authType; // 认证方式，枚举值：SHA256_RSA、 SHA256_ECDSA、CERT等
        @Index(value = 10, name = "auth_credential")
        private String authCredential; // 凭证内容：公钥值、证书内容等
        @Index(value = 11, name = "expired_time")
        private long expiredTime; // 合约过期时间 optional
        @Index(value = 12, name = "status")
        private String status; // APPLIED已申请 APPROVED已授权 REJECTED已拒绝 TERMINATED已 解除
        @Index(value = 13, name = "contract_id")
        private String contractId;
    }

    @Data
    class IncAck implements Serializable {

        private static final long serialVersionUID = 7003417744216186566L;
        @Index(value = 1, name = "contract_id")
        private String contractId;
        @Index(value = 2, name = "status")
        private String status; // APPLIED已申请 APPROVED已授权 REJECTED已拒绝 TERMINATED已 解除
        @Index(value = 3, name = "auth_type")
        private String authType; // 认证方式，枚举值：SHA256_RSA、 SHA256_ECDSA、CERT等
        @Index(value = 4, name = "auth_credential")
        private String authCredential; // 凭证内容：公钥值、证书内容等
        @Index(value = 5, name = "expired_time")
        private long expiredTime; // 合约过期时间 optional
    }

    @Data
    class IncOption implements Serializable {

        private static final long serialVersionUID = 7885756365042387199L;
    }
}
