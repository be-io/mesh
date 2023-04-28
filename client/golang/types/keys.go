/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

const (
	RootCaPrivateKey KeyKind = 1 << iota // 已连通
	RootCaPublicKey
	RootCaCrtKey
	IssuePrivateKey
	IssuePublicKey
	IssueCrtKey
)

type KeyKind int

type Keys struct {
	CNO     string  `index:"0" json:"cno" xml:"cno" yaml:"cno"`
	PNO     string  `index:"5" json:"pno" xml:"pno" yaml:"pno"`
	KNO     string  `index:"10" json:"kno" xml:"kno" yaml:"kno"`
	Kind    KeyKind `index:"15" json:"kind" xml:"kind" yaml:"kind"`
	Csr     string  `index:"20" json:"csr" xml:"csr" yaml:"csr"`
	Key     string  `index:"25" json:"key" xml:"key" yaml:"key"`
	Status  int     `index:"30" json:"status" xml:"status" yaml:"status"`
	Version int     `index:"35" json:"version" xml:"version" yaml:"version"`
}

type KeyCsr struct {
	CNO      string   `index:"0" json:"cno" xml:"cno" yaml:"cno"`
	PNO      string   `index:"5" json:"pno" xml:"pno" yaml:"pno"`
	Domain   string   `index:"10" json:"domain" xml:"domain" yaml:"domain"`
	Subject  string   `index:"15" json:"subject" xml:"subject" yaml:"subject"`
	Length   int      `index:"20" json:"length" xml:"length" yaml:"length"`
	ExpireAt Time     `index:"25" json:"expire_at" xml:"expire_at" yaml:"expire_at"`
	Mail     string   `index:"30" json:"mail" xml:"mail" yaml:"mail"`
	IsCA     bool     `index:"35" json:"is_ca" xml:"is_ca" yaml:"is_ca"`
	CaCert   string   `index:"40" json:"ca_cert" xml:"ca_cert" yaml:"ca_cert"`
	CaKey    string   `index:"45" json:"ca_key" xml:"ca_key" yaml:"ca_key"`
	IPs      []string `index:"50" json:"ips" xml:"ips" yaml:"ips"`
}

type KeysSet []*Keys

func (that KeysSet) Get(kind KeyKind) string {
	for _, key := range that {
		if key.Kind == kind {
			return key.Key
		}
	}
	return ""
}
