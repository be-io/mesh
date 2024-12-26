/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package main

import (
	"github.com/opendatav/mesh/client/golang/mpc"
	"github.com/opendatav/mesh/client/golang/plugin"
	_ "github.com/opendatav/mesh/cmd"
)

//go:generate go run client/golang/dyn/mpc.go -i plugin -e plugin/raft -m github.com/opendatav/mesh

func main() {
	container := plugin.LoadC("mesh")
	container.Start(mpc.Context())
}
