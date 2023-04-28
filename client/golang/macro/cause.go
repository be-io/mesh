/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package macro

type Cause struct {
	Name string `index:"0" json:"name" xml:"name" yaml:"name" comment:"Cause name"`
	Pos  string `index:"5" json:"pos" xml:"pos" yaml:"pos" comment:"Cause position"`
	Text string `index:"10" json:"text" xml:"text" yaml:"text" comment:"Cause descriptor"`
	Buff []byte `index:"15" json:"buff" xml:"buff" yaml:"buff" comment:"Cause stack"`
}

func Errors(err error) *Cause {
	return &Cause{
		Name: err.Error(),
		Pos:  "",
		Text: err.Error(),
	}
}
