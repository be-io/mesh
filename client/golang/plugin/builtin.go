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
	GDB         Name = "cayley"
	Metabase    Name = "metabase"
	NSQ         Name = "nsq"
	PRSIM       Name = "prsim"
	KMS         Name = "kms"
	Panel       Name = "panel"
	Redis       Name = "redis"
	Cache       Name = "cache"
	Raft        Name = "raft"
)
