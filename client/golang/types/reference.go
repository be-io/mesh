/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

type Reference struct {
	URN       string `index:"0" json:"urn" xml:"urn" yaml:"urn" comment:""`
	Namespace string `index:"5" json:"namespace" xml:"namespace" yaml:"namespace" comment:""`
	Name      string `index:"10" json:"name" xml:"name" yaml:"name" comment:""`
	Version   string `index:"15" json:"version" xml:"version" yaml:"version" comment:""`
	Proto     string `index:"20" json:"proto" xml:"proto" yaml:"proto" comment:""`
	Codec     string `index:"25" json:"codec" xml:"codec" yaml:"codec" comment:""`
	Flags     int64  `index:"30" json:"flags" xml:"flags" yaml:"flags" comment:"Service flag 1 asyncable 2 encrypt 4 communal"`
	Timeout   int64  `index:"35" json:"timeout" xml:"timeout" yaml:"timeout" comment:""`
	Retries   int    `index:"40" json:"retries" xml:"retries" yaml:"retries" comment:""`
	Node      string `index:"45" json:"node" xml:"node" yaml:"node" comment:""`
	Inst      string `index:"50" json:"inst" xml:"inst" yaml:"inst" comment:""`
	Zone      string `index:"55" json:"zone" xml:"zone" yaml:"zone" comment:""`
	Cluster   string `index:"60" json:"cluster" xml:"cluster" yaml:"cluster" comment:""`
	Cell      string `index:"65" json:"cell" xml:"cell" yaml:"cell" comment:""`
	Group     string `index:"70" json:"group" xml:"group" yaml:"group" comment:""`
	Address   string `index:"75" json:"address" xml:"address" yaml:"address" comment:""`
}

func (that *Reference) GetURN() string {
	return that.URN
}

func (that *Reference) GetProto() string {
	return that.Proto
}

func (that *Reference) GetCodec() string {
	return that.Codec
}

func (that *Reference) GetTimeout() int64 {
	return that.Timeout
}

func (that *Reference) GetRetries() int {
	return that.Retries
}
