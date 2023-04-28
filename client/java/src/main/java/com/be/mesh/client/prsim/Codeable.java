/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.prsim;

/**
 * See https://blue-elephant.yuque.com/blue-elephant/dg3hry/srnhhr.
 *
 * <pre>
 *
 * 1.5版本后，算子、算法、业务应用均往该规范靠.
 *
 * 错误码格式：
 *     E      0       00       000         0000
 *    前缀  异常等级  异常类型  模块或应用    业务码含义
 *
 * 执行成功通用的错误码：E0000000000.
 *
 * 总共10位错误码:
 * 1. 1位是等级，为0的等级直接对客
 * 2. 23位是异常类型
 * 3. 456位是模块或应用标志
 * 4. 78910位是自定义的业务码含义.
 *
 * 错误message可以是开发便于排查的，在网关层转化为产品定义的message，可配置.
 *
 * 错误码	模块
 * 000	   mesh
 * 001	   omega
 * 002	   asset
 * 003	   theta
 * 004	   cube
 * 005	   edge
 * 006     tensor
 * 007     base
 * </pre>
 *
 * @author coyzeng@gmail.com
 */
public interface Codeable {

    /**
     * Get the code.
     */
    String getCode();

    /**
     * Get the message.
     */
    String getMessage();

    /**
     * Error code is matchable.
     *
     * @param codeable code interface
     * @return true is matched
     */
    default boolean is(Codeable codeable) {
        return is(codeable.getCode());
    }

    /**
     * Error code is matchable.
     *
     * @param code code
     * @return true is matched
     */
    default boolean is(String code) {
        return null != code && code.equals(getCode());
    }

    /**
     * Return the code level.
     *
     * @return code level
     */
    default int level() {
        return getCode().length() > 1 ? getCode().charAt(1) : 9;
    }

    /**
     * Return the code kind.
     *
     * @return code kind
     */
    default String kind() {
        return getCode().length() > 3 ? getCode().substring(2, 4) : "99";
    }

    /**
     * Return the code group.
     *
     * @return code group
     */
    default String group() {
        return getCode().length() > 6 ? getCode().substring(4, 7) : "999";
    }
}
