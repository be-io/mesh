/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {mpi} from "@/macro/mpi";
import {index} from "@/macro/idx";
import {spi} from "@/macro/spi";
import type {Context} from "@/prsim/context";
import {Status} from "@/cause";


@spi("mesh")
export abstract class Sequence {

    /**
     * 生成全网唯一序列号
     *
     * @param ctx ctx
     * @param kind 序列号类型，各个业务保持唯一.
     * @param length length
     * @return 唯一序列号
     */
    @mpi("mesh.sequence.next", String)
    next(ctx: Context, @index(0, 'kind') kind: string, @index(1, 'length') length: number): Promise<string> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * 获取序列号段.
     *
     * @param ctx ctx
     * @param kind 类型
     * @param size 号段大小
     * @param length length
     * @return 返回序列号列表
     */
    @mpi("mesh.sequence.section", [Array, String])
    section(ctx: Context, @index(0, 'kind') kind: string, @index(1, 'size') size: number, @index(2, 'length') length: number): Promise<string[]> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

}