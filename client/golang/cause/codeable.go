/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package cause

// Codeable
/**
 *
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
 * 总共10位错误码:
 * 1. 1位是等级，为0的等级直接对客
 * 2. 23位是异常类型
 * 3. 456位是模块或应用标志
 * 4. 78910位是自定义的业务码含义.
 *
 * 错误message可以是开发便于排查的，在网关层转化为产品定义的message，可配置.
 *
 * 错误码	模块
 * 000	    mesh
 * 001	    omega
 * 002	    asset
 * 003	    theta
 * 004	    cube
 * 005	    edge
 * 006      tensor
 * 007      base
 * 008      ark
 * 009      jewel
 * 010      pandora
 * 011      graph
 * </pre>
 *
 * @author coyzeng@gmail.com
 */
type Codeable interface {

	// GetCode Get the code.
	GetCode() string

	// GetMessage Get the message.
	GetMessage() string
}
