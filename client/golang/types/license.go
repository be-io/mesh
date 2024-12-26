/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

import (
	"bytes"
	"github.com/opendatav/mesh/client/golang/cause"
	"github.com/opendatav/mesh/client/golang/codec"
	"sort"
)

const (
	LicenseCPUID     LicenseFactor = "CPUID"
	LicenseMAC       LicenseFactor = "MAC"
	LicenseKey       LicenseFactor = "KEY"
	LicenseSignature               = "signature"
	LicenseVersion                 = "1.0.0"
)

type LicenseFactor string

type License struct {
	Version        string   `index:"0" json:"version" xml:"version" yaml:"version" comment:"License version"`
	Level          int64    `index:"5" json:"level" xml:"level" yaml:"level" comment:"License level"`
	Name           string   `index:"10" json:"name" xml:"name" yaml:"name" comment:"License name"`
	CreateBy       string   `index:"15" json:"create_by" xml:"create_by" yaml:"create_by" comment:"License creator name"`
	CreateAt       int64    `index:"20" json:"create_at" xml:"create_at" yaml:"create_at" comment:"License create time"`
	ActiveAt       int64    `index:"25" json:"active_at" xml:"active_at" yaml:"active_at" comment:"License create by"`
	Factors        []string `index:"30" json:"factors" xml:"factors" yaml:"factors" comment:"License factors"`
	Signature      string   `index:"35" json:"signature" xml:"signature" yaml:"signature" comment:"License signature"`
	NodeId         string   `index:"40" json:"node_id" xml:"node_id" yaml:"node_id" comment:"License node identity"`
	InstId         string   `index:"45" json:"inst_id" xml:"inst_id" yaml:"inst_id" comment:"License institution identity"`
	Server         string   `index:"50" json:"server" xml:"server" yaml:"server" comment:"License server"`
	Crt            string   `index:"55" json:"crt" xml:"crt" yaml:"crt" comment:"License certification"`
	Group          []string `index:"60" json:"group" xml:"group" yaml:"group" comment:"License group"`
	Replicas       int64    `index:"65" json:"replicas" xml:"replicas" yaml:"replicas" comment:"License max replicas"`
	MaxCooperators int64    `index:"70" json:"max_cooperators" xml:"max_cooperators" yaml:"max_cooperators" comment:"License max cooperators"`
	MaxTenants     int64    `index:"75" json:"max_tenants" xml:"max_tenants" yaml:"max_tenants" comment:"License max tenants"`
	MaxUsers       int64    `index:"80" json:"max_users" xml:"max_users" yaml:"max_users" comment:"License max users"`
	MaxMills       int64    `index:"85" json:"max_mills" xml:"max_mills" yaml:"max_mills" comment:"License max mills"`
	WhiteURN       []string `index:"90" json:"white_urns" xml:"white_urns" yaml:"white_urns" comment:"License white urns"`
	BlackURN       []string `index:"95" json:"black_urns" xml:"black_urns" yaml:"black_urns" comment:"License black urns"`
	SuperURN       []string `index:"100" json:"super_urns" xml:"super_urns" yaml:"super_urns" comment:"License supervise urns"`
}

func (that *License) SortedMapString(codec codec.Codec) ([]byte, error) {
	buff, err := codec.Encode(that)
	if nil != err {
		return nil, cause.Error(err)
	}
	var kv map[string]interface{}
	if _, err = codec.Decode(buff, &kv); nil != err {
		return nil, cause.Error(err)
	}
	var keys []string
	for key, value := range kv {
		if nil == value {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	signature := &bytes.Buffer{}
	for _, key := range keys {
		if LicenseSignature == key || nil == kv[key] {
			continue
		}
		vff, err := codec.Encode(kv[key])
		if nil != err {
			return nil, cause.Error(err)
		}
		signature.Write(vff.Bytes())
	}
	return signature.Bytes(), nil
}
