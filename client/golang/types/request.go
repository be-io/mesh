/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

type Request struct {
	Version string      `index:"0" json:"version,omitempty" yaml:"version,omitempty" xml:"version,omitempty"`
	Method  string      `index:"1" json:"method,omitempty" yaml:"method,omitempty" xml:"method,omitempty"`
	Content interface{} `index:"2" json:"content,omitempty" yaml:"content,omitempty" xml:"content,omitempty"`
}
