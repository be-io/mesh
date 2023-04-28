/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

const (
	EXPRESSION = "EXPRESSION"
	VALUE      = "VALUE"
	SCRIPT     = "SCRIPT"
)

type Script struct {
	Code       string            `index:"0" xml:"code" json:"code" yaml:"code" comment:"Script code"`
	Name       string            `index:"5" xml:"name" json:"name" yaml:"name" comment:"Script name"`
	Desc       string            `index:"10" xml:"desc" json:"desc" yaml:"desc" comment:"Script desc"`
	Kind       string            `index:"15" xml:"kind" json:"kind" yaml:"kind" comment:"Script kind"`
	Expr       string            `index:"20" xml:"expr" json:"expr" yaml:"expr" comment:"Script expression"`
	Attachment map[string]string `index:"25" xml:"attachment" json:"attachment" yaml:"attachment" comment:"Script attachment"`
}
