/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

type Service struct {
	URN       string            `index:"0" json:"urn" xml:"urn" yaml:"urn"`                   // URN is service urn
	Namespace string            `index:"5" json:"namespace" xml:"namespace" yaml:"namespace"` // Namespace Service topic
	Name      string            `index:"10" json:"name" xml:"name" yaml:"name"`
	Version   string            `index:"15" json:"version" xml:"version" yaml:"version"`
	Proto     string            `index:"20" json:"proto" xml:"proto" yaml:"proto"`
	Codec     string            `index:"25" json:"codec" xml:"codec" yaml:"codec"`
	Flags     int64             `index:"30" json:"flags" xml:"flags" yaml:"flags" comment:"Service flag 1 asyncable 2 encrypt 4 communal"`
	Timeout   int64             `index:"35" json:"timeout" xml:"timeout" yaml:"timeout"`
	Retries   int               `index:"40" json:"retries" xml:"retries" yaml:"retries"`
	Node      string            `index:"45" json:"node" xml:"node" yaml:"node"`
	Inst      string            `index:"50" json:"inst" xml:"inst" yaml:"inst"`
	Zone      string            `index:"55" json:"zone" xml:"zone" yaml:"zone"`
	Cluster   string            `index:"60" json:"cluster" xml:"cluster" yaml:"cluster"`
	Cell      string            `index:"65" json:"cell" xml:"cell" yaml:"cell"`
	Group     string            `index:"70" json:"group" xml:"group" yaml:"group"`
	Sets      string            `index:"75" json:"sets" xml:"sets" yaml:"sets"`
	Address   string            `index:"80" json:"address" xml:"address" yaml:"address"`
	Kind      string            `index:"85" json:"kind" xml:"kind" yaml:"kind"`
	Lang      string            `index:"90" json:"lang" xml:"lang" yaml:"lang"`
	Attrs     map[string]string `index:"95" json:"attrs" xml:"attrs" yaml:"attrs"` // Attrs Service arguments.
}

func (that *Service) GetURN() string {
	return that.URN
}

func (that *Service) GetProto() string {
	return that.Proto
}

func (that *Service) GetCodec() string {
	return that.Codec
}

func (that *Service) GetTimeout() int64 {
	return that.Timeout
}

func (that *Service) GetRetries() int {
	return that.Retries
}
