/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package macro

import (
	"os"
	"strconv"
)

func Env(backoff string, keys ...string) string {
	for _, key := range keys {
		if "" != os.Getenv(key) {
			return os.Getenv(key)
		}
	}
	return backoff
}

func Name() string {
	return Env("mesh", "MESH_NAME", "MESH-NAME", "MESH.NAME", "mesh_name", "mesh-name", "mesh.name")
}

func Direct() string {
	return Env("", "MESH_DIRECT", "MESH-DIRECT", "mesh.direct", "mesh_direct")
}

func Subset() string {
	return Env("", "MESH_SUBSET", "MESH-SUBSET", "mesh.subset", "mesh_subset")
}

func Mode() string {
	return Env("", "MESH_MODE", "MESH-MODE", "mesh.mode", "mesh_mode")
}

func Proxy() string {
	return Env("", "MESH_PROXY", "mesh_proxy", "mesh.proxy", "MESH-PROXY")
}

func Address() string {
	return Env("127.0.0.1", "MESH_ADDRESS", "MESH-ADDRESS", "mesh_address", "mesh.address", "mesh-address")
}

func Runtime() string {
	return Env("", "MESH_RUNTIME", "MESH-RUNTIME", "mesh_runtime", "mesh.runtime", "mesh-runtime")
}

func LHome() string {
	return Env("", "MESH_LOG", "MESH-LOG", "MESH.LOG", "mesh_log", "mesh-log", "mesh.log")
}

func MDC() string {
	return Env("", "MESH_MDC", "MESH-MDC", "MESH.MDC", "mesh_mdc", "mesh-mdc", "mesh.mdc")
}

func SPA() string {
	return Env("", "MESH_SPA", "MESH-SPA", "MESH.SPA", "mesh_spa", "mesh-spa", "mesh.spa")
}

type Modes int64

var (
	Disable         = Modes(1)
	Failfast        = Modes(2)
	Nolog           = Modes(4)
	JsonLogFormat   = Modes(8)
	RCache          = Modes(16)
	PHeader         = Modes(32)
	Metrics         = Modes(64)
	RLog            = Modes(128)
	MGrpc           = Modes(256)
	PermitCirculate = Modes(512)
	NoStdColor      = Modes(1024)
	NoTeeReport     = Modes(2048)
	DisableTee      = Modes(4096)
	NoStaticFile    = Modes(8192)
	OpenTelemetry   = Modes(16384)
	EableSPIFirst   = Modes(32768)
)

func (that Modes) Enable() bool {
	return mode.Get().Match(that)
}

func (that Modes) Match(mode Modes) bool {
	return (int64(that) & int64(mode)) == int64(mode)
}

var mode = new(Once[Modes]).With(func() Modes {
	m := Mode()
	if "" == m {
		return Failfast | JsonLogFormat
	}
	n, err := strconv.ParseInt(m, 10, 64)
	if nil != err {
		return Failfast | JsonLogFormat
	}
	return Modes(n)
})

func SetMode(m Modes) {
	mode.Set(m)
}

func WithMode(m Modes) {
	mode.Set(mode.Get() | m)
}
