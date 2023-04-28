/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

type CommerceLicense struct {
	Cipher   string   `index:"0" json:"cipher" xml:"cipher" yaml:"cipher" comment:"License cipher"`
	Explain  *License `index:"5" json:"explain" xml:"explain" yaml:"explain" comment:"License explain"`
	CreateAt Time     `index:"10" json:"create_at" xml:"create_at" yaml:"create_at" comment:"License create time"`
}

type CommerceEnviron struct {
	Cipher  string   `index:"0" json:"cipher" xml:"cipher" yaml:"cipher" comment:"Node cipher"`
	Explain *Environ `index:"5" json:"explain" xml:"explain" yaml:"explain" comment:"Node explain"`
	NodeKey string   `index:"10" json:"node_key" xml:"node_key" yaml:"node_key" comment:"Node private key"`
}
