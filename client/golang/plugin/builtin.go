/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package plugin

type Name string

const (
	ServerCar   Name = "server-car"
	NodeCar     Name = "node-car"
	PanelCar    Name = "panel-car"
	WorkloadCar Name = "operator-car"
	ProxyCar    Name = "proxy-car"
	// Proxy -------------------------
	Proxy       Name = "proxy"
	Http        Name = "http"
	DNS         Name = "dns"
	GDB         Name = "gdb"
	Wasm        Name = "wasm"
	Metabase    Name = "metabase"
	NSQ         Name = "nsq"
	PRSIM       Name = "prsim"
	MYSQL       Name = "mysql"
	KMS         Name = "kms"
	OAuth2      Name = "oauth2"
	Graphics    Name = "graphics"
	Panel       Name = "panel"
	GoProxy     Name = "goproxy"
	Redis       Name = "redis"
	Packet      Name = "packet"
	Cache       Name = "cache"
	Raft        Name = "raft"
	Telemetry   Name = "telemetry"
	Shadowsocks Name = "shadowsocks"
)
