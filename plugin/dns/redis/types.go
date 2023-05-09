/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package redis

import "net"

type Zone struct {
	Name      string
	Locations map[string]struct{}
}

type Record struct {
	A     []A     `json:"a,omitempty"`
	AAAA  []AAAA  `json:"aaaa,omitempty"`
	TXT   []TXT   `json:"txt,omitempty"`
	CNAME []CNAME `json:"cname,omitempty"`
	NS    []NS    `json:"ns,omitempty"`
	MX    []MX    `json:"mx,omitempty"`
	SRV   []SRV   `json:"srv,omitempty"`
	CAA   []CAA   `json:"caa,omitempty"`
	SOA   SOA     `json:"soa,omitempty"`
}

type A struct {
	Ttl uint32 `json:"ttl,omitempty"`
	Ip  net.IP `json:"ip"`
}

type AAAA struct {
	Ttl uint32 `json:"ttl,omitempty"`
	Ip  net.IP `json:"ip"`
}

type TXT struct {
	Ttl  uint32 `json:"ttl,omitempty"`
	Text string `json:"text"`
}

type CNAME struct {
	Ttl  uint32 `json:"ttl,omitempty"`
	Host string `json:"host"`
}

type NS struct {
	Ttl  uint32 `json:"ttl,omitempty"`
	Host string `json:"host"`
}

type MX struct {
	Ttl        uint32 `json:"ttl,omitempty"`
	Host       string `json:"host"`
	Preference uint16 `json:"preference"`
}

type SRV struct {
	Ttl      uint32 `json:"ttl,omitempty"`
	Priority uint16 `json:"priority"`
	Weight   uint16 `json:"weight"`
	Port     uint16 `json:"port"`
	Target   string `json:"target"`
}

type SOA struct {
	Ttl     uint32 `json:"ttl,omitempty"`
	Ns      string `json:"ns"`
	MBox    string `json:"MBox"`
	Refresh uint32 `json:"refresh"`
	Retry   uint32 `json:"retry"`
	Expire  uint32 `json:"expire"`
	MinTtl  uint32 `json:"minttl"`
}

type CAA struct {
	Flag  uint8  `json:"flag"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}
