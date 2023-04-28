/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

type Inbound struct {
	Arguments   []interface{}     `index:"0" xml:"arguments" json:"arguments" yaml:"arguments" comment:"Invoke parameters"`
	Attachments map[string]string `index:"1" json:"attachments" xml:"attachments" yaml:"attachments" comment:"Invoke attachments"`
}
