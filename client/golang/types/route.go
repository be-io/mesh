/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

import (
	"context"
	"github.com/opendatav/mesh/client/golang/codec"
	"github.com/opendatav/mesh/client/golang/log"
	"github.com/opendatav/mesh/client/golang/macro"
	"net/url"
	"strings"
	"time"
)

// EdgeStatus 1 连通 2 组网 4 手动组网 8 自动发现 16 入网申请 32 逻辑删除 64 虚拟组网
type EdgeStatus int

const (
	Connected     EdgeStatus = 1 << iota // 已连通
	Weaved                               // 已组网
	FromMan                              // 手动组网
	FromDiscovery                        // 自动发现
	FromWeaving                          // 入网申请
	Disabled                             // 逻辑禁用
	Removed                              // 逻辑删除
	Approving                            // 入网审批
	Virtual                              // 虚拟组网
	Reject                               // 入网被拒绝 For UnionBank PPC
)

func (that EdgeStatus) Is(code int32) bool {
	return (int(that) & int(code)) == int(that)
}

func (that EdgeStatus) Or(code int32) int32 {
	return int32(that) | code
}

func (that EdgeStatus) OrElse(prev int32, code int32) int32 {
	if that.Is(prev) {
		return int32(that) | code
	}
	return code
}

func (that EdgeStatus) Not(code int32) int32 {
	return code & (^int32(that))
}

// Route is the node edge model
type Route struct {
	NodeId      string `index:"0" json:"node_id" xml:"node_id" yaml:"node_id" comment:"节点编号"`
	InstId      string `index:"5" json:"inst_id" xml:"inst_id" yaml:"inst_id" comment:"机构编号"`
	Name        string `index:"10" json:"name" xml:"name" yaml:"name" comment:"节点名称"`
	InstName    string `index:"15" json:"inst_name" xml:"inst_name" yaml:"inst_name" comment:"机构名称"`
	Address     string `index:"20" json:"address" xml:"address" yaml:"address" comment:"节点地址"`
	Describe    string `index:"25" json:"describe" xml:"describe" yaml:"describe" comment:"节点描述"`
	HostRoot    string `index:"30" json:"host_root" xml:"host_root" yaml:"host_root" comment:"Host root certifications"`
	HostKey     string `index:"35" json:"host_key" xml:"host_key" yaml:"host_key" comment:"Host private certifications key"`
	HostCrt     string `index:"40" json:"host_crt" xml:"host_crt" yaml:"host_crt" comment:"Host certification"`
	GuestRoot   string `index:"45" json:"guest_root" xml:"guest_root" yaml:"guest_root" comment:"Guest root certifications"`
	GuestKey    string `index:"50" json:"guest_key" xml:"guest_key" yaml:"guest_key" comment:"Guest private certifications key"`
	GuestCrt    string `index:"55" json:"guest_crt" xml:"guest_crt" yaml:"guest_crt" comment:"Guest certification"`
	Status      int32  `index:"60" json:"status" xml:"status" yaml:"status" comment:"1 连通 2 组网 4 手动组网 8自动发现 16入网申请"`
	Version     int32  `index:"65" json:"version" xml:"version" yaml:"version" comment:"版本"`
	AuthCode    string `index:"70" json:"auth_code" xml:"auth_code" yaml:"auth_code" comment:"授权码"`
	ExpireAt    int64  `index:"75" json:"expire_at" xml:"expire_at" yaml:"expire_at" comment:"生效截止时间"`
	Extra       string `index:"80" json:"extra" xml:"extra" yaml:"extra" comment:"补充信息说明"`
	CreateAt    string `index:"85" json:"create_at" xml:"create_at" yaml:"create_at" comment:"创建时间"`
	UpdateAt    string `index:"90" json:"update_at" xml:"update_at" yaml:"update_at" comment:"更新时间"`
	CreateBy    string `index:"95" json:"create_by" xml:"create_by" yaml:"create_by" comment:"创建人"`
	UpdateBy    string `index:"100" json:"update_by" xml:"update_by" yaml:"update_by" comment:"更新人"`
	Group       string `index:"105" json:"group" xml:"group" yaml:"group" comment:"联盟节点信息"`
	Upstream    int64  `index:"110" json:"upstream" xml:"upstream" yaml:"upstream" comment:"上行流量"`
	Downstream  int64  `index:"115" json:"downstream" xml:"downstream" yaml:"downstream" comment:"下行流量"`
	StaticIP    string `index:"120" json:"static_ip" xml:"static_ip" yaml:"static_ip" comment:"静态出口IP"`
	Proxy       string `index:"125" json:"proxy" xml:"proxy" yaml:"proxy" comment:"Proxy endpoint in transport"`
	Concurrency int64  `index:"130" json:"concurrency" xml:"concurrency" yaml:"concurrency" comment:"MPC concurrency"`
}

func (that *Route) GetCertificate(ctx context.Context) *RouteCertificate {
	return &RouteCertificate{
		HostRoot:  that.HostRoot,
		HostKey:   that.HostKey,
		HostCrt:   that.HostCrt,
		GuestRoot: that.GuestRoot,
		GuestKey:  that.GuestKey,
		GuestCrt:  that.GuestCrt,
	}
}

func (that *Route) Copy() *Route {
	return &Route{
		NodeId:      that.NodeId,
		InstId:      that.InstId,
		Name:        that.Name,
		InstName:    that.InstName,
		Address:     that.Address,
		Describe:    that.Describe,
		HostRoot:    that.HostRoot,
		HostKey:     that.HostKey,
		HostCrt:     that.HostCrt,
		GuestRoot:   that.GuestRoot,
		GuestKey:    that.GuestKey,
		GuestCrt:    that.GuestCrt,
		Status:      that.Status,
		Version:     that.Version,
		AuthCode:    that.AuthCode,
		ExpireAt:    that.ExpireAt,
		Extra:       that.Extra,
		CreateAt:    that.CreateAt,
		UpdateAt:    that.UpdateAt,
		CreateBy:    that.CreateBy,
		UpdateBy:    that.UpdateBy,
		Group:       that.Group,
		Upstream:    that.Upstream,
		Downstream:  that.Downstream,
		StaticIP:    that.StaticIP,
		Proxy:       that.Proxy,
		Concurrency: that.Concurrency,
	}
}

func (that *Route) ID(ctx context.Context) *NodeID {
	nodeId, err := FromNodeID(that.NodeId)
	if nil == err {
		return nodeId
	}
	log.Debug(ctx, err.Error())
	return &NodeID{SEQ: that.NodeId}
}

func (that *Route) URC() URC {
	return URC(that.Address)
}

// RouteCertificate is point to point certifications
type RouteCertificate struct {
	HostRoot  string `index:"0" json:"host_root" xml:"host_root" yaml:"host_root" comment:"Host root certifications"`
	HostKey   string `index:"5" json:"host_key" xml:"host_key" yaml:"host_key" comment:"Host private certifications key"`
	HostCrt   string `index:"10" json:"host_crt" xml:"host_crt" yaml:"host_crt" comment:"Host certification"`
	GuestRoot string `index:"15" json:"guest_root" xml:"guest_root" yaml:"guest_root" comment:"Guest root certifications"`
	GuestKey  string `index:"20" json:"guest_key" xml:"guest_key" yaml:"guest_key" comment:"Guest private certifications key"`
	GuestCrt  string `index:"25" json:"guest_crt" xml:"guest_crt" yaml:"guest_crt" comment:"Guest certification"`
}

func (that *RouteCertificate) Override0(nc *Route) *RouteCertificate {
	cw := new(RouteCertificate)
	if nil != nc && "" != nc.HostRoot {
		cw.HostRoot = nc.HostRoot
	} else {
		cw.HostRoot = that.HostRoot
	}
	if nil != nc && "" != nc.GuestRoot {
		cw.GuestRoot = nc.GuestRoot
	} else {
		cw.GuestRoot = that.GuestRoot
	}
	if nil != nc && "" != nc.HostKey && "" != nc.HostCrt {
		cw.HostKey = nc.HostKey
		cw.HostCrt = nc.HostCrt
	} else {
		cw.HostKey = that.HostKey
		cw.HostCrt = that.HostCrt
	}
	if nil != nc && "" != nc.GuestKey && "" != nc.GuestCrt {
		cw.GuestKey = nc.GuestKey
		cw.GuestCrt = nc.GuestCrt
	} else {
		cw.GuestKey = that.GuestKey
		cw.GuestCrt = that.GuestCrt
	}
	return cw
}

func (that *RouteCertificate) Override(nc *RouteCertificate) *RouteCertificate {
	cw := new(RouteCertificate)
	if nil != nc && "" != nc.HostRoot {
		cw.HostRoot = nc.HostRoot
	} else {
		cw.HostRoot = that.HostRoot
	}
	if nil != nc && "" != nc.GuestRoot {
		cw.GuestRoot = nc.GuestRoot
	} else {
		cw.GuestRoot = that.GuestRoot
	}
	if nil != nc && "" != nc.HostKey && "" != nc.HostCrt {
		cw.HostKey = nc.HostKey
		cw.HostCrt = nc.HostCrt
	} else {
		cw.HostKey = that.HostKey
		cw.HostCrt = that.HostCrt
	}
	if nil != nc && "" != nc.GuestKey && "" != nc.GuestCrt {
		cw.GuestKey = nc.GuestKey
		cw.GuestCrt = nc.GuestCrt
	} else {
		cw.GuestKey = that.GuestKey
		cw.GuestCrt = that.GuestCrt
	}
	return cw
}

func (that *RouteCertificate) Decode(ctx context.Context, certification string, decoder codec.Codec) {
	if _, err := decoder.DecodeString(certification, that); nil != err {
		log.Error(ctx, "Decode certification %s with unexpected cause, %s", certification, err.Error())
	}
}

func (that *RouteCertificate) Encode(ctx context.Context, encoder codec.Codec) string {
	if buff, err := encoder.Encode(that); nil != err {
		log.Error(ctx, "Encode certification with unexpected cause, %s", err.Error())
		return ""
	} else {
		return buff.String()
	}
}

type RouteAttribute struct {
	Upstream    int64  `index:"0" json:"upstream" xml:"upstream" yaml:"upstream" comment:"上行流量"`
	Downstream  int64  `index:"5" json:"downstream" xml:"downstream" yaml:"downstream" comment:"下行流量"`
	StaticIP    string `index:"10" json:"static_ip" xml:"static_ip" yaml:"static_ip" comment:"静态出口IP"`
	DirectAddr  string `index:"15" json:"direct_addr" xml:"direct_addr" yaml:"direct_addr" comment:"Direct address in old arch."`
	Proxy       string `index:"30" json:"proxy" xml:"proxy" yaml:"proxy" comment:"Proxy endpoint in transport"`
	Concurrency int64  `index:"35" json:"concurrency" xml:"concurrency" yaml:"concurrency" comment:"MPC concurrency"`
	TPC         string `index:"40" json:"tpc" xml:"tpc" yaml:"tpc" comment:"MPC tpc"`
	Credential  string `index:"45" json:"credential" xml:"credential" yaml:"credential" comment:"Cert credential"`
}

func (that *RouteAttribute) Encode(ctx context.Context, route *Route, encoder codec.Codec) (string, error) {
	that.Upstream = route.Upstream
	that.Downstream = route.Downstream
	that.StaticIP = route.StaticIP
	that.Proxy = route.Proxy
	that.Concurrency = route.Concurrency
	if strings.HasPrefix(route.Extra, "{") {
		if _, err := encoder.DecodeString(route.Extra, that); nil != err {
			log.Error(ctx, "Encode frontend route extra, %s", err.Error())
		}
	} else {
		that.DirectAddr = route.Extra
	}
	return encoder.EncodeString(that)
}

func (that *RouteAttribute) Decode(ctx context.Context, extra string, encoder codec.Codec) *RouteAttribute {
	if "" == extra || !strings.HasPrefix(extra, "{") {
		that.DirectAddr = extra
		return that
	}
	if _, err := encoder.DecodeString(extra, that); nil != err {
		log.Error(ctx, err.Error())
	}
	return that
}

func (that *RouteAttribute) Compat(ctx context.Context, encoder codec.Codec) string {
	if "" == that.Proxy {
		return that.DirectAddr
	}
	txt, err := encoder.EncodeString(map[string]string{"direct_addr": that.DirectAddr, "proxy": that.Proxy})
	if nil != err {
		log.Error(ctx, "Compat frontend route extra, %s", err.Error())
		return that.DirectAddr
	}
	return txt
}

type DirectRoute struct {
	Routes []*Route   `index:"0" json:"routes" xml:"routes" yaml:"routes" comment:"Routes"`
	Direct *Principal `index:"5" json:"direct" xml:"direct" yaml:"direct" comment:"Direct"`
}

type Server interface {
	URL() *url.URL
}

// VIP name:matcher=label:[host]
type VIP struct {
	Name    string   `json:"name"`
	Matcher string   `json:"matcher"`
	Label   string   `json:"label"`
	Hosts   []string `json:"hosts"`
}

func (that *VIP) Matches(k *VIPKey, s Server) bool {
	switch that.Name {
	case "ip":
		return k.IP == that.Matcher
	case "src_id":
		return k.SrcID == that.Matcher || k.SrcID == NodeSEQ(that.Matcher)
	case "dst_id":
		return k.DstID == that.Matcher || k.DstID == NodeSEQ(that.Matcher)
	case "version":
		return k.Version == that.Matcher
	case "gray":
		return k.Gray == that.Matcher || k.Gray == NodeSEQ(that.Matcher)
	case "all":
		return true
	default:
		return false
	}
}

type VIPKey struct {
	IP      string `json:"ip"`
	SrcID   string `json:"src_id"`
	DstID   string `json:"dst_id"`
	Version string `json:"version"`
	Gray    string `json:"gray"`
}

// Trick name:kind+Service=proto://address/pattern,pattern,
type Trick struct {
	Name     string   `json:"name"`
	Kind     string   `json:"kind"`
	Service  string   `json:"service"`
	Proto    string   `json:"proto"`
	Patterns []string `json:"patterns"`
	Address  string   `json:"address"`
	Plugins  string   `json:"plugins"`
}

func LocRoute(environ *Environ, instId string) *Route {
	return &Route{
		NodeId:      environ.NodeId,
		InstId:      instId,
		Name:        environ.InstName,
		InstName:    environ.InstName,
		Address:     macro.Runtime(),
		Describe:    "",
		HostRoot:    "",
		HostKey:     "",
		HostCrt:     "",
		GuestRoot:   "",
		GuestKey:    "",
		GuestCrt:    "",
		Status:      0,
		Version:     0,
		AuthCode:    "",
		ExpireAt:    time.Now().AddDate(10, 0, 0).UnixMilli(),
		Extra:       "",
		CreateAt:    "",
		UpdateAt:    "",
		CreateBy:    "",
		UpdateBy:    "",
		Group:       "",
		Upstream:    0,
		Downstream:  0,
		StaticIP:    "",
		Proxy:       "",
		Concurrency: 0,
	}
}
