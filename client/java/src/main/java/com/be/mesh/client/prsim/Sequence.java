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

/**
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
public interface Sequence {

    default String next(@Index(0) String kind) {
        return this.next(kind, 8);
    }

    /**
     * 生成全网唯一序列号
     *
     * @param kind   序列号类型，各个业务保持唯一.
     * @param length 长度
     * @return 唯一序列号
     */
    @MPI("mesh.sequence.next")
    String next(@Index(0) String kind, int length);

    /**
     * 获取序列号段.
     *
     * @param kind   类型
     * @param size   号段大小
     * @param length 长度
     * @return 返回序列号列表
     */
    @MPI("mesh.sequence.section")
    String[] section(@Index(0) String kind, @Index(1) int size, int length);

}
