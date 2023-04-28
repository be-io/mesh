/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

type Timeout struct {
	TaskId   string  `index:"0" json:"task_id" xml:"task_id" yaml:"task_id"`
	Binding  *Topic  `index:"5" json:"binding" xml:"binding" yaml:"binding"`
	Status   int64   `index:"10" json:"status" xml:"status" yaml:"status"`
	CreateAt int64   `index:"15" json:"create_at" xml:"create_at" yaml:"create_at"`
	InvokeAt int64   `index:"20" json:"invoke_at" xml:"invoke_at" yaml:"invoke_at"`
	Entity   *Entity `index:"25" json:"entity" xml:"entity" yaml:"entity"`
}
