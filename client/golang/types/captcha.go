/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

type Captcha struct {
	MNO  string `index:"0" json:"mno" yaml:"mno" xml:"mno" comment:""`
	Kind string `index:"5" json:"kind" yaml:"kind" xml:"kind" comment:""`
	Mime []byte `index:"10" json:"mime" yaml:"mime" xml:"mime" comment:""`
	Text string `index:"15" json:"text" yaml:"text" xml:"text" comment:""`
}
