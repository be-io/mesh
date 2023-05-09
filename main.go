/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package main

import (
	"github.com/be-io/mesh/client/golang/mpc"
	"github.com/be-io/mesh/client/golang/plugin"
	_ "github.com/be-io/mesh/cmd"
)

//go:generate go run client/golang/dyn/mpc.go -i plugin -e plugin/raft -m github.com/be-io/mesh

func main() {
	container := plugin.LoadC("mesh")
	container.Start(mpc.Context())
}
