/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

type Institution struct {
	NodeId   string `index:"0" json:"node_id" yaml:"node_id" xml:"node_id" comment:""`
	InstId   string `index:"5" json:"inst_id" yaml:"inst_id" xml:"inst_id" comment:""`
	InstName string `index:"10" json:"inst_name" yaml:"inst_name" xml:"inst_name" comment:""`
	Status   int32  `index:"15" json:"status" yaml:"status" xml:"status" comment:""`
}

type GenInstReq struct {
	NodeId   string `index:"0" json:"node_id" yaml:"node_id" xml:"node_id" comment:""`
	InstType string `index:"5" json:"inst_type" yaml:"inst_type" xml:"inst_type" comment:""`
}
