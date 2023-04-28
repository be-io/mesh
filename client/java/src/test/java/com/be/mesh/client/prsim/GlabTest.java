/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.prsim;

import com.be.mesh.client.annotate.Index;
import com.be.mesh.client.annotate.MPI;
import com.be.mesh.client.mpc.ServiceProxy;
import lombok.Data;
import lombok.extern.slf4j.Slf4j;
import org.testng.annotations.Test;

import java.util.Date;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
public class GlabTest {

    @Test
    public void testGlab() {
        System.setProperty("mesh.address", "10.99.51.33:570");
        GlabDataProductV2Dto.DataProductQueryDetailParam param = new GlabDataProductV2Dto.DataProductQueryDetailParam();
        param.setAuthorityCode("641fbcc48328406291970c48ffccf4ba");
        param.setInstId("JG0100000100000000");
        I i = ServiceProxy.proxy(I.class);
        GlabDataProductV2Dto.DataProductQueryDetail ret = i.queryDataProductDetail(param);
        log.info("{}", ret);
    }

    public static interface I {
        @MPI(value = "gman.management.dataProduct.detail")
        GlabDataProductV2Dto.DataProductQueryDetail queryDataProductDetail(@Index(value = 0, name = "request") GlabDataProductV2Dto.DataProductQueryDetailParam request);
    }

    public static final class GlabDataProductV2Dto {

        @Data
        public static class DataProductQueryDetail {
            /**
             * 产品id
             */
            private String dataProductId;
            /**
             * 产品名称
             */
            private String dataProductName;
            /**
             * 产品类型
             */
            private String dataProductType;
            /**
             * 产品描述
             */
            private String remark;
            /**
             *
             */
            private Date startTime;

            private Date endTime;
            /**
             * 授权时间
             */
            private Date createTime;
            /**
             * 算法配置
             */
            private String algorithmSetting;
            /**
             * 资产信息
             */
            private String dataAsset;
            /**
             * 这是啥
             */
            private String extInfo;
            /**
             * 数据方的机构信息
             */
            private String instId;
        }

        @Data
        public static class DataProductQueryDetailParam {
            /**
             * 授权方机构ID
             */
            private String instId;

            /**
             * 数据服务授权码
             */
            private String authorityCode;
        }
    }
}
